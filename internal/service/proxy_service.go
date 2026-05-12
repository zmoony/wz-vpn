package service

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/zmoony/pi-gateway/internal/domain"
	"github.com/zmoony/pi-gateway/internal/platform/nginx"
	"github.com/zmoony/pi-gateway/internal/store"
)

var ErrProxyHostNotFound = errors.New("proxy host not found")

type ProxyService struct {
	Hosts   store.ProxyHostStore
	Manager nginx.Manager
}

type UpsertProxyHostInput struct {
	Name                string `json:"name"`
	ServerName          string `json:"serverName"`
	UpstreamURL         string `json:"upstreamUrl"`
	CertificateCertPath string `json:"certificateCertPath"`
	CertificateKeyPath  string `json:"certificateKeyPath"`
	Enabled             bool   `json:"enabled"`
	Description         string `json:"description"`
}

func (s ProxyService) List(ctx context.Context) ([]domain.ProxyHost, error) {
	return s.Hosts.List(ctx)
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
		Name:                strings.TrimSpace(input.Name),
		ServerName:          strings.TrimSpace(input.ServerName),
		UpstreamURL:         strings.TrimSpace(input.UpstreamURL),
		CertificateCertPath: strings.TrimSpace(input.CertificateCertPath),
		CertificateKeyPath:  strings.TrimSpace(input.CertificateKeyPath),
		Enabled:             input.Enabled,
		Description:         strings.TrimSpace(input.Description),
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
	item.CertificateCertPath = strings.TrimSpace(input.CertificateCertPath)
	item.CertificateKeyPath = strings.TrimSpace(input.CertificateKeyPath)
	item.Enabled = input.Enabled
	item.Description = strings.TrimSpace(input.Description)

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
	if strings.TrimSpace(input.CertificateCertPath) == "" || strings.TrimSpace(input.CertificateKeyPath) == "" {
		return errors.New("certificate paths are required")
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
