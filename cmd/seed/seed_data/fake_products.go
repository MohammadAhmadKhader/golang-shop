package seed_data

import (
	"log"

	"github.com/brianvoe/gofakeit/v6"
	"main.go/internal/database"
	"main.go/pkg/models"
)

func SeedProducts() {
		for i := 0; i < 100 ; i++ {
			images := []models.Image{}
			prodDesc := gofakeit.ProductDescription()
			imagesNumber := gofakeit.UintRange(1,2)

			isMain := true
			for j := 0;j < int(imagesNumber); j++ {
				images = append(images, models.Image{
					IsMain:&isMain,
					ImageUrl: gofakeit.ImageURL(600,600),
					ImagePublicId: gofakeit.UUID(),
				})
				isMain = false
			}
			
			product := models.Product{
				Name: gofakeit.ProductName(),
				Description: &prodDesc,
				Quantity: gofakeit.UintRange(10,300),
				CategoryID: gofakeit.UintRange(1,17),
				Price: gofakeit.Float64Range(5,350),
				Images: images,	
			}

			DB := database.InitDB()
			err := DB.Create(&product).Error
			if err != nil {
				log.Fatal(err)
			}
		}

}