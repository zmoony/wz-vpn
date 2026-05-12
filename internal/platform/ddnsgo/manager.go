package ddnsgo

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zmoony/pi-gateway/internal/domain"
	"github.com/zmoony/pi-gateway/internal/platform"
)

type Manager interface {
	WriteManagedConfig(ctx context.Context, items []domain.DDNSConfig) error
	ReadStatus(ctx context.Context) (domain.DDNSRuntimeStatus, error)
	Sync(ctx context.Context) error
}

type managedConfigFile struct {
	ManagedBy string              `json:"managedBy"`
	Items     []managedConfigItem `json:"items"`
}

type managedConfigItem struct {
	Provider             string `json:"provider"`
	AccessKeyID          string `json:"accessKeyId"`
	AccessKeySecret      string `json:"accessKeySecret"`
	Domain               string `json:"domain"`
	Subdomain            string `json:"subdomain"`
	Enabled              bool   `json:"enabled"`
	CheckIntervalSeconds int    `json:"checkIntervalSeconds"`
}

type SystemManager struct {
	ConfigPath    string
	StatusPath    string
	ReloadCommand string
	Runner        platform.CommandRunner
}

func (m SystemManager) WriteManagedConfig(_ context.Context, items []domain.DDNSConfig) error {
	file := managedConfigFile{
		ManagedBy: "pi-gateway",
		Items:     make([]managedConfigItem, 0, len(items)),
	}
	for _, item := range items {
		file.Items = append(file.Items, managedConfigItem{
			Provider:             item.Provider,
			AccessKeyID:          item.AccessKeyID,
			AccessKeySecret:      item.AccessKeySecretEnc,
			Domain:               item.Domain,
			Subdomain:            item.Subdomain,
			Enabled:              item.Enabled,
			CheckIntervalSeconds: item.CheckIntervalSeconds,
		})
	}

	if err := os.MkdirAll(filepath.Dir(m.ConfigPath), 0o755); err != nil {
		return fmt.Errorf("create ddns-go config dir: %w", err)
	}

	content, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal ddns-go managed config: %w", err)
	}

	if err := os.WriteFile(m.ConfigPath, content, 0o600); err != nil {
		return fmt.Errorf("write ddns-go managed config: %w", err)
	}

	return nil
}

func (m SystemManager) ReadStatus(_ context.Context) (domain.DDNSRuntimeStatus, error) {
	if m.StatusPath == "" {
		return domain.DDNSRuntimeStatus{LastStatus: "unconfigured"}, nil
	}

	content, err := os.ReadFile(m.StatusPath)
	if err != nil {
		if os.IsNotExist(err) {
			return domain.DDNSRuntimeStatus{LastStatus: "status-unavailable"}, nil
		}
		return domain.DDNSRuntimeStatus{}, fmt.Errorf("read ddns-go status: %w", err)
	}

	var status domain.DDNSRuntimeStatus
	if err := json.Unmarshal(content, &status); err != nil {
		return domain.DDNSRuntimeStatus{}, fmt.Errorf("decode ddns-go status: %w", err)
	}

	if status.LastStatus == "" {
		status.LastStatus = "unknown"
	}
	return status, nil
}

func (m SystemManager) Sync(ctx context.Context) error {
	if strings.TrimSpace(m.ReloadCommand) == "" {
		return nil
	}
	parts := strings.Fields(m.ReloadCommand)
	if len(parts) == 0 {
		return nil
	}
	_, err := m.Runner.Run(ctx, parts[0], parts[1:]...)
	if err != nil {
		return fmt.Errorf("reload ddns-go: %w", err)
	}
	return nil
}
