package service

import (
	"context"
	"meli/notifications/dto"
	"meli/notifications/repository"
	"meli/notifications/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=mock/notify.go -package=mock . NotifyService
type NotifyService interface {
	Save(notify *dto.NotifyDTO, ctx context.Context) error
}

type notifyServiceImpl struct {
	repository repository.NotifyRepository
}

func NewNotifyService(repository repository.NotifyRepository) NotifyService {
	return &notifyServiceImpl{
		repository,
	}
}

func (s *notifyServiceImpl) Save(notification *dto.NotifyDTO, ctx context.Context) error {
	if notification.ID == "" {
		notification.ID = uuid.New().String()
	}
	err := s.repository.Save(notification, ctx)
	if err != nil {
		utils.Log.Error("error on try save on database", zap.Error(err))
		return err
	}
	return nil
}
