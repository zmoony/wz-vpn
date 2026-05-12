package service

import (
	"context"
	"testing"

	"github.com/zmoony/pi-gateway/internal/domain"
)

type ddnsStoreStub struct {
	items []domain.DDNSConfig
}

func (s *ddnsStoreStub) List(context.Context) ([]domain.DDNSConfig, error) { return s.items, nil }
func (s *ddnsStoreStub) FindByID(_ context.Context, id int64) (*domain.DDNSConfig, error) {
	for _, item := range s.items {
		if item.ID == id {
			copy := item
			return &copy, nil
		}
	}
	return nil, nil
}
func (s *ddnsStoreStub) FindByDomain(_ context.Context, domainName, subdomain string) (*domain.DDNSConfig, error) {
	for _, item := range s.items {
		if item.Domain == domainName && item.Subdomain == subdomain {
			copy := item
			return &copy, nil
		}
	}
	return nil, nil
}
func (s *ddnsStoreStub) Create(_ context.Context, item *domain.DDNSConfig) error {
	item.ID = int64(len(s.items) + 1)
	s.items = append(s.items, *item)
	return nil
}
func (s *ddnsStoreStub) Update(_ context.Context, item *domain.DDNSConfig) error {
	for i := range s.items {
		if s.items[i].ID == item.ID {
			s.items[i] = *item
		}
	}
	return nil
}

type ddnsManagerStub struct {
	lastWrite []domain.DDNSConfig
	syncCount int
}

func (s *ddnsManagerStub) WriteManagedConfig(_ context.Context, items []domain.DDNSConfig) error {
	s.lastWrite = append([]domain.DDNSConfig(nil), items...)
	return nil
}
func (s *ddnsManagerStub) ReadStatus(context.Context) (domain.DDNSRuntimeStatus, error) {
	return domain.DDNSRuntimeStatus{LastStatus: "ok"}, nil
}
func (s *ddnsManagerStub) Sync(context.Context) error {
	s.syncCount++
	return nil
}

func TestDDNSServiceCreateWritesManagedConfigAndSyncs(t *testing.T) {
	store := &ddnsStoreStub{}
	manager := &ddnsManagerStub{}
	service := DDNSService{Configs: store, Manager: manager}

	item, err := service.Create(context.Background(), UpsertDDNSConfigInput{
		Provider:             "aliyun",
		AccessKeyID:          "ak",
		AccessKeySecret:      "secret",
		Domain:               "example.com",
		Subdomain:            "home",
		Enabled:              true,
		CheckIntervalSeconds: 300,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if item.LastStatus != "synced" {
		t.Fatalf("expected synced status, got %s", item.LastStatus)
	}
	if len(manager.lastWrite) != 1 || manager.lastWrite[0].Domain != "example.com" {
		t.Fatalf("expected managed config write, got %#v", manager.lastWrite)
	}
	if manager.syncCount != 1 {
		t.Fatalf("expected sync count 1, got %d", manager.syncCount)
	}
}

func TestDDNSServiceListReturnsRuntimeStatus(t *testing.T) {
	store := &ddnsStoreStub{
		items: []domain.DDNSConfig{{ID: 1, Domain: "example.com", Subdomain: "home"}},
	}
	manager := &ddnsManagerStub{}
	service := DDNSService{Configs: store, Manager: manager}

	items, status, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected one item, got %d", len(items))
	}
	if status.LastStatus != "ok" {
		t.Fatalf("expected runtime status ok, got %s", status.LastStatus)
	}
}
