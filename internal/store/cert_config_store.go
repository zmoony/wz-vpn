package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/zmoony/pi-gateway/internal/domain"
)

type CertificateConfigStore interface {
	List(ctx context.Context) ([]domain.CertificateConfig, error)
	FindByID(ctx context.Context, id int64) (*domain.CertificateConfig, error)
	FindByRootDomain(ctx context.Context, rootDomain string) (*domain.CertificateConfig, error)
	Create(ctx context.Context, item *domain.CertificateConfig) error
	Update(ctx context.Context, item *domain.CertificateConfig) error
}

type SQLiteCertificateConfigStore struct {
	DB *sql.DB
}

func (s SQLiteCertificateConfigStore) List(ctx context.Context) ([]domain.CertificateConfig, error) {
	rows, err := s.DB.QueryContext(ctx, `
		select id, root_domain, provider, access_key_id, access_key_secret_encrypted, install_dir,
		       enabled_auto_renew, fullchain_path, private_key_path, last_issue_status, last_issue_error,
		       last_issued_at, created_at, updated_at
		from cert_configs
		order by id asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.CertificateConfig
	for rows.Next() {
		var item domain.CertificateConfig
		var enabled int
		if err := rows.Scan(
			&item.ID, &item.RootDomain, &item.Provider, &item.AccessKeyID, &item.AccessKeySecretEnc,
			&item.InstallDir, &enabled, &item.FullchainPath, &item.PrivateKeyPath, &item.LastIssueStatus,
			&item.LastIssueError, &item.LastIssuedAt, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		item.EnabledAutoRenew = enabled == 1
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s SQLiteCertificateConfigStore) FindByID(ctx context.Context, id int64) (*domain.CertificateConfig, error) {
	return s.findOne(ctx, "where id = ?", id)
}

func (s SQLiteCertificateConfigStore) FindByRootDomain(ctx context.Context, rootDomain string) (*domain.CertificateConfig, error) {
	return s.findOne(ctx, "where root_domain = ?", rootDomain)
}

func (s SQLiteCertificateConfigStore) findOne(ctx context.Context, filter string, value any) (*domain.CertificateConfig, error) {
	row := s.DB.QueryRowContext(ctx, `
		select id, root_domain, provider, access_key_id, access_key_secret_encrypted, install_dir,
		       enabled_auto_renew, fullchain_path, private_key_path, last_issue_status, last_issue_error,
		       last_issued_at, created_at, updated_at
		from cert_configs `+filter, value)

	var item domain.CertificateConfig
	var enabled int
	if err := row.Scan(
		&item.ID, &item.RootDomain, &item.Provider, &item.AccessKeyID, &item.AccessKeySecretEnc,
		&item.InstallDir, &enabled, &item.FullchainPath, &item.PrivateKeyPath, &item.LastIssueStatus,
		&item.LastIssueError, &item.LastIssuedAt, &item.CreatedAt, &item.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	item.EnabledAutoRenew = enabled == 1
	return &item, nil
}

func (s SQLiteCertificateConfigStore) Create(ctx context.Context, item *domain.CertificateConfig) error {
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now
	enabled := 0
	if item.EnabledAutoRenew {
		enabled = 1
	}
	result, err := s.DB.ExecContext(ctx, `
		insert into cert_configs (
			root_domain, provider, access_key_id, access_key_secret_encrypted, install_dir,
			enabled_auto_renew, fullchain_path, private_key_path, last_issue_status, last_issue_error,
			last_issued_at, created_at, updated_at
		)
		values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, item.RootDomain, item.Provider, item.AccessKeyID, item.AccessKeySecretEnc, item.InstallDir,
		enabled, item.FullchainPath, item.PrivateKeyPath, item.LastIssueStatus, item.LastIssueError,
		item.LastIssuedAt, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return err
	}
	item.ID, err = result.LastInsertId()
	return err
}

func (s SQLiteCertificateConfigStore) Update(ctx context.Context, item *domain.CertificateConfig) error {
	item.UpdatedAt = time.Now()
	enabled := 0
	if item.EnabledAutoRenew {
		enabled = 1
	}
	_, err := s.DB.ExecContext(ctx, `
		update cert_configs
		set root_domain = ?, provider = ?, access_key_id = ?, access_key_secret_encrypted = ?, install_dir = ?,
		    enabled_auto_renew = ?, fullchain_path = ?, private_key_path = ?, last_issue_status = ?, last_issue_error = ?,
		    last_issued_at = ?, updated_at = ?
		where id = ?
	`, item.RootDomain, item.Provider, item.AccessKeyID, item.AccessKeySecretEnc, item.InstallDir,
		enabled, item.FullchainPath, item.PrivateKeyPath, item.LastIssueStatus, item.LastIssueError,
		item.LastIssuedAt, item.UpdatedAt, item.ID)
	return err
}
