package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lucaswiix/notifications-tracking-app/domain"
	"github.com/lucaswiix/notifications-tracking-app/domain/mock"
	"github.com/stretchr/testify/assert"
)

var (
	userID = "34065a27-cb0c-42f9-9bea-258c5806aaa5"
	nType  = "web"
	ctx    = context.TODO()
)

func TestTrackByUserIDSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)

	notificationClient := mock.NewMockNotificationClient(ctrl)

	notification := &domain.Notification{
		Message: "msg",
		Title:   "titulo",
		Image:   "cat.png",
		Type:    nType,
		UserID:  userID,
	}
	json, _ := json.Marshal(notification)

	notificationClient.EXPECT().ConsumeByUserID(ctx, userID, nType).Return(json, nil)

	target := NewNotificationUsecase(notificationClient)

	notification, err := target.TrackByUserID(ctx, userID, nType)

	assert.NoError(t, err)
	assert.Equal(t, notification, notification)
}

func TestTrackByUserIDError(t *testing.T) {
	ctrl := gomock.NewController(t)

	notificationClient := mock.NewMockNotificationClient(ctrl)

	notificationClient.EXPECT().ConsumeByUserID(ctx, userID, nType).Return(nil, errors.New("error"))

	target := NewNotificationUsecase(notificationClient)

	_, err := target.TrackByUserID(ctx, userID, nType)

	assert.EqualError(t, err, "error")
}
