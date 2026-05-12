package service

import (
	"context"
	"testing"
	"time"

	"github.com/zmoony/pi-gateway/internal/domain"
)

type authUserStoreStub struct {
	user *domain.User
}

func (s *authUserStoreStub) FindByUsername(context.Context, string) (*domain.User, error) {
	return s.user, nil
}

func (s *authUserStoreStub) Create(context.Context, *domain.User) error { return nil }
func (s *authUserStoreStub) Count(context.Context) (int64, error)       { return 1, nil }

func TestHashAndComparePassword(t *testing.T) {
	hash, err := HashPassword("secret123")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if hash == "secret123" {
		t.Fatalf("expected hashed password to differ from raw password")
	}

	if err := ComparePassword(hash, "secret123"); err != nil {
		t.Fatalf("ComparePassword() error = %v", err)
	}
}

func TestLoginReturnsJWTForValidCredentials(t *testing.T) {
	hash, err := HashPassword("secret123")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	service := AuthService{
		Users: &authUserStoreStub{
			user: &domain.User{ID: 1, Username: "admin", PasswordHash: hash},
		},
		JWTSecret:  "jwt-secret",
		SessionTTL: 2 * time.Hour,
	}

	token, user, err := service.Login(context.Background(), "admin", "secret123")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if token == "" {
		t.Fatalf("expected non-empty token")
	}
	if user == nil || user.Username != "admin" {
		t.Fatalf("expected matching user, got %#v", user)
	}
}
