package test_utils

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"main.go/constants"
	"main.go/internal/database"
	"main.go/pkg/models"
	"main.go/services/auth"
)

func GetRoutePath(path string) string {
	return constants.Prefix + path
}

func GenSuperAdminCookie(w http.ResponseWriter, r *http.Request) error {
	user := &models.User{
		ModelBasicsTrackedDel: models.ModelBasicsTrackedDel{ID: 17},
		Email:                 "texteemail@gmail.com",
	}
	token, err := auth.CreateJWT(*user, w, r)
	if err != nil {
		return err
	}
	SetCookieForTesting(w, r, user, token)
	return nil
}

// the difference between this and the one used in production, this uses "Get" method, the one for production uses "New" method
//
// theoretically it should not be causing any problem to use any of them
// but "New" has issue and not setting a new cookie on testing server, for that reason this function is created.
func SetCookieForTesting(w http.ResponseWriter, r *http.Request, user *models.User, token string) (*sessions.Session, error) {
	session, err := auth.CookiesStore.Get(r, "session_token")
	if err != nil {
		return nil, err
	}

	session.Options = &sessions.Options{
		MaxAge:   auth.CookieMaxAge,
		Path:     constants.Prefix,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}

	session.Values["userId"] = user.ID
	session.Values["email"] = user.Email
	session.Values["token"] = token

	err = session.Save(r, w)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// this throws "multipart: NextPart: EOF" and return correct response on create product which is 201 and uploads the required file
func CreateImageFromData(keyName string, imagePath string) (*bytes.Buffer, *multipart.Writer, error) {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	file, err := os.Open(imagePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	formFile, err := writer.CreateFormFile(keyName, "testImage.jpg")
	if err != nil {
		return nil, nil, err
	}

	_, err = io.Copy(formFile, file)
	if err != nil {
		return nil, nil, err
	}

	return &b, writer, nil
}

func CreateTestProduct(productAdjuster func(prod *models.Product) *models.Product) (*models.Product, error) {
	var product models.Product
	product.Name = "test product"
	product.Quantity = 20
	product.Price = 400
	product.CategoryID = 1
	if productAdjuster != nil {
		productAdjuster(&product)
	}

	err := database.DB.Model(models.Product{}).Create(&product).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func CreateTestCategory(Adjuster func(prod *models.Category) *models.Category) (*models.Category, error) {
	var category models.Category
	category.Name = "test category"
	if Adjuster != nil {
		Adjuster(&category)
	}

	err := database.DB.Model(models.Category{}).Create(&category).Error
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func ExpectStatusCode(t *testing.T, rr *httptest.ResponseRecorder, expectedStatus int) {
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}
}

func AssertBodyType[TResponseBody any](t *testing.T, rr *httptest.ResponseRecorder, extraValidator func(response TResponseBody) bool) {
	var responseBody TResponseBody

	err := json.NewDecoder(rr.Body).Decode(&responseBody)
	if err != nil {
		t.Fatal("an error has occurred during decoding response body: ", err)
	}

	assert.NotZero(t, responseBody, "result should not be the zero value")

	dataType := reflect.TypeOf(responseBody)
	if extraValidator != nil {
		isValid := extraValidator(responseBody)

		if !isValid {
			log.Println("received: ", responseBody)
			t.Fatalf("extraValidator has failed to validate type of (%v)", dataType)
		}
	}
}

func CapStrLen(str string, maxLength int) string {
	if len(str) > maxLength {
		return str[:maxLength]
	}
	return str
}

func CreateRequestBody(t *testing.T, value any) *bytes.Buffer {
	jsonBody, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("could not marshal request body: %v", err)
	}

	return bytes.NewBuffer(jsonBody)
}

func ExpectEmptyJSON(t *testing.T, rr *httptest.ResponseRecorder) {
	responseData, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}

	assert.JSONEq(t, `{}`, string(responseData), "Expected an empty JSON object")
}