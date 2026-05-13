package service

import (
	"context"
	"testing"
	"time"

	"github.com/zmoony/pi-gateway/internal/domain"
)

type certStoreStub struct {
	items []domain.CertificateConfig
}

func (s *certStoreStub) List(context.Context) ([]domain.CertificateConfig, error) { return s.items, nil }
func (s *certStoreStub) FindByID(_ context.Context, id int64) (*domain.CertificateConfig, error) {
	for _, item := range s.items {
		if item.ID == id {
			copy := item
			return &copy, nil
		}
	}
	return nil, nil
}
func (s *certStoreStub) FindByRootDomain(_ context.Context, rootDomain string) (*domain.CertificateConfig, error) {
	for _, item := range s.items {
		if item.RootDomain == rootDomain {
			copy := item
			return &copy, nil
		}
	}
	return nil, nil
}
func (s *certStoreStub) Create(_ context.Context, item *domain.CertificateConfig) error {
	item.ID = int64(len(s.items) + 1)
	s.items = append(s.items, *item)
	return nil
}
func (s *certStoreStub) Update(_ context.Context, item *domain.CertificateConfig) error {
	for i := range s.items {
		if s.items[i].ID == item.ID {
			s.items[i] = *item
		}
	}
	return nil
}

type acmeManagerStub struct {
	issued []string
	renewed []string
}

func (s *acmeManagerStub) IssueAndInstall(_ context.Context, item domain.CertificateConfig) (string, string, error) {
	s.issued = append(s.issued, item.RootDomain)
	return item.InstallDir + "/fullchain.pem", item.InstallDir + "/privkey.pem", nil
}
func (s *acmeManagerStub) Renew(_ context.Context, item domain.CertificateConfig) error {
	s.renewed = append(s.renewed, item.RootDomain)
	return nil
}
func (s *acmeManagerStub) ReadInstalledStatus(_ context.Context, item domain.CertificateConfig) (domain.CertificateConfig, error) {
	notAfter := time.Now().Add(45 * 24 * time.Hour)
	item.NotAfter = &notAfter
	item.DaysRemaining = 45
	return item, nil
}

func TestCertificateServiceCreateIssuesWildcardCertificate(t *testing.T) {
	store := &certStoreStub{}
	manager := &acmeManagerStub{}
	service := CertificateService{
		Configs:            store,
		Manager:            manager,
		DefaultInstallRoot: "/opt/pi-gateway/data/certs",
	}

	item, err := service.Create(context.Background(), UpsertCertificateConfigInput{
		RootDomain:        "example.com",
		Provider:          "aliyun",
		AccessKeyID:       "ak",
		AccessKeySecret:   "secret",
		EnabledAutoRenew:  true,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if item.LastIssueStatus != "issued" {
		t.Fatalf("expected issued status, got %s", item.LastIssueStatus)
	}
	if len(manager.issued) != 1 || manager.issued[0] != "example.com" {
		t.Fatalf("expected issue call, got %#v", manager.issued)
	}
	if item.DaysRemaining != 45 {
		t.Fatalf("expected days remaining 45, got %d", item.DaysRemaining)
	}
}

func TestCertificateServiceRenewCallsManager(t *testing.T) {
	store := &certStoreStub{
		items: []domain.CertificateConfig{{ID: 1, RootDomain: "example.com", InstallDir: "/opt/pi-gateway/data/certs/example.com"}},
	}
	manager := &acmeManagerStub{}
	service := CertificateService{
		Configs:            store,
		Manager:            manager,
		DefaultInstallRoot: "/opt/pi-gateway/data/certs",
	}

	item, err := service.Renew(context.Background(), 1)
	if err != nil {
		t.Fatalf("Renew() error = %v", err)
	}
	if item.LastIssueStatus != "renewed" {
		t.Fatalf("expected renewed status, got %s", item.LastIssueStatus)
	}
	if len(manager.renewed) != 1 || manager.renewed[0] != "example.com" {
		t.Fatalf("expected renew call, got %#v", manager.renewed)
	}
}
