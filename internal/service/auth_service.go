package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/zmoony/pi-gateway/internal/domain"
	"github.com/zmoony/pi-gateway/internal/store"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
	Users          store.UserStore
	JWTSecret      string
	SessionTTL     time.Duration
	DefaultUser    string
	DefaultPassRaw string
}

func (s AuthService) EnsureDefaultAdmin(ctx context.Context) error {
	count, err := s.Users.Count(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hash, err := HashPassword(s.DefaultPassRaw)
	if err != nil {
		return err
	}

	return s.Users.Create(ctx, &domain.User{
		Username:     s.DefaultUser,
		PasswordHash: hash,
	})
}

func (s AuthService) Login(ctx context.Context, username, password string) (string, *domain.User, error) {
	user, err := s.Users.FindByUsername(ctx, username)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, ErrInvalidCredentials
	}

	if err := ComparePassword(user.PasswordHash, password); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	claims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(s.SessionTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return "", nil, err
	}

	return signed, user, nil
}

func (s AuthService) ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(s.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func ComparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
