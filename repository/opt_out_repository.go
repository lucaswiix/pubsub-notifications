package repository

import (
	"context"
	"fmt"
	"meli/notifications/utils"

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

	return r.redisDB.Set(ctx, fmt.Sprintf("opt-out:%s", userID), optOutVal, 0).Err()
}

func (r *implOptOutRepository) Del(userID string, ctx context.Context) error {

	return r.redisDB.Del(ctx, fmt.Sprintf("opt-out:%s", userID)).Err()
}

func (r *implOptOutRepository) Is(userID string, ctx context.Context) (bool, error) {
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
