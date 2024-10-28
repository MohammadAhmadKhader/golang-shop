package payloads_test

import (
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"main.go/pkg/payloads"
	"main.go/pkg/test_utils"
	"main.go/pkg/utils"
)

func TestAddressPayload_Create(t *testing.T) {
	t.Run("Should return an error when payload is empty", func(t *testing.T) {
		state := ""
		zipCode := ""
		ca := payloads.CreateAddress{
			FullName:      "",
			StreetAddress: "",
			City:          "",
			State:         &state,
			Country:       "",
			ZipCode:       &zipCode,
		}

		err := utils.ValidateStruct(ca)
		if err != nil {
			err, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("an error has occurred during type assertion validationErrors")
			}
			errsCount := test_utils.GetErrorsCount(err)
			expectedErrCount := 6
			assert.Equal(t, expectedErrCount, errsCount, test_utils.ExpectedErrsCountMsg(expectedErrCount, errsCount))
		}
	})

	t.Run("Should return errors for each field that they are less than min length", func(t *testing.T) {
		state := "abc"
		zipCode := "ab"
		ca := payloads.CreateAddress{
			FullName:      "abc",
			StreetAddress: "abc",
			City:          "ab",
			State:         &state,
			Country:       "abc",
			ZipCode:       &zipCode,
		}

		err := utils.ValidateStruct(ca)
		if err != nil {
			err, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("an error has occurred during type assertion validationErrors")
			}
			errsCount := test_utils.GetErrorsCount(err)
			expectedErrCount := 6
			assert.Equal(t, expectedErrCount, errsCount, test_utils.ExpectedErrsCountMsg(expectedErrCount, errsCount))
		}
	})

	t.Run("Should return errors for each field that they are more than max length", func(t *testing.T) {
		state := test_utils.StringLong33
		zipCode := test_utils.StringLong13
		ca := payloads.CreateAddress{
			FullName:      "abc",
			StreetAddress: "abc",
			City:          test_utils.StringLong33,
			State:         &state,
			Country:       test_utils.StringLong33,
			ZipCode:       &zipCode,
		}

		err := utils.ValidateStruct(ca)
		if err != nil {
			err, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("an error has occurred during type assertion validationErrors")
			}
			errsCount := test_utils.GetErrorsCount(err)
			expectedErrCount := 6
			assert.Equal(t, expectedErrCount, errsCount, test_utils.ExpectedErrsCountMsg(expectedErrCount, errsCount))
		}
	})

	t.Run("Should accept string fields at max lengths", func(t *testing.T) {
		state := test_utils.StringLong32
		zipCode := test_utils.StringLong12
		ca := payloads.CreateAddress{
			FullName:      test_utils.StringLong32,
			StreetAddress: test_utils.StringLong64,
			City:          test_utils.StringLong32,
			State:         &state,
			Country:       test_utils.StringLong32,
			ZipCode:       &zipCode,
		}

		err := utils.ValidateStruct(ca)
		if err != nil {
			err, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("an error has occurred during type assertion validationErrors")
			}
			errsCount := test_utils.GetErrorsCount(err)
			expectedErrCount := 0
			assert.Equal(t, expectedErrCount, errsCount, test_utils.ExpectedErrsCountMsg(expectedErrCount, errsCount))
		}
	})

	t.Run("Should accept string fields at max lengths", func(t *testing.T) {
		state := "abcd"
		zipCode := "abc"
		ca := payloads.CreateAddress{
			FullName:      "abcd",
			StreetAddress: "abcd",
			City:          "abc",
			State:         &state,
			Country:       "abcd",
			ZipCode:       &zipCode,
		}

		err := utils.ValidateStruct(ca)
		if err != nil {
			err, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("an error has occurred during type assertion validationErrors")
			}
			errsCount := test_utils.GetErrorsCount(err)
			expectedErrCount := 0
			assert.Equal(t, expectedErrCount, errsCount, test_utils.ExpectedErrsCountMsg(expectedErrCount, errsCount))
		}
	})

	t.Run("Strings must be trimmed after using trim function", func(t *testing.T) {
		state := " abcd "
		zipCode := " abc "
		ca := payloads.CreateAddress{
			FullName:      " abcd   ",
			StreetAddress: " abcd  ",
			City:          "  abc  ",
			State:         &state,
			Country:       " abcd ",
			ZipCode:       &zipCode,
		}

		err := utils.ValidateStruct(ca)
		if err != nil {
			err, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("an error has occurred during type assertion validationErrors")
			}
			errsCount := test_utils.GetErrorsCount(err)
			expectedErrCount := 0
			assert.Equal(t, expectedErrCount, errsCount, test_utils.ExpectedErrsCountMsg(expectedErrCount, errsCount))
		}

		trimmedStrs := ca.TrimStrs()
		test_utils.ValidateTrimming(t, *trimmedStrs)
	})
}

