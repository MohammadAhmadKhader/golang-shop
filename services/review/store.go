package review

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
	Generic *generic.GenericRepository[models.Review]
}

func NewStore(DB *gorm.DB) *Store {
	return &Store{
		DB:      DB,
		Generic: &generic.GenericRepository[models.Review]{DB: DB},
	}
}

func (reviewStore *Store) GetReviewById(Id uint) (*models.Review, error) {
	notFoundMsg := "review with id: '%v' is not found"
	review, err := reviewStore.Generic.GetOne(Id, notFoundMsg)
	if err != nil {
		return nil, err
	}

	return &review, err
}

func (reviewStore *Store) GetAllReviews(page, limit int) ([]models.Review, int64, error) {
	reviews,count, errs := reviewStore.Generic.GetAll(page, limit)
	if len(errs) != 0 {
		return nil, 0, errs[0]
	}

	return reviews, count, nil
}

func (reviewStore *Store) UpdateReview(id uint, updatePayload *models.Review, excluder types.Excluder) (*models.Review, error) {
	notFoundMsg :="review with id: '%v' was not found"
	_, err := reviewStore.Generic.GetOne(id, notFoundMsg); 
	if err != nil {
		return nil, err
	}

	uCols := excluder.Exclude(constants.CommentCols)
	review, err := reviewStore.Generic.Update(id, updatePayload, uCols)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (reviewStore *Store) CreateReview(createPayload *models.Review) (*models.Review, error) {
	review, err := reviewStore.Generic.Create(createPayload, constants.CommentCols)
	if err != nil {
		return nil, err
	}

	return review, nil
}

//func (reviewStore *Store) SoftDeleteReview(Id uint) error {
//	notFoundErr := fmt.Errorf("review with id: '%v' is not found", Id)
//	err := reviewStore.Generic.SoftDelete(Id, notFoundErr)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func (reviewStore *Store) HardDeleteReview(Id uint, userId uint) error {
	notFoundErr := fmt.Errorf("review with id: '%v' is not found", Id)
	err := reviewStore.Generic.HardDeleteWithUserId(Id, userId, notFoundErr)
	if err != nil {
		return err
	}

	return nil
}