package payloads

import (
	"slices"
	"strings"

	"main.go/pkg/models"
)

type UserLogin struct {
	Email    string `json:"email" validate:"required,email,max=64"`
	Password string `json:"password" validate:"required,min=6,max=24"`
}

type UserSignUp struct {
	Name     string `json:"name" validate:"required,min=4,max=32,alphanumWithSpaces"`
	Email    string `json:"email" validate:"required,email,max=64"`
	Password string `json:"password" validate:"required,min=6,max=24"`
}

type ResetPassword struct {
	OldPassword string `json:"oldPassword" validate:"required,min=6,max=24"`
	NewPassword string `json:"newPassword" validate:"required,min=6,max=24"`
	ConfirmNewPassword string `json:"confirmNewPassword" validate:"required,min=6,max=24,eqfield=NewPassword"`
}

type AssignRolePayload struct {
	RoleId  uint `json:"roleId" validate:"required,min=1"`
}

func (usi *UserLogin) TrimStrs() *UserLogin {
	if usi != nil {
		usi.Email = strings.Trim(usi.Email, " ")
		usi.Password= strings.Trim(usi.Password, " ")
	}
	return usi
}

func (usu *UserSignUp) TrimStrs() *UserSignUp {
	if usu != nil {
		usu.Name = strings.Trim(usu.Name, " ")
		usu.Email = strings.Trim(usu.Email, " ")
		usu.Password = strings.Trim(usu.Password, " ")
	}

	return usu
}

func (rp *ResetPassword) TrimStrs() *ResetPassword {
	if rp != nil {
		rp.OldPassword = strings.Trim(rp.OldPassword, " ")
		rp.NewPassword = strings.Trim(rp.NewPassword, " ")
		rp.ConfirmNewPassword = strings.Trim(rp.ConfirmNewPassword, " ")
	}
	
	return rp
}

// Avatar must be set manually inside the route, if it was set with (ToModel) as the other fields might cause injection from the frontend.
// The purpose of this field (Avatar) is to satisfy (Exclude) operation which ensures that this field will be updated in the database.
type UpdateProfile struct {
	Name     string `json:"name,omitempty" validate:"omitempty,min=4,max=32"`
	Email    string `json:"email,omitempty" validate:"omitempty,email,max=64"`
	MobileNumber string `json:"mobileNumber,omitempty" validate:"omitempty,min=8,max=32"`
	Avatar string `json:"avatar,omitempty"`
}

func (up *UpdateProfile) TrimStrs() *UpdateProfile {
	if up != nil {
		up.Name = strings.Trim(up.Name, " ")
		up.Email = strings.Trim(up.Email, " ")
		up.MobileNumber = strings.Trim(up.MobileNumber, " ")
	}
	
	return up
}

func (up *UpdateProfile) Exclude(selectedFields []string) []string {
	removedCols := map[string]any{}
	if up.Name == "" {
		removedCols["Name"] = 1
	}
	if up.Email == "" {
		removedCols["Email"] = 1
	}
	if up.MobileNumber == "" {
		removedCols["MobileNumber"] = 1
	}
	if up.Avatar == "" {
		removedCols["Avatar"] = 1
	}
	
	selectedFields = slices.DeleteFunc(selectedFields, func(element string) bool {
		_, exists := removedCols[element]
		return exists
	})
	return selectedFields
}

func (up *UpdateProfile) ToModel() *models.User {
	if up != nil {
		up.Name = strings.Trim(up.Name, " ")
		up.Email = strings.Trim(up.Email, " ")
		up.MobileNumber = strings.Trim(up.MobileNumber, " ")
	}
	
	return &models.User{
		Name: up.Name,
		Email: up.Email,
		MobileNumber: &up.MobileNumber,
	}
}

// must be called after TrimStrs
func (up *UpdateProfile) IsEmpty() bool {
	return (up == nil) || (up.Email == "" && up.MobileNumber == "" && up.Name == "") 
}

