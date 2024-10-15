package payloads

import (
	"strings"

	"main.go/pkg/models"
)

type CreateRole struct {
	Name string `json:"name" validate:"required,min=2,max=32,alphanumWithSpaces"`
}

type UpdateRole struct {
	Name string `json:"name" validate:"required,min=2,max=32,alphanumWithSpaces"`
}

func (ur *UpdateRole) TrimStrs() *UpdateRole {
	if ur != nil {
		ur.Name = strings.Trim(ur.Name, " ")
	}

	return ur
}

func (ur *UpdateRole) ToModel() *models.Role {
	if ur != nil {
		return &models.Role{
			Name: ur.Name,
		}
	}

	return nil
}

func (cr *CreateRole) TrimStrs() *CreateRole {
	if cr != nil {
		cr.Name = strings.Trim(cr.Name, " ")
	}

	return cr
}

func (cr *CreateRole) ToModel() *models.Role {
	if cr != nil {
		return &models.Role{
			Name: cr.Name,
		}
	}

	return nil
}