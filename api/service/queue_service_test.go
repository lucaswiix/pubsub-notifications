package service

import (
	"errors"
	"testing"

	"github.com/lucaswiix/meli/notifications/dto"
	"github.com/lucaswiix/meli/notifications/repository/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSendSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockQueueRepository(ctrl)

	repo.EXPECT().Send(gomock.Any(), gomock.Any()).Return(nil)

	target := NewQueueService(repo)

	notification := &dto.NotifyDTO{
		Title:   "Celular",
		Message: "celular promocao",
		Image:   "cat.jpg",
		Type:    "web",
	}

	err := target.Send(notification, ctx)

	assert.NoError(t, err)
	assert.Equal(t, notification.ID, notification.ID)

}

func TestSendError(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockQueueRepository(ctrl)

	notification := &dto.NotifyDTO{
		Title:   "Celular",
		Message: "celular promocao",
		Image:   "cat.jpg",
		Type:    "web",
	}

	repo.EXPECT().Send(gomock.Any(), gomock.Any()).Return(errors.New("error send"))

	target := NewQueueService(repo)

	err := target.Send(notification, ctx)

	assert.EqualError(t, err, "error send")
}
