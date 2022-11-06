package usecase

import (
	"context"
	"errors"
	"github.com/lucaswiix/meli/notifications/dto"
	"github.com/lucaswiix/meli/notifications/service/mock"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	userID         = "fc2d1670-be5b-4235-bc55-9452bce74a0a"
	notificationID = "c487a2d6-2280-49c0-92cf-738d8cc71366"
	ctx            = context.TODO()
	schedulerDate  = time.Now().Local().Add(time.Hour * time.Duration(1)).Round(0 * time.Second).Format("2006-01-02 15:04:05")
)

func TestSendNotificationSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)

	notifyService := mock.NewMockNotifyService(ctrl)
	queueService := mock.NewMockQueueService(ctrl)
	optOutService := mock.NewMockOptOutService(ctrl)

	notification := &dto.NotifyDTO{
		Message:       "message",
		Title:         "titulo",
		Image:         "cat.png",
		Type:          "web",
		ToUserID:      userID,
		SchedulerDate: schedulerDate,
	}

	optOutService.EXPECT().Is(userID, ctx).Return(false, nil)
	notifyService.EXPECT().Save(notification, ctx).DoAndReturn(func(notification *dto.NotifyDTO, ctx context.Context) error {
		notification.ID = notificationID
		return nil
	}).AnyTimes()
	queueService.EXPECT().Send(notification, ctx).Return(nil)

	target := NewNotificationUsecase(notifyService, queueService, optOutService)

	err := target.SendNotification(notification, ctx)

	assert.NoError(t, err)

}

func TestSendNotificationIsOptOutError(t *testing.T) {
	ctrl := gomock.NewController(t)

	optOutService := mock.NewMockOptOutService(ctrl)

	notification := &dto.NotifyDTO{
		Message:       "message",
		Title:         "titulo",
		Image:         "cat.png",
		Type:          "web",
		ToUserID:      userID,
		SchedulerDate: schedulerDate,
	}

	optOutService.EXPECT().Is(userID, ctx).Return(true, nil)

	target := NewNotificationUsecase(nil, nil, optOutService)

	err := target.SendNotification(notification, ctx)

	assert.EqualError(t, err, "opt-out user")
}

func TestSendNotificationFunctionError(t *testing.T) {
	ctrl := gomock.NewController(t)

	optOutService := mock.NewMockOptOutService(ctrl)

	notification := &dto.NotifyDTO{
		Message:       "message",
		Title:         "titulo",
		Image:         "cat.png",
		Type:          "web",
		ToUserID:      userID,
		SchedulerDate: schedulerDate,
	}

	optOutService.EXPECT().Is(userID, ctx).Return(false, errors.New("errors"))

	target := NewNotificationUsecase(nil, nil, optOutService)

	err := target.SendNotification(notification, ctx)

	assert.EqualError(t, err, "errors")
}

func TestSendNotificationOnStorageError(t *testing.T) {
	ctrl := gomock.NewController(t)

	notifyService := mock.NewMockNotifyService(ctrl)
	optOutService := mock.NewMockOptOutService(ctrl)

	notification := &dto.NotifyDTO{
		Message:       "message",
		Title:         "titulo",
		Image:         "cat.png",
		Type:          "web",
		ToUserID:      userID,
		SchedulerDate: schedulerDate,
	}

	optOutService.EXPECT().Is(userID, ctx).Return(false, nil)
	notifyService.EXPECT().Save(notification, ctx).Return(errors.New("error"))

	target := NewNotificationUsecase(notifyService, nil, optOutService)

	err := target.SendNotification(notification, ctx)

	assert.EqualError(t, err, "error")

}

func TestSendNotificationSetStatusFailedWhenErrorToSendToQueue(t *testing.T) {
	ctrl := gomock.NewController(t)

	notifyService := mock.NewMockNotifyService(ctrl)
	queueService := mock.NewMockQueueService(ctrl)
	optOutService := mock.NewMockOptOutService(ctrl)

	notification := &dto.NotifyDTO{
		Message:       "message",
		Title:         "titulo",
		Image:         "cat.png",
		Type:          "web",
		ToUserID:      userID,
		SchedulerDate: schedulerDate,
	}

	optOutService.EXPECT().Is(userID, ctx).Return(false, nil)
	notifyService.EXPECT().Save(notification, ctx).DoAndReturn(func(notification *dto.NotifyDTO, ctx context.Context) error {
		notification.ID = notificationID
		return nil
	}).AnyTimes()
	queueService.EXPECT().Send(notification, ctx).Return(errors.New("error rbq"))

	target := NewNotificationUsecase(notifyService, queueService, optOutService)

	err := target.SendNotification(notification, ctx)

	assert.NoError(t, err)
	assert.Equal(t, notification.Status, dto.Failed)

}

func TestSendNotificationOnUpdateError(t *testing.T) {
	ctrl := gomock.NewController(t)

	notifyService := mock.NewMockNotifyService(ctrl)
	queueService := mock.NewMockQueueService(ctrl)
	optOutService := mock.NewMockOptOutService(ctrl)

	notification := &dto.NotifyDTO{
		Message:       "message",
		Title:         "titulo",
		Image:         "cat.png",
		Type:          "web",
		ToUserID:      userID,
		SchedulerDate: schedulerDate,
	}

	first := optOutService.EXPECT().Is(userID, ctx).Return(false, nil)
	second := notifyService.EXPECT().Save(notification, ctx).DoAndReturn(func(notification *dto.NotifyDTO, ctx context.Context) error {
		notification.ID = notificationID
		return nil
	})
	third := queueService.EXPECT().Send(notification, ctx).Return(nil)
	fourth := notifyService.EXPECT().Save(notification, ctx).Return(errors.New("error update"))

	gomock.InOrder(
		first,
		second,
		third,
		fourth,
	)
	target := NewNotificationUsecase(notifyService, queueService, optOutService)

	err := target.SendNotification(notification, ctx)

	assert.EqualError(t, err, "error update")

}
