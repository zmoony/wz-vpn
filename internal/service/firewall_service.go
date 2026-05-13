package service

import (
	"context"
	"errors"
	"net"
	"strconv"
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
	ForwardRules store.FirewallForwardRuleStore
	Settings     store.SettingStore
	Manager      nftables.Manager
	PendingTTL   time.Duration
	DefaultWGInterface string

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

type UpsertFirewallForwardInput struct {
	Enabled bool   `json:"enabled"`
	LanCIDR string `json:"lanCidr"`
}

type UpsertFirewallForwardRuleInput struct {
	Name            string `json:"name"`
	SourceCIDR      string `json:"sourceCidr"`
	DestinationCIDR string `json:"destinationCidr"`
	Protocol        string `json:"protocol"`
	DestinationPort int    `json:"destinationPort"`
	Enabled         bool   `json:"enabled"`
	Priority        int    `json:"priority"`
	Description     string `json:"description"`
}

const (
	firewallForwardEnabledKey = "firewall.forward.enabled"
	firewallForwardCIDRKey    = "firewall.forward.lan_cidr"
)

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
	forwardRules, err := s.ForwardRules.List(ctx)
	if err != nil {
		return "", err
	}
	forward, err := s.GetForwardConfig(ctx)
	if err != nil {
		return "", err
	}
	return s.Manager.RenderRules(rules, forwardRules, forward), nil
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
	forwardRules, err := s.ForwardRules.List(ctx)
	if err != nil {
		return nil, "", err
	}
	forward, err := s.GetForwardConfig(ctx)
	if err != nil {
		return nil, "", err
	}
	candidate := s.Manager.RenderRules(rules, forwardRules, forward)
	rollbackText := s.Manager.RenderRules(nil, nil, domain.FirewallForwardConfig{})
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

func (s *FirewallService) GetForwardConfig(ctx context.Context) (domain.FirewallForwardConfig, error) {
	enabledRaw, err := s.Settings.Get(ctx, firewallForwardEnabledKey)
	if err != nil {
		return domain.FirewallForwardConfig{}, err
	}
	lanCIDR, err := s.Settings.Get(ctx, firewallForwardCIDRKey)
	if err != nil {
		return domain.FirewallForwardConfig{}, err
	}
	return domain.FirewallForwardConfig{
		Enabled:     strings.EqualFold(enabledRaw, "true"),
		WGInterface: s.DefaultWGInterface,
		LanCIDR:     strings.TrimSpace(lanCIDR),
	}, nil
}

func (s *FirewallService) ListForwardRules(ctx context.Context) ([]domain.FirewallForwardRule, error) {
	return s.ForwardRules.List(ctx)
}

func (s *FirewallService) CreateForwardRule(ctx context.Context, input UpsertFirewallForwardRuleInput) (*domain.FirewallForwardRule, error) {
	item, err := s.buildForwardRule(ctx, input, nil)
	if err != nil {
		return nil, err
	}
	if err := s.ForwardRules.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *FirewallService) UpdateForwardRule(ctx context.Context, id int64, input UpsertFirewallForwardRuleInput) (*domain.FirewallForwardRule, error) {
	current, err := s.ForwardRules.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, ErrFirewallRuleNotFound
	}
	item, err := s.buildForwardRule(ctx, input, current)
	if err != nil {
		return nil, err
	}
	if err := s.ForwardRules.Update(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *FirewallService) DeleteForwardRule(ctx context.Context, id int64) error {
	current, err := s.ForwardRules.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if current == nil {
		return ErrFirewallRuleNotFound
	}
	return s.ForwardRules.Delete(ctx, id)
}

func (s *FirewallService) UpdateForwardConfig(ctx context.Context, input UpsertFirewallForwardInput) (domain.FirewallForwardConfig, error) {
	if input.Enabled {
		if strings.TrimSpace(input.LanCIDR) == "" {
			return domain.FirewallForwardConfig{}, errors.New("lan cidr is required when forward is enabled")
		}
		if _, _, err := net.ParseCIDR(strings.TrimSpace(input.LanCIDR)); err != nil {
			return domain.FirewallForwardConfig{}, errors.New("lan cidr must be a valid CIDR")
		}
	}
	if err := s.Settings.Set(ctx, firewallForwardEnabledKey, strings.ToLower(strconv.FormatBool(input.Enabled))); err != nil {
		return domain.FirewallForwardConfig{}, err
	}
	if err := s.Settings.Set(ctx, firewallForwardCIDRKey, strings.TrimSpace(input.LanCIDR)); err != nil {
		return domain.FirewallForwardConfig{}, err
	}
	return s.GetForwardConfig(ctx)
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

func (s *FirewallService) buildForwardRule(ctx context.Context, input UpsertFirewallForwardRuleInput, current *domain.FirewallForwardRule) (*domain.FirewallForwardRule, error) {
	if err := validateFirewallForwardRuleInput(input); err != nil {
		return nil, err
	}
	priority := input.Priority
	if priority <= 0 {
		if current != nil {
			priority = current.Priority
		} else {
			next, err := s.ForwardRules.NextPriority(ctx)
			if err != nil {
				return nil, err
			}
			priority = next
		}
	}
	item := &domain.FirewallForwardRule{
		Name:            strings.TrimSpace(input.Name),
		SourceCIDR:      strings.TrimSpace(input.SourceCIDR),
		DestinationCIDR: strings.TrimSpace(input.DestinationCIDR),
		Protocol:        strings.ToLower(strings.TrimSpace(input.Protocol)),
		DestinationPort: input.DestinationPort,
		Enabled:         input.Enabled,
		Priority:        priority,
		Description:     strings.TrimSpace(input.Description),
	}
	if current != nil {
		item.ID = current.ID
		item.CreatedAt = current.CreatedAt
		item.UpdatedAt = current.UpdatedAt
	}
	return item, nil
}

func validateFirewallForwardRuleInput(input UpsertFirewallForwardRuleInput) error {
	if strings.TrimSpace(input.Name) == "" {
		return errors.New("forward rule name is required")
	}
	source := strings.TrimSpace(input.SourceCIDR)
	destination := strings.TrimSpace(input.DestinationCIDR)
	if source == "" || destination == "" {
		return errors.New("source cidr and destination cidr are required")
	}
	if _, _, err := net.ParseCIDR(source); err != nil {
		return errors.New("source cidr must be a valid CIDR")
	}
	if _, _, err := net.ParseCIDR(destination); err != nil {
		return errors.New("destination cidr must be a valid CIDR")
	}
	protocol := strings.ToLower(strings.TrimSpace(input.Protocol))
	switch protocol {
	case "", "any":
		if input.DestinationPort > 0 {
			return errors.New("destination port requires tcp or udp protocol")
		}
	case "tcp", "udp":
		if input.DestinationPort < 0 || input.DestinationPort > 65535 {
			return errors.New("destination port must be between 0 and 65535")
		}
	default:
		return errors.New("forward protocol must be any, tcp, or udp")
	}
	return nil
}
