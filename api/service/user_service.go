package service

import (
	"context"
	"time"

	"github.com/lucaswiix/meli/notifications/dto"
	"github.com/lucaswiix/meli/notifications/repository"
	"go.elastic.co/apm"
)

//go:generate mockgen -destination=mock/opt_out.go -package=mock . UserService
type UserService interface {
	PutOptOut(userID string, ctx context.Context) error
	Put(user dto.User, ctx context.Context) error
	Del(userID string, ctx context.Context) error
	Get(userID string, ctx context.Context) (dto.User, error)
}

type UserServiceImpl struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserService {
	return &UserServiceImpl{
		repository,
	}
}

func (s *UserServiceImpl) PutOptOut(userID string, ctx context.Context) error {
	span, ctx := apm.StartSpan(ctx, "SetUser", "service")
	defer span.End()
	user := dto.User{
		ID:        userID,
		IsOptOut:  true,
		CreatedAt: time.Now().Format(time.RFC822),
	}
	return s.repository.Put(user, ctx)
}

func (s *UserServiceImpl) Put(user dto.User, ctx context.Context) error {
	span, ctx := apm.StartSpan(ctx, "SetUser", "service")
	defer span.End()
	return s.repository.Put(user, ctx)
}

func (s *UserServiceImpl) Del(userID string, ctx context.Context) error {
	span, ctx := apm.StartSpan(ctx, "DelUser", "service")
	defer span.End()
	return s.repository.Del(userID, ctx)
}

func (s *UserServiceImpl) Get(userID string, ctx context.Context) (dto.User, error) {
	span, ctx := apm.StartSpan(ctx, "IsUser", "service")
	defer span.End()
	user, err := s.repository.Get(userID, ctx)
	if err != nil {
		return user, err
	}
	return user, nil
}
