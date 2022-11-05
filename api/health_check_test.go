package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_HealthCheck_APIUp_Return200(t *testing.T) {
	t.Run("sucess with return 200", func(t *testing.T) {
		//prepare
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/healthcheck", nil)

		router := gin.Default()
		RegisterHealthCheckHandlers(router)

		//execute
		router.ServeHTTP(w, req)

		//check
		body, _ := io.ReadAll(w.Body)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `{"message":"ok"}`, string(body))
	})
}
