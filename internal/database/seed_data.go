package database

import (
	"gorm.io/gorm"
	"main.go/pkg/models"
)

func SeedData(db *gorm.DB) error {
	roles := []models.Role{
		{Name: "Admin"},
		{Name: "RegularUser"},
		{Name: "SuperAdmin"},
	}

	// Roles seed data
	for _, role := range roles {
		var existingRole models.Role
		if err := db.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&role).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}