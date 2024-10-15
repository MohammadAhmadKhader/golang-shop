package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"main.go/config"
	"main.go/constants"
	"main.go/pkg/models"
)

const GenericErrMessage = "An unexpected error has occurred, please try again later!"

func Trim(str *string) *string {
	ts := strings.Trim(*str, "")
	return &ts
}

func IsEmptyStr(str string) bool {
	return len(str) == 0
}

func IsDefaultFloat64(str string) bool {
	return len(str) == 0
}

func ConvertStrToUint(str string) (*uint, error) {
	uint64Val, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return nil, err
	}
	uintVal := uint(uint64Val)

	return &uintVal, nil
}

func ConvertStrToFloat64(str string) (*float64, error) {
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return nil, err
	}

	return &val, nil
}

func CalculateOffset(page, limit int) int {
	offset := (limit * page) - limit
	return offset
}

// payload here must a pointer, this function parse the body and set the result of parsing to the payload
func ParseBodyJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("request body is empty")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func ParseFormFileJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("request body is empty")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

// ! we dont need to capture stack trace for every failing request
func WriteError(w http.ResponseWriter, status int, err error) {
	stackTrace := logCaptureStackTrace()
	returnedErr := err.Error()

	errObj := make(map[string]any, 0)
	var UnmarshalTypeErr *json.UnmarshalTypeError
	var validationErrMsg validator.ValidationErrors

	if config.Envs.Env == "production" {
		if status >= 500 {
			returnedErr = GenericErrMessage
			errObj["error"] = GenericErrMessage
			errObj["statusCode"] = status

		} else if errors.As(err, &validationErrMsg) {
			errObj["error"] = validationErrMsgHandler(err.(validator.ValidationErrors))
			errObj["statusCode"] = status

		} else if errors.As(err, &UnmarshalTypeErr) {
			errObj["error"] = unmarshalErrMsgHandler(err.(*json.UnmarshalTypeError))
			errObj["statusCode"] = 400

		} else {
			errObj["error"] = returnedErr
			errObj["statusCode"] = status
		}

	} else {
		errObj["error"] = returnedErr
		errObj["statusCode"] = status
		errObj["stackTrace"] = stackTrace
	}

	WriteJSON(w, status, errObj)
}

// it validate the parse the body based on the passed type
func ValidateAndParseBody[TPayload any](r *http.Request) (*TPayload, error) {
	var payload TPayload
	if err := ParseBodyJSON(r, &payload); err != nil {
		return nil, err
	}

	if err := Validate.Struct(payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

func ValidateAndParseFormData[TPayload any](r *http.Request, formValuesGetter func() (*TPayload, error)) (*TPayload, error) {
	if err := r.ParseMultipartForm(2 << 20); err != nil {
		return nil, err
	}

	payload, err := formValuesGetter()
	if err != nil {
		return nil, err
	}

	if err := Validate.Struct(payload); err != nil {
		return nil, err
	}

	return payload, nil
}

func HandleMultipleFilesUpload(r *http.Request, sizeInMB int64, keyName string) ([]*multipart.FileHeader, error) {
	if err := r.ParseMultipartForm(sizeInMB << 20); err != nil {
		return nil, fmt.Errorf("an error has occurred during parsing files")
	}

	files := r.MultipartForm.File[keyName]

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("unable to open file")
		}

		defer file.Close()
	}

	return files, nil
}

func HandleOneFileUpload(r *http.Request, sizeInMB int64, keyName string) (multipart.File, *multipart.FileHeader, error) {
	if err := r.ParseMultipartForm(sizeInMB << 20); err != nil {
		return nil, nil, fmt.Errorf("an error has occurred during parsing files")
	}

	file, fileHeader, err := r.FormFile(keyName)
	if err != nil {
		return nil, nil, fmt.Errorf("no file was found")
	}
	defer file.Close()

	return file, fileHeader, nil
}

func RoutePath(method, path string) string {
	return method + " " + constants.Prefix + path
}

func GetUserIdFromTokenPayload(r *http.Request) (*uint, error) {
	tokenClaims, ok := r.Context().Value(constants.TokenPayload).(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("no token claims found")
	}

	userID, ok := tokenClaims["userId"].(string)
	if !ok {
		return nil, fmt.Errorf("userId not found in token claims")
	}
	userIdInt, err := strconv.Atoi(userID)
	if err != nil {
		return nil, err
	}

	userIDAsUint := uint(userIdInt)

	return &userIDAsUint, nil
}

func GetUserEmailFromTokenPayload(r *http.Request) (*string, error) {
	tokenClaims, ok := r.Context().Value(constants.TokenPayload).(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("no token claims found")
	}

	email, ok := tokenClaims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("email not found in token claims")
	}

	return &email, nil
}

func GetValidateId(r *http.Request, pathKey string) (*uint, error) {
	IdAsString := r.PathValue(pathKey)
	Id, err := strconv.ParseUint(IdAsString, 10, 64)
	if err != nil {
		return nil, err
	}
	if Id == 0 {
		return nil, fmt.Errorf("id can not be 0")
	}
	uintId := uint(Id)

	return &uintId, nil
}

// TODO: this used to prettify middlewares and make them more readable
func MWCaller(funcs ...http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		for _, hf := range funcs {
			hf(w, r)
		}
	}
}

func GetUserIdCtx(r *http.Request) (*uint, error) {
	user, ok := r.Context().Value(constants.UserKey).(*models.User)
	if !ok {
		return nil, fmt.Errorf("user id was not found in context")
	}

	return &user.ID, nil
}
