package service

import (
	"context"
	"errors"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/zmoony/pi-gateway/internal/domain"
	"github.com/zmoony/pi-gateway/internal/platform/nftables"
	"github.com/zmoony/pi-gateway/internal/store"
)

var (
	ErrFirewallRuleNotFound = errors.New("firewall rule not found")
	ErrFirewallApplyPending = errors.New("firewall apply already pending confirmation")
)

type FirewallService struct {
	Rules        store.FirewallRuleStore
	Manager      nftables.Manager
	PendingTTL   time.Duration

	mu           sync.Mutex
	pendingTimer *time.Timer
}

type UpsertFirewallRuleInput struct {
	Name           string `json:"name"`
	Kind           string `json:"kind"`
	TemplateKey    string `json:"templateKey"`
	Protocol       string `json:"protocol"`
	Port           int    `json:"port"`
	PortRangeStart int    `json:"portRangeStart"`
	PortRangeEnd   int    `json:"portRangeEnd"`
	SourceCIDR     string `json:"sourceCidr"`
	Enabled        bool   `json:"enabled"`
	Priority       int    `json:"priority"`
	Description    string `json:"description"`
}

func (s *FirewallService) ResumePending(ctx context.Context) error {
	state, err := s.Manager.LoadPendingState()
	if err != nil || state == nil || !state.Pending {
		return err
	}
	if state.ExpiresAt == nil || time.Now().After(*state.ExpiresAt) {
		return s.Manager.Rollback(ctx)
	}
	s.scheduleRollback(state)
	return nil
}

func (s *FirewallService) List(ctx context.Context) ([]domain.FirewallRule, error) {
	return s.Rules.List(ctx)
}

func (s *FirewallService) Create(ctx context.Context, input UpsertFirewallRuleInput) (*domain.FirewallRule, error) {
	item, err := s.buildRule(ctx, input, nil)
	if err != nil {
		return nil, err
	}
	if err := s.Rules.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *FirewallService) Update(ctx context.Context, id int64, input UpsertFirewallRuleInput) (*domain.FirewallRule, error) {
	current, err := s.Rules.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, ErrFirewallRuleNotFound
	}
	item, err := s.buildRule(ctx, input, current)
	if err != nil {
		return nil, err
	}
	if err := s.Rules.Update(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *FirewallService) Delete(ctx context.Context, id int64) error {
	current, err := s.Rules.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if current == nil {
		return ErrFirewallRuleNotFound
	}
	return s.Rules.Delete(ctx, id)
}

func (s *FirewallService) Preview(ctx context.Context) (string, error) {
	rules, err := s.Rules.List(ctx)
	if err != nil {
		return "", err
	}
	return s.Manager.RenderInputRules(rules), nil
}

func (s *FirewallService) Apply(ctx context.Context) (*domain.FirewallPendingState, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, err := s.Manager.LoadPendingState()
	if err != nil {
		return nil, "", err
	}
	if state != nil && state.Pending {
		return state, "", ErrFirewallApplyPending
	}

	rules, err := s.Rules.List(ctx)
	if err != nil {
		return nil, "", err
	}
	candidate := s.Manager.RenderInputRules(rules)
	rollbackText := s.Manager.RenderInputRules(nil)
	state, err = s.Manager.Apply(ctx, candidate, rollbackText, s.PendingTTL)
	if err != nil {
		return nil, candidate, err
	}
	s.scheduleRollback(state)
	return state, candidate, nil
}

func (s *FirewallService) Confirm(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.pendingTimer != nil {
		s.pendingTimer.Stop()
		s.pendingTimer = nil
	}
	return s.Manager.Confirm(ctx)
}

func (s *FirewallService) PendingState() (*domain.FirewallPendingState, error) {
	return s.Manager.LoadPendingState()
}

func (s *FirewallService) scheduleRollback(state *domain.FirewallPendingState) {
	if s.pendingTimer != nil {
		s.pendingTimer.Stop()
	}
	if state == nil || state.ExpiresAt == nil {
		return
	}
	delay := time.Until(*state.ExpiresAt)
	if delay < 0 {
		delay = 0
	}
	s.pendingTimer = time.AfterFunc(delay, func() {
		_ = s.Manager.Rollback(context.Background())
		s.mu.Lock()
		s.pendingTimer = nil
		s.mu.Unlock()
	})
}

func (s *FirewallService) buildRule(ctx context.Context, input UpsertFirewallRuleInput, current *domain.FirewallRule) (*domain.FirewallRule, error) {
	if err := validateFirewallInput(input); err != nil {
		return nil, err
	}
	priority := input.Priority
	if priority <= 0 {
		if current != nil {
			priority = current.Priority
		} else {
			next, err := s.Rules.NextPriority(ctx)
			if err != nil {
				return nil, err
			}
			priority = next
		}
	}

	item := &domain.FirewallRule{
		Name:           strings.TrimSpace(input.Name),
		Kind:           strings.ToLower(strings.TrimSpace(input.Kind)),
		TemplateKey:    strings.ToLower(strings.TrimSpace(input.TemplateKey)),
		Protocol:       strings.ToLower(strings.TrimSpace(input.Protocol)),
		Port:           input.Port,
		PortRangeStart: input.PortRangeStart,
		PortRangeEnd:   input.PortRangeEnd,
		SourceCIDR:     strings.TrimSpace(input.SourceCIDR),
		Action:         "accept",
		Enabled:        input.Enabled,
		Priority:       priority,
		Description:    strings.TrimSpace(input.Description),
	}
	if current != nil {
		item.ID = current.ID
		item.CreatedAt = current.CreatedAt
		item.UpdatedAt = current.UpdatedAt
	}
	return item, nil
}

func validateFirewallInput(input UpsertFirewallRuleInput) error {
	if strings.TrimSpace(input.Name) == "" {
		return errors.New("rule name is required")
	}
	kind := strings.ToLower(strings.TrimSpace(input.Kind))
	switch kind {
	case "template":
		switch strings.ToLower(strings.TrimSpace(input.TemplateKey)) {
		case "ssh", "http", "https", "wireguard":
		default:
			return errors.New("invalid template key")
		}
	case "custom":
		protocol := strings.ToLower(strings.TrimSpace(input.Protocol))
		switch protocol {
		case "tcp", "udp":
			if input.Port <= 0 && (input.PortRangeStart <= 0 || input.PortRangeEnd < input.PortRangeStart) {
				return errors.New("custom tcp/udp rules require a port or valid port range")
			}
		case "icmp", "icmpv6":
		default:
			return errors.New("invalid custom protocol")
		}
	default:
		return errors.New("rule kind must be template or custom")
	}

	if source := strings.TrimSpace(input.SourceCIDR); source != "" {
		if _, _, err := net.ParseCIDR(source); err != nil && net.ParseIP(source) == nil {
			return errors.New("source cidr must be a valid IP or CIDR")
		}
	}
	return nil
}
