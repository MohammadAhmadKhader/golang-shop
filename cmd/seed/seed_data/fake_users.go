package seed_data

import (
	"fmt"
	"log"

	"github.com/brianvoe/gofakeit/v6"
	"main.go/internal/database"
	"main.go/pkg/models"
	"main.go/services/auth"
)

func SeedUsers() {
	for i := 0; i < 100; i++ {
		pw, _ := auth.HashPassword(gofakeit.Password(true, true, true, true, true, 6))
		user := models.User{
			Name:     gofakeit.Username(),
			Email:    gofakeit.Email(),
			Password: pw,
		}

		DB := database.InitDB()
		err := DB.Create(&user).Error
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Username: %v @@@@ email: %v", user.Name, user.Email)
	}
}
