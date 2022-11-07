package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/lucaswiix/meli/notifications/dto"
	"github.com/lucaswiix/meli/notifications/usecase"
	"github.com/lucaswiix/meli/notifications/utils"
	"go.elastic.co/apm"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type notifyHandler struct {
	notificationUsecase usecase.NotificationUseCase
}

func RegisterNotifyHandlers(handler *gin.Engine, notificationUsecase usecase.NotificationUseCase) {
	ah := &notifyHandler{
		notificationUsecase,
	}

	handler.POST("/api/notify", ah.Sent)
}

func (h *notifyHandler) Sent(c *gin.Context) {
	span, ctx := apm.StartSpan(c.Request.Context(), "SentNotification", "request")
	defer span.End()
	var notifyDTO dto.NotifyDTO
	userID := c.GetHeader("x-user-id")
	notifyDTO.ToUserID = userID
	status := http.StatusOK

	defer func() {
		notificationStatus.WithLabelValues(notifyDTO.ToUserID, strconv.Itoa(status)).Inc()
	}()

	if userID == "" {
		status = http.StatusBadRequest
		utils.Log.Debug(fmt.Sprintf("error while parsing request parameters api/notify/save rid %s ", c.GetString("request_id")))
		c.JSON(status, gin.H{"message": "header x-user-id is required"})
		return
	}

	if err := c.ShouldBind(&notifyDTO); err != nil {
		status = http.StatusBadRequest
		utils.Log.Debug(fmt.Sprintf("error while parsing request parameters api/notify/save rid %s ", c.GetString("request_id")), zap.Error(err))
		c.JSON(status, gin.H{"message": utils.ValidateErrors(err)})
		return
	}

	err := h.notificationUsecase.SendNotification(&notifyDTO, ctx)
	if err != nil {
		status = http.StatusInternalServerError

		if err == utils.ErrOptOutUser {
			status = http.StatusForbidden

			utils.Log.Debug(fmt.Sprintf("error on send notification to opt-out user %s", c.GetString("request_id")), zap.Error(err))
			c.JSON(status, gin.H{
				"message": err.Error(),
			})
			return
		}

		utils.Log.Error(fmt.Sprintf("error on try save database %s", c.GetString("request_id")), zap.Error(err))
		c.JSON(status, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(status, gin.H{
		"message": "notification sent",
	})

}
