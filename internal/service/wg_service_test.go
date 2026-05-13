package service

import (
	"context"
	"strings"
	"testing"

	"github.com/zmoony/pi-gateway/internal/domain"
)

type wgPeerStoreStub struct {
	peers []domain.WireGuardPeer
}

func (s *wgPeerStoreStub) List(context.Context) ([]domain.WireGuardPeer, error) {
	return s.peers, nil
}
func (s *wgPeerStoreStub) FindByID(context.Context, int64) (*domain.WireGuardPeer, error) {
	return nil, nil
}
func (s *wgPeerStoreStub) FindByName(context.Context, string) (*domain.WireGuardPeer, error) {
	return nil, nil
}
func (s *wgPeerStoreStub) Create(context.Context, *domain.WireGuardPeer) error { return nil }
func (s *wgPeerStoreStub) Update(context.Context, *domain.WireGuardPeer) error { return nil }
func (s *wgPeerStoreStub) Delete(context.Context, int64) error { return nil }

type wireGuardManagerStub struct{}

func (wireGuardManagerStub) ListRuntimePeers(context.Context, string) ([]domain.WireGuardRuntimePeer, error) {
	return nil, nil
}
func (wireGuardManagerStub) GenerateKeyPair(context.Context) (string, string, error) {
	return "private-key", "public-key", nil
}
func (wireGuardManagerStub) GeneratePreSharedKey(context.Context) (string, error) {
	return "psk-key", nil
}
func (wireGuardManagerStub) ApplyPeer(context.Context, string, domain.WireGuardPeerSpec) error {
	return nil
}
func (wireGuardManagerStub) RemovePeer(context.Context, string, string) error { return nil }

type qrStub struct{}

func (qrStub) Encode(context.Context, string) (string, error) { return "qr", nil }

func TestAllocateIPv4SkipsUsedAddresses(t *testing.T) {
	service := WireGuardService{
		Peers: &wgPeerStoreStub{
			peers: []domain.WireGuardPeer{
				{ClientIPv4: "10.66.66.2"},
				{ClientIPv4: "10.66.66.3"},
			},
		},
		SubnetV4: "10.66.66.0/24",
	}

	ip, err := service.allocateIPv4(context.Background(), service.SubnetV4)
	if err != nil {
		t.Fatalf("allocateIPv4() error = %v", err)
	}

	if ip != "10.66.66.4" {
		t.Fatalf("expected 10.66.66.4, got %s", ip)
	}
}

func TestCreateBuildsClientConfig(t *testing.T) {
	store := &wgPeerStoreStub{}
	service := WireGuardService{
		Peers:          store,
		Manager:        wireGuardManagerStub{},
		QRGenerator:    qrStub{},
		InterfaceName:  "wg0",
		ServerEndpoint: "vpn.example.com:51820",
		ServerPublicKey: "server-public-key",
		DefaultDNS:     "192.168.1.1",
		SubnetV4:       "10.66.66.0/24",
		ServerCIDRV4:   "10.66.66.1/24",
	}

	result, err := service.Create(context.Background(), CreatePeerInput{Name: "iphone"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if !strings.Contains(result.ClientConf, "Endpoint = vpn.example.com:51820") {
		t.Fatalf("expected client conf to include endpoint, got: %s", result.ClientConf)
	}
	if !strings.Contains(result.ClientConf, "PublicKey = server-public-key") {
		t.Fatalf("expected client conf to include server public key, got: %s", result.ClientConf)
	}
	if result.Peer.ClientIPv4 != "10.66.66.2" {
		t.Fatalf("expected first allocated IP to be 10.66.66.2, got %s", result.Peer.ClientIPv4)
	}
}
