package service

import (
	"context"

	"github.com/lucaswiix/meli/notifications/dto"
	"github.com/lucaswiix/meli/notifications/repository"
	"github.com/lucaswiix/meli/notifications/utils"
	"go.elastic.co/apm"

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
	span, ctx := apm.StartSpan(ctx, "SaveNotification", "service")
	defer span.End()
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
