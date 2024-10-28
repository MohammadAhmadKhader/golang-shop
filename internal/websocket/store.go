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
	err := utils.ValidateStruct(message)
	if err != nil {
		return err
	}

	err = s.DB.Create(&message).Error
	if err != nil {
		return err
	}

	GlobalManager.BroadcastCUMessage(message, []uint{message.From,message.To}, MessageCreated)

	return nil
}

func (s *Store) UpdateMessage(id uint, changes models.Message, excluder types.Excluder) (*models.Message, error){
	var message models.Message
	err := s.DB.First(&message,id).Error
	if err != nil {
		return nil, err
	}

	messageColsCopy := utils.CopyCols(selectedFields)
	fields := excluder.Exclude(messageColsCopy)
	err = s.DB.Model(&message).Select(fields).Updates(changes).Error
	if err != nil {
		return nil, err
	}

	GlobalManager.BroadcastCUMessage(message, []uint{message.From,message.To}, MessageUpdated)

	return &message, nil
}

func (s *Store) DeleteMessage(id uint) error {
	var msg models.Message
	err := s.DB.First(&msg, id).Error
	if err != nil {
		return fmt.Errorf("message with id: '%v' was not found", id)
	}
	
	err = s.DB.Delete(msg, id).Error
	if err != nil {
		return err
	}

	GlobalManager.BroadcastDMessage(DeleteMessagePayload{Id: msg.ID}, []uint{msg.From,msg.To})

	return nil
}

func (s *Store) UpdateMessageStatus(id uint, status string) (error){
	var message models.Message
	err := s.DB.First(&message,id).Error
	if err != nil {
		return err
	}

	res := s.DB.Model(&message).Update("status",status)
	if res.Error != nil {
		return err
	}

	GlobalManager.BroadcastCUMessage(message, []uint{message.From,message.To}, MessageStatusUpdated)

	return nil
}