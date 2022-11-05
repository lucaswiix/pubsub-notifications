package usecase

import (
	"context"
	"encoding/json"

	"github.com/lucaswiix/notifications-tracking-app/domain"
)

type notificationUsecase struct {
	pc domain.NotificationClient
}

func NewNotificationUsecase(pClient domain.NotificationClient) domain.NotificationUsecase {
	return &notificationUsecase{pc: pClient}
}

func (p *notificationUsecase) TrackByUserID(ctx context.Context, id, nType string) (*domain.Notification, error) {
	bytes, err := p.pc.ConsumeByUserID(ctx, id, nType)
	if err != nil {
		return nil, err
	}

	var res domain.Notification
	err = json.Unmarshal(bytes, &res)
	return &res, err
}
