package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type SettingStore interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
}

type SQLiteSettingStore struct {
	DB *sql.DB
}

func (s SQLiteSettingStore) Get(ctx context.Context, key string) (string, error) {
	row := s.DB.QueryRowContext(ctx, `select value from settings where key = ?`, key)
	var value string
	if err := row.Scan(&value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	return value, nil
}

func (s SQLiteSettingStore) Set(ctx context.Context, key string, value string) error {
	_, err := s.DB.ExecContext(ctx, `
		insert into settings (key, value, updated_at)
		values (?, ?, ?)
		on conflict(key) do update set value = excluded.value, updated_at = excluded.updated_at
	`, key, value, time.Now())
	return err
}
