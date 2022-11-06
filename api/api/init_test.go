package api

import (
	"os"
	"testing"

	"github.com/lucaswiix/meli/notifications/utils"
)

func TestMain(m *testing.M) {
	utils.InitLogger()
	utils.InitValidation()
	os.Exit(m.Run())
}
