package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/zmoony/pi-gateway/internal/domain"
	"github.com/zmoony/pi-gateway/internal/platform/ddnsgo"
	"github.com/zmoony/pi-gateway/internal/store"
)

var ErrDDNSConfigNotFound = errors.New("ddns config not found")

type DDNSService struct {
	Configs store.DDNSConfigStore
	Manager ddnsgo.Manager
}

type UpsertDDNSConfigInput struct {
	Provider             string `json:"provider"`
	AccessKeyID          string `json:"accessKeyId"`
	AccessKeySecret      string `json:"accessKeySecret"`
	Domain               string `json:"domain"`
	Subdomain            string `json:"subdomain"`
	Enabled              bool   `json:"enabled"`
	CheckIntervalSeconds int    `json:"checkIntervalSeconds"`
}

func (s DDNSService) List(ctx context.Context) ([]domain.DDNSConfig, domain.DDNSRuntimeStatus, error) {
	items, err := s.Configs.List(ctx)
	if err != nil {
		return nil, domain.DDNSRuntimeStatus{}, err
	}
	status, err := s.Manager.ReadStatus(ctx)
	if err != nil {
		return items, domain.DDNSRuntimeStatus{}, err
	}
	return items, status, nil
}

func (s DDNSService) Create(ctx context.Context, input UpsertDDNSConfigInput) (*domain.DDNSConfig, error) {
	if err := validateDDNSInput(input); err != nil {
		return nil, err
	}
	if existing, err := s.Configs.FindByDomain(ctx, input.Domain, input.Subdomain); err != nil {
		return nil, err
	} else if existing != nil {
		return nil, errors.New("ddns record already exists")
	}

	item := domain.DDNSConfig{
		Provider:             normalizedProvider(input.Provider),
		AccessKeyID:          strings.TrimSpace(input.AccessKeyID),
		AccessKeySecretEnc:   strings.TrimSpace(input.AccessKeySecret),
		Domain:               strings.TrimSpace(input.Domain),
		Subdomain:            strings.TrimSpace(input.Subdomain),
		Enabled:              input.Enabled,
		CheckIntervalSeconds: input.CheckIntervalSeconds,
		LastStatus:           "configured",
	}
	if err := s.Configs.Create(ctx, &item); err != nil {
		return nil, err
	}
	if err := s.writeAndSync(ctx); err != nil {
		item.LastStatus = "error"
		item.LastError = err.Error()
		_ = s.Configs.Update(ctx, &item)
		return &item, err
	}
	now := time.Now()
	item.LastStatus = "synced"
	item.LastError = ""
	item.LastSyncedAt = &now
	if err := s.Configs.Update(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (s DDNSService) Update(ctx context.Context, id int64, input UpsertDDNSConfigInput) (*domain.DDNSConfig, error) {
	item, err := s.Configs.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrDDNSConfigNotFound
	}
	if err := validateDDNSInput(input); err != nil {
		return nil, err
	}
	if existing, err := s.Configs.FindByDomain(ctx, input.Domain, input.Subdomain); err != nil {
		return nil, err
	} else if existing != nil && existing.ID != id {
		return nil, errors.New("ddns record already exists")
	}

	item.Provider = normalizedProvider(input.Provider)
	item.AccessKeyID = strings.TrimSpace(input.AccessKeyID)
	item.AccessKeySecretEnc = strings.TrimSpace(input.AccessKeySecret)
	item.Domain = strings.TrimSpace(input.Domain)
	item.Subdomain = strings.TrimSpace(input.Subdomain)
	item.Enabled = input.Enabled
	item.CheckIntervalSeconds = input.CheckIntervalSeconds
	if err := s.Configs.Update(ctx, item); err != nil {
		return nil, err
	}

	if err := s.writeAndSync(ctx); err != nil {
		item.LastStatus = "error"
		item.LastError = err.Error()
		_ = s.Configs.Update(ctx, item)
		return item, err
	}
	now := time.Now()
	item.LastStatus = "synced"
	item.LastError = ""
	item.LastSyncedAt = &now
	if err := s.Configs.Update(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s DDNSService) Sync(ctx context.Context) error {
	if err := s.writeAndSync(ctx); err != nil {
		return err
	}
	items, err := s.Configs.List(ctx)
	if err != nil {
		return err
	}
	now := time.Now()
	for i := range items {
		items[i].LastStatus = "synced"
		items[i].LastError = ""
		items[i].LastSyncedAt = &now
		if err := s.Configs.Update(ctx, &items[i]); err != nil {
			return err
		}
	}
	return nil
}

func (s DDNSService) writeAndSync(ctx context.Context) error {
	items, err := s.Configs.List(ctx)
	if err != nil {
		return err
	}
	if err := s.Manager.WriteManagedConfig(ctx, items); err != nil {
		return err
	}
	return s.Manager.Sync(ctx)
}

func validateDDNSInput(input UpsertDDNSConfigInput) error {
	if normalizedProvider(input.Provider) == "" {
		return errors.New("provider is required")
	}
	if strings.TrimSpace(input.AccessKeyID) == "" || strings.TrimSpace(input.AccessKeySecret) == "" {
		return errors.New("access key id and secret are required")
	}
	if strings.TrimSpace(input.Domain) == "" || strings.TrimSpace(input.Subdomain) == "" {
		return errors.New("domain and subdomain are required")
	}
	if input.CheckIntervalSeconds <= 0 {
		return errors.New("check interval must be greater than 0")
	}
	return nil
}

func normalizedProvider(provider string) string {
	value := strings.TrimSpace(strings.ToLower(provider))
	if value == "" {
		return "aliyun"
	}
	return value
}
