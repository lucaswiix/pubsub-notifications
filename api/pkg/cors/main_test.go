package cors_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lucaswiix/meli/notifications/pkg/cors"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupRouterCors() *gin.Engine {

	c := cors.NewCors()

	router := gin.New()
	router.Use(c.CORS())

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "get")
	})
	router.POST("/", func(c *gin.Context) {
		c.String(http.StatusOK, "post")
	})
	router.PATCH("/", func(c *gin.Context) {
		c.String(http.StatusOK, "patch")
	})
	router.OPTIONS("/", func(c *gin.Context) {
		c.String(http.StatusOK, "options")
	})
	return router
}

func TestCors(t *testing.T) {
	router := setupRouterCors()

	w := performRequest(router, "GET", "https://gist.github.com")
	assert.Equal(t, 200, w.Code)

	w = performRequest(router, "GET", "https://github.com")
	assert.Equal(t, 200, w.Code)

	w = performRequest(router, "OPTIONS", "https://github.com")
	assert.Equal(t, 204, w.Code)
}

func performRequest(r *gin.Engine, method, origin string) *httptest.ResponseRecorder {
	return performRequestWithHeaders(r, method, origin, http.Header{})
}

func performRequestWithHeaders(r *gin.Engine, method, origin string, header http.Header) *httptest.ResponseRecorder {
	req, _ := http.NewRequestWithContext(context.Background(), method, "/", nil)
	req.Host = header.Get("Host")
	header.Del("Host")
	if len(origin) > 0 {
		header.Set("Origin", origin)
	}
	req.Header = header
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
