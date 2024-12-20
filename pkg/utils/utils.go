package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	appErrs "main.go/errors"

	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"main.go/config"
	"main.go/constants"
	"main.go/pkg/models"
)

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

func WriteError(w http.ResponseWriter, status int, err error) {
	errObj := make(map[string]any, 0)
	var UnmarshalTypeErr *json.UnmarshalTypeError
	var mySqlError *mysql.MySQLError

	if config.Envs.Env == "production" {
		if status >= 500 {
			errObj["error"] = appErrs.ErrGenericMessage.Error()
			errObj["statusCode"] = status
			logCaptureStackTrace()

		} else if errors.As(err, &validator.ValidationErrors{}) {
			errObj["error"] = validationErrMsgHandler(err.(validator.ValidationErrors))
			errObj["statusCode"] = status

		} else if errors.As(err, &mySqlError) {
			errObj["error"] = appErrs.ErrGenericMessage.Error()
			errObj["statusCode"] = 500
			logCaptureStackTrace()

		} else if errors.As(err, &UnmarshalTypeErr) {
			errObj["error"] = unmarshalErrMsgHandler(err.(*json.UnmarshalTypeError))
			errObj["statusCode"] = 400

		} else {
			errObj["error"] = err.Error()
			errObj["statusCode"] = status
		}

	} else {
		errObj["error"] = err.Error()
		errObj["statusCode"] = status
		errObj["stackTrace"] = logCaptureStackTrace()
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

func ValidateStruct[TPayload any](model TPayload) error {
	if err := Validate.Struct(model); err != nil {
		return err
	}

	return nil
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
		return nil, nil, appErrs.ErrNoFileFound
	}
	defer file.Close()

	return file, fileHeader, nil
}

func RoutePath(method, path string) string {
	return method + " " + constants.Prefix + path
}

func GetUserIdFromToken(r *http.Request) (*uint, error) {
	tokenClaims, ok := r.Context().Value(constants.TokenPayload).(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("no token claims found")
	}

	userID, ok := tokenClaims["userId"].(float64)
	if !ok {
		return nil, fmt.Errorf("userId not found in token claims")
	}

	userIDAsUint := uint(userID)

	return &userIDAsUint, nil
}

func GetEmailFromToken(r *http.Request) (*string, error) {
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

func GetValidateId(r *http.Request, pathKey string) (id *uint, receivedStr string, err error) {
	IdAsString := r.PathValue(pathKey)
	var uintId uint
	Id, err := strconv.ParseUint(IdAsString, 10, 64)
	if err != nil {
		return &uintId, IdAsString, err
	}
	if Id == 0 {
		return &uintId, "",fmt.Errorf("id can not be 0")
	}
	uintId = uint(Id)

	return &uintId, "",nil
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

func GetUserCtx(r *http.Request) (*models.User, error) {
	user, ok := r.Context().Value(constants.UserKey).(*models.User)
	if !ok {
		return nil, fmt.Errorf("user id was not found in context")
	}

	return user, nil
}

func GetResourceCtx[TModel any](r *http.Request, modelName string) (*TModel, error) {
	model, ok := r.Context().Value(constants.ResourceKey).(*TModel)
	if !ok {
		return nil, fmt.Errorf("%v was not found in context", modelName)
	}

	return model, nil
}

func TruncateToTwoDecimals(value float64) float64 {
	return float64(int(value*100)) / 100
}

func GetFilesCount(r *http.Request, keyName string) (int, error) {
	if err := r.ParseMultipartForm(1); err != nil {
		return 0, err
	}
	files, ok := r.MultipartForm.File["images"]
	if !ok {
		return 0, fmt.Errorf("failed to access files")
	}
	return len(files), nil
}

func CopyCols(selectedFields []string) []string {
	colsCopy := make([]string, len(selectedFields))
	copy(colsCopy, selectedFields)
	return colsCopy
}