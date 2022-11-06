package usecase

import (
	"github.com/lucaswiix/meli/notifications/utils"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	utils.InitLogger()
	os.Exit(m.Run())
}
