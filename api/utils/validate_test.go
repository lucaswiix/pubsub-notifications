package utils

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	_ "time/tzdata"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

var (
	allMsgErr     = "Name is required\nEmail must me valid email\nUUID invalid user uuid format\nImage only accept png image format\nType must be one of webDate must be set to future\nMin must be 5 length at least\nerror in filed max"
	schedulerDate = time.Now().Local().Add(time.Hour * time.Duration(1)).Round(0 * time.Second).Format("2006-01-02 15:04:05")
)

type ValidateStructTest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"email"`
	UUID  string `json:"uuid" validate:"uuid"`
	Image string `json:"image" validate:"endswith=png"`
	Type  string `json:"type" validate:"Enum=web"`
	Date  string `json:"date" validate:"IsAfterNow"`
	Min   string `json:"age" validate:"min=5"`
	Max   string `json:"max" validate:"max=2"`
}

func TestMain(m *testing.M) {
	time.Local, _ = time.LoadLocation("America/Sao_Paulo")
	os.Exit(m.Run())
}
func TestValidateErrorsMessage(t *testing.T) {
	t.Run("json syntax error", func(t *testing.T) {

		var jsonErr *json.SyntaxError

		strData := ValidateErrors(jsonErr)

		assert.Equal(t, "invalid json input", strData)
	})
	t.Run("validate errors message", func(t *testing.T) {
		validate := validator.New()
		validate.RegisterValidation("Enum", Enum)
		validate.RegisterValidation("IsAfterNow", IsAfterNow)

		required := ValidateStructTest{
			Date: "2010-10-02 10:10:00",
			Max:  "123",
		}
		err := validate.Struct(required)

		strData := ValidateErrors(err)

		assert.Equal(t, allMsgErr, strData)
	})
}

func TestIsAfterNowTrue(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("IsAfterNow", IsAfterNow)
	err := validate.Var(schedulerDate, "IsAfterNow")
	assert.NoError(t, err)
}

func TestIsAfterNowEmptyFalse(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("IsAfterNow", IsAfterNow)
	err := validate.Var("", "IsAfterNow")
	assert.NoError(t, err)
}

func TestIsAfterNowParseErrorFalse(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("IsAfterNow", IsAfterNow)
	err := validate.Var("12312", "IsAfterNow")
	assert.Error(t, err)
}
func TestIsAfterNowFalse(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("IsAfterNow", IsAfterNow)
	err := validate.Var("2010-10-02 10:10:00", "IsAfterNow")
	assert.Error(t, err)
}

func TestValidateRequests(t *testing.T) {

}
func TestEnumTrue(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("Enum", Enum)
	err := validate.Var("web", "Enum=web")
	assert.NoError(t, err)
}

func TestEnumFalse(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("Enum", Enum)
	err := validate.Var("sms", "Enum=web")
	assert.Error(t, err)
}

func TestInitValidator(t *testing.T) {
	InitValidation()
	validate := binding.Validator.Engine().(*validator.Validate)
	err := validate.Var("web", "Enum=web")
	assert.NoError(t, err)
}
