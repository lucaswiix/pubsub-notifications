package api

import (
	"fmt"
	"github.com/lucaswiix/meli/notifications/dto"
	"github.com/lucaswiix/meli/notifications/service"
	"github.com/lucaswiix/meli/notifications/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type optOutHandler struct {
	optOutService service.OptOutService
}

func RegisterOptOutHandlers(handler *gin.Engine, optOutService service.OptOutService) {
	ah := &optOutHandler{
		optOutService,
	}

	handler.POST("/opt-out", ah.Set)
	handler.DELETE("/opt-out/:id", ah.Del)
}

// @Summary Set user opt out
// @Tags SetOptOut
// @Schemes
// @Accept json
// @Param body body dto.NotifyDTO true "Send Notification"
// @Produce json
// @Success 200 {object} object{uuid=string,type_name=string} "Resposta de successo quando é atualizado um tenant"
// @Success 201 {object} object{uuid=string,type_name=string} "Resposta de successo quando é criado um novo tenant"
// @Failure 400 {object} object{message=string} "Resposta de erro quando identifica que o atributo type_name está vazio ou inválido ou quando o account não está autorizado a delegar um tenant"
// @Failure 500 {object} object{message=string} "Resposta de erro durante o processo de criar/alterar um tenant"
// @Router /set-opt-out [post]
func (h *optOutHandler) Set(c *gin.Context) {
	ctx := c.Request.Context()
	var optOut dto.OptOut

	if err := c.ShouldBind(&optOut); err != nil {
		utils.Log.Debug(fmt.Sprintf("error while parsing request parameters api/optout/set rid %s ", c.GetString("request_id")), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"message": utils.ValidateErrors(err)})
		return
	}

	err := h.optOutService.Set(optOut.UserID, ctx)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("error on try save database %s", c.GetString("request_id")), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user opt-out set",
	})

}

// @Summary Set user opt out
// @Tags SetOptOut
// @Schemes
// @Accept json
// @Param body body dto.NotifyDTO true "Send Notification"
// @Produce json
// @Success 200 {object} object{uuid=string,type_name=string} "Resposta de successo quando é atualizado um tenant"
// @Success 201 {object} object{uuid=string,type_name=string} "Resposta de successo quando é criado um novo tenant"
// @Failure 400 {object} object{message=string} "Resposta de erro quando identifica que o atributo type_name está vazio ou inválido ou quando o account não está autorizado a delegar um tenant"
// @Failure 500 {object} object{message=string} "Resposta de erro durante o processo de criar/alterar um tenant"
// @Router /opt-out [delete]
func (h *optOutHandler) Del(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("id")

	err := h.optOutService.Del(userID, ctx)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("error on try delete opt-out user in database %s", c.GetString("request_id")), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})

}
