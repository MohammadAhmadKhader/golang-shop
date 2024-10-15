package category

import (

	"gorm.io/gorm"
	"main.go/constants"
	"main.go/pkg/models"
	"main.go/services/generic"
)

type Store struct {
	DB *gorm.DB
	Generic *generic.GenericRepository[models.Category]
}

func NewStore(DB *gorm.DB) *Store {
	return &Store{
		DB: DB,
		Generic: &generic.GenericRepository[models.Category]{DB: DB},
	}
}

func (cateStore *Store) GetCategoryById(Id uint) (*models.Category, error) {
	notFoundMsg := "category with id: '%v' was not found"
	category, err := cateStore.Generic.GetOne(Id, notFoundMsg)
	if err != nil {
		return nil, err
	}

	return &category, err
}

func (cateStore *Store) GetAllCategories(page, limit int) ([]models.Category, int64, error) {
	categories,count, errs := cateStore.Generic.GetAll(page, limit)
	if len(errs) != 0 {
		return nil, 0, errs[0]
	}

	return categories, count, nil
}

func (cateStore *Store) CreateCategory(category *models.Category) (*models.Category, error) {
	category, err := cateStore.Generic.Create(category, constants.CategoryCols)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (cateStore *Store) UpdateCategory(id uint, category *models.Category) (*models.Category, error) {
	notFoundMsg := "category with id: '%v' was not found"
	_, err := cateStore.Generic.GetOne(id, notFoundMsg); 
	if err != nil {
		return nil, err
	}
	
	uCategory, err := cateStore.Generic.Update(id, category, constants.CategoryCols)
	if err != nil {
		return nil, err
	}

	return uCategory, nil
}