package api

import (
	"fmt"
	"net/http"

	"github.com/lucaswiix/meli/notifications/dto"
	"github.com/lucaswiix/meli/notifications/usecase"
	"github.com/lucaswiix/meli/notifications/utils"

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

	handler.POST("/notify", ah.Sent)
}

// @Summary Send notification
// @Tags Tenants
// @Schemes
// @Accept json
// @Param body body dto.NotifyDTO true "Send Notification"
// @Produce json
// @Success 200 {object} object{uuid=string,type_name=string} "Resposta de successo quando é atualizado um tenant"
// @Success 201 {object} object{uuid=string,type_name=string} "Resposta de successo quando é criado um novo tenant"
// @Failure 400 {object} object{message=string} "Resposta de erro quando identifica que o atributo type_name está vazio ou inválido ou quando o account não está autorizado a delegar um tenant"
// @Failure 500 {object} object{message=string} "Resposta de erro durante o processo de criar/alterar um tenant"
// @Router /tenants [post]
func (h *notifyHandler) Sent(c *gin.Context) {
	ctx := c.Request.Context()
	var notifyDTO dto.NotifyDTO
	userID := c.GetHeader("x-user-id")
	notifyDTO.ToUserID = userID

	if userID == "" {
		utils.Log.Debug(fmt.Sprintf("error while parsing request parameters api/notify/save rid %s ", c.GetString("request_id")))
		c.JSON(http.StatusBadRequest, gin.H{"message": "header x-user-id is required"})
		return
	}

	if err := c.ShouldBind(&notifyDTO); err != nil {
		utils.Log.Debug(fmt.Sprintf("error while parsing request parameters api/notify/save rid %s ", c.GetString("request_id")), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"message": utils.ValidateErrors(err)})
		return
	}

	err := h.notificationUsecase.SendNotification(&notifyDTO, ctx)
	if err != nil {
		if err == utils.ErrOptOutUser {
			utils.Log.Debug(fmt.Sprintf("error on send notification to opt-out user %s", c.GetString("request_id")), zap.Error(err))
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
			})
			return
		}
		utils.Log.Error(fmt.Sprintf("error on try save database %s", c.GetString("request_id")), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "notification sent",
	})

}
