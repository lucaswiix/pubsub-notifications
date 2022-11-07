package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
	InitLogger()
	assert.NotNil(t, Log)
}
