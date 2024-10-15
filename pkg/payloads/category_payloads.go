package payloads

import (
	"strings"

	"main.go/pkg/models"
)

// in this unique case both create and update have same conditions, but for the sake of separate of concern

type CreateCategory struct {
	Name string `json:"name" validate:"required,min=3,max=32,alphanumWithSpaces"`
}

type UpdateCategory struct {
	Name string `json:"name" validate:"required,min=3,max=32,alphanumWithSpaces"`
}

func (uc *UpdateCategory) TrimStrs() *UpdateCategory {
	if uc != nil {
		uc.Name = strings.Trim(uc.Name, " ")
	}
	
	return uc
}

func (uc *UpdateCategory) ToModel() *models.Category {
	if uc != nil {
		return &models.Category{
			Name:uc.Name,
		}
	}
	
	return nil
}

func (cc *CreateCategory) TrimStrs() *CreateCategory {
	if cc != nil {
		cc.Name = strings.Trim(cc.Name, " ")
	}
	
	return cc
}

func (cc *CreateCategory) ToModel() *models.Category {
	if cc != nil {
		return &models.Category{
			Name:cc.Name,
		}
	}
	
	return nil
}