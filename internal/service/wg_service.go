package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/zmoony/pi-gateway/internal/domain"
	"github.com/zmoony/pi-gateway/internal/platform/wireguard"
	"github.com/zmoony/pi-gateway/internal/store"
)

var ErrPeerNotFound = errors.New("wireguard peer not found")

type WireGuardService struct {
	Peers          store.WireGuardPeerStore
	Settings       store.SettingStore
	Manager        wireguard.Manager
	QRGenerator    wireguard.QRCodeGenerator
	InterfaceName  string
	ServerEndpoint string
	ServerPublicKey string
	DefaultDNS     string
	SubnetV4       string
	ServerCIDRV4   string
}

type CreatePeerInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreatePeerResult struct {
	Peer       domain.WireGuardPeer `json:"peer"`
	ClientConf string               `json:"clientConf"`
	QRCode     string               `json:"qrCode"`
}

func (s WireGuardService) List(ctx context.Context) ([]domain.WireGuardPeer, error) {
	peers, err := s.Peers.List(ctx)
	if err != nil {
		return nil, err
	}

	runtimePeers, err := s.Manager.ListRuntimePeers(ctx, s.interfaceName(ctx))
	if err != nil {
		return peers, nil
	}

	runtimeIndex := map[string]domain.WireGuardRuntimePeer{}
	for _, peer := range runtimePeers {
		runtimeIndex[peer.PublicKey] = peer
	}

	for i := range peers {
		if runtimePeer, ok := runtimeIndex[peers[i].PublicKey]; ok {
			peers[i].LastHandshakeAt = runtimePeer.LastHandshakeAt
			peers[i].RxBytes = runtimePeer.RxBytes
			peers[i].TxBytes = runtimePeer.TxBytes
			peers[i].Online = runtimePeer.Online
		}
	}

	return peers, nil
}

func (s WireGuardService) Create(ctx context.Context, input CreatePeerInput) (*CreatePeerResult, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, errors.New("peer name is required")
	}

	existing, err := s.Peers.FindByName(ctx, input.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("peer name already exists")
	}

	allocatedIP, err := s.allocateIPv4(ctx, s.subnetV4(ctx))
	if err != nil {
		return nil, err
	}

	privateKey, publicKey, err := s.Manager.GenerateKeyPair(ctx)
	if err != nil {
		return nil, err
	}

	peer := domain.WireGuardPeer{
		Name:                input.Name,
		ClientIPv4:          allocatedIP,
		PublicKey:           publicKey,
		PrivateKeyEncrypted: privateKey,
		DNS:                 s.defaultDNS(ctx),
		AllowedIPs:          "192.168.1.0/24,10.66.66.1/32",
		Endpoint:            s.serverEndpoint(ctx),
		PersistentKeepalive: 25,
		Enabled:             true,
		Description:         input.Description,
	}

	if err := s.Peers.Create(ctx, &peer); err != nil {
		return nil, err
	}

	spec := domain.WireGuardPeerSpec{
		Name:                peer.Name,
		ClientIPv4:          peer.ClientIPv4,
		PublicKey:           peer.PublicKey,
		PrivateKey:          privateKey,
		DNS:                 peer.DNS,
		AllowedIPs:          peer.AllowedIPs,
		Endpoint:            peer.Endpoint,
		PersistentKeepalive: peer.PersistentKeepalive,
		Enabled:             peer.Enabled,
		Description:         peer.Description,
	}
	if err := s.Manager.ApplyPeer(ctx, s.interfaceName(ctx), spec); err != nil {
		return nil, err
	}

	clientConf := s.buildClientConfig(ctx, spec)
	qrCode, err := s.QRGenerator.Encode(ctx, clientConf)
	if err != nil {
		return nil, err
	}

	return &CreatePeerResult{
		Peer:       peer,
		ClientConf: clientConf,
		QRCode:     qrCode,
	}, nil
}

