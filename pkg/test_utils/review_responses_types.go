package test_utils

import (

	"main.go/pkg/models"
)

type ReviewCreate struct {
	Review models.Review `json:"review"`
}

type ReviewGetAll struct {
	Pagination
	Reviews []models.Review `json:"reviews"`
}

type ReviewUpdate struct {
	Review models.Review `json:"review"`
}

type ReviewGetOne struct {
	Review models.Review `json:"review"`
}
