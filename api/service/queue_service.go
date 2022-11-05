package service

import (
	"context"
	"meli/notifications/dto"
	"meli/notifications/repository"
	"meli/notifications/utils"

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

	err := s.repository.Send(notification, ctx)
	if err != nil {
		utils.Log.Error("error on try sent notification", zap.Error(err))
		return err
	}
	return nil
}