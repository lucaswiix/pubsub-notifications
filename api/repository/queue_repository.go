package repository

import (
	"context"
	"fmt"
	"github.com/lucaswiix/meli/notifications/dto"
	"github.com/lucaswiix/meli/notifications/utils"
	"time"

	"github.com/goccy/go-json"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=mock/queue.go -package=mock . QueueRepository
type QueueRepository interface {
	Send(notification *dto.NotifyDTO, ctx context.Context) error
}

type implQueueRepository struct {
	ch *amqp.Channel
}

func NewQueueRepository(ch *amqp.Channel) QueueRepository {
	return &implQueueRepository{ch}
}

func (r *implQueueRepository) Send(notification *dto.NotifyDTO, ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cleanNotify := dto.CleanNotify{
		Title:   notification.Title,
		Image:   notification.Image,
		Message: notification.Message,
		UserID:  notification.ToUserID,
		Type:    notification.Type,
	}
	json, err := json.Marshal(cleanNotify)
	if err != nil {
		utils.Log.Error("failed to marshal notification interface", zap.Error(err))
		return err
	}

	headers := make(amqp.Table)
	if notification.SchedulerDate != "" {
		now := time.Now()
		scheduledDate, err := time.ParseInLocation("2006-01-02 15:04:05", notification.SchedulerDate, time.Local)
		if err != nil {
			utils.Log.Error("failed when try parse date", zap.Error(err))
			return err
		}
		subDate := scheduledDate.Sub(now).Milliseconds()
		headers["x-delay"] = subDate
		utils.Log.Info(fmt.Sprintf("scheduler message to %v miliseconds", subDate))
	}

	err = r.ch.PublishWithContext(ctx,
		"notifications",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(json),
			MessageId:   notification.ToUserID,
			Type:        notification.Type,
			Headers:     headers,
		})
	if err != nil {
		utils.Log.Error("failed to publish a message", zap.Error(err))
		return err
	}
	utils.Log.Debug(fmt.Sprintf(" [x] Sent %s\n", json))

	return nil
}
