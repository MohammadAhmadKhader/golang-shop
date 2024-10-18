package middlewares

import (

	"main.go/internal/database"
	"main.go/pkg/models"
)

type UserLookup struct {}

func (u *UserLookup) GetUserById(Id uint) (*models.User, error) {
	var user models.User
	err := database.DB.Where("id = ?", Id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserLookup) GetUserRolesByUserId(Id uint) ([]models.UserRoles, error) {
	var roles []models.UserRoles
	err := database.DB.Where("user_id = ?", Id).Preload("Role").Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func NewUserLookup() *UserLookup {
	return &UserLookup{}
}