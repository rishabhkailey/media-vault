package db

import (
	"context"
	"fmt"
	"time"

	gsredis "github.com/go-session/redis/v3"
	"github.com/go-session/session/v3"

	"github.com/go-redis/redis/v8"
	"github.com/rishabhkailey/media-service/internal/config"
)

type RedisStore struct {
	Client *redis.Client
}

func NewRedisClient(config config.RedisCacheConfig) (*redis.Client, error) {
	redis := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.Db,
	})
	if err := redis.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return redis, nil
}

func NewRedisStore(config config.RedisCacheConfig) (*RedisStore, error) {
	redis := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.Db,
	})
	if err := redis.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	store := &RedisStore{Client: redis}
	return store, nil
}

func NewRedisSessionStore(config config.RedisCacheConfig) session.ManagerStore {
	return gsredis.NewRedisStore(&gsredis.Options{
		Addr:     fmt.Sprintf("%v:%v", config.Host, config.Port),
		DB:       config.Db,
		Password: config.Password,
	})
}

func mediaTypeKey(fileName string) string {
	return fmt.Sprintf("mediaType:%s", fileName)
}

func (store *RedisStore) GetMediaType(ctx context.Context, fileName string) (mediaType string, err error) {
	mediaType, err = store.Client.Get(ctx, mediaTypeKey(fileName)).Result()
	if err != nil {
		if err == redis.Nil {
			return mediaType, fmt.Errorf("[GetMediaType]: key doesn't exist: %w", err)
		}
		return mediaType, fmt.Errorf("[GetMediaType]: failed to get nonce: %w", err)
	}
	return mediaType, err
}

func (store *RedisStore) SetMediaType(ctx context.Context, fileName string, mediaType string) error {
	expire := 1 * time.Hour
	err := store.Client.Set(ctx, mediaTypeKey(fileName), mediaType, expire).Err()
	if err != nil {
		return fmt.Errorf("[SetMediaType]: failed to save nonce: %w", err)
	}
	return nil
}
