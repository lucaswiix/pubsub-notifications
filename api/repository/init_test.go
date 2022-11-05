package repository

import (
	"meli/notifications/utils"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	utils.InitLogger()
	os.Exit(m.Run())
}
