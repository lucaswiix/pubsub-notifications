package usecase

import (
	"os"
	"testing"

	"github.com/lucaswiix/notifications-tracking-app/utils"
)

func TestMain(m *testing.M) {
	utils.InitLogger()
	os.Exit(m.Run())
}
