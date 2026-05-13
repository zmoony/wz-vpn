package service

import (
	"context"
	"testing"

	"github.com/zmoony/pi-gateway/internal/config"
)

func TestSettingsGetUsesStoredWireGuardOverrides(t *testing.T) {
	service := SettingsService{
		Settings: &settingStoreStub{
			values: map[string]string{
				settingWireGuardInterface: "wg-test",
				settingWireGuardEndpoint:  "vpn.home.test:51820",
			},
		},
		Config: config.Config{
			WireGuardIface:    "wg0",
			WireGuardSubnetV4: "10.66.66.0/24",
			WireGuardServerV4: "10.66.66.1/24",
			WireGuardDNS:      "192.168.1.1",
			WireGuardEndpoint: "vpn.example.com:51820",
			WireGuardPubKey:   "server-public-key",
		},
	}

	result, err := service.Get(context.Background())
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if result.WireGuard.Interface != "wg-test" {
		t.Fatalf("expected stored interface, got %q", result.WireGuard.Interface)
	}
	if result.WireGuard.Endpoint != "vpn.home.test:51820" {
		t.Fatalf("expected stored endpoint, got %q", result.WireGuard.Endpoint)
	}
	if result.WireGuard.PublicKey != "server-public-key" {
		t.Fatalf("expected default public key fallback, got %q", result.WireGuard.PublicKey)
	}
}

func TestSettingsUpdatePersistsWireGuardValues(t *testing.T) {
	store := &settingStoreStub{values: map[string]string{}}
	service := SettingsService{
		Settings: store,
		Config: config.Config{
			WireGuardIface:    "wg0",
			WireGuardSubnetV4: "10.66.66.0/24",
			WireGuardServerV4: "10.66.66.1/24",
			WireGuardDNS:      "192.168.1.1",
			WireGuardEndpoint: "vpn.example.com:51820",
			WireGuardPubKey:   "server-public-key",
		},
	}

	result, err := service.Update(context.Background(), UpdateSettingsInput{
		WireGuard: WireGuardSettingsSection{
			Interface: "wg-home",
			SubnetV4:  "10.77.77.0/24",
			ServerV4:  "10.77.77.1/24",
			DNS:       "192.168.50.1",
			Endpoint:  "vpn.home.test:51820",
			PublicKey: "public-key-2",
		},
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if result.WireGuard.Interface != "wg-home" {
		t.Fatalf("expected updated interface, got %q", result.WireGuard.Interface)
	}
	if store.values[settingWireGuardSubnetV4] != "10.77.77.0/24" {
		t.Fatalf("expected subnet to be stored, got %q", store.values[settingWireGuardSubnetV4])
	}
}
