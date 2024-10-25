package test_utils

import "main.go/pkg/models"

type ProductCreate struct {
	Product models.Product `json:"product"`
}

type ProductGetAll struct {
	Pagination
	Products []models.Product `json:"products"`
}

type ProductUpdate struct {
	Product models.Product `json:"product"`
}

type ProductGetOne struct {
	Product models.Product `json:"product"`
}

type ProductRestore struct {
	Message string `json:"message"`
	Product models.Product `json:"product"`
}




