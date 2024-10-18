package image

import (
	"context"
	"fmt"
	"sync"

	"gorm.io/gorm"
	"main.go/constants"
	"main.go/pkg/models"
	"main.go/pkg/utils"
	"main.go/services/generic"
	"main.go/types"
)

type Store struct {
	DB      *gorm.DB
	Generic *generic.GenericRepository[models.Image]
}

func NewStore(DB *gorm.DB) *Store {
	return &Store{
		DB:      DB,
		Generic: &generic.GenericRepository[models.Image]{DB: DB},
	}
}

var (
	notFoundMsg = "image with id: '%v' was not found"
)

func (imageStore *Store) GetImageById(id uint) (*models.Image, error) {
	image, err := imageStore.Generic.GetOne(id, notFoundMsg)
	if err != nil {
		return nil, err
	}

	return &image, nil
}

func (imageStore *Store) CreateImage(image *models.Image) (*models.Image, error) {
	image, err := imageStore.Generic.Create(image, constants.ImageCols)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (imageStore *Store) UpdateImageUrl(id uint, newImageUrl string) (error) {
	err := imageStore.DB.Model(&models.Image{}).Where("id = ?", id).Update("ImageUrl", newImageUrl).Error
    if err != nil {
        return err
    }

	return nil
}

func (imageStore *Store) GetCountOfProductImages(productId uint) (*int64, error) {
	var count int64
	err := imageStore.DB.Model(&models.Image{}).Where("product_id = ?", productId).Count(&count).Error
	if err != nil {
		return nil, err
	}

	return &count, nil
}

func (imageStore *Store) GetProductById(productId uint) (*models.Product, error) {
	var product models.Product
	err := imageStore.DB.Model(&models.Product{}).First(&product, productId).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (imageStore *Store) DeleteImageById(id uint) error {
	image, err := imageStore.Generic.GetOne(id, notFoundMsg)
	if err != nil {
		return err
	}
	if *image.IsMain == true {
		return fmt.Errorf("you can not delete a main product image, set another product image as main then try again")
	}
	imgHandler := utils.NewImagesHandler()

	err = imgHandler.DeleteOne(image.ImagePublicId, context.Background())
	if err != nil {
		return err
	}

	err = imageStore.DB.Unscoped().Delete(image).Error
	if err != nil {
		return err
	}

	return nil
}

func (imageStore *Store) SetImageAsMainTx(tx *gorm.DB, id, productId uint) error {
	isTrue := true
	result := tx.Where("id = ? AND product_id = ?", id, productId).Select("IsMain").Updates(&models.Image{IsMain: &isTrue})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("something went wrong, no images were updated")
	}

	return nil
}

func (imageStore *Store) SetImageAsNotMainTx(tx *gorm.DB, productId uint) error {
	isTrue := false
	err := tx.Where("is_main = 1 AND product_id = ?", productId).
	Select("IsMain").Updates(&models.Image{IsMain: &isTrue}).Error
	if err != nil {
		return err
	}

	return nil
}

func (imageStore *Store) SwapMainStatus(id, productId uint) error {
	return imageStore.DB.Transaction(func(tx *gorm.DB) error {
		var wg sync.WaitGroup
		errChan := make(chan error, 2)

		wg.Add(1)
		go func(){
			defer wg.Done()
			err := imageStore.SetImageAsMainTx(tx, id, productId)
			if err != nil {
				errChan <- err
			}
		}()
		
		wg.Add(1)
		go func(){
			defer wg.Done()
			err := imageStore.SetImageAsNotMainTx(tx, productId)
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
				return err
			}
		}

		return nil
	})
}

func (imageStore *Store) CreateManyImages(uploadResults []*types.UploadResponse, productId *uint) ([]models.Image, error) {
	var images = make([]models.Image, 0, len(uploadResults))
	isMain := false
	for _, upResult := range uploadResults {
		newImg := models.Image{
			ProductID: *productId,
			ImageUrl: upResult.URL,
			IsMain: &isMain,
			ImagePublicId: upResult.PublicID,
		}
		
		images = append(images, newImg)
	}

	err := imageStore.DB.Create(images).Error
	if err != nil {
		return nil, err
	}

	return images, nil
}