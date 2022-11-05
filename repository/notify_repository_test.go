package repository

import (
	"context"
	"encoding/json"
	"meli/notifications/dto"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

var (
	userUUID       = "fc2d1670-be5b-4235-bc55-9452bce74a0a"
	notificationID = "c487a2d6-2280-49c0-92cf-738d8cc71366"
	ctx            = context.TODO()
	regexUUID      = `^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$`
)

func TestNotifySaveSuccess(t *testing.T) {
	db, mock := redismock.NewClientMock()

	repo := NewNotifyRepository(db)

	notification := &dto.NotifyDTO{
		ID:       notificationID,
		Message:  "teste",
		Title:    "lucas",
		Image:    "cat.jpg",
		Type:     "normal",
		ToUserID: userUUID,
	}

	json, _ := json.Marshal(notification)
	mock.Regexp().ExpectSet(notificationID, json, 0).SetVal("")

	err := repo.Save(notification, ctx)

	assert.NoError(t, err)
}
