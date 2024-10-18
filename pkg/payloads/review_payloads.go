package payloads

import (
	"slices"
	"strings"

	"main.go/pkg/models"
)

type CreateReview struct {
	Rate    uint8  `json:"rate" validate:"required,oneof=1 2 3 4 5"`
	Comment string `json:"comment" validate:"required,min=2,max=256,alphanumWithSpaces"`
}

type UpdateReview struct {
	Rate    uint8  `json:"rate,omitempty" validate:"omitempty,oneof=1 2 3 4 5"`
	Comment string `json:"comment,omitempty" validate:"min=2,max=256,alphanumWithSpaces,omitempty"`
}

func (ur *UpdateReview) IsEmpty() bool {
	return ur.Rate == 0 && ur.Comment == ""
}

func (ur *UpdateReview) TrimStrs() *UpdateReview {
	if ur != nil {
		ur.Comment = strings.Trim(ur.Comment, " ")
	}
	
	return ur
}

func (cr *CreateReview) TrimStrs() *CreateReview {
	if cr!= nil {
		cr.Comment = strings.Trim(cr.Comment, " ")
	}
	
	return cr
}

func (ur *UpdateReview) ToModel(userId uint, productId uint) *models.Review {
	if ur != nil {
		return &models.Review{
			Rate:    ur.Rate,
			Comment: ur.Comment,
			UserID: userId,
			ProductID: productId,
		}
	}

	return nil
}

func (cr *CreateReview) ToModel(userId uint, productId uint) *models.Review {
	if cr != nil {
		return &models.Review{
			Rate:    cr.Rate,
			Comment: cr.Comment,
			UserID: userId,
			ProductID: productId,
		}
	}
	
	return nil
}

func (ur *UpdateReview) Exclude(selectedFields []string) []string {
	removedCols := map[string]any{}
	if ur.Comment == "" {
		removedCols["Comment"] = 1
	}
	if ur.Rate == 0 {
		removedCols["Rate"] = 1
	}
	
	selectedFields = slices.DeleteFunc(selectedFields, func(element string) bool {
		_, exists := removedCols[element]
		return exists
	})
	return selectedFields
}