package seed_data

import (
	"log"

	"github.com/brianvoe/gofakeit/v6"
	"main.go/internal/database"
	"main.go/pkg/models"
)

func SeedReviews() {
	
	for i := 0; i < 100 ; i++ {
		review := models.Review{
			ProductID: gofakeit.UintRange(1, 107),
			UserID: gofakeit.UintRange(30, 107),
			Comment: gofakeit.Comment(),
			Rate: uint8(gofakeit.UintRange(1, 5)),
		}

		DB := database.InitDB()
		err := DB.Create(&review).Error
		if err != nil {
			log.Fatal(err)
		}
	}
	
}