func TestAddressPayload_Update(t *testing.T) {
	t.Run("Should return an error when payload is empty", func(t *testing.T) {
		ca := payloads.UpdateAddress{
			FullName:      "",
			StreetAddress: "",
			City:          "",
			State:         nil,
			Country:       "",
			ZipCode:       nil,
		}

		err := utils.ValidateStruct(ca)
		if err != nil {
			err, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("an error has occurred during type assertion validationErrors")
			}
			errsCount := test_utils.GetErrorsCount(err)
			expectedErrCount := 0
			fmt.Println(err.Error())
			assert.Equal(t, expectedErrCount, errsCount, test_utils.ExpectedErrsCountMsg(expectedErrCount, errsCount))
		}

		isEmpty := ca.IsEmpty()
		assert.Equal(t, true, isEmpty, "expected IsEmpty to return true value")
	})

	t.Run("Should return errors for each field that they are less than min length", func(t *testing.T) {
		state := "abc"
		zipCode := "ab"
		ca := payloads.UpdateAddress{
			FullName:      "abc",
			StreetAddress: "abc",
			City:          "ab",
			State:         &state,
			Country:       "abc",
			ZipCode:       &zipCode,
		}

		err := utils.ValidateStruct(ca)
		if err != nil {
			err, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("an error has occurred during type assertion validationErrors")
			}
			errsCount := test_utils.GetErrorsCount(err)
			expectedErrCount := 6
			assert.Equal(t, expectedErrCount, errsCount, test_utils.ExpectedErrsCountMsg(expectedErrCount, errsCount))
		}
	})

	t.Run("Should return errors for each field that they are more than max length", func(t *testing.T) {
		state := test_utils.StringLong33
		zipCode := test_utils.StringLong13
		ca := payloads.UpdateAddress{
			FullName:      "abc",
			StreetAddress: "abc",
			City:          test_utils.StringLong33,
			State:         &state,
			Country:       test_utils.StringLong33,
			ZipCode:       &zipCode,
		}

		err := utils.ValidateStruct(ca)
		if err != nil {
			err, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("an error has occurred during type assertion validationErrors")
			}
			errsCount := test_utils.GetErrorsCount(err)
			expectedErrCount := 6
			assert.Equal(t, expectedErrCount, errsCount, test_utils.ExpectedErrsCountMsg(expectedErrCount, errsCount))
		}
	})

	t.Run("Should accept string fields at max lengths", func(t *testing.T) {
		state := test_utils.StringLong32
		zipCode := test_utils.StringLong12
		ca := payloads.UpdateAddress{
			FullName:      test_utils.StringLong32,
			StreetAddress: test_utils.StringLong64,
			City:          test_utils.StringLong32,
			State:         &state,
			Country:       test_utils.StringLong32,
			ZipCode:       &zipCode,
		}

		err := utils.ValidateStruct(ca)
		if err != nil {
			err, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("an error has occurred during type assertion validationErrors")
			}
			errsCount := test_utils.GetErrorsCount(err)
			expectedErrCount := 0
			assert.Equal(t, expectedErrCount, errsCount, test_utils.ExpectedErrsCountMsg(expectedErrCount, errsCount))
		}
	})

	t.Run("Should accept string fields at max lengths", func(t *testing.T) {
		state := "abcd"
		zipCode := "abc"
		ca := payloads.UpdateAddress{
			FullName:      "abcd",
			StreetAddress: "abcd",
			City:          "abc",
			State:         &state,
			Country:       "abcd",
			ZipCode:       &zipCode,
		}

		err := utils.ValidateStruct(ca)
		if err != nil {
			err, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("an error has occurred during type assertion validationErrors")
			}
			errsCount := test_utils.GetErrorsCount(err)
			expectedErrCount := 0
			assert.Equal(t, expectedErrCount, errsCount, test_utils.ExpectedErrsCountMsg(expectedErrCount, errsCount))
		}
	})

	t.Run("Strings must be trimmed after using trim function", func(t *testing.T) {
		state := " abcd "
		zipCode := " abc "
		ca := payloads.UpdateAddress{
			FullName:      " abcd   ",
			StreetAddress: " abcd  ",
			City:          "  abc  ",
			State:         &state,
			Country:       " abcd ",
			ZipCode:       &zipCode,
		}

		err := utils.ValidateStruct(ca)
		if err != nil {
			err, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("an error has occurred during type assertion validationErrors")
			}
			errsCount := test_utils.GetErrorsCount(err)
			expectedErrCount := 0
			assert.Equal(t, expectedErrCount, errsCount, test_utils.ExpectedErrsCountMsg(expectedErrCount, errsCount))
		}

		trimmedStrs := ca.TrimStrs()
		test_utils.ValidateTrimming(t, *trimmedStrs)
	})
}
