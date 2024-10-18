package payloads

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"slices"
	"strings"

	"main.go/pkg/models"
)

type CreateProduct struct {
	Name        string         `json:"name" validate:"required,max=32,min=3,alphanumWithSpaces"`
	Quantity    uint           `json:"quantity" validate:"required,max=10000,min=0"`
	Image       multipart.File `json:"image"`
	Description string        `json:"description" validate:"omitempty,max=256,min=4,alphanumWithSpaces"`
	CategoryID  uint           `json:"categoryId" validate:"required,min=1"`
	Price       float64        `json:"price" validate:"gt=0.0"`
}


type UpdateProduct struct {
	Name        string         `json:"name" validate:"omitempty,max=32,min=3,alphanumWithSpaces"`
	Quantity    uint           `json:"quantity" validate:"omitempty,min=0,max=10000"`
	Description string        `json:"description" validate:"omitempty,max=256,min=4,alphanumWithSpaces"`
	CategoryID  uint           `json:"categoryId" validate:"omitempty,min=1"`
	Price       float64        `json:"price" validate:"omitempty,gt=0.0"`
}

func (cp *CreateProduct) ToModelWithImage(url string) *models.Product {
	if cp != nil {
		return &models.Product{
			Name:        cp.Name,
			Quantity:    cp.Quantity,
			Description: &cp.Description,
			CategoryID:  cp.CategoryID,
			Price:       cp.Price,
		}
	}

	return nil
}

func (cp *CreateProduct) TrimStrs() *CreateProduct {
	if cp != nil {
		cp.Name = strings.Trim(cp.Name, " ")
		cp.Description = strings.Trim(cp.Description, " ")
	}

	return cp
}

func (up *UpdateProduct) TrimStrs() *UpdateProduct {
	if up != nil {
		up.Name = strings.Trim(up.Name, " ")
		up.Description = strings.Trim(up.Description, " ")
	}

	return up
}

func (up *UpdateProduct) ToModel() *models.Product {
	if up != nil {
		return &models.Product{
			Name:        up.Name,
			Quantity:    up.Quantity,
			Description: &up.Description,
			CategoryID:  up.CategoryID,
			Price:       up.Price,
		}
	}
	return nil
}

// ! Handle url in the model with image
func (up *UpdateProduct) ToModelWithImage(url string) *models.Product {
	if up != nil {
		return &models.Product{
			Name:        up.Name,
			Quantity:    up.Quantity,
			Description: &up.Description,
			CategoryID:  up.CategoryID,
			Price:       up.Price,
		}
	}
	return nil
}

func (up *UpdateProduct) Exclude(selectedFields []string) []string {
	removedCols := map[string]any{}
	if up.Name == "" {
		removedCols["Name"] = 1
	}
	if up.Price == 0.0 {
		removedCols["Price"] = 1
	}
	if up.CategoryID == 0 {
		removedCols["CategoryID"] = 1
	}
	if up.Quantity == 0 {
		removedCols["Quantity"] = 1
	}
	if up.Description == "" {
		removedCols["Description"] = 1
	}

	selectedFields = slices.DeleteFunc(selectedFields, func(element string) bool {
		_, exists := removedCols[element]
		return exists
	})
	return selectedFields
}

func NewCreatePayload(r *http.Request, uIntConvertor func(s string) (*uint, error), floatConvertor func(s string) (*float64, error)) (*CreateProduct, error) {
	qty, err := uIntConvertor(r.FormValue("quantity"))
	if err != nil {
		return nil, fmt.Errorf("invalid quantity")
	}

	categoryId, err := uIntConvertor(r.FormValue("categoryId"))
	if err != nil {
		return nil, fmt.Errorf("invalid category id")
	}

	price, err := floatConvertor(r.FormValue("price"))
	if err != nil {
		return nil, fmt.Errorf("invalid price")
	}
	
	payload := &CreateProduct{
		Name: r.FormValue("name"),
		Quantity: *qty,
		Description: r.FormValue("description"),
		CategoryID: *categoryId,
		Price: *price,
	}
	
	return payload, nil
}

func (up *UpdateProduct) IsEmpty() bool {
	return up.Name == "" && up.Quantity == 0 && up.Description == "" && up.Price == 0 && up.CategoryID == 0 
}