package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/erupshis/key_keeper/internal/common/auth/models"
	"github.com/erupshis/key_keeper/internal/common/auth/storage"
	"github.com/erupshis/key_keeper/internal/common/hasher"
	"github.com/erupshis/key_keeper/internal/common/jwtgenerator"
)

const (
	TokenHeader = "Authorization"
	TokenType   = "Bearer "
	UserID      = "user_id"
)

type Config struct {
	Storage storage.BaseAuthStorage
	JWT     *jwtgenerator.JwtGenerator
	Hasher  *hasher.Hasher
}

type Manager struct {
	storage storage.BaseAuthStorage
	jwt     *jwtgenerator.JwtGenerator
	hasher  *hasher.Hasher
}

func NewManager(cfg *Config) *Manager {
	return &Manager{
		storage: cfg.Storage,
		jwt:     cfg.JWT,
		hasher:  cfg.Hasher,
	}
}

func (m *Manager) Login(ctx context.Context, user *models.User) (string, error) {
	userData, err := m.storage.GetUserByLogin(ctx, user.Login)
	if err != nil {
		return "", fmt.Errorf("check user in db by login: %w", err)
	}

	hashedPwd, err := m.hasher.HashMsg([]byte(user.Password))
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	if userData.Password != hashedPwd {
		return "", ErrMismatchPassword
	}

	token, err := m.jwt.BuildJWTString(user.ID)
	if err != nil {
		return "", fmt.Errorf("create session token: %w", err)
	}
	return addBearerPrefix(token), nil
}

func (m *Manager) Register(ctx context.Context, user *models.User) error {
	userData, err := m.storage.GetUserByLogin(ctx, user.Login)
	if err != nil && !errors.Is(err, storage.ErrUserNotFound) {
		return fmt.Errorf("check user in db by login: %w", err)
	}

	if userData != nil {
		return ErrLoginOccupied
	}

	user.Password, err = m.hasher.HashMsg([]byte(user.Password))
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	if err = m.storage.AddUser(ctx, user); err != nil {
		return fmt.Errorf("add new user: %w", err)
	}

	return nil
}
