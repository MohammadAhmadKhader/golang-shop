package product

import (
	"fmt"

	"gorm.io/gorm"
	"main.go/constants"
	"main.go/pkg/models"
	"main.go/pkg/utils"
	"main.go/services/generic"
	"main.go/services/image"
	"main.go/types"
)

type Store struct {
	DB      *gorm.DB
	Generic *generic.GenericRepository[models.Product]
}

func NewStore(DB *gorm.DB) *Store {
	return &Store{
		DB:      DB,
		Generic: &generic.GenericRepository[models.Product]{DB: DB},
	}
}

var (
	notFoundMsg = "product with id: '%v' was not found"
)

func (prodStore *Store) GetProductById(Id uint) ([]types.RowGetProductById, error) {
	var qRows []types.RowGetProductById
	err := prodStore.DB.Model(&models.Product{}).Select(getProductByIdQ, Id).
	Joins(getProductByIdJoins).Where("products.id", Id).Group(groupByGetProductById).Scan(&qRows).Error
	
	if err != nil {
		return nil, fmt.Errorf(notFoundMsg, Id)
	}

	return qRows, err
}

func (prodStore *Store) GetAllProducts(page, limit int, filter func(db *gorm.DB, filters []types.FilterCondition) ([]models.Product, error)) ([]models.Product, int64, error) {
	products, count, errs := prodStore.Generic.GetAll(page, limit)
	if len(errs) != 0 {
		return nil, 0, errs[0]
	}

	return products, count, nil
}

func (prodStore *Store) CreateProduct(product *models.Product) (*models.Product, error) {
	products, errs := prodStore.Generic.Create(product, constants.ProductCols)
	if errs != nil {
		return nil, errs
	}

	return products, nil
}

func (prodStore *Store) UpdateProduct(id uint, changes *models.Product, excluder types.Excluder) (*models.Product, error) {
	prodColsCopy := utils.CopyCols(constants.ProductCols)
	fields := excluder.Exclude(prodColsCopy)
	product, errs := prodStore.Generic.FindThenUpdate(id, changes, fields, notFoundMsg)
	if errs != nil {
		return nil, errs
	}

	return product, nil
}

func (prodStore *Store) CreateImageTx(tx *gorm.DB, uploadResp *types.UploadResponse, productId uint, isMain bool) (*models.Image, error) {
	imageStore := image.NewStore(tx)
	newImage := &models.Image{
		ProductID:     productId,
		IsMain:        &isMain,
		ImageUrl:      uploadResp.URL,
		ImagePublicId: uploadResp.PublicID,
	}

	err := imageStore.Generic.CreateTx(newImage, tx)
	if err != nil {
		return nil, err
	}

	return newImage, nil
}

func (prodStore *Store) CreateProductWithImage(product *models.Product, uploadResp *types.UploadResponse) (*models.Product, error) {
	var returnedProduct models.Product
	var returnedImg models.Image

	err := prodStore.DB.Transaction(func(tx *gorm.DB) error {
		err := prodStore.Generic.CreateTx(product, tx)
		if err != nil {
			return err
		}

		newImg, err := prodStore.CreateImageTx(tx, uploadResp, product.ID, true)
		if err != nil {
			return err
		}

		returnedImg = *newImg
		returnedProduct = *product

		return nil
	})

	if err != nil {
		return nil, err
	}
	returnedProduct.Image = &returnedImg

	return &returnedProduct, nil
}