func (s WireGuardService) Toggle(ctx context.Context, id int64) (*domain.WireGuardPeer, error) {
	peer, err := s.Peers.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if peer == nil {
		return nil, ErrPeerNotFound
	}

	peer.Enabled = !peer.Enabled
	if peer.Enabled {
		err = s.Manager.ApplyPeer(ctx, s.interfaceName(ctx), domain.WireGuardPeerSpec{
			Name:                peer.Name,
			ClientIPv4:          peer.ClientIPv4,
			PublicKey:           peer.PublicKey,
			PrivateKey:          peer.PrivateKeyEncrypted,
			DNS:                 peer.DNS,
			AllowedIPs:          peer.AllowedIPs,
			Endpoint:            peer.Endpoint,
			PersistentKeepalive: peer.PersistentKeepalive,
			Enabled:             peer.Enabled,
			Description:         peer.Description,
		})
	} else {
		err = s.Manager.RemovePeer(ctx, s.interfaceName(ctx), peer.PublicKey)
	}
	if err != nil {
		return nil, err
	}

	if err := s.Peers.Update(ctx, peer); err != nil {
		return nil, err
	}

	return peer, nil
}

func (s WireGuardService) Delete(ctx context.Context, id int64) error {
	peer, err := s.Peers.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if peer == nil {
		return ErrPeerNotFound
	}

	if err := s.Manager.RemovePeer(ctx, s.interfaceName(ctx), peer.PublicKey); err != nil {
		return err
	}

	return s.Peers.Delete(ctx, id)
}

func (s WireGuardService) allocateIPv4(ctx context.Context, subnet string) (string, error) {
	_, network, err := net.ParseCIDR(subnet)
	if err != nil {
		return "", fmt.Errorf("parse wireguard subnet: %w", err)
	}

	used := map[string]struct{}{}
	peers, err := s.Peers.List(ctx)
	if err != nil {
		return "", err
	}
	for _, peer := range peers {
		used[peer.ClientIPv4] = struct{}{}
	}

	base := network.IP.To4()
	for host := 2; host < 255; host++ {
		ip := net.IPv4(base[0], base[1], base[2], byte(host)).String()
		if _, exists := used[ip]; !exists {
			return ip, nil
		}
	}

	return "", errors.New("no available IPv4 address left in subnet")
}

func (s WireGuardService) buildClientConfig(ctx context.Context, spec domain.WireGuardPeerSpec) string {
	serverAddress := s.serverCIDRV4(ctx)
	if ip, _, err := net.ParseCIDR(serverAddress); err == nil {
		serverAddress = ip.String()
	}

	lines := []string{
		"[Interface]",
		"PrivateKey = " + spec.PrivateKey,
		"Address = " + spec.ClientIPv4 + "/32",
		"DNS = " + spec.DNS,
		"",
		"[Peer]",
		"PublicKey = " + s.serverPublicKey(ctx),
	}
	if spec.PresharedKey != "" {
		lines = append(lines, "PresharedKey = "+spec.PresharedKey)
	}
	lines = append(lines,
		"AllowedIPs = "+spec.AllowedIPs,
		"Endpoint = "+spec.Endpoint,
		"PersistentKeepalive = "+strconv.Itoa(spec.PersistentKeepalive),
		"# Server = "+serverAddress,
	)

	return strings.Join(lines, "\n")
}

func (s WireGuardService) interfaceName(ctx context.Context) string {
	return s.getSetting(ctx, settingWireGuardInterface, s.InterfaceName)
}

func (s WireGuardService) subnetV4(ctx context.Context) string {
	return s.getSetting(ctx, settingWireGuardSubnetV4, s.SubnetV4)
}

func (s WireGuardService) serverCIDRV4(ctx context.Context) string {
	return s.getSetting(ctx, settingWireGuardServerV4, s.ServerCIDRV4)
}

func (s WireGuardService) defaultDNS(ctx context.Context) string {
	return s.getSetting(ctx, settingWireGuardDNS, s.DefaultDNS)
}

func (s WireGuardService) serverEndpoint(ctx context.Context) string {
	return s.getSetting(ctx, settingWireGuardEndpoint, s.ServerEndpoint)
}

func (s WireGuardService) serverPublicKey(ctx context.Context) string {
	return s.getSetting(ctx, settingWireGuardPublicKey, s.ServerPublicKey)
}

func (s WireGuardService) getSetting(ctx context.Context, key, fallback string) string {
	if s.Settings == nil {
		return fallback
	}
	value, err := s.Settings.Get(ctx, key)
	if err != nil || strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}
