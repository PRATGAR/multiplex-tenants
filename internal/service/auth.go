package service

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/pbkdf2"

	"tenantly/internal/config"
	"tenantly/internal/middleware"
	"tenantly/internal/model"
	"tenantly/internal/store"
	"tenantly/internal/tokens"
)

type Auth struct {
	users *store.Users
	rdb   *redis.Client
	salt  string
	ttl   time.Duration
}

func NewAuth(users *store.Users, rdb *redis.Client, salt string, ttl time.Duration) *Auth {
	return &Auth{users: users, rdb: rdb, salt: salt, ttl: ttl}
}

func HashPassword(salt, password string) string {
	sum := pbkdf2.Key([]byte(password), []byte(salt), config.PasswordIterations, 32, sha256.New)
	return hex.EncodeToString(sum)
}

func (a *Auth) Login(ctx context.Context, tenantSlug, email, password string) (string, error) {
	u, err := a.users.ByLogin(ctx, tenantSlug, email)
	if err != nil {
		return "", ErrBadCredentials
	}
	candidate := HashPassword(a.salt, password)
	if subtle.ConstantTimeCompare([]byte(u.PasswordHash), []byte(candidate)) != 1 {
		return "", ErrBadCredentials
	}
	token, err := tokens.Secret()
	if err != nil {
		return "", err
	}
	payload, err := json.Marshal(model.Session{UserID: u.ID, TenantID: u.TenantID, Role: u.Role})
	if err != nil {
		return "", err
	}
	if err := a.rdb.Set(ctx, middleware.SessionPrefix+token, payload, a.ttl).Err(); err != nil {
		return "", err
	}
	return token, nil
}

func (a *Auth) Logout(ctx context.Context, token string) error {
	return a.rdb.Del(ctx, middleware.SessionPrefix+token).Err()
}
