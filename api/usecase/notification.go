package usecase

import (
	"context"
	"fmt"

	"github.com/lucaswiix/meli/notifications/dto"
	"github.com/lucaswiix/meli/notifications/service"
	"github.com/lucaswiix/meli/notifications/utils"
	"go.elastic.co/apm"

	"go.uber.org/zap"
)

//go:generate mockgen -destination=mock/notification.go -package=mock . NotificationUseCase
type NotificationUseCase interface {
	SendNotification(notify *dto.NotifyDTO, ctx context.Context) error
}

type notificationUsecaseImpl struct {
	queueService  service.QueueService
	optOutService service.UserService
}

func NewNotificationUsecase(queueService service.QueueService, userService service.UserService) NotificationUseCase {
	return &notificationUsecaseImpl{
		queueService,
		userService,
	}
}
func (s *notificationUsecaseImpl) SendNotification(notification *dto.NotifyDTO, ctx context.Context) error {
	span, ctx := apm.StartSpan(ctx, "SentNotification", "usecase")
	defer span.End()

	if err := IsOptOut(notification, ctx, s.optOutService); err != nil {
		return err
	}

	if err := s.queueService.Send(notification, ctx); err != nil {
		utils.Log.Error("failed to send to queue", zap.Error(err))
		return err
	}

	utils.Log.Info(fmt.Sprintf("notification sent to user %s type %s status %s", notification.ToUserID, notification.Type, notification.Status))
	return nil
}

func IsOptOut(notification *dto.NotifyDTO, ctx context.Context, userService service.UserService) error {
	user, err := userService.Get(notification.ToUserID, ctx)
	if err != nil {
		utils.Log.Error("error on try to get opt-out user", zap.Error(err))
		return err
	}
	if user.ID != "" && user.IsOptOut {
		utils.Log.Debug("is opt-out user, skipped", zap.Error(err))
		return utils.ErrOptOutUser
	}
	return nil
}
