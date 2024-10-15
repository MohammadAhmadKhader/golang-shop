package payloads

import "slices"

type CreateImage struct {
	ProductID     uint   `json:"productId" validate:"required"`
	ImageUrl      string `json:"imageUrl" validate:"required"`
	IsMain        *bool   `json:"isMain" validate:"required"`
	ImagePublicId string `json:"imagePublicId" validate:"required"`
}

type UpdateImage struct {
	ProductID     uint   `json:"productId"`
	ImageUrl      string `json:"imageUrl"`
	IsMain        *bool   `json:"isMain" validate:"required"`
	ImagePublicId string `json:"imagePublicId"`
}

func (ui *UpdateImage) Exclude(selectedFields []string) []string {
	removedCols := map[string]any{}
	if ui.ImageUrl == "" {
		removedCols["ImageUrl"] = 1
	}
	if ui.ProductID == 0 {
		removedCols["ProductID"] = 1
	}
	if ui.ImagePublicId == "" {
		removedCols["ImagePublicId"] = 1
	}
	if ui.IsMain == nil {
		removedCols["IsMain"] = 1
	}

	selectedFields = slices.DeleteFunc(selectedFields, func(element string) bool {
		_, exists := removedCols[element]
		return exists
	})
	return selectedFields
}