package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lucaswiix/meli/notifications/dto"
	"github.com/lucaswiix/meli/notifications/service/mock"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPutUserOptOutSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userService := mock.NewMockUserService(ctrl)

	user := &dto.OptOut{
		UserID: userID,
	}

	userService.EXPECT().PutOptOut(userID, gomock.Any()).Return(nil)

	router := gin.Default()

	RegisterUserHandlers(router, userService)

	requestBody, _ := json.Marshal(user)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/user/opt-out", strings.NewReader(string(requestBody)))
	req.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_Set_Opt_Out_User_Validate_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := &dto.OptOut{
		UserID: "invalid-uuid",
	}

	router := gin.Default()

	RegisterUserHandlers(router, nil)

	requestBody, _ := json.Marshal(user)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/user/opt-out", strings.NewReader(string(requestBody)))
	req.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	body, _ := io.ReadAll(w.Body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"message":"UserID invalid user uuid format\n"}`, string(body))

}
func Test_Set_Opt_Out_User_Service_Error_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userService := mock.NewMockUserService(ctrl)

	user := &dto.OptOut{
		UserID: userID,
	}

	userService.EXPECT().PutOptOut(user.UserID, gomock.Any()).Return(errors.New("error"))

	router := gin.Default()

	RegisterUserHandlers(router, userService)

	requestBody, _ := json.Marshal(user)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/user/opt-out", strings.NewReader(string(requestBody)))
	req.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	body, _ := io.ReadAll(w.Body)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, `{"message":"error"}`, string(body))

}

func Test_Del_Opt_Out_User_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userService := mock.NewMockUserService(ctrl)

	userService.EXPECT().Del(userID, gomock.Any()).Return(nil)

	router := gin.Default()

	RegisterUserHandlers(router, userService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/user/opt-out/%s", userID), nil)
	req.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func Test_Del_Opt_Out_Service_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userService := mock.NewMockUserService(ctrl)

	userService.EXPECT().Del(userID, gomock.Any()).Return(errors.New("error"))

	router := gin.Default()

	RegisterUserHandlers(router, userService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/user/opt-out/%s", userID), nil)
	req.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	body, _ := io.ReadAll(w.Body)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, `{"message":"error"}`, string(body))

}
