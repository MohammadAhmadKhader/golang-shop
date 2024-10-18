package generic

import (
	"fmt"
	"sync"

	"gorm.io/gorm"
	"main.go/pkg/utils"
)

type GenericRepository[TModel any] struct {
	DB *gorm.DB
}

// the error message (notFoundMsg) expects only one param to be passed which is the Id
func (g GenericRepository[TModel]) GetOne(Id uint, notFoundMsg string) (TModel, error) {
	var model TModel
	if err := g.DB.Where("id = ?", Id).First(&model).Error; err != nil {
		return model, fmt.Errorf(notFoundMsg, Id)
	}

	return model, nil
}

func (g GenericRepository[TModel]) GetForUpdate(model *TModel, notFoundMsg string) (*TModel, error) {
	if err := g.DB.Model(model).First(model, 4).Error; err != nil {
		return model, fmt.Errorf(notFoundMsg, 4)
	}

	return model, nil
}


// this function meant to be used with any model that contains user id inside it, it will search based on both the resource id and user id.
// 
// * if the model does not contain user id it will throw an error.
func (g GenericRepository[TModel]) GetOneWithUserId(Id uint, userId uint,notFoundMsg string) (TModel, error) {
	var model TModel
	if err := g.DB.Where("id = ? AND user_id = ?", Id, userId).First(&model).Error; err != nil {
		return model, fmt.Errorf(notFoundMsg, Id)
	}

	return model, nil
}

func (g GenericRepository[TModel]) GetAll(page, limit int) ([]TModel, int64, []error) {
	var models []TModel
	var model TModel
	var count int64
	var wg sync.WaitGroup
	wg.Add(2)

	errors := make([]error, 0)

	// get models
	go func() {
		defer wg.Done()
		offset := utils.CalculateOffset(page, limit)
		if err := g.DB.Find(&models).Order("created_at DESC").Offset(offset).Limit(limit).Error; err != nil {
			errors = append(errors, err)
		}
	}()

	//get count
	go func() {
		defer wg.Done()
		if err := g.DB.Model(&model).Count(&count).Error; err != nil {
			errors = append(errors, err)
		}
	}()

	wg.Wait()

	return models, count, errors
}

func (g GenericRepository[TModel]) Create(model *TModel, selectedFields []string) (*TModel, error) {
	result := g.DB.Select(selectedFields).Create(model)
	if result.Error != nil {
		return nil, result.Error
	}

	return model, nil
}

func (prodStore GenericRepository[TModel]) CreateTx(model *TModel, tx *gorm.DB) error {
	err := tx.Create(model).Error
	if err != nil {
		return err
	}

	return nil
}

func (prodStore GenericRepository[TModel]) UpdateTx(model *TModel, tx *gorm.DB, columns []string) error {
	err := tx.Model(model).Select(columns).Updates(model).Error
	if err != nil {
		return err
	}

	return nil
}

func (g GenericRepository[TModel]) UpdateAndReturn(id uint, model *TModel, selectedFields []string) (*TModel, error) {
	result := g.DB.Model(model).Select(selectedFields).Where("id = ?", id).Updates(model)
	if result.Error != nil {
		return nil, result.Error
	}

	if err := g.DB.First(model, id).Error; err != nil {
		return nil, err
	}

	return model, nil
}

func (g GenericRepository[TModel]) Update(id uint, model *TModel, selectedFields []string) (*TModel, error) {
	result := g.DB.Model(model).Select(selectedFields).Where("id = ?", id).Updates(model)
	if result.Error != nil {
		return nil, result.Error
	}
	
	return model, nil
}

func (g GenericRepository[TModel]) UpdateWithUserId(id uint, model *TModel, selectedFields []string, userId uint) (*TModel, error) {
	result := g.DB.Model(model).Select(selectedFields).Where("id = ? AND user_id = ?", id, userId).Updates(model)
	if result.Error != nil {
		return nil, result.Error
	}

	if err := g.DB.First(model, id).Error; err != nil {
		return nil, err
	}

	return model, nil
}

func (g GenericRepository[TModel]) SoftDelete(Id uint, notFoundMsg string) error {
	var modelToDelete TModel
	if err := g.DB.Model(&modelToDelete).First(&modelToDelete, Id).Error; err != nil {
		return fmt.Errorf(notFoundMsg, Id)
	}

	if err := g.DB.Delete(&modelToDelete).Error; err != nil {
		return err
	}

	return nil
}

func (g GenericRepository[TModel]) SoftDeleteWithUserId(Id uint, userId uint ,notFoundMsg string) error {
	var modelToDelete TModel
	if err := g.DB.Model(&modelToDelete).First(&modelToDelete, Id).Where("user_id = ?", userId).Error; err != nil {
		return fmt.Errorf(notFoundMsg, Id)
	}

	if err := g.DB.Delete(&modelToDelete).Error; err != nil {
		return err
	}

	return nil
}

func (g GenericRepository[TModel]) FindThenHardDelete(Id uint, notFoundMsg string) error {
	var modelToDelete TModel
	if err := g.DB.Model(&modelToDelete).Unscoped().First(&modelToDelete, Id).Error; err != nil {
		return fmt.Errorf(notFoundMsg, Id)
	}

	if err := g.DB.Unscoped().Delete(&modelToDelete).Error; err != nil {
		return err
	}

	return nil
}

func (g GenericRepository[TModel]) HardDelete(Id uint) error {
	var modelToDelete TModel
	if err := g.DB.Unscoped().Where("id = ?", Id).Delete(&modelToDelete).Error; err != nil {
		return err
	}

	return nil
}

// * This store function applies only when soft delete is applied on the route
func (g *GenericRepository[TModel]) Restore(id uint, notFoundMsg string) (*TModel, error) {
	var item TModel
	if err := g.DB.Unscoped().First(&item, id).Where("DeletedAt != NULL").Error; err != nil {
		return nil, fmt.Errorf(notFoundMsg, id)
	}

	if err := g.DB.Model(&item).Update("DeletedAt", nil).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (g *GenericRepository[TModel]) RestoreWithUserId(id uint, userId uint, notFoundMsg string) (*TModel, error) {
	var item TModel
	if err := g.DB.Unscoped().First(&item, id).Where("user_id = ? & DeletedAt != NULL", userId).Error; err != nil {
		return nil, fmt.Errorf(notFoundMsg, id)
	}

	if err := g.DB.Model(&item).Update("DeletedAt", nil).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

// * This store function applies only when soft delete is applied on the route
func (g *GenericRepository[TModel]) GetAllDeleted(page, limit int) ([]TModel, int64, []error) {
	var models []TModel
	var model TModel
	var count int64
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(2)

	errors := make([]error, 0)

	// get models
	go func() {
		defer wg.Done()
		offset := utils.CalculateOffset(page, limit)
		if err := g.DB.Unscoped().Where("deleted_at is NOT NULL").Order("deleted_at DESC").
		Offset(offset).Limit(limit).Find(&models).Error; err != nil {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
		}
	}()

	//get count
	go func() {
		defer wg.Done()
		if err := g.DB.Unscoped().Where("deleted_at is NOT NULL").Model(&model).Count(&count).Error; err != nil {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
		}
	}()

	wg.Wait()

	return models, count, errors
}
