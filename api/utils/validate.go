package utils

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	JsonSyntaxErr = "invalid json input"
)

func InitValidation() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("Enum", Enum)
		_ = v.RegisterValidation("IsAfterNow", IsAfterNow)
	}
}

func Enum(
	fl validator.FieldLevel,
) bool {
	enumString := fl.Param()                    // get string male_female
	value := fl.Field().String()                // the actual field
	enumSlice := strings.Split(enumString, "_") // convert to slice
	for _, v := range enumSlice {
		if value == v {
			return true
		}
	}
	return false
}

func IsAfterNow(
	fl validator.FieldLevel,
) bool {
	value := fl.Field().String() // the actual field
	if value == "" {
		return true
	}
	now := time.Now()
	datanow := now.Format(time.RFC822)
	Log.Info(datanow)

	schedulerDate, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
	dataFuture := schedulerDate.Format(time.RFC822)
	Log.Info(dataFuture)

	if err != nil {
		return false
	}
	isAfter := schedulerDate.After(now)
	return isAfter
}

// custom validation error messages
func ValidateErrors(requestError error) string {
	var jsonErr *json.SyntaxError
	if requestError == jsonErr {
		return JsonSyntaxErr
	}
	return validate(requestError.(validator.ValidationErrors))
}

func validate(errors validator.ValidationErrors) string {
	resultErrors := ""
	for _, err := range errors {
		switch err.Tag() {
		case "required":
			resultErrors += err.Field() + " is required\n"
		case "email":
			resultErrors += err.Field() + " must me valid email\n"
		case "min":
			resultErrors += err.Field() + " must be " + err.Param() + " length at least\n"
		case "Enum":
			replacer := *strings.NewReplacer("_", ",")
			resultErrors += err.Field() + " must be one of " + replacer.Replace(err.Param())
		case "IsAfterNow":
			resultErrors += err.Field() + " must be set to future\n"
		case "uuid":
			resultErrors += err.Field() + " invalid user uuid format\n"
		case "endswith":
			resultErrors += err.Field() + " only accept png image format\n"

		default:
			resultErrors += "error in filed " + err.Tag()
		}
	}
	return resultErrors
}
