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

type userHandler struct {
	userService service.UserService
}

func RegisterUserHandlers(handler *gin.Engine, userService service.UserService) {
	ah := &userHandler{
		userService,
	}

	handler.POST("/api/user/opt-out", ah.Put)
	handler.DELETE("/api/user/opt-out/:id", ah.Del)
}

func (h *userHandler) Put(c *gin.Context) {
	span, ctx := apm.StartSpan(c.Request.Context(), "PutUser", "request")
	defer span.End()
	var optOut dto.OptOut

	if err := c.ShouldBind(&optOut); err != nil {
		utils.Log.Debug(fmt.Sprintf("error while parsing request parameters api/optout/set rid %s ", c.GetString("request_id")), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"message": utils.ValidateErrors(err)})
		return
	}

	err := h.userService.PutOptOut(optOut.UserID, ctx)
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

func (h *userHandler) Del(c *gin.Context) {
	span, ctx := apm.StartSpan(c.Request.Context(), "DeleteUser", "request")
	defer span.End()
	userID := c.Param("id")

	err := h.userService.Del(userID, ctx)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("error on try delete opt-out user in database %s", c.GetString("request_id")), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})

}
