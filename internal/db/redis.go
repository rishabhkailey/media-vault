package db

import (
	"context"
	"fmt"
	"time"

	oredis "github.com/go-oauth2/redis/v4"
	gsredis "github.com/go-session/redis/v3"
	"github.com/go-session/session/v3"

	"github.com/go-redis/redis/v8"
	"github.com/rishabhkailey/media-service/internal/config"
)

type RedisStore struct {
	*oredis.TokenStore
	Client *redis.Client
}

func NewRedisTokenStore(config config.RedisCacheConfig) (*RedisStore, error) {
	redis := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.Db,
	})
	if err := redis.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	store := &RedisStore{oredis.NewRedisStoreWithCli(redis), redis}
	return store, nil
}

func NewRedisSessionStore(config config.RedisCacheConfig) session.ManagerStore {
	return gsredis.NewRedisStore(&gsredis.Options{
		Addr:     fmt.Sprintf("%v:%v", config.Host, config.Port),
		DB:       config.Db,
		Password: config.Password,
	})
}

func getNonceKey(userID string) string {
	return fmt.Sprintf("code:%v:nonce", userID)
}

func (store *RedisStore) CreateNonceByUserID(ctx context.Context, userID string, nonce string) error {
	// todo - expire time? for new request new nonce will be generated and for refresh as well i'm assuming new code verifier, state and nonce will be generated.
	// i think for refresh we don't verify nonce, state and code_verifier
	expire := 10 * time.Minute
	err := store.Client.Set(ctx, getNonceKey(userID), nonce, expire).Err()
	if err != nil {
		return fmt.Errorf("[CreateNonce]: failed to save nonce: %w", err)
	}
	return nil
}

func (store *RedisStore) GetNonceByUserID(ctx context.Context, userID string) (nonce string, err error) {
	nonce, err = store.Client.Get(ctx, getNonceKey(userID)).Result()
	if err != nil {
		if err == redis.Nil {
			return nonce, fmt.Errorf("[GetNonceByCode]: key doesn't exist: %w", err)
		}
		return nonce, fmt.Errorf("[GetNonceByCode]: failed to get nonce: %w", err)
	}
	return nonce, err
}

func (store *RedisStore) GetAndDeleteNonceByUserID(ctx context.Context, userID string) (nonce string, err error) {
	nonce, err = store.GetNonceByUserID(ctx, userID)
	if err != nil {
		// in case of failure we will not be deleting key if client want's to retry the request
		// it will be automatticaly expire after 10 minutes
		return nonce, err
	}
	err = store.Client.Del(ctx, getNonceKey(userID)).Err()
	if err != nil {
		return nonce, fmt.Errorf("[GetAndDeleteNonceByCode]: key deletion failed: %w", err)
	}
	return nonce, err
}
