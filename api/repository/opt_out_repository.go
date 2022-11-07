package repository

import (
	"context"
	"fmt"

	"github.com/lucaswiix/meli/notifications/utils"
	"go.elastic.co/apm"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=mock/opt_out.go -package=mock . OptOutRepository
type OptOutRepository interface {
	Set(userID string, ctx context.Context) error
	Del(userID string, ctx context.Context) error
	Is(userID string, ctx context.Context) (bool, error)
}

type implOptOutRepository struct {
	redisDB *redis.Client
}

func NewOptOutRepository(redisDB *redis.Client) OptOutRepository {
	return &implOptOutRepository{redisDB}
}

const (
	optOutVal = "true"
)

func (r *implOptOutRepository) Set(userID string, ctx context.Context) error {
	span, ctx := apm.StartSpan(ctx, "SetOptOut", "repository")
	defer span.End()
	return r.redisDB.Set(ctx, fmt.Sprintf("opt-out:%s", userID), optOutVal, 0).Err()
}

func (r *implOptOutRepository) Del(userID string, ctx context.Context) error {
	span, ctx := apm.StartSpan(ctx, "DelOptOut", "repository")
	defer span.End()
	return r.redisDB.Del(ctx, fmt.Sprintf("opt-out:%s", userID)).Err()
}

func (r *implOptOutRepository) Is(userID string, ctx context.Context) (bool, error) {
	span, ctx := apm.StartSpan(ctx, "IsOptOut", "repository")
	defer span.End()
	_, err := r.redisDB.Get(ctx, fmt.Sprintf("opt-out:%s", userID)).Result()
	switch {
	case err == redis.Nil:
		return false, nil
	case err != nil && err != redis.Nil:
		utils.Log.Error("error on search for opt-out user", zap.Error(err))
		return false, err
	default:
		return true, nil
	}
}
