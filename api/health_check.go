package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
)

func RegisterHealthCheckHandlers(handler *gin.Engine) {
	handler.GET("/healthcheck", HealthCheck)
}

// @Summary Verificar se o serviço do platform-account-api se encontra ativo
// @Tags Healthcheck
// @Produce application/json
// @Success 200 {object} object{message=string} "Resposta de Successo"
// @Router /healthcheck [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
	txn := nrgin.Transaction(c)
	txn.Ignore()
}
