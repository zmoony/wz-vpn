package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/zmoony/pi-gateway/internal/domain"
)

type DDNSConfigStore interface {
	List(ctx context.Context) ([]domain.DDNSConfig, error)
	FindByID(ctx context.Context, id int64) (*domain.DDNSConfig, error)
	FindByDomain(ctx context.Context, domainName, subdomain string) (*domain.DDNSConfig, error)
	Create(ctx context.Context, item *domain.DDNSConfig) error
	Update(ctx context.Context, item *domain.DDNSConfig) error
}

type SQLiteDDNSConfigStore struct {
	DB *sql.DB
}

func (s SQLiteDDNSConfigStore) List(ctx context.Context) ([]domain.DDNSConfig, error) {
	rows, err := s.DB.QueryContext(ctx, `
		select id, provider, access_key_id, access_key_secret_encrypted, domain, subdomain, enabled,
		       check_interval_seconds, last_known_ipv6, last_status, last_error, last_synced_at,
		       created_at, updated_at
		from ddns_configs
		order by id asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.DDNSConfig
	for rows.Next() {
		var item domain.DDNSConfig
		var enabled int
		if err := rows.Scan(
			&item.ID, &item.Provider, &item.AccessKeyID, &item.AccessKeySecretEnc, &item.Domain,
			&item.Subdomain, &enabled, &item.CheckIntervalSeconds, &item.LastKnownIPv6,
			&item.LastStatus, &item.LastError, &item.LastSyncedAt, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		item.Enabled = enabled == 1
		items = append(items, item)
	}

	return items, rows.Err()
}

func (s SQLiteDDNSConfigStore) FindByID(ctx context.Context, id int64) (*domain.DDNSConfig, error) {
	row := s.DB.QueryRowContext(ctx, `
		select id, provider, access_key_id, access_key_secret_encrypted, domain, subdomain, enabled,
		       check_interval_seconds, last_known_ipv6, last_status, last_error, last_synced_at,
		       created_at, updated_at
		from ddns_configs where id = ?
	`, id)
	return scanDDNS(row)
}

func (s SQLiteDDNSConfigStore) FindByDomain(ctx context.Context, domainName, subdomain string) (*domain.DDNSConfig, error) {
	row := s.DB.QueryRowContext(ctx, `
		select id, provider, access_key_id, access_key_secret_encrypted, domain, subdomain, enabled,
		       check_interval_seconds, last_known_ipv6, last_status, last_error, last_synced_at,
		       created_at, updated_at
		from ddns_configs where domain = ? and subdomain = ?
	`, domainName, subdomain)
	return scanDDNS(row)
}

func scanDDNS(row *sql.Row) (*domain.DDNSConfig, error) {
	var item domain.DDNSConfig
	var enabled int
	if err := row.Scan(
		&item.ID, &item.Provider, &item.AccessKeyID, &item.AccessKeySecretEnc, &item.Domain,
		&item.Subdomain, &enabled, &item.CheckIntervalSeconds, &item.LastKnownIPv6,
		&item.LastStatus, &item.LastError, &item.LastSyncedAt, &item.CreatedAt, &item.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	item.Enabled = enabled == 1
	return &item, nil
}

func (s SQLiteDDNSConfigStore) Create(ctx context.Context, item *domain.DDNSConfig) error {
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now
	enabled := 0
	if item.Enabled {
		enabled = 1
	}
	result, err := s.DB.ExecContext(ctx, `
		insert into ddns_configs (
			provider, access_key_id, access_key_secret_encrypted, domain, subdomain, enabled,
			check_interval_seconds, last_known_ipv6, last_status, last_error, last_synced_at, created_at, updated_at
		)
		values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, item.Provider, item.AccessKeyID, item.AccessKeySecretEnc, item.Domain, item.Subdomain, enabled,
		item.CheckIntervalSeconds, item.LastKnownIPv6, item.LastStatus, item.LastError, item.LastSyncedAt,
		item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return err
	}
	item.ID, err = result.LastInsertId()
	return err
}

func (s SQLiteDDNSConfigStore) Update(ctx context.Context, item *domain.DDNSConfig) error {
	item.UpdatedAt = time.Now()
	enabled := 0
	if item.Enabled {
		enabled = 1
	}
	_, err := s.DB.ExecContext(ctx, `
		update ddns_configs
		set provider = ?, access_key_id = ?, access_key_secret_encrypted = ?, domain = ?, subdomain = ?, enabled = ?,
		    check_interval_seconds = ?, last_known_ipv6 = ?, last_status = ?, last_error = ?, last_synced_at = ?, updated_at = ?
		where id = ?
	`, item.Provider, item.AccessKeyID, item.AccessKeySecretEnc, item.Domain, item.Subdomain, enabled,
		item.CheckIntervalSeconds, item.LastKnownIPv6, item.LastStatus, item.LastError, item.LastSyncedAt, item.UpdatedAt, item.ID)
	return err
}
