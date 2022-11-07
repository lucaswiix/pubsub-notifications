package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lucaswiix/meli/notifications/dto"
	"github.com/lucaswiix/meli/notifications/service/mock"

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

	queueService := mock.NewMockQueueService(ctrl)
	userService := mock.NewMockUserService(ctrl)

	notification := &dto.NotifyDTO{
		Message:       "message",
		Title:         "titulo",
		Image:         "cat.png",
		Type:          "web",
		ToUserID:      userID,
		SchedulerDate: schedulerDate,
	}

	user := dto.User{
		ID:        userID,
		CreatedAt: time.Now().Format(time.RFC822),
		IsOptOut:  false,
	}

	userService.EXPECT().Get(userID, ctx).Return(user, nil)
	queueService.EXPECT().Send(notification, ctx).Return(nil)

	target := NewNotificationUsecase(queueService, userService)

	err := target.SendNotification(notification, ctx)

	assert.NoError(t, err)

}

func TestSendNotificationIsOptOutError(t *testing.T) {
	ctrl := gomock.NewController(t)

	userService := mock.NewMockUserService(ctrl)

	notification := &dto.NotifyDTO{
		Message:       "message",
		Title:         "titulo",
		Image:         "cat.png",
		Type:          "web",
		ToUserID:      userID,
		SchedulerDate: schedulerDate,
	}

	user := dto.User{
		ID:        userID,
		IsOptOut:  true,
		CreatedAt: time.Now().Format(time.RFC822),
	}

	userService.EXPECT().Get(userID, ctx).Return(user, nil)

	target := NewNotificationUsecase(nil, userService)

	err := target.SendNotification(notification, ctx)

	assert.EqualError(t, err, "opt-out user")
}

func TestSendNotificationFunctionError(t *testing.T) {
	ctrl := gomock.NewController(t)

	userService := mock.NewMockUserService(ctrl)

	notification := &dto.NotifyDTO{
		Message:       "message",
		Title:         "titulo",
		Image:         "cat.png",
		Type:          "web",
		ToUserID:      userID,
		SchedulerDate: schedulerDate,
	}

	userService.EXPECT().Get(userID, ctx).Return(false, errors.New("errors"))

	target := NewNotificationUsecase(nil, userService)

	err := target.SendNotification(notification, ctx)

	assert.EqualError(t, err, "errors")
}

func TestSendNotificationSetStatusFailedWhenErrorToSendToQueue(t *testing.T) {
	ctrl := gomock.NewController(t)

	queueService := mock.NewMockQueueService(ctrl)
	userService := mock.NewMockUserService(ctrl)

	notification := &dto.NotifyDTO{
		Message:       "message",
		Title:         "titulo",
		Image:         "cat.png",
		Type:          "web",
		ToUserID:      userID,
		SchedulerDate: schedulerDate,
	}

	user := dto.User{
		ID:        userID,
		IsOptOut:  false,
		CreatedAt: time.Now().Format(time.RFC822),
	}

	userService.EXPECT().Get(userID, ctx).Return(user, nil)
	queueService.EXPECT().Send(notification, ctx).Return(errors.New("error rbq"))

	target := NewNotificationUsecase(queueService, userService)

	err := target.SendNotification(notification, ctx)

	assert.Error(t, err)
	assert.EqualError(t, err, "error rbq")

}
