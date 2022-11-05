package api

import (
	"fmt"
	"log"
	"meli/notifications/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gin-gonic/gin"
)

var (
	routesFileYaml string = `
routes:
  - path: "/"
    methods:
      - name: "GET"
        skip_auth: true
`
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func init() {
	gin.SetMode(gin.TestMode)
}

func createRouteFile() {

	fmt.Println("basepath", basepath)

	f, err := os.OpenFile(fmt.Sprintf("%s/../routesTest.yaml", basepath), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(routesFileYaml)); err != nil {
		log.Fatal(err)
	}
}

func removeRouteFile() {
	err := os.RemoveAll(fmt.Sprintf("%s/../routesTest.yaml", basepath))
	if err != nil {
		log.Fatal(err)
	}
}

func TestInitWebServer(t *testing.T) {
	t.Setenv("ROUTES_CONFIG_FULL_PATH", fmt.Sprintf("%s/../routesTest.yaml", basepath))
	t.Setenv("CASBIN_POLICY_FULL_PATH", fmt.Sprintf("%s/../gen_auto_policy.csv", basepath))
	t.Setenv("CASBIN_MODEL_FULL_PATH", fmt.Sprintf("%s/../keymatch_model.conf", basepath))
	createRouteFile()
	utils.InitLogger()
	t.Run("success", func(t *testing.T) {
		server := InitWebServer()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
	})

	t.Cleanup(removeRouteFile)
}
