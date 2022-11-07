package service

import (
	"context"

	"github.com/lucaswiix/meli/notifications/dto"
	"github.com/lucaswiix/meli/notifications/repository"
	"github.com/lucaswiix/meli/notifications/utils"
	"go.elastic.co/apm"

	"go.uber.org/zap"
)

//go:generate mockgen -destination=mock/queue.go -package=mock . QueueService
type QueueService interface {
	Send(notification *dto.NotifyDTO, ctx context.Context) error
}

type queueServiceImpl struct {
	repository repository.QueueRepository
}

func NewQueueService(repository repository.QueueRepository) QueueService {
	return &queueServiceImpl{
		repository,
	}
}

func (s *queueServiceImpl) Send(notification *dto.NotifyDTO, ctx context.Context) error {
	span, ctx := apm.StartSpan(ctx, "SendToQueue", "service")
	defer span.End()
	err := s.repository.Send(notification, ctx)
	if err != nil {
		utils.Log.Error("error on try sent notification", zap.Error(err))
		return err
	}
	return nil
}
