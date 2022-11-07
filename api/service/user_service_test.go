package service

import (
	"errors"
	"testing"
	"time"

	"github.com/lucaswiix/meli/notifications/dto"
	"github.com/lucaswiix/meli/notifications/repository/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	userID = "c487a2d6-2280-49c0-92cf-738d8cc71366"
)

func TestPutSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockUserRepository(ctrl)

	user := dto.User{
		ID:        userID,
		IsOptOut:  true,
		CreatedAt: time.Now().Format(time.RFC822),
	}

	repo.EXPECT().Put(user, gomock.Any()).Return(nil)

	target := NewUserService(repo)

	err := target.Put(user, ctx)

	assert.NoError(t, err)

}

func TestDelSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockUserRepository(ctrl)

	repo.EXPECT().Del(userID, gomock.Any()).Return(nil)

	target := NewUserService(repo)

	err := target.Del(userID, ctx)

	assert.NoError(t, err)
}

func TestGetSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockUserRepository(ctrl)

	user := dto.User{
		ID:        userID,
		IsOptOut:  true,
		CreatedAt: time.Now().Format(time.RFC822),
	}

	repo.EXPECT().Get(userID, gomock.Any()).Return(user, nil)

	target := NewUserService(repo)

	user, err := target.Get(userID, ctx)

	assert.NoError(t, err)
	assert.Equal(t, user.IsOptOut, true)
}

func TestGetError(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockUserRepository(ctrl)

	emptyUser := dto.User{}
	repo.EXPECT().Get(userID, gomock.Any()).Return(emptyUser, errors.New("error"))

	target := NewUserService(repo)

	user, err := target.Get(userID, ctx)

	assert.EqualError(t, err, "error")
	assert.Equal(t, user, emptyUser)
}
