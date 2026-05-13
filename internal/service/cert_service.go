package service

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"time"

	"github.com/zmoony/pi-gateway/internal/domain"
	"github.com/zmoony/pi-gateway/internal/platform/acme"
	"github.com/zmoony/pi-gateway/internal/store"
)

var ErrCertificateConfigNotFound = errors.New("certificate config not found")

type CertificateService struct {
	Configs           store.CertificateConfigStore
	Manager           acme.Manager
	DefaultInstallRoot string
}

type UpsertCertificateConfigInput struct {
	RootDomain       string `json:"rootDomain"`
	Provider         string `json:"provider"`
	AccessKeyID      string `json:"accessKeyId"`
	AccessKeySecret  string `json:"accessKeySecret"`
	InstallDir       string `json:"installDir"`
	EnabledAutoRenew bool   `json:"enabledAutoRenew"`
}

func (s CertificateService) List(ctx context.Context) ([]domain.CertificateConfig, error) {
	items, err := s.Configs.List(ctx)
	if err != nil {
		return nil, err
	}
	for i := range items {
		item, statusErr := s.Manager.ReadInstalledStatus(ctx, items[i])
		if statusErr == nil {
			items[i] = item
		}
	}
	return items, nil
}

func (s CertificateService) Create(ctx context.Context, input UpsertCertificateConfigInput) (*domain.CertificateConfig, error) {
	if err := validateCertificateInput(input); err != nil {
		return nil, err
	}
	rootDomain := strings.TrimSpace(input.RootDomain)
	if existing, err := s.Configs.FindByRootDomain(ctx, rootDomain); err != nil {
		return nil, err
	} else if existing != nil {
		return nil, errors.New("certificate root domain already exists")
	}

	item := domain.CertificateConfig{
		RootDomain:         rootDomain,
		Provider:           normalizedProvider(input.Provider),
		AccessKeyID:        strings.TrimSpace(input.AccessKeyID),
		AccessKeySecretEnc: strings.TrimSpace(input.AccessKeySecret),
		InstallDir:         s.resolveInstallDir(input),
		EnabledAutoRenew:   input.EnabledAutoRenew,
		LastIssueStatus:    "configured",
	}
	if err := s.Configs.Create(ctx, &item); err != nil {
		return nil, err
	}
	if err := s.issue(ctx, &item); err != nil {
		_ = s.Configs.Update(ctx, &item)
		return &item, err
	}
	if err := s.Configs.Update(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (s CertificateService) Update(ctx context.Context, id int64, input UpsertCertificateConfigInput) (*domain.CertificateConfig, error) {
	item, err := s.Configs.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrCertificateConfigNotFound
	}
	if err := validateCertificateInput(input); err != nil {
		return nil, err
	}
	rootDomain := strings.TrimSpace(input.RootDomain)
	if existing, err := s.Configs.FindByRootDomain(ctx, rootDomain); err != nil {
		return nil, err
	} else if existing != nil && existing.ID != id {
		return nil, errors.New("certificate root domain already exists")
	}

	item.RootDomain = rootDomain
	item.Provider = normalizedProvider(input.Provider)
	item.AccessKeyID = strings.TrimSpace(input.AccessKeyID)
	if secret := strings.TrimSpace(input.AccessKeySecret); secret != "" {
		item.AccessKeySecretEnc = secret
	}
	item.InstallDir = s.resolveInstallDir(input)
	item.EnabledAutoRenew = input.EnabledAutoRenew
	if err := s.Configs.Update(ctx, item); err != nil {
		return nil, err
	}
	itemWithStatus, statusErr := s.Manager.ReadInstalledStatus(ctx, *item)
	if statusErr == nil {
		*item = itemWithStatus
	}
	return item, nil
}

func (s CertificateService) Renew(ctx context.Context, id int64) (*domain.CertificateConfig, error) {
	item, err := s.Configs.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrCertificateConfigNotFound
	}
	if err := s.Manager.Renew(ctx, *item); err != nil {
		item.LastIssueStatus = "error"
		item.LastIssueError = err.Error()
		_ = s.Configs.Update(ctx, item)
		return item, err
	}
	item.LastIssueStatus = "renewed"
	item.LastIssueError = ""
	now := time.Now()
	item.LastIssuedAt = &now
	itemWithStatus, statusErr := s.Manager.ReadInstalledStatus(ctx, *item)
	if statusErr == nil {
		*item = itemWithStatus
	}
	if err := s.Configs.Update(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s CertificateService) issue(ctx context.Context, item *domain.CertificateConfig) error {
	fullchain, keyPath, err := s.Manager.IssueAndInstall(ctx, *item)
	now := time.Now()
	item.LastIssuedAt = &now
	if err != nil {
		item.LastIssueStatus = "error"
		item.LastIssueError = err.Error()
		return err
	}
	item.FullchainPath = fullchain
	item.PrivateKeyPath = keyPath
	item.LastIssueStatus = "issued"
	item.LastIssueError = ""
	itemWithStatus, statusErr := s.Manager.ReadInstalledStatus(ctx, *item)
	if statusErr == nil {
		*item = itemWithStatus
	}
	return nil
}

func (s CertificateService) resolveInstallDir(input UpsertCertificateConfigInput) string {
	if dir := strings.TrimSpace(input.InstallDir); dir != "" {
		return dir
	}
	return filepath.Join(s.DefaultInstallRoot, strings.TrimSpace(input.RootDomain))
}

func validateCertificateInput(input UpsertCertificateConfigInput) error {
	if strings.TrimSpace(input.RootDomain) == "" {
		return errors.New("root domain is required")
	}
	if normalizedProvider(input.Provider) == "" {
		return errors.New("provider is required")
	}
	if strings.TrimSpace(input.AccessKeyID) == "" || strings.TrimSpace(input.AccessKeySecret) == "" {
		return errors.New("access key id and secret are required")
	}
	return nil
}
