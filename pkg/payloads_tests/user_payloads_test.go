package payloads_test

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"main.go/pkg/payloads"
	"main.go/pkg/test_utils"
	"main.go/pkg/utils"
)

func TestUserPayload_UserLogin(t *testing.T) {
	t.Run("Should return an error when payload is empty", func(t *testing.T) {
		ul := payloads.UserLogin{
			Email:    "",
			Password: "",
		}

		err := utils.ValidateStruct(ul)
		if err != nil {
			err, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("an error has occurred during type assertion validationErrors")
			}
			errsCount := test_utils.GetErrorsCount(err)
			expectedErrCount := 2
			assert.Equal(t, expectedErrCount, errsCount, test_utils.ExpectedErrsCountMsg(expectedErrCount, errsCount))
		}
	})

	t.Run("Should return errors for each field that they are less than min length", func(t *testing.T) {
		ul := payloads.UserLogin{
			Email:    "email@gmail.com",
			Password: test_utils.StringLessThanMinPassword,
		}

		err := utils.ValidateStruct(ul)
		if err != nil {
			err, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("an error has occurred during type assertion validationErrors")
			}
			errsCount := test_utils.GetErrorsCount(err)
			expectedErrCount := 1
			assert.Equal(t, expectedErrCount, errsCount, test_utils.ExpectedErrsCountMsg(expectedErrCount, errsCount))
		}
	})

	t.Run("Should return errors for each field that they are more than max length", func(t *testing.T) {
		ul := payloads.UserLogin{
			Email:    test_utils.EmailMoreThanMaxLong65,
			Password: test_utils.StringMoreThanMaxPassword,
		}

		err := utils.ValidateStruct(ul)
		if err != nil {
			err, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("an error has occurred during type assertion validationErrors")
			}

			errsCount := test_utils.GetErrorsCount(err)
			expectedErrCount := 2
			assert.Equal(t, expectedErrCount, errsCount, test_utils.ExpectedErrsCountMsg(expectedErrCount, errsCount))
		}
	})

	t.Run("Should accept string fields at max lengths", func(t *testing.T) {
		ul := payloads.UserLogin{
			Email:    test_utils.EmailLong64,
			Password: test_utils.StringMaxPassword,
		}

		err := utils.ValidateStruct(ul)
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

	t.Run("Should accept string fields at min lengths", func(t *testing.T) {
		ul := payloads.UserLogin{
			Email:    "email@gmail.com",
			Password: test_utils.StringMinPassword,
		}

		err := utils.ValidateStruct(ul)
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
		ul := payloads.UserLogin{
			Email:    "email@gmail.com",
			Password: "abcdef",
		}

		err := utils.ValidateStruct(ul)
		if err != nil {
			t.Fatal("an error has occurred during validation")
		}

		ul.TrimStrs()
		test_utils.ValidateTrimming(t, ul)
	})
}

func TestUserPayload_UserSignUp(t *testing.T) {
	t.Run("Should return an error when payload is empty", func(t *testing.T) {
		ul := payloads.UserSignUp{
			Name:     "",
			Email:    "",
			Password: "",
		}

		err := utils.ValidateStruct(ul)
		if err != nil {
			err, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("an error has occurred during type assertion validationErrors")
			}
			errsCount := test_utils.GetErrorsCount(err)
			expectedErrCount := 3
			assert.Equal(t, expectedErrCount, errsCount, test_utils.ExpectedErrsCountMsg(expectedErrCount, errsCount))
		}
	})

	t.Run("Should return errors for each field that they are less than min length", func(t *testing.T) {
		ul := payloads.UserSignUp{
			Name:     "use",
			Email:    "email@gmail.com",
			Password: "abcde",
		}

		err := utils.ValidateStruct(ul)
		if err != nil {
			err, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("an error has occurred during type assertion validationErrors")
			}
			errsCount := test_utils.GetErrorsCount(err)
			expectedErrCount := 2
			assert.Equal(t, expectedErrCount, errsCount, test_utils.ExpectedErrsCountMsg(expectedErrCount, errsCount))
		}
	})

	t.Run("Should return errors for each field that they are more than max length", func(t *testing.T) {
		ul := payloads.UserSignUp{
			Name:     test_utils.StringLong33,
			Email:    test_utils.EmailMoreThanMaxLong65,
			Password: test_utils.StringLong25,
		}

		err := utils.ValidateStruct(ul)
		if err != nil {
			err, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("an error has occurred during type assertion validationErrors")
			}

			errsCount := test_utils.GetErrorsCount(err)
			expectedErrCount := 3
			assert.Equal(t, expectedErrCount, errsCount, test_utils.ExpectedErrsCountMsg(expectedErrCount, errsCount))
		}
	})

	t.Run("Should accept string fields at max lengths", func(t *testing.T) {
		ul := payloads.UserSignUp{
			Name:     test_utils.StringLong24,
			Email:    test_utils.EmailLong64,
			Password: test_utils.StringLong24,
		}

		err := utils.ValidateStruct(ul)
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

	t.Run("Should accept string fields at min lengths", func(t *testing.T) {
		ul := payloads.UserSignUp{
			Name:     "user",
			Email:    "email@gmail.com",
			Password: "abcdef",
		}

		err := utils.ValidateStruct(ul)
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
		usu := payloads.UserSignUp{
			Name:     " user ",
			Email:    "email@gmail.com",
			Password: " abcdef ",
		}

		err := utils.ValidateStruct(usu)
		if err != nil {
			t.Fatal("an error has occurred during validation")
		}

		usu.TrimStrs()
		test_utils.ValidateTrimming(t, usu)
	})
}

func TestUserPayload_ResetPassword(t *testing.T) {
	t.Run("Should return an error when payload is empty", func(t *testing.T) {
		ul := payloads.ResetPassword{
			OldPassword:        "",
			NewPassword:        "",
			ConfirmNewPassword: "",
		}

		err := utils.ValidateStruct(ul)
		if err != nil {
			err, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("an error has occurred during type assertion validationErrors")
			}
			errsCount := test_utils.GetErrorsCount(err)
			expectedErrCount := 3
			assert.Equal(t, expectedErrCount, errsCount, test_utils.ExpectedErrsCountMsg(expectedErrCount, errsCount))
		}
	})

	t.Run("Should return errors for each field that they are less than min length", func(t *testing.T) {
		ul := payloads.ResetPassword{
			OldPassword:        test_utils.StringMinPassword,
			NewPassword:        test_utils.StringMinPassword,
			ConfirmNewPassword: test_utils.StringMinPassword,
		}

		err := utils.ValidateStruct(ul)
		if err != nil {
			err, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("an error has occurred during type assertion validationErrors")
			}
			errsCount := test_utils.GetErrorsCount(err)
			expectedErrCount := 2
			assert.Equal(t, expectedErrCount, errsCount, test_utils.ExpectedErrsCountMsg(expectedErrCount, errsCount))
		}
	})

	t.Run("Should return errors for each field that they are more than max length", func(t *testing.T) {
		ul := payloads.ResetPassword{
			OldPassword:        test_utils.StringMoreThanMaxPassword,
			NewPassword:        test_utils.StringMoreThanMaxPassword,
			ConfirmNewPassword: test_utils.StringMoreThanMaxPassword,
		}

		err := utils.ValidateStruct(ul)
		if err != nil {
			err, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("an error has occurred during type assertion validationErrors")
			}

			errsCount := test_utils.GetErrorsCount(err)
			expectedErrCount := 3
			assert.Equal(t, expectedErrCount, errsCount, test_utils.ExpectedErrsCountMsg(expectedErrCount, errsCount))
		}
	})

	t.Run("Should accept string fields at max lengths", func(t *testing.T) {
		ul := payloads.ResetPassword{
			OldPassword:        test_utils.StringMaxPassword,
			NewPassword:        test_utils.StringMaxPassword,
			ConfirmNewPassword: test_utils.StringMaxPassword,
		}

		err := utils.ValidateStruct(ul)
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

	t.Run("Should accept string fields at min lengths", func(t *testing.T) {
		ul := payloads.ResetPassword{
			OldPassword:        test_utils.StringMinPassword,
			NewPassword:        test_utils.StringMinPassword,
			ConfirmNewPassword: test_utils.StringMinPassword,
		}

		err := utils.ValidateStruct(ul)
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

	t.Run("Should return an error when confirmNewPassword does not equal new password", func(t *testing.T) {
		ul := payloads.ResetPassword{
			OldPassword:        test_utils.StringMinPassword,
			NewPassword:        test_utils.StringMinPassword,
			ConfirmNewPassword: "random string",
		}

		err := utils.ValidateStruct(ul)
		if err != nil {
			err, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("an error has occurred during type assertion validationErrors")

				assert.Contains(t, err.Error(), "eqfield", "error must be from tag 'eqfield'")
			}
			errsCount := test_utils.GetErrorsCount(err)
			expectedErrCount := 1
			assert.Equal(t, expectedErrCount, errsCount, test_utils.ExpectedErrsCountMsg(expectedErrCount, errsCount))
		}
	})

	t.Run("Strings must be trimmed after using trim function", func(t *testing.T) {
		rp := payloads.ResetPassword{
			OldPassword:        " abced ",
			NewPassword:        " abced ",
			ConfirmNewPassword: " abced ",
		}

		err := utils.ValidateStruct(rp)
		if err != nil {
			t.Fatal("an error has occurred during validation")
		}

		rp.TrimStrs()
		test_utils.ValidateTrimming(t, rp)
	})
}

func TestUserPayload_UpdateProfile(t *testing.T) {
	t.Run("Should return true for IsEmpty Method", func(t *testing.T) {
		ul := payloads.UpdateProfile{
			Name:         "",
			Email:        "",
			MobileNumber: "",
		}

		err := utils.ValidateStruct(ul)
		if err != nil {
			t.Fatal("an error has occurred during type assertion validationErrors")

		}

		assert.Equal(t, true, ul.IsEmpty(), "Expected 'IsEmpty' Method on 'UpdateProfile' to return true")
	})

	t.Run("Should return errors for each field that they are less than min length", func(t *testing.T) {
		ul := payloads.UpdateProfile{
			Name:        "abc",
			Email:        "email@gmail.com",
			MobileNumber: "0598727",
		}

		err := utils.ValidateStruct(ul)
		if err != nil {
			err, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("an error has occurred during type assertion validationErrors")
			}
			errsCount := test_utils.GetErrorsCount(err)
			expectedErrCount := 2
			assert.Equal(t, expectedErrCount, errsCount, test_utils.ExpectedErrsCountMsg(expectedErrCount, errsCount))
		}
	})

	t.Run("Should return errors for each field that they are more than max length", func(t *testing.T) {
		ul := payloads.UpdateProfile{
			Name:        test_utils.StringLong33,
			Email:       test_utils.EmailMoreThanMaxLong65,
			MobileNumber:test_utils.MobileNumberMoreThanMax,
		}

		err := utils.ValidateStruct(ul)
		if err != nil {
			err, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("an error has occurred during type assertion validationErrors")
			}

			errsCount := test_utils.GetErrorsCount(err)
			expectedErrCount := 3
			assert.Equal(t, expectedErrCount, errsCount, test_utils.ExpectedErrsCountMsg(expectedErrCount, errsCount))
		}
	})

	t.Run("Should accept string fields at max lengths", func(t *testing.T) {
		ul := payloads.UpdateProfile{
			Name:        test_utils.StringLong32,
			Email:       test_utils.EmailLong64,
			MobileNumber:test_utils.MobileNumberMax,
		}

		err := utils.ValidateStruct(ul)
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

	t.Run("Should accept string fields at min lengths", func(t *testing.T) {
		ul := payloads.UpdateProfile{
			Name:         "abcd",
			Email:        "email@gmail.com",
			MobileNumber: "05927183",
		}

		err := utils.ValidateStruct(ul)
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
		up := payloads.UpdateProfile{
			Name:         " abcd ",
			Email:        "email@gmail.com",
			MobileNumber: "  0598372 ",
		}

		err := utils.ValidateStruct(up)
		if err != nil {
			t.Fatal("an error has occurred during validation")
		}

		up.TrimStrs()
		test_utils.ValidateTrimming(t, up)
	})
}
