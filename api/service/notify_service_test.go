package service

import (
	"context"
	"errors"
	"github.com/lucaswiix/meli/notifications/dto"
	"github.com/lucaswiix/meli/notifications/repository/mock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	ctx            = context.TODO()
	notificationID = "c487a2d6-2280-49c0-92cf-738d8cc71366"
)

func TestNotifySaveSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockNotifyRepository(ctrl)

	repo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

	target := NewNotifyService(repo)

	notification := &dto.NotifyDTO{
		Title:   "Celular",
		Message: "celular promocao",
		Image:   "cat.jpg",
		Type:    "web",
	}

	err := target.Save(notification, ctx)

	assert.NoError(t, err)
	assert.Equal(t, notification.ID, notification.ID)

}

func TestNotifyUpdateSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockNotifyRepository(ctrl)

	notification := &dto.NotifyDTO{
		ID:      notificationID,
		Title:   "Celular",
		Message: "celular promocao",
		Image:   "cat.jpg",
		Type:    "web",
	}

	repo.EXPECT().Save(notification, gomock.Any()).Return(nil)

	target := NewNotifyService(repo)

	err := target.Save(notification, ctx)

	assert.NoError(t, err)
	assert.Equal(t, notification.ID, notification.ID)

}
func TestNotifySaveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockNotifyRepository(ctrl)

	notification := &dto.NotifyDTO{
		Title:   "Celular",
		Message: "celular promocao",
		Image:   "cat.jpg",
		Type:    "web",
	}

	repo.EXPECT().Save(notification, gomock.Any()).Return(errors.New("error"))

	target := NewNotifyService(repo)

	err := target.Save(notification, ctx)

	assert.EqualError(t, err, "error")
}
