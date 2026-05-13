package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/zmoony/pi-gateway/internal/domain"
)

type FirewallRuleStore interface {
	List(ctx context.Context) ([]domain.FirewallRule, error)
	FindByID(ctx context.Context, id int64) (*domain.FirewallRule, error)
	Create(ctx context.Context, item *domain.FirewallRule) error
	Update(ctx context.Context, item *domain.FirewallRule) error
	Delete(ctx context.Context, id int64) error
	NextPriority(ctx context.Context) (int, error)
}

type SQLiteFirewallRuleStore struct {
	DB *sql.DB
}

func (s SQLiteFirewallRuleStore) List(ctx context.Context) ([]domain.FirewallRule, error) {
	rows, err := s.DB.QueryContext(ctx, `
		select id, name, kind, template_key, protocol, port, port_range_start, port_range_end,
		       source_cidr, action, enabled, priority, description, created_at, updated_at
		from firewall_rules
		order by priority asc, id asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.FirewallRule
	for rows.Next() {
		var item domain.FirewallRule
		var enabled int
		if err := rows.Scan(
			&item.ID, &item.Name, &item.Kind, &item.TemplateKey, &item.Protocol, &item.Port,
			&item.PortRangeStart, &item.PortRangeEnd, &item.SourceCIDR, &item.Action, &enabled,
			&item.Priority, &item.Description, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		item.Enabled = enabled == 1
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s SQLiteFirewallRuleStore) FindByID(ctx context.Context, id int64) (*domain.FirewallRule, error) {
	row := s.DB.QueryRowContext(ctx, `
		select id, name, kind, template_key, protocol, port, port_range_start, port_range_end,
		       source_cidr, action, enabled, priority, description, created_at, updated_at
		from firewall_rules where id = ?
	`, id)
	var item domain.FirewallRule
	var enabled int
	if err := row.Scan(
		&item.ID, &item.Name, &item.Kind, &item.TemplateKey, &item.Protocol, &item.Port,
		&item.PortRangeStart, &item.PortRangeEnd, &item.SourceCIDR, &item.Action, &enabled,
		&item.Priority, &item.Description, &item.CreatedAt, &item.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	item.Enabled = enabled == 1
	return &item, nil
}

func (s SQLiteFirewallRuleStore) Create(ctx context.Context, item *domain.FirewallRule) error {
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now
	enabled := 0
	if item.Enabled {
		enabled = 1
	}
	result, err := s.DB.ExecContext(ctx, `
		insert into firewall_rules (
			name, kind, template_key, protocol, port, port_range_start, port_range_end,
			source_cidr, action, enabled, priority, description, created_at, updated_at
		)
		values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, item.Name, item.Kind, item.TemplateKey, item.Protocol, item.Port, item.PortRangeStart,
		item.PortRangeEnd, item.SourceCIDR, item.Action, enabled, item.Priority, item.Description,
		item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return err
	}
	item.ID, err = result.LastInsertId()
	return err
}

func (s SQLiteFirewallRuleStore) Update(ctx context.Context, item *domain.FirewallRule) error {
	item.UpdatedAt = time.Now()
	enabled := 0
	if item.Enabled {
		enabled = 1
	}
	_, err := s.DB.ExecContext(ctx, `
		update firewall_rules
		set name = ?, kind = ?, template_key = ?, protocol = ?, port = ?, port_range_start = ?,
		    port_range_end = ?, source_cidr = ?, action = ?, enabled = ?, priority = ?, description = ?, updated_at = ?
		where id = ?
	`, item.Name, item.Kind, item.TemplateKey, item.Protocol, item.Port, item.PortRangeStart,
		item.PortRangeEnd, item.SourceCIDR, item.Action, enabled, item.Priority, item.Description,
		item.UpdatedAt, item.ID)
	return err
}

func (s SQLiteFirewallRuleStore) Delete(ctx context.Context, id int64) error {
	_, err := s.DB.ExecContext(ctx, `delete from firewall_rules where id = ?`, id)
	return err
}

func (s SQLiteFirewallRuleStore) NextPriority(ctx context.Context) (int, error) {
	row := s.DB.QueryRowContext(ctx, `select coalesce(max(priority), 0) + 1 from firewall_rules`)
	var next int
	if err := row.Scan(&next); err != nil {
		return 0, err
	}
	return next, nil
}
