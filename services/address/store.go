package address

import (
	"fmt"

	"gorm.io/gorm"
	"main.go/constants"
	"main.go/pkg/models"
	"main.go/services/generic"
	"main.go/types"
)

type Store struct {
	DB      *gorm.DB
	Generic *generic.GenericRepository[models.Address]
}

func NewStore(DB *gorm.DB) *Store {
	return &Store{
		DB:      DB,
		Generic: &generic.GenericRepository[models.Address]{DB: DB},
	}
}

func (addressStore *Store) GetById(id uint, userId uint) (*models.Address, error){
	notFoundErr := fmt.Errorf("address with id: '%v' was not found", id)
	address, err := addressStore.Generic.GetOneWithUserId(id, userId, notFoundErr)
	if err != nil {
		return nil, err
	}

	return &address, nil
}

func (addressStore *Store) CreateAddress(address *models.Address) (*models.Address, error){
	address, err := addressStore.Generic.Create(address, constants.AddressCols)
	if err != nil {
		return nil, err
	}
	
	return address, nil
}

func (addressStore *Store) UpdateAddress(id uint, userId uint ,address *models.Address, excluder types.Excluder) (*models.Address, error){
	fields := excluder.Exclude(constants.AddressCols)
	address, err := addressStore.Generic.Update(id, address, fields)
	if err != nil {
		return nil, err
	}
	
	return address, nil
}

func (addressStore *Store) DeleteAddress(id uint, userId uint) (error){
	notFoundErr := fmt.Errorf("address with id: '%v' was not found", id)
	err := addressStore.Generic.SoftDelete(id, notFoundErr)
	if err != nil {
		return err
	}
	
	return nil
}

func (addressStore *Store) GetAllAddresses(userId uint) ([]models.Address, error){
	var addresses []models.Address
	err := addressStore.DB.Model(&addresses).Where("user_id = ?", userId).Find(&addresses).Error
	if err != nil {
		return nil, err
	}

	return addresses,  err
}

func (addressStore *Store) GetUndeletedAddressesCount(userId uint) (*int64, error){
	var count int64
	err := addressStore.DB.Model(&models.Address{}).Where("user_id = ?", userId).Count(&count).Error
	if err != nil {
		return nil, err
	}

	return &count,  err
}