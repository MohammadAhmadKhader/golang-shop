package message

import (
	"time"

	"gorm.io/gorm"
	"main.go/pkg/models"
	"main.go/services/generic"
)

type Store struct {
	DB      *gorm.DB
	Generic *generic.GenericRepository[models.Message]
}

func NewStore(DB *gorm.DB) *Store {
	return &Store{
		DB:      DB,
		Generic: &generic.GenericRepository[models.Message]{DB: DB},
	}
}

// we must apply cursor
func (s *Store) GetById(userId uint, lastMessageId uint, cursor time.Time, limit int) ([]models.Message, error) {
	var messages []models.Message
	query := s.DB.Where("from = ?", userId)
	if !cursor.IsZero() {
		query = query.Where("created_at < ? OR (created_at < ?  AND id < ?)", cursor, cursor, lastMessageId)
	}
	
	err := query.Order("created_at DESC").Limit(limit).Find(messages).Error
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// we must apply cursor
func (s *Store) GetByUsersIds(from uint, to uint, lastMessageId uint, cursor time.Time, limit int) ([]models.Message, error) {
	var messages []models.Message
	query := s.DB.Where("(from = ? AND to = ?) OR (from = ? AND to = ?)", from, to, to, from)
	if !cursor.IsZero() {
		query = query.Where("created_at < ? OR (created_at < ?  AND id < ?)", cursor, cursor, lastMessageId)
	}
	err := query.Order("created_at DESC, id DESC").Limit(limit).Find(messages).Error
	if err != nil {
		return nil, err
	}

	return messages, nil
}
