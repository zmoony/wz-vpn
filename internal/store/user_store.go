package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/zmoony/pi-gateway/internal/domain"
)

type UserStore interface {
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Count(ctx context.Context) (int64, error)
}

type SQLiteUserStore struct {
	DB *sql.DB
}

func (s SQLiteUserStore) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	row := s.DB.QueryRowContext(ctx, `
		select id, username, password_hash, created_at, updated_at
		from users
		where username = ?
	`, username)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (s SQLiteUserStore) Create(ctx context.Context, user *domain.User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	result, err := s.DB.ExecContext(ctx, `
		insert into users (username, password_hash, created_at, updated_at)
		values (?, ?, ?, ?)
	`, user.Username, user.PasswordHash, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	user.ID, err = result.LastInsertId()
	return err
}

func (s SQLiteUserStore) Count(ctx context.Context) (int64, error) {
	row := s.DB.QueryRowContext(ctx, `select count(*) from users`)
	var count int64
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}
