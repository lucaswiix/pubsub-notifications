package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/lucaswiix/meli/notifications/dto"
	"github.com/lucaswiix/meli/notifications/usecase/mock"
	"github.com/lucaswiix/meli/notifications/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	userID        = "fc2d1670-be5b-4235-bc55-9452bce74a0a"
	schedulerDate = time.Now().Local().Add(time.Hour * time.Duration(1)).Round(0 * time.Second).Format("2006-01-02 15:04:05")
)

func Test_Send_Notification_With_Scheduler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	notificationUsecase := mock.NewMockNotificationUseCase(ctrl)

	notificationDTO := &dto.NotifyDTO{
		Message:       "mensagem",
		Title:         "Titulo",
		Image:         "cat.png",
		Type:          "web",
		ToUserID:      userID,
		SchedulerDate: schedulerDate,
	}

	notificationUsecase.EXPECT().SendNotification(notificationDTO, gomock.Any()).Return(nil)

	router := gin.Default()

	RegisterNotifyHandlers(router, notificationUsecase)

	requestBody, _ := json.Marshal(notificationDTO)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/notify", strings.NewReader(string(requestBody)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-user-id", userID)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_Send_Notification_Without_Scheduler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	notificationUsecase := mock.NewMockNotificationUseCase(ctrl)

	notificationDTO := &dto.NotifyDTO{
		Message:  "mensagem",
		Title:    "Titulo",
		Image:    "cat.png",
		Type:     "web",
		ToUserID: userID,
	}

	notificationUsecase.EXPECT().SendNotification(notificationDTO, gomock.Any()).Return(nil)

	router := gin.Default()

	RegisterNotifyHandlers(router, notificationUsecase)

	requestBody, _ := json.Marshal(notificationDTO)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/notify", strings.NewReader(string(requestBody)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-user-id", userID)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_Send_Notification_Without_UserID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	notificationDTO := &dto.NotifyDTO{
		Message: "mensagem",
		Title:   "Titulo",
		Image:   "cat.png",
		Type:    "web",
	}

	router := gin.Default()

	RegisterNotifyHandlers(router, nil)

	requestBody, _ := json.Marshal(notificationDTO)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/notify", strings.NewReader(string(requestBody)))
	req.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	body, _ := io.ReadAll(w.Body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"message":"header x-user-id is required"}`, string(body))

}

func Test_Send_Notification_Invalid_Fields_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	notificationDTO := &dto.NotifyDTO{
		Message:  "mensagem",
		Title:    "Titulo",
		Type:     "web",
		ToUserID: userID,
	}

	router := gin.Default()

	RegisterNotifyHandlers(router, nil)

	requestBody, _ := json.Marshal(notificationDTO)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/notify", strings.NewReader(string(requestBody)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-user-id", userID)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_Send_Notification_Is_Opt_Out_User_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	notificationUsecase := mock.NewMockNotificationUseCase(ctrl)

	notificationDTO := &dto.NotifyDTO{
		Message:  "mensagem",
		Title:    "Titulo",
		Image:    "cat.png",
		Type:     "web",
		ToUserID: userID,
	}

	notificationUsecase.EXPECT().SendNotification(notificationDTO, gomock.Any()).Return(utils.ErrOptOutUser)

	router := gin.Default()

	RegisterNotifyHandlers(router, notificationUsecase)

	requestBody, _ := json.Marshal(notificationDTO)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/notify", strings.NewReader(string(requestBody)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-user-id", userID)

	router.ServeHTTP(w, req)
	body, _ := io.ReadAll(w.Body)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Equal(t, `{"message":"opt-out user"}`, string(body))

}

func Test_Send_Notification_Generic_Send_Notification_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	notificationUsecase := mock.NewMockNotificationUseCase(ctrl)

	notificationDTO := &dto.NotifyDTO{
		Message:  "mensagem",
		Title:    "Titulo",
		Image:    "cat.png",
		Type:     "web",
		ToUserID: userID,
	}

	notificationUsecase.EXPECT().SendNotification(notificationDTO, gomock.Any()).Return(errors.New("error"))

	router := gin.Default()

	RegisterNotifyHandlers(router, notificationUsecase)

	requestBody, _ := json.Marshal(notificationDTO)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/notify", strings.NewReader(string(requestBody)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-user-id", userID)

	router.ServeHTTP(w, req)
	body, _ := io.ReadAll(w.Body)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, `{"message":"error"}`, string(body))

}
