package domain

import "context"

type Notification struct {
	Message string `json:"message"`
	Title   string `json:"title"`
	Image   string `json:"image"`
	Type    string `json:"type"`
	UserID  string `json:"to_user_id"`
}

//go:generate mockgen -destination=mock/notification_usecase.go -package=mock . NotificationUsecase

type NotificationUsecase interface {
	TrackByUserID(ctx context.Context, id, nType string) (*Notification, error)
}

//go:generate mockgen -destination=mock/notification_client.go -package=mock . NotificationClient
type NotificationClient interface {
	ConsumeByUserID(ctx context.Context, userID, nType string) ([]byte, error)
}
