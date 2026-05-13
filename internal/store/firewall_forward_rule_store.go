package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/zmoony/pi-gateway/internal/domain"
)

type FirewallForwardRuleStore interface {
	List(ctx context.Context) ([]domain.FirewallForwardRule, error)
	FindByID(ctx context.Context, id int64) (*domain.FirewallForwardRule, error)
	Create(ctx context.Context, item *domain.FirewallForwardRule) error
	Update(ctx context.Context, item *domain.FirewallForwardRule) error
	Delete(ctx context.Context, id int64) error
	NextPriority(ctx context.Context) (int, error)
}

type SQLiteFirewallForwardRuleStore struct {
	DB *sql.DB
}

func (s SQLiteFirewallForwardRuleStore) List(ctx context.Context) ([]domain.FirewallForwardRule, error) {
	rows, err := s.DB.QueryContext(ctx, `
		select id, name, source_cidr, destination_cidr, protocol, destination_port,
		       enabled, priority, description, created_at, updated_at
		from firewall_forward_rules
		order by priority asc, id asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.FirewallForwardRule
	for rows.Next() {
		var item domain.FirewallForwardRule
		var enabled int
		if err := rows.Scan(
			&item.ID, &item.Name, &item.SourceCIDR, &item.DestinationCIDR, &item.Protocol,
			&item.DestinationPort, &enabled, &item.Priority, &item.Description, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		item.Enabled = enabled == 1
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s SQLiteFirewallForwardRuleStore) FindByID(ctx context.Context, id int64) (*domain.FirewallForwardRule, error) {
	row := s.DB.QueryRowContext(ctx, `
		select id, name, source_cidr, destination_cidr, protocol, destination_port,
		       enabled, priority, description, created_at, updated_at
		from firewall_forward_rules
		where id = ?
	`, id)
	var item domain.FirewallForwardRule
	var enabled int
	if err := row.Scan(
		&item.ID, &item.Name, &item.SourceCIDR, &item.DestinationCIDR, &item.Protocol,
		&item.DestinationPort, &enabled, &item.Priority, &item.Description, &item.CreatedAt, &item.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	item.Enabled = enabled == 1
	return &item, nil
}

func (s SQLiteFirewallForwardRuleStore) Create(ctx context.Context, item *domain.FirewallForwardRule) error {
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now
	enabled := 0
	if item.Enabled {
		enabled = 1
	}
	result, err := s.DB.ExecContext(ctx, `
		insert into firewall_forward_rules (
			name, source_cidr, destination_cidr, protocol, destination_port,
			enabled, priority, description, created_at, updated_at
		)
		values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, item.Name, item.SourceCIDR, item.DestinationCIDR, item.Protocol, item.DestinationPort,
		enabled, item.Priority, item.Description, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return err
	}
	item.ID, err = result.LastInsertId()
	return err
}

func (s SQLiteFirewallForwardRuleStore) Update(ctx context.Context, item *domain.FirewallForwardRule) error {
	item.UpdatedAt = time.Now()
	enabled := 0
	if item.Enabled {
		enabled = 1
	}
	_, err := s.DB.ExecContext(ctx, `
		update firewall_forward_rules
		set name = ?, source_cidr = ?, destination_cidr = ?, protocol = ?, destination_port = ?,
		    enabled = ?, priority = ?, description = ?, updated_at = ?
		where id = ?
	`, item.Name, item.SourceCIDR, item.DestinationCIDR, item.Protocol, item.DestinationPort,
		enabled, item.Priority, item.Description, item.UpdatedAt, item.ID)
	return err
}

func (s SQLiteFirewallForwardRuleStore) Delete(ctx context.Context, id int64) error {
	_, err := s.DB.ExecContext(ctx, `delete from firewall_forward_rules where id = ?`, id)
	return err
}

func (s SQLiteFirewallForwardRuleStore) NextPriority(ctx context.Context) (int, error) {
	row := s.DB.QueryRowContext(ctx, `select coalesce(max(priority), 0) + 1 from firewall_forward_rules`)
	var next int
	if err := row.Scan(&next); err != nil {
		return 0, err
	}
	return next, nil
}
