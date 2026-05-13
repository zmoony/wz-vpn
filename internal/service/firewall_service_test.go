package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/zmoony/pi-gateway/internal/domain"
)

type firewallStoreStub struct {
	items []domain.FirewallRule
}

func (s *firewallStoreStub) List(context.Context) ([]domain.FirewallRule, error) { return s.items, nil }
func (s *firewallStoreStub) FindByID(_ context.Context, id int64) (*domain.FirewallRule, error) {
	for _, item := range s.items {
		if item.ID == id {
			copy := item
			return &copy, nil
		}
	}
	return nil, nil
}
func (s *firewallStoreStub) Create(_ context.Context, item *domain.FirewallRule) error {
	item.ID = int64(len(s.items) + 1)
	s.items = append(s.items, *item)
	return nil
}
func (s *firewallStoreStub) Update(_ context.Context, item *domain.FirewallRule) error {
	for i := range s.items {
		if s.items[i].ID == item.ID {
			s.items[i] = *item
		}
	}
	return nil
}
func (s *firewallStoreStub) Delete(_ context.Context, id int64) error {
	filtered := s.items[:0]
	for _, item := range s.items {
		if item.ID != id {
			filtered = append(filtered, item)
		}
	}
	s.items = filtered
	return nil
}
func (s *firewallStoreStub) NextPriority(context.Context) (int, error) { return len(s.items) + 1, nil }

type settingStoreStub struct {
	values map[string]string
}

func (s *settingStoreStub) Get(_ context.Context, key string) (string, error) {
	return s.values[key], nil
}

func (s *settingStoreStub) Set(_ context.Context, key string, value string) error {
	if s.values == nil {
		s.values = map[string]string{}
	}
	s.values[key] = value
	return nil
}

type nftManagerStub struct {
	appliedText string
	rollbackText string
	confirmed bool
	rolledBack bool
	state *domain.FirewallPendingState
}

func (s *nftManagerStub) RenderInputRules(rules []domain.FirewallRule) string {
	panic("old method should not be called")
}
func (s *nftManagerStub) RenderRules(rules []domain.FirewallRule, forward domain.FirewallForwardConfig) string {
	var names []string
	for _, rule := range rules {
		names = append(names, rule.Name)
	}
	text := strings.Join(names, "\n")
	if forward.Enabled && forward.LanCIDR != "" {
		if text != "" {
			text += "\n"
		}
		text += "forward:" + forward.LanCIDR
	}
	return text
}
func (s *nftManagerStub) Apply(_ context.Context, candidate string, rollbackText string, ttl time.Duration) (*domain.FirewallPendingState, error) {
	s.appliedText = candidate
	s.rollbackText = rollbackText
	now := time.Now()
	expires := now.Add(ttl)
	s.state = &domain.FirewallPendingState{Pending: true, AppliedAt: &now, ExpiresAt: &expires}
	return s.state, nil
}
func (s *nftManagerStub) Confirm(context.Context) error {
	s.confirmed = true
	s.state = nil
	return nil
}
func (s *nftManagerStub) Rollback(context.Context) error {
	s.rolledBack = true
	s.state = nil
	return nil
}
func (s *nftManagerStub) LoadPendingState() (*domain.FirewallPendingState, error) { return s.state, nil }

func TestFirewallPreviewIncludesOrderedRules(t *testing.T) {
	store := &firewallStoreStub{
		items: []domain.FirewallRule{
			{Name: "ssh", Priority: 1, Enabled: true},
			{Name: "http", Priority: 2, Enabled: true},
		},
	}
	manager := &nftManagerStub{}
	service := FirewallService{Rules: store, Settings: &settingStoreStub{}, Manager: manager, PendingTTL: 30 * time.Second}

	preview, err := service.Preview(context.Background())
	if err != nil {
		t.Fatalf("Preview() error = %v", err)
	}
	if !strings.Contains(preview, "ssh") || !strings.Contains(preview, "http") {
		t.Fatalf("expected preview to include rule names, got %q", preview)
	}
}

func TestFirewallApplyStartsPendingState(t *testing.T) {
	store := &firewallStoreStub{
		items: []domain.FirewallRule{{Name: "ssh", Priority: 1, Enabled: true}},
	}
	manager := &nftManagerStub{}
	settings := &settingStoreStub{values: map[string]string{
		firewallForwardEnabledKey: "true",
		firewallForwardCIDRKey:    "192.168.1.0/24",
	}}
	service := FirewallService{Rules: store, Settings: settings, Manager: manager, PendingTTL: 30 * time.Second, DefaultWGInterface: "wg0"}

	state, candidate, err := service.Apply(context.Background())
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if state == nil || !state.Pending {
		t.Fatalf("expected pending state, got %#v", state)
	}
	if candidate != "ssh\nforward:192.168.1.0/24" {
		t.Fatalf("expected rendered candidate, got %q", candidate)
	}
}

func TestFirewallConfirmClearsPendingState(t *testing.T) {
	manager := &nftManagerStub{
		state: &domain.FirewallPendingState{Pending: true},
	}
	service := FirewallService{Settings: &settingStoreStub{}, Manager: manager, PendingTTL: 30 * time.Second}

	if err := service.Confirm(context.Background()); err != nil {
		t.Fatalf("Confirm() error = %v", err)
	}
	if !manager.confirmed {
		t.Fatalf("expected confirm to be called")
	}
}

func TestFirewallUpdateForwardConfigStoresCIDR(t *testing.T) {
	settings := &settingStoreStub{values: map[string]string{}}
	service := FirewallService{Settings: settings, DefaultWGInterface: "wg0"}

	config, err := service.UpdateForwardConfig(context.Background(), UpsertFirewallForwardInput{
		Enabled: true,
		LanCIDR: "192.168.1.0/24",
	})
	if err != nil {
		t.Fatalf("UpdateForwardConfig() error = %v", err)
	}
	if !config.Enabled || config.LanCIDR != "192.168.1.0/24" || config.WGInterface != "wg0" {
		t.Fatalf("unexpected config: %#v", config)
	}
}
