package category_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"gorm.io/gorm"
	"main.go/internal/database"
	"main.go/pkg/models"
	"main.go/pkg/test_utils"
	"main.go/services"
)

func TestCategoriesHandler(t *testing.T) {
	server := http.NewServeMux()
	services.SetupAllServices(database.DB, server)
	categoryIdNotExistent := 9999999

	t.Run("Should get categories and return 200 status code", func(t *testing.T) {
		t.Parallel()
		req, err := http.NewRequest("GET", test_utils.GetRoutePath("/categories"), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, req)

		test_utils.AssertBodyType[test_utils.CategoryGetAll](t, rr, nil)
		test_utils.ExpectStatusCode(t, rr, http.StatusOK)
	})

	t.Run("Should fail to get a non-existent category and return 400 status code", func(t *testing.T) {
		t.Parallel()
		req, err := http.NewRequest("GET", test_utils.GetRoutePath("/categories/9999999"), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, req)

		test_utils.ExpectStatusCode(t, rr, http.StatusBadRequest)
	})

	t.Run("Should get a category by id successfully and return 200 status code", func(t *testing.T) {
		t.Parallel()
		categoryId := uint(15)
		req, err := http.NewRequest("GET", test_utils.GetRoutePath(fmt.Sprintf("/categories/%v",categoryId)), nil)
		if err != nil {
			t.Fatal(err)
		}
		
		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, req)

		test_utils.AssertBodyType(t, rr, func(response test_utils.CategoryGetOne) bool {
			return response.Category.ID == categoryId
		})
		test_utils.ExpectStatusCode(t, rr, http.StatusOK)
	})

	t.Run("Should failed to get category with id that does not exist", func(t *testing.T) {
		t.Parallel()
		req, err := http.NewRequest("GET", test_utils.GetRoutePath(fmt.Sprintf("/categories/%v",categoryIdNotExistent)), nil)
		if err != nil {
			t.Fatal(err)
		}
		
		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, req)

		test_utils.ExpectStatusCode(t, rr, http.StatusBadRequest)
	})

	t.Run("Should create category successfully and return 201 status code", func(t *testing.T) {
		t.Parallel()

		category := models.Category{
			Name: test_utils.CapStrLen(gofakeit.Name(), 32),
		}
		reqBody := test_utils.CreateRequestBody(t, category)
		req, err := http.NewRequest("POST", test_utils.GetRoutePath("/categories"), reqBody)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		err = test_utils.GenSuperAdminCookie(rr, req)
		if err != nil {
			t.Fatal(err)
		}

		server.ServeHTTP(rr, req)
		
		test_utils.AssertBodyType[test_utils.CategoryCreate](t, rr, func(r test_utils.CategoryCreate) bool {
			return r.Category.Name == category.Name
		})
		test_utils.ExpectStatusCode(t, rr, http.StatusCreated)
	})

	t.Run("Should soft delete category by id and return 204 status code", func(t *testing.T) {
		t.Parallel()
		category, err := test_utils.CreateTestCategory(nil)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("DELETE", test_utils.GetRoutePath(fmt.Sprintf("/categories/%v/soft-delete", category.ID)), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		err = test_utils.GenSuperAdminCookie(rr, req)
		if err != nil {
			t.Fatal(err)
		}

		server.ServeHTTP(rr, req)

		test_utils.ExpectEmptyJSON(t ,rr)
		test_utils.ExpectStatusCode(t, rr, http.StatusNoContent)
	})

	t.Run("Should restore a soft deleted category by id and return 200 status code", func(t *testing.T) {
		t.Parallel()

		category, err := test_utils.CreateTestCategory(func(category *models.Category) *models.Category {
			category.DeletedAt = &gorm.DeletedAt{Time: time.Now()}
			return category
		})
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("PATCH", test_utils.GetRoutePath(fmt.Sprintf("/categories/%v/restore", category.ID)), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		err = test_utils.GenSuperAdminCookie(rr, req)
		if err != nil {
			t.Fatal(err)
		}

		server.ServeHTTP(rr, req)

		test_utils.AssertBodyType(t, rr, func(response test_utils.CategoryRestore) bool {
			return response.Category.ID == category.ID
		})
		
		test_utils.ExpectStatusCode(t, rr, http.StatusOK)
	})

	t.Run("Should return an error that category is not found on attempt to restore", func(t *testing.T) {
		t.Parallel()
		req, err := http.NewRequest("PATCH", test_utils.GetRoutePath(fmt.Sprintf("/categories/%v/restore", categoryIdNotExistent)), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		err = test_utils.GenSuperAdminCookie(rr, req)
		if err != nil {
			t.Fatal(err)
		}

		server.ServeHTTP(rr, req)

		test_utils.ExpectStatusCode(t, rr, http.StatusBadRequest)
	})

	t.Run("Should hard delete category by id and return 204 status code", func(t *testing.T) {
		t.Parallel()

		category, err := test_utils.CreateTestCategory(nil)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("DELETE", test_utils.GetRoutePath(fmt.Sprintf("/categories/%v", category.ID)), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		err = test_utils.GenSuperAdminCookie(rr, req)
		if err != nil {
			t.Fatal(err)
		}

		server.ServeHTTP(rr, req)

		test_utils.ExpectEmptyJSON(t, rr)
		test_utils.ExpectStatusCode(t, rr, http.StatusNoContent)
	})

	t.Run("Should Update category successfully and return status code 202", func(t *testing.T) {
		changes := models.Category{
			Name: test_utils.CapStrLen(gofakeit.Name(), 32),
		}
		categoryId := uint(17)
		reqBody := test_utils.CreateRequestBody(t, changes)
		req, err := http.NewRequest("PUT", test_utils.GetRoutePath(fmt.Sprintf("/categories/%v", categoryId)), reqBody)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		err = test_utils.GenSuperAdminCookie(rr, req)
		if err != nil {
			t.Fatal(err)
		}

		server.ServeHTTP(rr, req)

		test_utils.AssertBodyType(t, rr, func(resp test_utils.CategoryUpdate) bool {
			r := resp.Category
			return r.Name == changes.Name && r.ID == categoryId
		})
		test_utils.ExpectStatusCode(t, rr, http.StatusAccepted)
	})

	t.Run("Should return an error that category does not exist on update attempt", func(t *testing.T) {
		t.Parallel()
		changes := models.Category{Name: test_utils.CapStrLen(gofakeit.Name(), 32)}
		reqBody := test_utils.CreateRequestBody(t, changes)
		req, err := http.NewRequest("PUT", test_utils.GetRoutePath(fmt.Sprintf("/categories/%v", categoryIdNotExistent)), reqBody)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		err = test_utils.GenSuperAdminCookie(rr, req)
		if err != nil {
			t.Fatal(err)
		}

		server.ServeHTTP(rr, req)

		test_utils.ExpectStatusCode(t, rr, http.StatusBadRequest)
	})

	t.Run("Should return an error that category does not exist on delete attempt", func(t *testing.T) {
		t.Parallel()
		req, err := http.NewRequest("DELETE", test_utils.GetRoutePath(fmt.Sprintf("/categories/%v", categoryIdNotExistent)), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		err = test_utils.GenSuperAdminCookie(rr, req)
		if err != nil {
			t.Fatal(err)
		}

		server.ServeHTTP(rr, req)

		test_utils.ExpectStatusCode(t, rr, http.StatusBadRequest)
	})
}