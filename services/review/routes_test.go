package review_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"main.go/internal/database"
	"main.go/pkg/models"
	"main.go/pkg/test_utils"
	"main.go/services"
)

func TestReviewsHandler(t *testing.T) {
	server := http.NewServeMux()
	services.SetupAllServices(database.DB, server)
	reviewIdNotExistent := 9999999

	t.Run("Should get all reviews with pagination and return 200 status code", func(t *testing.T) {
		t.Parallel()
		req, err := http.NewRequest("GET", test_utils.GetRoutePath("/reviews"), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		err = test_utils.GenSuperAdminCookie(rr, req)
		if err != nil {
			t.Fatal(err)
		}
		server.ServeHTTP(rr, req)

		test_utils.AssertBodyType[test_utils.Pagination](t, rr, nil)
		test_utils.ExpectStatusCode(t, rr, http.StatusOK)
	})

	t.Run("Should create review successfully and return 201 status code", func(t *testing.T) {
		userId := uint(61)
		productId := uint(55)

		review := models.Review{
			Comment: test_utils.CapStrLen(gofakeit.Name(), 32),
			Rate:    uint8(gofakeit.UintRange(1, 5)),
		}
		
		reqBody := test_utils.CreateRequestBody(t, review)
		req, err := http.NewRequest("POST", test_utils.GetRoutePath(fmt.Sprintf("/products/%v/reviews", productId)), reqBody)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		err = test_utils.GenCookieByUserId(rr, req, userId)
		if err != nil {
			t.Fatal(err)
		}
		server.ServeHTTP(rr, req)

		var createdRev models.Review
		test_utils.AssertBodyType(t, rr, func(r test_utils.ReviewCreate) bool {
			createdRev = r.Review
			return true
		})
		assert.Equal(t, productId, createdRev.ProductID, "created review and the review returned in response must be equal")
		test_utils.ExpectStatusCode(t, rr, http.StatusCreated)

		defer func() {
			err := test_utils.DeleteResourceById[models.Review](createdRev.ID)
			if err != nil {
				log.Fatal(err)
			}
		}()
	})

	t.Run("Should delete review by id and return 204 status code", func(t *testing.T) {
		t.Parallel()

		userId := uint(60)
		review, err := test_utils.CreateTestReview(func(rev *models.Review) *models.Review {
			rev.ProductID = 77
			rev.UserID = userId
			return rev
		})
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("DELETE", test_utils.GetRoutePath(fmt.Sprintf("/products/77/reviews/%v", review.ID)), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		err = test_utils.GenCookieByUserId(rr, req, userId)
		if err != nil {
			t.Fatal(err)
		}

		server.ServeHTTP(rr, req)

		test_utils.ExpectEmptyJSON(t, rr)
		test_utils.ExpectStatusCode(t, rr, http.StatusNoContent)

		t.Cleanup(func() {
			test_utils.DeleteResourceById[models.Review](review.ID)
		})
	})

	t.Run("Should Update review successfully and return status code 202", func(t *testing.T) {
		changes := models.Review{
			Comment: test_utils.CapStrLen(gofakeit.Name(), 32),
			Rate:    uint8(gofakeit.UintRange(1, 5)),
		}
		reviewId := uint(40)
		reqBody := test_utils.CreateRequestBody(t, changes)
		req, err := http.NewRequest("PUT", test_utils.GetRoutePath(fmt.Sprintf("/products/76/reviews/%v", reviewId)), reqBody)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		err = test_utils.GenCookieByUserId(rr, req, 46)
		if err != nil {
			t.Fatal(err)
		}

		server.ServeHTTP(rr, req)

		var revAfterUpdate models.Review

		test_utils.ExpectStatusCode(t, rr, http.StatusAccepted)
		test_utils.AssertBodyType(t, rr, func(response test_utils.ReviewUpdate) bool {
			revAfterUpdate = response.Review
			return true
		})
		assert.Equal(t, changes.Comment, revAfterUpdate.Comment, "returned review comment after update must be equal to the sent one")
		assert.Equal(t, changes.Rate, revAfterUpdate.Rate, "returned review rate after update must be equal to the sent one")
	})

	t.Run("Should return an error that review does not exist on update attempt", func(t *testing.T) {
		t.Parallel()
		userId := uint(33)
		changes := models.Review{
			Comment: test_utils.CapStrLen(gofakeit.Name(), 32),
		}
		reqBody := test_utils.CreateRequestBody(t, changes)
		req, err := http.NewRequest("PUT", test_utils.GetRoutePath(fmt.Sprintf("/products/22/reviews/%v", reviewIdNotExistent)), reqBody)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		err = test_utils.GenCookieByUserId(rr, req, userId)
		if err != nil {
			t.Fatal(err)
		}

		server.ServeHTTP(rr, req)
		test_utils.ExpectStatusCode(t, rr, http.StatusBadRequest)
	})

	t.Run("Should return an error that review does not exist on delete attempt", func(t *testing.T) {
		t.Parallel()
		req, err := http.NewRequest("DELETE", test_utils.GetRoutePath(fmt.Sprintf("/products/22/reviews/%v", reviewIdNotExistent)), nil)
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
