package role

import (

	"gorm.io/gorm"
	"main.go/constants"
	"main.go/pkg/models"
	"main.go/services/generic"
)

type Store struct {
	DB      *gorm.DB
	Generic *generic.GenericRepository[models.Role]
}

func NewStore(DB *gorm.DB) *Store {
	return &Store{
		DB:      DB,
		Generic: &generic.GenericRepository[models.Role]{DB: DB},
	}
}

func (roleStore *Store) GetAllRoles(page, limit int) ([]models.Role, int64, error) {
	roles,count, errs := roleStore.Generic.GetAll(page, limit)
	if len(errs) != 0 {
		return nil, 0, errs[0]
	}

	return roles, count, nil
}

func (roleStore *Store) CreateRole(role *models.Role) (*models.Role, error) {
	role, err := roleStore.Generic.Create(role, constants.RoleCols)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (roleStore *Store) UpdateRole(id uint, role *models.Role) (*models.Role, error) {
	notFoundMsg := "role with id: '%v' was not found"
	_, err := roleStore.Generic.GetOne(id, notFoundMsg); 
	if err != nil {
		return nil, err
	}

	updatedRole, err := roleStore.Generic.UpdateAndReturn(id, role, constants.RoleCols)
	if err != nil {
		return nil, err
	}

	return updatedRole, nil
}