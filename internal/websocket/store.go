package websocket

import (
	"fmt"

	"gorm.io/gorm"
	"main.go/pkg/models"
	"main.go/pkg/utils"
	"main.go/types"
)

type Store struct {
	DB *gorm.DB
}

func NewStore(DB *gorm.DB) *Store {
	return &Store{
		DB: DB,
	}
}

var selectedFields = []string{}


func (s *Store) CreateMessage(message models.Message) error {
	validatedMsg, err := utils.ValidateStruct(message)
	if err != nil {
		return err
	}

	err = s.DB.Create(validatedMsg).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) UpdateMessage(id uint, changes models.Message, exculuder types.Excluder) (*models.Message, error){
	var message models.Message
	err := s.DB.First(message,id).Error
	if err != nil {
		return nil, err
	}

	fields := exculuder.Exclude(selectedFields)
	err = s.DB.Model(&message).Select(fields).Updates(changes).Error
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (s *Store) DeleteMessage(id uint) error {
	var msg models.Message
	res := s.DB.Delete(msg, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("message with id: '%v' was not found", id)
	}

	return nil
}
