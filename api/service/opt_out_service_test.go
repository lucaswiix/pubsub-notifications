package service

import (
	"errors"
	"meli/notifications/repository/mock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	userID = "c487a2d6-2280-49c0-92cf-738d8cc71366"
)

func TestSetSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockOptOutRepository(ctrl)

	repo.EXPECT().Set(userID, gomock.Any()).Return(nil)

	target := NewOptOutService(repo)

	err := target.Set(userID, ctx)

	assert.NoError(t, err)

}

func TestDelSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockOptOutRepository(ctrl)

	repo.EXPECT().Del(userID, gomock.Any()).Return(nil)

	target := NewOptOutService(repo)

	err := target.Del(userID, ctx)

	assert.NoError(t, err)
}

func TestIsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockOptOutRepository(ctrl)

	repo.EXPECT().Is(userID, gomock.Any()).Return(true, nil)

	target := NewOptOutService(repo)

	isOpt, err := target.Is(userID, ctx)

	assert.NoError(t, err)
	assert.Equal(t, isOpt, true)
}

func TestIsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockOptOutRepository(ctrl)

	repo.EXPECT().Is(userID, gomock.Any()).Return(false, errors.New("error"))

	target := NewOptOutService(repo)

	isOpt, err := target.Is(userID, ctx)

	assert.EqualError(t, err, "error")
	assert.Equal(t, isOpt, false)
}
