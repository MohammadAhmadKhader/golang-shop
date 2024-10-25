package test_utils

import "main.go/pkg/models"

type CategoryCreate struct {
	Category models.Category `json:"category"`
}

type CategoryGetAll struct {
	Pagination
	Categories []models.Category `json:"categories"`
}

type CategoryUpdate struct {
	Category models.Category `json:"category"`
}

type CategoryGetOne struct {
	Category models.Category `json:"category"`
}

type CategoryRestore struct {
	Message string         `json:"message"`
	Category models.Category `json:"category"`
}
