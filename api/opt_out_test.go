package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"meli/notifications/dto"
	"meli/notifications/service/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Set_Opt_Out_User_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	optOutService := mock.NewMockOptOutService(ctrl)

	user := &dto.OptOut{
		UserID: userID,
	}

	optOutService.EXPECT().Set(user.UserID, gomock.Any()).Return(nil)

	router := gin.Default()

	RegisterOptOutHandlers(router, optOutService)

	requestBody, _ := json.Marshal(user)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/opt-out", strings.NewReader(string(requestBody)))
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

	RegisterOptOutHandlers(router, nil)

	requestBody, _ := json.Marshal(user)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/opt-out", strings.NewReader(string(requestBody)))
	req.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	body, _ := io.ReadAll(w.Body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"message":"UserID invalid user uuid format\n"}`, string(body))

}
func Test_Set_Opt_Out_User_Service_Error_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	optOutService := mock.NewMockOptOutService(ctrl)

	user := &dto.OptOut{
		UserID: userID,
	}

	optOutService.EXPECT().Set(user.UserID, gomock.Any()).Return(errors.New("error"))

	router := gin.Default()

	RegisterOptOutHandlers(router, optOutService)

	requestBody, _ := json.Marshal(user)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/opt-out", strings.NewReader(string(requestBody)))
	req.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	body, _ := io.ReadAll(w.Body)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, `{"message":"error"}`, string(body))

}

func Test_Del_Opt_Out_User_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	optOutService := mock.NewMockOptOutService(ctrl)

	optOutService.EXPECT().Del(userID, gomock.Any()).Return(nil)

	router := gin.Default()

	RegisterOptOutHandlers(router, optOutService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/opt-out/%s", userID), nil)
	req.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func Test_Del_Opt_Out_Service_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	optOutService := mock.NewMockOptOutService(ctrl)

	optOutService.EXPECT().Del(userID, gomock.Any()).Return(errors.New("error"))

	router := gin.Default()

	RegisterOptOutHandlers(router, optOutService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/opt-out/%s", userID), nil)
	req.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	body, _ := io.ReadAll(w.Body)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, `{"message":"error"}`, string(body))

}
