package message

import (
	"sync"
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

func (s *Store) GetById(userId uint, lastMessageId uint, cursor time.Time, limit int) ([]*models.Message, error) {
	var messages []*models.Message
	var user models.User
	errChan := make(chan error, 2)

	var wg sync.WaitGroup
	wg.Add(1)
	go func(){
		defer wg.Done()

		err := s.DB.First(&user,userId).Error
		if err != nil {
			errChan <- err
		}
	}()

	wg.Add(1)
	go func(){
		defer wg.Done()

		query := s.DB.Where("`from` = ?", userId)
		if !cursor.IsZero() {
			query = query.Where("created_at < ? OR (created_at < ?  AND id < ?)", cursor, cursor, lastMessageId)
		}

		err := query.Order("created_at DESC, id DESC").Limit(limit).Find(&messages).Error
		if err != nil {
			errChan <- err
		}
	}()

	go func(){
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}
	
	for _, message := range messages {
		if message.From == userId {
			message.FromUser = &user
		}
	}

	return messages, nil
}

func (s *Store) GetByUsersIds(userId uint, to uint, lastMessageId uint, cursor time.Time, limit int) ([]*models.Message, error) {
	var messages []*models.Message
	var user models.User
	errChan := make(chan error, 2)
	
	var wg sync.WaitGroup
	wg.Add(1)
	go func(){
		defer wg.Done()
		err := s.DB.First(&user,userId).Error
		if err != nil {
			errChan <- err
		}
	}()

	wg.Add(1)
	go func(){
		defer wg.Done()

		query := s.DB.Where("(`from` = ? AND `to` = ?) OR (`from` = ? AND `to` = ?)", userId, to, to, userId)
		if !cursor.IsZero() {
			query = query.Where("created_at < ? OR (created_at < ?  AND id < ?)", cursor, cursor, lastMessageId)
		}

		err := query.Order("created_at DESC, id DESC").Limit(limit).Find(&messages).Error
		if err != nil {
			errChan <- err
		}
	}()

	go func(){
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}
	
	for _, message := range messages {
		if message.From == userId {
			message.FromUser = &user
		}
		if message.To == userId {
			message.ToUser = &user
		}
	}

	return messages, nil
}
