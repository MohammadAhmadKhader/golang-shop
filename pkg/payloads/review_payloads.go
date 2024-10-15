package payloads

import (
	"slices"
	"strings"

	"main.go/pkg/models"
)

type CreateReview struct {
	Rate    uint8  `json:"rate" validate:"required,oneof=1 2 3 4 5"`
	Comment string `json:"comment" validate:"required,min=2,max=256,alphanumWithSpaces"`
	UserId  uint   `json:"userId" validate:"required,min=1"`
}

type UpdateReview struct {
	Rate    uint8  `json:"rate" validate:"oneof=1 2 3 4 5"`
	Comment string `json:"comment" validate:"required,min=2,max=256,alphanumWithSpaces"`
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

func (ur *UpdateReview) ToModel() *models.Review {
	if ur != nil {
		return &models.Review{
			Rate:    ur.Rate,
			Comment: ur.Comment,
		}
	}

	return nil
}

func (cr *CreateReview) ToModel() *models.Review {
	if cr != nil {
		return &models.Review{
			UserID: cr.UserId,
			Rate:    cr.Rate,
			Comment: cr.Comment,
		}
	}
	
	return nil
}

func (uc *UpdateReview) Exclude(selectedFields []string) []string {
	removedCols := map[string]any{}
	if uc.Comment == "" {
		removedCols["Comment"] = 1
	}
	if uc.Rate == 0 {
		removedCols["Rate"] = 1
	}
	
	selectedFields = slices.DeleteFunc(selectedFields, func(element string) bool {
		_, exists := removedCols[element]
		return exists
	})
	return selectedFields
}