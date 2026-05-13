package service

import (
	"context"
	"errors"
	"net"
	"strings"

	"github.com/zmoony/pi-gateway/internal/config"
	"github.com/zmoony/pi-gateway/internal/store"
)

type SettingsService struct {
	Settings store.SettingStore
	Config   config.Config
}

type WireGuardSettingsSection struct {
	Interface string `json:"interface"`
	SubnetV4  string `json:"subnetV4"`
	ServerV4  string `json:"serverV4"`
	DNS       string `json:"dns"`
	Endpoint  string `json:"endpoint"`
	PublicKey string `json:"publicKey"`
}

type RuntimeSettingsSection struct {
	AppEnv                 string `json:"appEnv"`
	ListenAddr             string `json:"listenAddr"`
	DataDir                string `json:"dataDir"`
	DatabasePath           string `json:"databasePath"`
	WebDistDir             string `json:"webDistDir"`
	NginxSitesDir          string `json:"nginxSitesDir"`
	DDNSGoConfigPath       string `json:"ddnsGoConfigPath"`
	CertsInstallRoot       string `json:"certsInstallRoot"`
	NftablesRulesPath      string `json:"nftablesRulesPath"`
	FirewallPendingSeconds int    `json:"firewallPendingSeconds"`
}

type SettingsView struct {
	WireGuard WireGuardSettingsSection `json:"wireGuard"`
	Runtime   RuntimeSettingsSection   `json:"runtime"`
}

type UpdateSettingsInput struct {
	WireGuard WireGuardSettingsSection `json:"wireGuard"`
}

func (s SettingsService) Get(ctx context.Context) (SettingsView, error) {
	return SettingsView{
		WireGuard: WireGuardSettingsSection{
			Interface: s.getSetting(ctx, settingWireGuardInterface, s.Config.WireGuardIface),
			SubnetV4:  s.getSetting(ctx, settingWireGuardSubnetV4, s.Config.WireGuardSubnetV4),
			ServerV4:  s.getSetting(ctx, settingWireGuardServerV4, s.Config.WireGuardServerV4),
			DNS:       s.getSetting(ctx, settingWireGuardDNS, s.Config.WireGuardDNS),
			Endpoint:  s.getSetting(ctx, settingWireGuardEndpoint, s.Config.WireGuardEndpoint),
			PublicKey: s.getSetting(ctx, settingWireGuardPublicKey, s.Config.WireGuardPubKey),
		},
		Runtime: RuntimeSettingsSection{
			AppEnv:                 s.Config.AppEnv,
			ListenAddr:             s.Config.ListenAddr,
			DataDir:                s.Config.DataDir,
			DatabasePath:           s.Config.DatabasePath,
			WebDistDir:             s.Config.WebDistDir,
			NginxSitesDir:          s.Config.NginxSitesDir,
			DDNSGoConfigPath:       s.Config.DDNSGoConfigPath,
			CertsInstallRoot:       s.Config.CertsInstallRoot,
			NftablesRulesPath:      s.Config.NftablesRulesPath,
			FirewallPendingSeconds: s.Config.FirewallPendingSeconds,
		},
	}, nil
}

func (s SettingsService) Update(ctx context.Context, input UpdateSettingsInput) (SettingsView, error) {
	if err := validateWireGuardSettings(input.WireGuard); err != nil {
		return SettingsView{}, err
	}
	pairs := map[string]string{
		settingWireGuardInterface: strings.TrimSpace(input.WireGuard.Interface),
		settingWireGuardSubnetV4:  strings.TrimSpace(input.WireGuard.SubnetV4),
		settingWireGuardServerV4:  strings.TrimSpace(input.WireGuard.ServerV4),
		settingWireGuardDNS:       strings.TrimSpace(input.WireGuard.DNS),
		settingWireGuardEndpoint:  strings.TrimSpace(input.WireGuard.Endpoint),
		settingWireGuardPublicKey: strings.TrimSpace(input.WireGuard.PublicKey),
	}
	for key, value := range pairs {
		if err := s.Settings.Set(ctx, key, value); err != nil {
			return SettingsView{}, err
		}
	}
	return s.Get(ctx)
}

func (s SettingsService) getSetting(ctx context.Context, key, fallback string) string {
	value, err := s.Settings.Get(ctx, key)
	if err != nil || strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}

func validateWireGuardSettings(input WireGuardSettingsSection) error {
	if strings.TrimSpace(input.Interface) == "" {
		return errors.New("wireguard interface is required")
	}
	if _, _, err := net.ParseCIDR(strings.TrimSpace(input.SubnetV4)); err != nil {
		return errors.New("wireguard subnet must be a valid CIDR")
	}
	if _, _, err := net.ParseCIDR(strings.TrimSpace(input.ServerV4)); err != nil {
		return errors.New("wireguard server address must be a valid CIDR")
	}
	if strings.TrimSpace(input.DNS) == "" {
		return errors.New("wireguard dns is required")
	}
	if strings.TrimSpace(input.Endpoint) == "" {
		return errors.New("wireguard endpoint is required")
	}
	if strings.TrimSpace(input.PublicKey) == "" {
		return errors.New("wireguard public key is required")
	}
	return nil
}
