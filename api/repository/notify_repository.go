package repository

import (
	"context"
	"encoding/json"
	"meli/notifications/dto"
	"meli/notifications/utils"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=mock/notify.go -package=mock . NotifyRepository
type NotifyRepository interface {
	Save(notification *dto.NotifyDTO, ctx context.Context) error
}

type implNotifyRepository struct {
	redisDB *redis.Client
}

func NewNotifyRepository(redisDB *redis.Client) NotifyRepository {
	return &implNotifyRepository{redisDB}
}

func (r *implNotifyRepository) Save(notification *dto.NotifyDTO, ctx context.Context) error {
	json, err := json.Marshal(notification)
	if err != nil {
		utils.Log.Error("failed to decode notification interface", zap.Error(err))
		return err
	}
	return r.redisDB.Set(ctx, notification.ID, json, 0).Err()
}
