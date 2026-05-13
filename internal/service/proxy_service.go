package service

import (
	"context"
	"errors"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/zmoony/pi-gateway/internal/domain"
	"github.com/zmoony/pi-gateway/internal/platform/nginx"
	"github.com/zmoony/pi-gateway/internal/store"
)

var ErrProxyHostNotFound = errors.New("proxy host not found")

type ProxyService struct {
	Hosts        store.ProxyHostStore
	Certificates store.CertificateConfigStore
	Manager      nginx.Manager
}

type UpsertProxyHostInput struct {
	Name                  string `json:"name"`
	ServerName            string `json:"serverName"`
	UpstreamURL           string `json:"upstreamUrl"`
	CertificateRootDomain string `json:"certificateRootDomain"`
	Enabled               bool   `json:"enabled"`
	Description           string `json:"description"`
}

func (s ProxyService) List(ctx context.Context) ([]domain.ProxyHost, error) {
	items, err := s.Hosts.List(ctx)
	if err != nil {
		return nil, err
	}
	if s.Certificates == nil {
		return items, nil
	}
	certs, err := s.Certificates.List(ctx)
	if err != nil {
		return items, nil
	}
	for i := range items {
		if strings.TrimSpace(items[i].CertificateRootDomain) != "" {
			continue
		}
		for _, cert := range certs {
			if certMatchesHost(cert, items[i]) {
				items[i].CertificateRootDomain = cert.RootDomain
				break
			}
		}
	}
	return items, nil
}

func (s ProxyService) Create(ctx context.Context, input UpsertProxyHostInput) (*domain.ProxyHost, error) {
	if err := validateProxyInput(input); err != nil {
		return nil, err
	}
	if existing, err := s.Hosts.FindByName(ctx, input.Name); err != nil {
		return nil, err
	} else if existing != nil {
		return nil, errors.New("proxy name already exists")
	}
	if existing, err := s.Hosts.FindByServerName(ctx, input.ServerName); err != nil {
		return nil, err
	} else if existing != nil {
		return nil, errors.New("server name already exists")
	}

	item := domain.ProxyHost{
		Name:                  strings.TrimSpace(input.Name),
		ServerName:            strings.TrimSpace(input.ServerName),
		UpstreamURL:           strings.TrimSpace(input.UpstreamURL),
		CertificateRootDomain: strings.TrimSpace(input.CertificateRootDomain),
		Enabled:               input.Enabled,
		Description:           strings.TrimSpace(input.Description),
	}
	if err := s.resolveCertificate(ctx, &item); err != nil {
		return nil, err
	}
	if err := s.Hosts.Create(ctx, &item); err != nil {
		return nil, err
	}
	s.apply(ctx, &item)
	if err := s.Hosts.Update(ctx, &item); err != nil {
		return nil, err
	}
	if item.LastApplyStatus == "error" {
		return &item, errors.New(item.LastApplyError)
	}
	return &item, nil
}

func (s ProxyService) Update(ctx context.Context, id int64, input UpsertProxyHostInput) (*domain.ProxyHost, error) {
	item, err := s.Hosts.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrProxyHostNotFound
	}
	if err := validateProxyInput(input); err != nil {
		return nil, err
	}
	if existing, err := s.Hosts.FindByName(ctx, input.Name); err != nil {
		return nil, err
	} else if existing != nil && existing.ID != id {
		return nil, errors.New("proxy name already exists")
	}
	if existing, err := s.Hosts.FindByServerName(ctx, input.ServerName); err != nil {
		return nil, err
	} else if existing != nil && existing.ID != id {
		return nil, errors.New("server name already exists")
	}

	item.Name = strings.TrimSpace(input.Name)
	item.ServerName = strings.TrimSpace(input.ServerName)
	item.UpstreamURL = strings.TrimSpace(input.UpstreamURL)
	item.CertificateRootDomain = strings.TrimSpace(input.CertificateRootDomain)
	item.Enabled = input.Enabled
	item.Description = strings.TrimSpace(input.Description)
	if err := s.resolveCertificate(ctx, item); err != nil {
		return nil, err
	}

	s.apply(ctx, item)
	if err := s.Hosts.Update(ctx, item); err != nil {
		return nil, err
	}
	if item.LastApplyStatus == "error" {
		return item, errors.New(item.LastApplyError)
	}
	return item, nil
}

func (s ProxyService) Delete(ctx context.Context, id int64) error {
	item, err := s.Hosts.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if item == nil {
		return ErrProxyHostNotFound
	}
	if err := s.Manager.Remove(ctx, item.ServerName); err != nil {
		return err
	}
	return s.Hosts.Delete(ctx, id)
}

func (s ProxyService) apply(ctx context.Context, item *domain.ProxyHost) {
	var err error
	if item.Enabled {
		err = s.Manager.Apply(ctx, *item)
	} else {
		err = s.Manager.Remove(ctx, item.ServerName)
	}
	now := time.Now()
	item.LastAppliedAt = &now
	if err != nil {
		item.LastApplyStatus = "error"
		item.LastApplyError = err.Error()
		return
	}
	item.LastApplyStatus = "applied"
	item.LastApplyError = ""
}

func validateProxyInput(input UpsertProxyHostInput) error {
	if strings.TrimSpace(input.Name) == "" {
		return errors.New("proxy name is required")
	}
	if strings.TrimSpace(input.ServerName) == "" {
		return errors.New("server name is required")
	}
	if strings.TrimSpace(input.CertificateRootDomain) == "" {
		return errors.New("certificate root domain is required")
	}
	upstream := strings.TrimSpace(input.UpstreamURL)
	if upstream == "" {
		return errors.New("upstream url is required")
	}
	parsed, err := url.Parse(upstream)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return errors.New("upstream url must be a valid absolute URL")
	}
	return nil
}

func (s ProxyService) resolveCertificate(ctx context.Context, item *domain.ProxyHost) error {
	if s.Certificates == nil {
		return errors.New("certificate store is not configured")
	}
	cert, err := s.Certificates.FindByRootDomain(ctx, strings.TrimSpace(item.CertificateRootDomain))
	if err != nil {
		return err
	}
	if cert == nil {
		return errors.New("selected certificate root domain does not exist")
	}
	item.CertificateRootDomain = cert.RootDomain
	item.CertificateCertPath = resolvedCertificateFullchain(*cert)
	item.CertificateKeyPath = resolvedCertificatePrivateKey(*cert)
	return nil
}

func resolvedCertificateFullchain(cert domain.CertificateConfig) string {
	if strings.TrimSpace(cert.FullchainPath) != "" {
		return strings.TrimSpace(cert.FullchainPath)
	}
	return filepath.Join(strings.TrimSpace(cert.InstallDir), "fullchain.pem")
}

func resolvedCertificatePrivateKey(cert domain.CertificateConfig) string {
	if strings.TrimSpace(cert.PrivateKeyPath) != "" {
		return strings.TrimSpace(cert.PrivateKeyPath)
	}
	return filepath.Join(strings.TrimSpace(cert.InstallDir), "privkey.pem")
}

func certMatchesHost(cert domain.CertificateConfig, host domain.ProxyHost) bool {
	fullchain := resolvedCertificateFullchain(cert)
	privateKey := resolvedCertificatePrivateKey(cert)
	return strings.EqualFold(strings.TrimSpace(host.CertificateCertPath), fullchain) &&
		strings.EqualFold(strings.TrimSpace(host.CertificateKeyPath), privateKey)
}
