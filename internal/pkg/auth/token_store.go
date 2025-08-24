package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type TokenStore struct {
	rdb *redis.Client
}

func NewTokenStore(rdb *redis.Client) *TokenStore {
	return &TokenStore{rdb: rdb}
}

// SaveRefreshToken simpan RT ke Redis
func (ts *TokenStore) SaveRefreshToken(ctx context.Context, userID uuid.UUID, jti string, ttl time.Duration) error {
	key := fmt.Sprintf("rt:%s:%s", userID.String(), jti)
	return ts.rdb.Set(ctx, key, "valid", ttl).Err()
}

// VerifyRefreshToken cek apakah RT masih valid di Redis
func (ts *TokenStore) VerifyRefreshToken(ctx context.Context, userID, jti string) (bool, error) {
	key := fmt.Sprintf("rt:%s:%s", userID, jti)
	val, err := ts.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil // tidak ada = invalid
	}
	if err != nil {
		return false, err
	}
	return val == "valid", nil
}

// RevokeRefreshToken hapus RT dari Redis (logout)
func (ts *TokenStore) RevokeRefreshToken(ctx context.Context, userID, jti string) error {
	key := fmt.Sprintf("rt:%s:%s", userID, jti)
	return ts.rdb.Del(ctx, key).Err()
}

// RevokeAllRefreshTokens hapus semua RT milik user
func (ts *TokenStore) RevokeAllRefreshTokens(ctx context.Context, userID string) error {
	pattern := fmt.Sprintf("rt:%s:*", userID)
	iter := ts.rdb.Scan(ctx, 0, pattern, 0).Iterator()
	var firstErr error
	for iter.Next(ctx) {
		if err := ts.rdb.Del(ctx, iter.Val()).Err(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	if err := iter.Err(); err != nil && firstErr == nil {
		firstErr = err
	}
	return firstErr
}
