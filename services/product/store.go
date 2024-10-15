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

func (prodStore *Store) GetProductById(Id uint) (*models.Product, error) {
	var product models.Product
	rows, err := prodStore.DB.Raw(qGetProductById, Id, Id).Rows()
	if err != nil {
		return nil, fmt.Errorf("product with id: %v was not found", Id)
	}

	Ids := map[string]any{}
	formater := func(Id uint, model string) string {
		return fmt.Sprintf("%v-%v", Id, model)
	}

	for rows.Next() {
		var image models.Image
		var review models.Review
		err := rows.Scan(scanGetProductById(&product, &image, &review))
		if err != nil {
			return nil, fmt.Errorf("something went wrong try again later")
		}

		_, existRev := Ids[formater(review.ID, "review")]
		if !existRev && review.ID != 0 {
			Ids[formater(review.ID, "review")] = formater(review.ID, "review")
			product.Reviews = append(product.Reviews, review)
		}
		_, existImg := Ids[formater(image.ID, "image")]
		if !existImg && image.ID != 0 {
			Ids[formater(image.ID, "image")] = formater(image.ID, "image")
			product.Images = append(product.Images, image)
		}

	}
	defer rows.Close()

	return &product, err
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

func (prodStore *Store) UpdateProduct(id uint, product *models.Product, excluder types.Excluder) (*models.Product, error) {
	notFoundMsg := "product with id: '%v' was not found"
	_, err := prodStore.Generic.GetOne(id, notFoundMsg)
	if err != nil {
		return nil, err
	}

	fields := excluder.Exclude(constants.ProductCols)
	products, errs := prodStore.Generic.Update(id, product, fields)
	if errs != nil {
		return nil, errs
	}

	return products, nil
}

func (prodStore *Store) CreateImageTx(tx *gorm.DB, uploadResp *utils.UploadResponse, productId uint, isMain bool) (*models.Image, error) {
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

func (prodStore *Store) CreateProductWithImage(product *models.Product, uploadResp *utils.UploadResponse) (*models.Product, error) {
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