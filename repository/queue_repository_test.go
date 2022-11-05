package repository

import (
	"meli/notifications/dto"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	userID = "c487a2d6-2280-49c0-92cf-738d8cc71366"
)

func TestSendParseSchedulerDateTimeError(t *testing.T) {

	notification := &dto.NotifyDTO{
		Message:       "msg",
		Type:          "web",
		Image:         "cat.png",
		ToUserID:      userID,
		Title:         "titulo",
		SchedulerDate: "20100821'T'04:20:15",
	}

	queueRepo := NewQueueRepository(nil)
	err := queueRepo.Send(notification, ctx)

	assert.Error(t, err)
}
