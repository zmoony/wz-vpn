package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/zmoony/pi-gateway/internal/domain"
)

type WireGuardPeerStore interface {
	List(ctx context.Context) ([]domain.WireGuardPeer, error)
	FindByID(ctx context.Context, id int64) (*domain.WireGuardPeer, error)
	FindByName(ctx context.Context, name string) (*domain.WireGuardPeer, error)
	Create(ctx context.Context, peer *domain.WireGuardPeer) error
	Update(ctx context.Context, peer *domain.WireGuardPeer) error
	Delete(ctx context.Context, id int64) error
}

type SQLiteWireGuardPeerStore struct {
	DB *sql.DB
}

func (s SQLiteWireGuardPeerStore) List(ctx context.Context) ([]domain.WireGuardPeer, error) {
	rows, err := s.DB.QueryContext(ctx, `
		select id, name, client_ipv4, client_ipv6, public_key, private_key_encrypted,
		       preshared_key_encrypted, dns, allowed_ips, endpoint, persistent_keepalive,
		       enabled, description, created_at, updated_at
		from wg_peers
		order by id asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var peers []domain.WireGuardPeer
	for rows.Next() {
		var peer domain.WireGuardPeer
		var enabled int
		if err := rows.Scan(
			&peer.ID,
			&peer.Name,
			&peer.ClientIPv4,
			&peer.ClientIPv6,
			&peer.PublicKey,
			&peer.PrivateKeyEncrypted,
			&peer.PresharedKeyEnc,
			&peer.DNS,
			&peer.AllowedIPs,
			&peer.Endpoint,
			&peer.PersistentKeepalive,
			&enabled,
			&peer.Description,
			&peer.CreatedAt,
			&peer.UpdatedAt,
		); err != nil {
			return nil, err
		}
		peer.Enabled = enabled == 1
		peers = append(peers, peer)
	}

	return peers, rows.Err()
}

func (s SQLiteWireGuardPeerStore) FindByID(ctx context.Context, id int64) (*domain.WireGuardPeer, error) {
	return s.findOne(ctx, `where id = ?`, id)
}

func (s SQLiteWireGuardPeerStore) FindByName(ctx context.Context, name string) (*domain.WireGuardPeer, error) {
	return s.findOne(ctx, `where name = ?`, name)
}

func (s SQLiteWireGuardPeerStore) findOne(ctx context.Context, filter string, value any) (*domain.WireGuardPeer, error) {
	query := `
		select id, name, client_ipv4, client_ipv6, public_key, private_key_encrypted,
		       preshared_key_encrypted, dns, allowed_ips, endpoint, persistent_keepalive,
		       enabled, description, created_at, updated_at
		from wg_peers
	` + filter

	row := s.DB.QueryRowContext(ctx, query, value)
	var peer domain.WireGuardPeer
	var enabled int
	if err := row.Scan(
		&peer.ID,
		&peer.Name,
		&peer.ClientIPv4,
		&peer.ClientIPv6,
		&peer.PublicKey,
		&peer.PrivateKeyEncrypted,
		&peer.PresharedKeyEnc,
		&peer.DNS,
		&peer.AllowedIPs,
		&peer.Endpoint,
		&peer.PersistentKeepalive,
		&enabled,
		&peer.Description,
		&peer.CreatedAt,
		&peer.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	peer.Enabled = enabled == 1
	return &peer, nil
}

func (s SQLiteWireGuardPeerStore) Create(ctx context.Context, peer *domain.WireGuardPeer) error {
	now := time.Now()
	peer.CreatedAt = now
	peer.UpdatedAt = now
	enabled := 0
	if peer.Enabled {
		enabled = 1
	}

	result, err := s.DB.ExecContext(ctx, `
		insert into wg_peers (
			name, client_ipv4, client_ipv6, public_key, private_key_encrypted,
			preshared_key_encrypted, dns, allowed_ips, endpoint, persistent_keepalive,
			enabled, description, created_at, updated_at
		)
		values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		peer.Name,
		peer.ClientIPv4,
		peer.ClientIPv6,
		peer.PublicKey,
		peer.PrivateKeyEncrypted,
		peer.PresharedKeyEnc,
		peer.DNS,
		peer.AllowedIPs,
		peer.Endpoint,
		peer.PersistentKeepalive,
		enabled,
		peer.Description,
		peer.CreatedAt,
		peer.UpdatedAt,
	)
	if err != nil {
		return err
	}

	peer.ID, err = result.LastInsertId()
	return err
}

func (s SQLiteWireGuardPeerStore) Update(ctx context.Context, peer *domain.WireGuardPeer) error {
	peer.UpdatedAt = time.Now()
	enabled := 0
	if peer.Enabled {
		enabled = 1
	}

	_, err := s.DB.ExecContext(ctx, `
		update wg_peers
		set name = ?, client_ipv4 = ?, client_ipv6 = ?, public_key = ?, private_key_encrypted = ?,
		    preshared_key_encrypted = ?, dns = ?, allowed_ips = ?, endpoint = ?, persistent_keepalive = ?,
		    enabled = ?, description = ?, updated_at = ?
		where id = ?
	`,
		peer.Name,
		peer.ClientIPv4,
		peer.ClientIPv6,
		peer.PublicKey,
		peer.PrivateKeyEncrypted,
		peer.PresharedKeyEnc,
		peer.DNS,
		peer.AllowedIPs,
		peer.Endpoint,
		peer.PersistentKeepalive,
		enabled,
		peer.Description,
		peer.UpdatedAt,
		peer.ID,
	)
	return err
}

func (s SQLiteWireGuardPeerStore) Delete(ctx context.Context, id int64) error {
	_, err := s.DB.ExecContext(ctx, `delete from wg_peers where id = ?`, id)
	return err
}
