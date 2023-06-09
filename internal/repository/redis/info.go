package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vgekko/ani-go/internal/entity"
)

const notExists = iota

type InfoRepo struct {
	db  *redis.Client
	ttl time.Duration
}

func NewInfoRepo(db *redis.Client, ttl time.Duration) *InfoRepo {
	return &InfoRepo{db, ttl}
}

func (r *InfoRepo) Cache(ctx context.Context, key string, value entity.TitleInfos) error {
	err := r.db.Set(ctx, key, value, r.ttl).Err()
	if err != nil {
		return fmt.Errorf("InfoRedisRepo.Cache: %w", err)
	}

	return nil
}

func (r *InfoRepo) Lookup(ctx context.Context, key string) bool {
	data, _ := r.db.Exists(ctx, key).Result()

	return data != notExists
}

func (r *InfoRepo) FromCache(ctx context.Context, key string) (entity.TitleInfos, error) {
	var ti entity.TitleInfos
	var asBytes []byte

	if err := r.db.Get(ctx, key).Scan(&asBytes); err != nil {
		return nil, fmt.Errorf("InfoRedisRepo.FromCache: %w", err)
	}

	if err := json.Unmarshal(asBytes, &ti); err != nil {
		return nil, fmt.Errorf("InfoRedisRepo.FromCache: %w", err)
	}

	return ti, nil
}
