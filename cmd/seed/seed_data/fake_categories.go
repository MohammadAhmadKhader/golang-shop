package seed_data

import (

	"github.com/brianvoe/gofakeit/v6"
	"main.go/internal/database"
	"main.go/pkg/models"
)

func SeedCategories() {
	for i := 0; i < 17; i++ {
		category := models.Category{
			Name: gofakeit.ProductCategory(),
		}

		DB := database.InitDB()
		err := DB.Create(&category).Error
		if err != nil {
			continue
		}
	}

}