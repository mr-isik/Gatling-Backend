package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/mr-isik/gatling-backend/internal/domain"
	"github.com/mr-isik/gatling-backend/internal/repository"
	"github.com/redis/go-redis/v9"
)

type runStateRepository struct {
	client *redis.Client
}

func NewRunStateRepository(client *redis.Client) repository.RunStateRepository {
	return &runStateRepository{client: client}
}

func (r *runStateRepository) SetState(ctx context.Context, runID string, state domain.RunStatus, ttl time.Duration) error {
	key := fmt.Sprintf("run:state:%s", runID)
	return r.client.Set(ctx, key, string(state), ttl).Err()
}

func (r *runStateRepository) GetState(ctx context.Context, runID string) (domain.RunStatus, error) {
	key := fmt.Sprintf("run:state:%s", runID)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", domain.ErrNotFound
	}
	if err != nil {
		return "", err
	}
	return domain.RunStatus(val), nil
}

func (r *runStateRepository) DeleteState(ctx context.Context, runID string) error {
	key := fmt.Sprintf("run:state:%s", runID)
	return r.client.Del(ctx, key).Err()
}
