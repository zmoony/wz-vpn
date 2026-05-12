package wireguard

import (
	"context"
	"fmt"
	"strings"

	"github.com/zmoony/pi-gateway/internal/domain"
	"github.com/zmoony/pi-gateway/internal/platform"
)

type Manager interface {
	ListRuntimePeers(ctx context.Context, iface string) ([]domain.WireGuardRuntimePeer, error)
	GenerateKeyPair(ctx context.Context) (privateKey string, publicKey string, err error)
	GeneratePreSharedKey(ctx context.Context) (string, error)
	ApplyPeer(ctx context.Context, iface string, spec domain.WireGuardPeerSpec) error
	RemovePeer(ctx context.Context, iface string, publicKey string) error
}

type SystemManager struct {
	Runner platform.CommandRunner
}

func (m SystemManager) ListRuntimePeers(ctx context.Context, iface string) ([]domain.WireGuardRuntimePeer, error) {
	output, err := m.Runner.Run(ctx, "wg", "show", iface, "dump")
	if err != nil {
		return nil, err
	}

	return ParseDump(output), nil
}

func (m SystemManager) GenerateKeyPair(ctx context.Context) (string, string, error) {
	privateKey, err := m.Runner.Run(ctx, "wg", "genkey")
	if err != nil {
		return "", "", err
	}

	privateKey = strings.TrimSpace(privateKey)
	publicKey, err := m.Runner.RunInput(ctx, privateKey, "wg", "pubkey")
	if err != nil {
		return "", "", fmt.Errorf("derive public key: %w", err)
	}

	return privateKey, strings.TrimSpace(publicKey), nil
}

func (m SystemManager) GeneratePreSharedKey(ctx context.Context) (string, error) {
	output, err := m.Runner.Run(ctx, "wg", "genpsk")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output), nil
}

func (m SystemManager) ApplyPeer(ctx context.Context, iface string, spec domain.WireGuardPeerSpec) error {
	args := []string{
		"set", iface,
		"peer", spec.PublicKey,
		"allowed-ips", spec.ClientIPv4 + "/32",
	}
	if spec.PersistentKeepalive > 0 {
		args = append(args, "persistent-keepalive", fmt.Sprintf("%d", spec.PersistentKeepalive))
	}

	_, err := m.Runner.Run(ctx, "wg", args...)
	return err
}

func (m SystemManager) RemovePeer(ctx context.Context, iface string, publicKey string) error {
	_, err := m.Runner.Run(ctx, "wg", "set", iface, "peer", publicKey, "remove")
	return err
}
