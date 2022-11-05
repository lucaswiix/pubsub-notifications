package service

import (
	"context"
	"meli/notifications/repository"
)

//go:generate mockgen -destination=mock/opt_out.go -package=mock . OptOutService
type OptOutService interface {
	Set(userID string, ctx context.Context) error
	Del(userID string, ctx context.Context) error
	Is(userID string, ctx context.Context) (bool, error)
}

type OptOutServiceImpl struct {
	repository repository.OptOutRepository
}

func NewOptOutService(repository repository.OptOutRepository) OptOutService {
	return &OptOutServiceImpl{
		repository,
	}
}

func (s *OptOutServiceImpl) Set(userID string, ctx context.Context) error {
	return s.repository.Set(userID, ctx)
}

func (s *OptOutServiceImpl) Del(userID string, ctx context.Context) error {
	return s.repository.Del(userID, ctx)
}

func (s *OptOutServiceImpl) Is(userID string, ctx context.Context) (bool, error) {
	isOptOut, err := s.repository.Is(userID, ctx)
	if err != nil {
		return false, err
	}
	return isOptOut, nil
}
