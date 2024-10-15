package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"

	"github.com/go-playground/validator/v10"
)

const allowedSize int64 = 10 * 1024 * 1024 // 10 MB
var allowedMimeTypes = []string{"image/jpeg", "image/png", "image/jpg"}

var Validate = validator.New()

func init(){
	Validate.RegisterValidation("alphanumWithSpaces", isAlphanumericWithSpaces)
}

func isAlphanumericWithSpaces(fl validator.FieldLevel) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9 ]+$`)
	
	return regex.MatchString(fl.Field().String())
}


func ValidateFile(fileHeader *multipart.FileHeader) error {
	if fileHeader.Size > allowedSize {
		return errors.New("file size can not exceed 10 MB")
	}

	buf := make([]byte, 512)
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("error during open file process: %v", err)
	}

	_, err = io.ReadFull(file, buf)
	if err != nil {
		return fmt.Errorf("error during reading: %v", err)
	}
	fileMimeType := http.DetectContentType(buf)
	for _, mimeType := range allowedMimeTypes {
		if mimeType == fileMimeType {
			return nil
		}
	}

	return fmt.Errorf("invalid file type: %v", fileMimeType)
}