package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
)

func IsDuplicateKeyErr(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return true
	}

	return false
}

func validationErrMsgHandler(errors validator.ValidationErrors) string {
	var message string

	switch errors[0].Tag() {
	case "required":
		message = fmt.Sprintf("%s is required", errors[0].Field())
	case "email":
		message = "invalid email address"
	case "gte":
		message = fmt.Sprintf("%s must be greater than or equal to %s", errors[0].Field(), errors[0].Param())
	case "lte":
		message = fmt.Sprintf("%s must be less than or equal to %s", errors[0].Field(), errors[0].Param())
	case "gt":
		message = fmt.Sprintf("%s must be greater than %s", errors[0].Field(), errors[0].Param())
	case "lt":
		message = fmt.Sprintf("%s must be less than %s", errors[0].Field(), errors[0].Param())
	case "alphanumWithSpaces":
		message = fmt.Sprintf("%s only characters allowed are (a-z) and (A-Z) and (0-9)", errors[0].Field())
	case "eqfield":
		if errors[0].Field() == "ConfirmNewPassword" && errors[0].Param() == "NewPassword" {
			message = fmt.Sprintf("%s must equal %s", strings.ToLower(errors[0].Field()), strings.ToLower(errors[0].Param()))
		}	
	case "min":
		if errors[0].Kind() == reflect.String {
			message = fmt.Sprintf("%s minimum length allowed is %s", errors[0].Field(), errors[0].Param())
		} else {
			message = fmt.Sprintf("%s minimum allowed is %s", errors[0].Field(), errors[0].Param())
		}
	case "max":
		if errors[0].Kind() == reflect.String {
			message = fmt.Sprintf("%s maximum length allowed is %s", errors[0].Field(), errors[0].Param())
		} else {
			message = fmt.Sprintf("%s maximum allowed is %s", errors[0].Field(), errors[0].Param())
		}
	default:
		message = fmt.Sprintf("%s is invalid", errors[0].Field())
	}

	return message
}

func unmarshalErrMsgHandler(error *json.UnmarshalTypeError) string {
	errMsg := fmt.Sprintf("%v is type %v can't be equal to %v", error.Field, error.Type, error.Value)
	return errMsg
}

func logCaptureStackTrace() string {
	buf := make([]byte, 1<<16)
	n := runtime.Stack(buf, false)
	var fullSrackTrace = string(buf[:n])

	lines := strings.Split(fullSrackTrace, "\n")

	filteredStackTraceAsString := strings.Join(lines, "\n")
	log.Printf("\nstack Trace: %v\n", filteredStackTraceAsString)

	return filteredStackTraceAsString
}