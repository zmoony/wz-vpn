package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/zmoony/pi-gateway/internal/domain"
)

type ProxyHostStore interface {
	List(ctx context.Context) ([]domain.ProxyHost, error)
	FindByID(ctx context.Context, id int64) (*domain.ProxyHost, error)
	FindByName(ctx context.Context, name string) (*domain.ProxyHost, error)
	FindByServerName(ctx context.Context, serverName string) (*domain.ProxyHost, error)
	Create(ctx context.Context, host *domain.ProxyHost) error
	Update(ctx context.Context, host *domain.ProxyHost) error
	Delete(ctx context.Context, id int64) error
}

type SQLiteProxyHostStore struct {
	DB *sql.DB
}

func (s SQLiteProxyHostStore) List(ctx context.Context) ([]domain.ProxyHost, error) {
	rows, err := s.DB.QueryContext(ctx, `
		select id, name, server_name, upstream_url, certificate_cert_path, certificate_key_path,
		       enabled, description, last_apply_status, last_apply_error, last_applied_at,
		       created_at, updated_at
		from proxy_hosts
		order by id asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.ProxyHost
	for rows.Next() {
		var item domain.ProxyHost
		var enabled int
		if err := rows.Scan(
			&item.ID, &item.Name, &item.ServerName, &item.UpstreamURL, &item.CertificateCertPath,
			&item.CertificateKeyPath, &enabled, &item.Description, &item.LastApplyStatus,
			&item.LastApplyError, &item.LastAppliedAt, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		item.Enabled = enabled == 1
		items = append(items, item)
	}

	return items, rows.Err()
}

func (s SQLiteProxyHostStore) FindByID(ctx context.Context, id int64) (*domain.ProxyHost, error) {
	return s.findOne(ctx, "where id = ?", id)
}

func (s SQLiteProxyHostStore) FindByName(ctx context.Context, name string) (*domain.ProxyHost, error) {
	return s.findOne(ctx, "where name = ?", name)
}

func (s SQLiteProxyHostStore) FindByServerName(ctx context.Context, serverName string) (*domain.ProxyHost, error) {
	return s.findOne(ctx, "where server_name = ?", serverName)
}

func (s SQLiteProxyHostStore) findOne(ctx context.Context, filter string, value any) (*domain.ProxyHost, error) {
	query := `
		select id, name, server_name, upstream_url, certificate_cert_path, certificate_key_path,
		       enabled, description, last_apply_status, last_apply_error, last_applied_at,
		       created_at, updated_at
		from proxy_hosts
	` + filter

	row := s.DB.QueryRowContext(ctx, query, value)
	var item domain.ProxyHost
	var enabled int
	if err := row.Scan(
		&item.ID, &item.Name, &item.ServerName, &item.UpstreamURL, &item.CertificateCertPath,
		&item.CertificateKeyPath, &enabled, &item.Description, &item.LastApplyStatus,
		&item.LastApplyError, &item.LastAppliedAt, &item.CreatedAt, &item.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	item.Enabled = enabled == 1
	return &item, nil
}

func (s SQLiteProxyHostStore) Create(ctx context.Context, host *domain.ProxyHost) error {
	now := time.Now()
	host.CreatedAt = now
	host.UpdatedAt = now
	enabled := 0
	if host.Enabled {
		enabled = 1
	}

	result, err := s.DB.ExecContext(ctx, `
		insert into proxy_hosts (
			name, server_name, upstream_url, certificate_cert_path, certificate_key_path,
			enabled, description, last_apply_status, last_apply_error, last_applied_at, created_at, updated_at
		)
		values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, host.Name, host.ServerName, host.UpstreamURL, host.CertificateCertPath, host.CertificateKeyPath,
		enabled, host.Description, host.LastApplyStatus, host.LastApplyError, host.LastAppliedAt, host.CreatedAt, host.UpdatedAt)
	if err != nil {
		return err
	}

	host.ID, err = result.LastInsertId()
	return err
}

func (s SQLiteProxyHostStore) Update(ctx context.Context, host *domain.ProxyHost) error {
	host.UpdatedAt = time.Now()
	enabled := 0
	if host.Enabled {
		enabled = 1
	}
	_, err := s.DB.ExecContext(ctx, `
		update proxy_hosts
		set name = ?, server_name = ?, upstream_url = ?, certificate_cert_path = ?, certificate_key_path = ?,
		    enabled = ?, description = ?, last_apply_status = ?, last_apply_error = ?, last_applied_at = ?, updated_at = ?
		where id = ?
	`, host.Name, host.ServerName, host.UpstreamURL, host.CertificateCertPath, host.CertificateKeyPath,
		enabled, host.Description, host.LastApplyStatus, host.LastApplyError, host.LastAppliedAt, host.UpdatedAt, host.ID)
	return err
}

func (s SQLiteProxyHostStore) Delete(ctx context.Context, id int64) error {
	_, err := s.DB.ExecContext(ctx, `delete from proxy_hosts where id = ?`, id)
	return err
}
