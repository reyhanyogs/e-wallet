package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/reyhanyogs/e-wallet/domain"
	"github.com/reyhanyogs/e-wallet/internal/config"
)

type redisCacheRepository struct {
	rdb *redis.Client
}

func NewRedisClient(config *config.Config) domain.CacheRepository {
	return &redisCacheRepository{
		rdb: redis.NewClient(&redis.Options{
			Addr:     config.Redis.Addr,
			Password: config.Redis.Pass,
			DB:       0,
		}),
	}
}

// Get implements domain.CacheRepository.
func (r *redisCacheRepository) Get(key string) ([]byte, error) {
	val, err := r.rdb.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	return []byte(val), nil
}

// Set implements domain.CacheRepository.
func (r *redisCacheRepository) Set(key string, entry []byte) error {
	return r.rdb.Set(context.Background(), key, entry, 15*time.Minute).Err()
}
