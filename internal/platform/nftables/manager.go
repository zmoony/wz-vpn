package nftables

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zmoony/pi-gateway/internal/domain"
	"github.com/zmoony/pi-gateway/internal/platform"
)

type Manager interface {
	RenderInputRules(rules []domain.FirewallRule) string
	Apply(ctx context.Context, candidate string, rollbackText string, ttl time.Duration) (*domain.FirewallPendingState, error)
	Confirm(ctx context.Context) error
	Rollback(ctx context.Context) error
	LoadPendingState() (*domain.FirewallPendingState, error)
}

type SystemManager struct {
	Runner        platform.CommandRunner
	RulesPath     string
	StatePath     string
	BackupsDir    string
}

func (m SystemManager) RenderInputRules(rules []domain.FirewallRule) string {
	lines := []string{
		"table inet pi_gateway {",
		"    chain input {",
		"        type filter hook input priority 0;",
		"        policy drop;",
		"        iifname \"lo\" accept",
		"        ct state established,related accept",
		"        icmpv6 type { nd-neighbor-solicit, nd-neighbor-advert, nd-router-solicit, nd-router-advert, echo-request, echo-reply, packet-too-big, time-exceeded, parameter-problem } accept",
	}

	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}
		if rendered := renderRule(rule); rendered != "" {
			lines = append(lines, "        "+rendered)
		}
	}

	lines = append(lines,
		"    }",
		"}",
	)

	return strings.Join(lines, "\n") + "\n"
}

func (m SystemManager) Apply(ctx context.Context, candidate string, rollbackText string, ttl time.Duration) (*domain.FirewallPendingState, error) {
	if err := os.MkdirAll(filepath.Dir(m.RulesPath), 0o755); err != nil {
		return nil, fmt.Errorf("create nftables rules dir: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(m.StatePath), 0o755); err != nil {
		return nil, fmt.Errorf("create nftables state dir: %w", err)
	}
	if err := os.MkdirAll(m.BackupsDir, 0o755); err != nil {
		return nil, fmt.Errorf("create nftables backup dir: %w", err)
	}

	backupPath := filepath.Join(m.BackupsDir, fmt.Sprintf("input-%d.nft", time.Now().UnixNano()))
	originalContent, hadOriginal, err := readNftFile(m.RulesPath)
	if err != nil {
		return nil, err
	}
	if !hadOriginal {
		originalContent = []byte(rollbackText)
	}
	if err := os.WriteFile(backupPath, originalContent, 0o600); err != nil {
		return nil, fmt.Errorf("write nftables backup: %w", err)
	}

	if err := os.WriteFile(m.RulesPath, []byte(candidate), 0o644); err != nil {
		return nil, fmt.Errorf("write nftables rules: %w", err)
	}
	if _, err := m.Runner.Run(ctx, "nft", "-f", m.RulesPath); err != nil {
		_ = os.WriteFile(m.RulesPath, originalContent, 0o644)
		return nil, fmt.Errorf("apply nftables rules: %w", err)
	}

	now := time.Now()
	expires := now.Add(ttl)
	state := &domain.FirewallPendingState{
		Pending:      true,
		BackupPath:   backupPath,
		ExpiresAt:    &expires,
		AppliedAt:    &now,
		OriginalPath: m.RulesPath,
	}
	if err := m.writeState(state); err != nil {
		return nil, err
	}
	return state, nil
}

func (m SystemManager) Confirm(ctx context.Context) error {
	state, err := m.LoadPendingState()
	if err != nil {
		return err
	}
	if state == nil || !state.Pending {
		return nil
	}
	if state.BackupPath != "" {
		_ = os.Remove(state.BackupPath)
	}
	if err := os.Remove(m.StatePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove firewall state: %w", err)
	}
	return nil
}

func (m SystemManager) Rollback(ctx context.Context) error {
	state, err := m.LoadPendingState()
	if err != nil {
		return err
	}
	if state == nil || !state.Pending {
		return nil
	}

	backup, err := os.ReadFile(state.BackupPath)
	if err != nil {
		return fmt.Errorf("read firewall backup: %w", err)
	}
	if err := os.WriteFile(m.RulesPath, backup, 0o644); err != nil {
		return fmt.Errorf("restore firewall rules: %w", err)
	}
	if _, err := m.Runner.Run(ctx, "nft", "-f", m.RulesPath); err != nil {
		return fmt.Errorf("reapply backup firewall rules: %w", err)
	}
	_ = os.Remove(state.BackupPath)
	if err := os.Remove(m.StatePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove firewall state: %w", err)
	}
	return nil
}

func (m SystemManager) LoadPendingState() (*domain.FirewallPendingState, error) {
	content, err := os.ReadFile(m.StatePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read firewall state: %w", err)
	}
	var state domain.FirewallPendingState
	if err := json.Unmarshal(content, &state); err != nil {
		return nil, fmt.Errorf("decode firewall state: %w", err)
	}
	return &state, nil
}

func (m SystemManager) writeState(state *domain.FirewallPendingState) error {
	content, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal firewall state: %w", err)
	}
	if err := os.WriteFile(m.StatePath, content, 0o600); err != nil {
		return fmt.Errorf("write firewall state: %w", err)
	}
	return nil
}

func readNftFile(path string) ([]byte, bool, error) {
	content, err := os.ReadFile(path)
	if err == nil {
		return content, true, nil
	}
	if os.IsNotExist(err) {
		return nil, false, nil
	}
	return nil, false, fmt.Errorf("read nftables file: %w", err)
}

func renderRule(rule domain.FirewallRule) string {
	protocol := strings.ToLower(strings.TrimSpace(rule.Protocol))
	portStart := rule.PortRangeStart
	portEnd := rule.PortRangeEnd
	port := rule.Port

	if rule.Kind == "template" {
		switch strings.ToLower(strings.TrimSpace(rule.TemplateKey)) {
		case "ssh":
			protocol, port = "tcp", 22
		case "http":
			protocol, port = "tcp", 80
		case "https":
			protocol, port = "tcp", 443
		case "wireguard":
			protocol, port = "udp", 51820
		default:
			return ""
		}
	}

	parts := make([]string, 0, 4)
	if source := strings.TrimSpace(rule.SourceCIDR); source != "" {
		if strings.Contains(source, ":") {
			parts = append(parts, "ip6 saddr "+source)
		} else if parsed := net.ParseIP(source); parsed != nil {
			parts = append(parts, "ip saddr "+source)
		} else {
			parts = append(parts, "ip saddr "+source)
		}
	}

	switch protocol {
	case "tcp", "udp":
		parts = append(parts, protocol)
		switch {
		case port > 0:
			parts = append(parts, fmt.Sprintf("dport %d", port))
		case portStart > 0 && portEnd >= portStart:
			parts = append(parts, fmt.Sprintf("dport %d-%d", portStart, portEnd))
		default:
			return ""
		}
	case "icmp", "icmpv6":
		parts = append(parts, protocol)
	default:
		return ""
	}

	parts = append(parts, "accept")
	return strings.Join(parts, " ")
}
