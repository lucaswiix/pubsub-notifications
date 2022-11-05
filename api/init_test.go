package api

import (
	"meli/notifications/utils"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	utils.InitLogger()
	utils.InitValidation()
	os.Exit(m.Run())
}
