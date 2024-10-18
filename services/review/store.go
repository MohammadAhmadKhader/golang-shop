package review

import (

	"gorm.io/gorm"
	"main.go/constants"
	"main.go/pkg/models"
	"main.go/services/generic"
	"main.go/types"
)

var (
	notFoundMsg = "review with id: '%v' is not found"
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

func (reviewStore *Store) UpdateReview(id,userId uint, updatePayload *models.Review,excluder types.Excluder) (*models.Review, error) {
	uCols := excluder.Exclude(constants.CommentUpdateCols)
	review, err := reviewStore.Generic.FindThenUpdateWithAuth(id, updatePayload, uCols,notFoundMsg ,userId)
	if err != nil {
		return nil, err
	}

	return review, nil
}


func (reviewStore *Store) CreateReview(createPayload *models.Review) (*models.Review, error) {
	review, err := reviewStore.Generic.Create(createPayload, constants.CommentCreateCols)
	if err != nil {
		return nil, err
	}

	return review, nil
}

//func (reviewStore *Store) HardDelete(Id uint) error {
//	err := reviewStore.Generic.HardDelete(Id, notFoundMsg)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func (reviewStore *Store) HardDelete(Id uint, userId uint) error {
	_,err := reviewStore.Generic.FindThenDeleteWithAuth(Id, notFoundMsg, userId)
	if err != nil {
		return err
	}

	return nil
}

func (reviewStore *Store) GetProductById(productId uint) (*models.Product, error) {
	var product models.Product
	err := reviewStore.DB.First(&product, productId).Error
	if err != nil {
		return nil, err
	}

	return &product, err
}