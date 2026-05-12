package service

import (
	"context"
	"errors"
	"testing"

	"github.com/zmoony/pi-gateway/internal/domain"
)

type proxyStoreStub struct {
	items []domain.ProxyHost
}

func (s *proxyStoreStub) List(context.Context) ([]domain.ProxyHost, error) { return s.items, nil }
func (s *proxyStoreStub) FindByID(_ context.Context, id int64) (*domain.ProxyHost, error) {
	for _, item := range s.items {
		if item.ID == id {
			copy := item
			return &copy, nil
		}
	}
	return nil, nil
}
func (s *proxyStoreStub) FindByName(_ context.Context, name string) (*domain.ProxyHost, error) {
	for _, item := range s.items {
		if item.Name == name {
			copy := item
			return &copy, nil
		}
	}
	return nil, nil
}
func (s *proxyStoreStub) FindByServerName(_ context.Context, serverName string) (*domain.ProxyHost, error) {
	for _, item := range s.items {
		if item.ServerName == serverName {
			copy := item
			return &copy, nil
		}
	}
	return nil, nil
}
func (s *proxyStoreStub) Create(_ context.Context, host *domain.ProxyHost) error {
	host.ID = int64(len(s.items) + 1)
	s.items = append(s.items, *host)
	return nil
}
func (s *proxyStoreStub) Update(_ context.Context, host *domain.ProxyHost) error {
	for i := range s.items {
		if s.items[i].ID == host.ID {
			s.items[i] = *host
		}
	}
	return nil
}
func (s *proxyStoreStub) Delete(_ context.Context, id int64) error {
	filtered := s.items[:0]
	for _, item := range s.items {
		if item.ID != id {
			filtered = append(filtered, item)
		}
	}
	s.items = filtered
	return nil
}

type proxyManagerStub struct {
	applied []string
	removed []string
	err     error
}

func (s *proxyManagerStub) Apply(_ context.Context, host domain.ProxyHost) error {
	s.applied = append(s.applied, host.ServerName)
	return s.err
}
func (s *proxyManagerStub) Remove(_ context.Context, serverName string) error {
	s.removed = append(s.removed, serverName)
	return s.err
}
func (s *proxyManagerStub) Render(domain.ProxyHost) string { return "" }

func TestProxyServiceCreateAppliesEnabledHost(t *testing.T) {
	store := &proxyStoreStub{}
	manager := &proxyManagerStub{}
	service := ProxyService{Hosts: store, Manager: manager}

	result, err := service.Create(context.Background(), UpsertProxyHostInput{
		Name:                "blog",
		ServerName:          "blog.example.com",
		UpstreamURL:         "http://127.0.0.1:3000",
		CertificateCertPath: "/certs/fullchain.pem",
		CertificateKeyPath:  "/certs/privkey.pem",
		Enabled:             true,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if result.LastApplyStatus != "applied" {
		t.Fatalf("expected applied status, got %s", result.LastApplyStatus)
	}
	if len(manager.applied) != 1 || manager.applied[0] != "blog.example.com" {
		t.Fatalf("expected nginx apply call, got %#v", manager.applied)
	}
}

func TestProxyServiceDeleteRemovesConfig(t *testing.T) {
	store := &proxyStoreStub{
		items: []domain.ProxyHost{{ID: 1, ServerName: "blog.example.com"}},
	}
	manager := &proxyManagerStub{}
	service := ProxyService{Hosts: store, Manager: manager}

	if err := service.Delete(context.Background(), 1); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if len(manager.removed) != 1 || manager.removed[0] != "blog.example.com" {
		t.Fatalf("expected remove call, got %#v", manager.removed)
	}
}

func TestProxyServiceSurfacesManagerFailure(t *testing.T) {
	store := &proxyStoreStub{}
	manager := &proxyManagerStub{err: errors.New("nginx -t failed")}
	service := ProxyService{Hosts: store, Manager: manager}

	result, err := service.Create(context.Background(), UpsertProxyHostInput{
		Name:                "blog",
		ServerName:          "blog.example.com",
		UpstreamURL:         "http://127.0.0.1:3000",
		CertificateCertPath: "/certs/fullchain.pem",
		CertificateKeyPath:  "/certs/privkey.pem",
		Enabled:             true,
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	if result == nil || result.LastApplyStatus != "error" {
		t.Fatalf("expected result to record error status, got %#v", result)
	}
}
