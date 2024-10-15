package payloads

import (
	"slices"
	"strings"

	"main.go/pkg/models"
)

type CreateAddress struct {
	FullName string `json:"fullName" validate:"required,min=4,max=32"`
	City          string  `json:"city" validate:"required,min=3,max=32"`
	StreetAddress string  `json:"streetAddress" validate:"required,min=4,max=64"`
	State         *string `json:"state" validate:"min=4,max=32"`
	ZipCode       *string `json:"zipCode" validate:"min=3,max=12"`
	Country       string  `json:"country" validate:"required,min=4,max=32"`
}

type UpdateAddress struct {
	FullName string `json:"fullName" validate:"min=4,max=32"`
	City          string  `json:"city" validate:"min=3,max=32"`
	StreetAddress string  `json:"streetAddress" validate:"min=4,max=64"`
	State         *string `json:"state" validate:"min=4,max=32"`
	ZipCode       *string `json:"zipCode" validate:"min=3,max=12"`
	Country       string  `json:"country" validate:"min=4,max=32"`
}

func (u *UpdateAddress) IsEmpty() bool {
	return u.City == "" && u.Country == "" && (u.State == nil || *u.State == "") && u.StreetAddress == "" && (u.ZipCode == nil || *u.ZipCode == "") && u.FullName == ""
}

func (u *UpdateAddress) ToModel() *models.Address {
	return &models.Address{
		FullName: u.FullName,
		City: u.City,
		StreetAddress: u.StreetAddress,
		State: u.State,
		ZipCode: u.ZipCode,
		Country: u.Country,
	}
}

func (u *CreateAddress) ToModel(userId uint) *models.Address {
	return &models.Address{
		FullName: u.FullName,
		City: u.City,
		StreetAddress: u.StreetAddress,
		State: u.State,
		ZipCode: u.ZipCode,
		Country: u.Country,
		UserID: userId,
	}
}

func (c *CreateAddress) TrimStrs() *CreateAddress {
	if c != nil {
		c.City = strings.Trim(c.City, " ")
		c.StreetAddress = strings.Trim(c.StreetAddress, " ")
		c.Country = strings.Trim(c.Country, " ")
		c.FullName = strings.Trim(c.FullName, " ")
		if c.State != nil {
			*c.State = strings.Trim(*c.State, " ")
		}
		if c.ZipCode != nil {
			*c.ZipCode = strings.Trim(*c.ZipCode, " ")
		}
	}
	return c
}

func (u *UpdateAddress) TrimStrs() *UpdateAddress {
	if u != nil {
		u.City = strings.Trim(u.City, " ")
		u.StreetAddress = strings.Trim(u.StreetAddress, " ")
		u.Country = strings.Trim(u.Country, " ")
		u.FullName = strings.Trim(u.FullName, " ")
		if u.State != nil {
			*u.State = strings.Trim(*u.State, " ")
		}
		if u.ZipCode != nil {
			*u.ZipCode = strings.Trim(*u.ZipCode, " ")
		}
	}
	return u
}

func (u *UpdateAddress) Exclude(selectedFields []string) []string {
	deletedCols := map[string]any{}
	if u.City == "" {
		deletedCols[u.City] = 1
	}
	if u.StreetAddress == "" {
		deletedCols[u.StreetAddress] = 1
	}
	if u.Country == "" {
		deletedCols[u.Country] = 1
	}
	if u.FullName == "" {
		deletedCols[u.FullName] = 1
	}
	if u.State == nil || *u.State == "" {
		deletedCols[*u.State] = 1
	}
	if u.ZipCode == nil || *u.ZipCode == "" {
		deletedCols[*u.ZipCode] = 1
	}
	
	selectedFields = slices.DeleteFunc(selectedFields, func(element string) bool {
		_, exists := deletedCols[element]
		return exists
	})

	return selectedFields
}