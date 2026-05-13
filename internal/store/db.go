package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func Open(databasePath string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(databasePath), 0o755); err != nil {
		return nil, fmt.Errorf("create database dir: %w", err)
	}

	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	if err := migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	statements := []string{
		`create table if not exists users (
			id integer primary key autoincrement,
			username text not null unique,
			password_hash text not null,
			created_at datetime not null,
			updated_at datetime not null
		);`,
		`create table if not exists wg_peers (
			id integer primary key autoincrement,
			name text not null unique,
			client_ipv4 text not null unique,
			client_ipv6 text not null default '',
			public_key text not null unique,
			private_key_encrypted text not null,
			preshared_key_encrypted text not null default '',
			dns text not null default '',
			allowed_ips text not null default '',
			endpoint text not null default '',
			persistent_keepalive integer not null default 25,
			enabled integer not null default 1,
			description text not null default '',
			created_at datetime not null,
			updated_at datetime not null
		);`,
		`create table if not exists settings (
			key text primary key,
			value text not null,
			updated_at datetime not null
		);`,
		`create table if not exists proxy_hosts (
			id integer primary key autoincrement,
			name text not null unique,
			server_name text not null unique,
			upstream_url text not null,
			certificate_cert_path text not null,
			certificate_key_path text not null,
			enabled integer not null default 1,
			description text not null default '',
			last_apply_status text not null default '',
			last_apply_error text not null default '',
			last_applied_at datetime,
			created_at datetime not null,
			updated_at datetime not null
		);`,
		`create table if not exists ddns_configs (
			id integer primary key autoincrement,
			provider text not null,
			access_key_id text not null,
			access_key_secret_encrypted text not null,
			domain text not null,
			subdomain text not null,
			enabled integer not null default 1,
			check_interval_seconds integer not null default 300,
			last_known_ipv6 text not null default '',
			last_status text not null default '',
			last_error text not null default '',
			last_synced_at datetime,
			created_at datetime not null,
			updated_at datetime not null,
			unique(domain, subdomain)
		);`,
		`create table if not exists cert_configs (
			id integer primary key autoincrement,
			root_domain text not null unique,
			provider text not null,
			access_key_id text not null,
			access_key_secret_encrypted text not null,
			install_dir text not null,
			enabled_auto_renew integer not null default 1,
			fullchain_path text not null default '',
			private_key_path text not null default '',
			last_issue_status text not null default '',
			last_issue_error text not null default '',
			last_issued_at datetime,
			created_at datetime not null,
			updated_at datetime not null
		);`,
		`create table if not exists firewall_rules (
			id integer primary key autoincrement,
			name text not null,
			kind text not null,
			template_key text not null default '',
			protocol text not null default '',
			port integer not null default 0,
			port_range_start integer not null default 0,
			port_range_end integer not null default 0,
			source_cidr text not null default '',
			action text not null default 'accept',
			enabled integer not null default 1,
			priority integer not null default 1,
			description text not null default '',
			created_at datetime not null,
			updated_at datetime not null
		);`,
	}

	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return fmt.Errorf("migrate database: %w", err)
		}
	}

	return nil
}
