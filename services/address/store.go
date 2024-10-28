package address

import (
	"gorm.io/gorm"
	"main.go/constants"
	"main.go/pkg/models"
	"main.go/pkg/utils"
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

var (
	notFoundMsg = "address with id: '%v' was not found"
)

func (addressStore *Store) GetById(id uint, userId uint) (*models.Address, error){
	address, err := addressStore.Generic.GetOneWithUserId(id, userId, notFoundMsg)
	if err != nil {
		return nil, err
	}

	return &address, nil
}

func (addressStore *Store) GetAddressById(id uint) (*models.Address, error){
	address, err := addressStore.Generic.GetOne(id, notFoundMsg)
	if err != nil {
		return nil, err
	}

	return &address, nil
}

func (addressStore *Store) CreateAddress(address *models.Address) (*models.Address, error){
	address, err := addressStore.Generic.Create(address, constants.AddressCreateCols)
	if err != nil {
		return nil, err
	}
	
	return address, nil
}

func (addressStore *Store) UpdateAddress(id uint, address *models.Address, excluder types.Excluder) (*models.Address, error){
	addressFieldsCopy := utils.CopyCols(constants.AddressUpdateCols)
	fields := excluder.Exclude(addressFieldsCopy)
	address, err := addressStore.Generic.UpdateAndReturn(id, address, fields)
	if err != nil {
		return nil, err
	}
	
	return address, nil
}

func (addressStore *Store) DeleteAddress(id uint, userId uint) (error){
	err := addressStore.Generic.SoftDelete(id, notFoundMsg)
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