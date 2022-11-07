package api

import (
	"fmt"
	"net/http"

	"github.com/lucaswiix/meli/notifications/dto"
	"github.com/lucaswiix/meli/notifications/service"
	"github.com/lucaswiix/meli/notifications/utils"
	"go.elastic.co/apm"

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

	handler.POST("/api/user/opt-out", ah.Set)
	handler.DELETE("/api/user/opt-out/:id", ah.Del)
}

func (h *optOutHandler) Set(c *gin.Context) {
	span, ctx := apm.StartSpan(c.Request.Context(), "SetOptOut", "request")
	defer span.End()
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

func (h *optOutHandler) Del(c *gin.Context) {
	span, ctx := apm.StartSpan(c.Request.Context(), "DeleteOptOut", "request")
	defer span.End()
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
