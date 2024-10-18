package user

import (
	"fmt"

	"gorm.io/gorm"
	"main.go/constants"
	"main.go/pkg/models"
	"main.go/pkg/payloads"
	"main.go/types"

	"main.go/services/auth"
	"main.go/services/generic"
)

type Store struct {
	DB      *gorm.DB
	Generic *generic.GenericRepository[models.User]
}

func NewStore(DB *gorm.DB) *Store {
	return &Store{
		DB:      DB,
		Generic: &generic.GenericRepository[models.User]{DB: DB},
	}
}

func (userStore *Store) GetUserById(Id uint) (*models.User, error) {
	notFoundMsg := "user with id: '%v' is not found"
	user, err := userStore.Generic.GetOne(Id, notFoundMsg)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (userStore *Store) GetUserWithRolesById(Id uint) (*models.User, error) {
	var userWithRoles models.User
	err := userStore.DB.Joins("Role").First(&userWithRoles, Id).Error
	if err != nil {
		return nil, err
	}

	return &userWithRoles, nil
}

func (userStore *Store) GetUserByEmail(email string) (*models.User, error) {
	notFoundErr := fmt.Errorf("user with email: '%v' is not found", email)
	var user models.User
	if err := userStore.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, notFoundErr
	}

	return &user, nil
}

func (userStore *Store) CreateUser(user payloads.UserSignUp) (*models.User, error) {
	user.TrimStrs()

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	var newUser *models.User
	err = userStore.DB.Transaction(func(tx *gorm.DB) error {
		newUser = &models.User{
			Name:     user.Name,
			Email:    user.Email,
			Password: hashedPassword,
		}
		err := tx.Create(newUser).Error
		if err != nil {
			return err
		}

		regularUser := string(types.RegularUser)
		assignedRole := models.Role{Name: regularUser}
		err = tx.Where("name = ?", regularUser).First(&assignedRole).Error
		if err != nil {
			return err
		}

		assignRole := models.UserRoles{
			UserID: newUser.ID,
			RoleID: assignedRole.ID,
		}

		err = tx.Create(&assignRole).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (userStore *Store) UpdatePassword(newHashedPassword, email string) error {
	if err := userStore.DB.Model(&models.User{}).Where("email = ?", email).Update("Password", newHashedPassword).Error; err != nil {
		return err
	}

	return nil
}

func (userStore *Store) UpdateProfile(id uint, user *models.User, excluder types.Excluder) (*models.User, error) {
	colsToUpdate := excluder.Exclude(constants.UserUpdateCols)
	updatedUser, err := userStore.Generic.UpdateAndReturn(id, user, colsToUpdate)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (userStore *Store) RemoveUserRole(roleId, userId uint) (error) {
	var userRole = models.UserRoles{
		RoleID: roleId,
		UserID: userId,
	}
	err := userStore.DB.Model(&userRole).First(&userRole).Error
	if err != nil {
		return fmt.Errorf("user role with keys (roleId-userId) '%v'-'%v' was not found", roleId, userId)
	}

	err = userStore.DB.Model(&userRole).Unscoped().Delete(&userRole).Error
	if err != nil {
		return err
	}

	return nil
}

func (userStore *Store) AssignUserRole(roleId, userId uint) (*models.UserRoles, error) {
	var userRole = models.UserRoles{
		RoleID: roleId,
		UserID: userId,
	}

	user, err := userStore.GetUserById(userId)
	if err != nil{
		return nil, fmt.Errorf("user with id:'%v' was not found", userId)
	}

	var role models.Role
	err = userStore.DB.First(&role, roleId).Error
	if err != nil{
		return nil, fmt.Errorf("role with id:'%v' was not found", roleId)
	}

	err = userStore.DB.Model(&userRole).Create(&userRole).Error
	if err != nil {
		return nil, err
	}
	userRole.Role = &role
	userRole.User = user

	return &userRole, nil
}