package product_test

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

func TestProductsHandler(t *testing.T) {
	server := http.NewServeMux()
	services.SetupAllServices(database.DB, server)
	productIdNotExist := 9999999

	t.Run("Should get products and return 200 status code", func(t *testing.T) {
		t.Parallel()
		req, err := http.NewRequest("GET", test_utils.GetRoutePath("/products"), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, req)

		test_utils.AssertBodyType[test_utils.ProductGetAll](t, rr, nil)
		test_utils.ExpectStatusCode(t, rr, http.StatusOK)
	})

	t.Run("Should fail to get a non-existent product and return 400 status code", func(t *testing.T) {
		t.Parallel()
		req, err := http.NewRequest("GET", test_utils.GetRoutePath("/products/9999999"), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, req)

		
		test_utils.ExpectStatusCode(t, rr, http.StatusBadRequest)
	})

	t.Run("Should get a product by id successfully and return 200 status code", func(t *testing.T) {
		t.Parallel()
		productId := uint(77)
		req, err := http.NewRequest("GET", test_utils.GetRoutePath(fmt.Sprintf("/products/%v",productId)), nil)
		if err != nil {
			t.Fatal(err)
		}

		
		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, req)

		test_utils.AssertBodyType(t, rr, func(response test_utils.ProductGetOne) bool {
			return response.Product.ID == productId
		})
		test_utils.ExpectStatusCode(t, rr, http.StatusOK)
	})

	t.Run("Should failed to get product with id that does not exist", func(t *testing.T) {
		t.Parallel()
		req, err := http.NewRequest("GET", test_utils.GetRoutePath(fmt.Sprintf("/products/%v",productIdNotExist)), nil)
		if err != nil {
			t.Fatal(err)
		}

		
		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, req)

		test_utils.ExpectStatusCode(t, rr, http.StatusBadRequest)
	})

	t.Run("Should create product successfully and return 201 status code", func(t *testing.T) {
		t.Skip()
		t.Parallel()

		bytes, writer, err := test_utils.CreateImageFromData("image", "../../testsdata/pexels-photo-1667088.jpeg")
		if err != nil {
			t.Fatal(err)
		}
		writer.WriteField("quantity", "23")
		writer.WriteField("categoryId", "1")
		writer.WriteField("price", "400")
		writer.WriteField("name", "product name")

		writer.Close()
		req, err := http.NewRequest("POST", test_utils.GetRoutePath("/products"), bytes)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())
		rr := httptest.NewRecorder()
		err = test_utils.GenSuperAdminCookie(rr, req)
		if err != nil {
			t.Fatal(err)
		}

		server.ServeHTTP(rr, req)
		
		test_utils.AssertBodyType[test_utils.ProductCreate](t, rr, func(r test_utils.ProductCreate) bool {
			return r.Product.Image.IsMain != nil && *r.Product.Image.IsMain == true && r.Product.Image.ImageUrl != ""
		})
		test_utils.ExpectStatusCode(t, rr, http.StatusCreated)
	})

	t.Run("Should soft delete product by id and return 204 status code", func(t *testing.T) {
		t.Parallel()
		prod, err := test_utils.CreateTestProduct(nil)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("DELETE", test_utils.GetRoutePath(fmt.Sprintf("/products/%v/soft-delete", prod.ID)), nil)
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

	t.Run("Should restore a soft deleted product by id and return 200 status code", func(t *testing.T) {
		t.Parallel()

		prod, err := test_utils.CreateTestProduct(func(prod *models.Product) *models.Product {
			prod.DeletedAt = &gorm.DeletedAt{Time: time.Now()}
			return prod
		})
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("PATCH", test_utils.GetRoutePath(fmt.Sprintf("/products/%v/restore", prod.ID)), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		err = test_utils.GenSuperAdminCookie(rr, req)
		if err != nil {
			t.Fatal(err)
		}

		server.ServeHTTP(rr, req)

		test_utils.AssertBodyType(t, rr, func(response test_utils.ProductRestore) bool {
			return response.Product.ID == prod.ID
		})
		
		test_utils.ExpectStatusCode(t, rr, http.StatusOK)
	})

	t.Run("Should return an error that product is not found on attempt to restore", func(t *testing.T) {
		t.Parallel()
		req, err := http.NewRequest("PATCH", test_utils.GetRoutePath(fmt.Sprintf("/products/%v/restore", productIdNotExist)), nil)
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

	t.Run("Should hard delete product by id and return 204 status code", func(t *testing.T) {
		t.Parallel()

		prod, err := test_utils.CreateTestProduct(nil)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("DELETE", test_utils.GetRoutePath(fmt.Sprintf("/products/%v", prod.ID)), nil)
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

	t.Run("Should Update product successfully and return status code 202", func(t *testing.T) {
		desc := test_utils.CapStrLen(gofakeit.ProductDescription(), 256)
		changes := models.Product{
			Name: test_utils.CapStrLen(gofakeit.Name(), 32),
			Quantity: gofakeit.UintRange(1, 200),
			Price: gofakeit.Float64Range(30,500),
			Description: &desc,
		}
		productId := uint(111)
		reqBody := test_utils.CreateRequestBody(t, changes)
		req, err := http.NewRequest("PUT", test_utils.GetRoutePath(fmt.Sprintf("/products/%v", productId)), reqBody)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		err = test_utils.GenSuperAdminCookie(rr, req)
		if err != nil {
			t.Fatal(err)
		}

		server.ServeHTTP(rr, req)

		test_utils.AssertBodyType(t, rr, func(resp test_utils.ProductUpdate) bool {
			r := resp.Product
			return r.Name == changes.Name && r.ID == productId && 
			r.Quantity == changes.Quantity && r.Price == changes.Price &&
			r.Description != nil && *r.Description == *changes.Description
		})
		test_utils.ExpectStatusCode(t, rr, http.StatusAccepted)
	})

	t.Run("Should return an error that product does not exist on update attempt", func(t *testing.T) {
		t.Parallel()
		changes := models.Product{Name: test_utils.CapStrLen(gofakeit.Name(), 32)}
		reqBody := test_utils.CreateRequestBody(t, changes)
		req, err := http.NewRequest("PUT", test_utils.GetRoutePath(fmt.Sprintf("/products/%v", productIdNotExist)), reqBody)
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

	t.Run("Should return an error that product does not exist on delete attempt", func(t *testing.T) {
		t.Parallel()
		req, err := http.NewRequest("DELETE", test_utils.GetRoutePath(fmt.Sprintf("/products/%v", productIdNotExist)), nil)
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