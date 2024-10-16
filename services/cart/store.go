package cart

import (
	"fmt"

	"gorm.io/gorm"
	"main.go/pkg/models"
	"main.go/pkg/payloads"
	"main.go/services/generic"
)

type Store struct {
	DB      *gorm.DB
	Generic *generic.GenericRepository[models.CartItem]
}

func NewStore(DB *gorm.DB) *Store {
	return &Store{
		DB:      DB,
		Generic: &generic.GenericRepository[models.CartItem]{DB: DB},
	}
}

//// ! must be reworked
//func (cartStore *Store) GetCartByUserId(userId uint) (*models.Cart, error) {
//	var cart models.Cart
//	err := cartStore.DB.Where("userId = ?", userId).First(&cart).Error
//	if err != nil {
//		return nil, err
//	}
//
//	return &cart, err
//}

func (cartStore *Store) GetCartItemById(Id uint) (*models.CartItem, error) {
	notFoundMsg := "cart item with id: '%v' was not found"
	cartItem ,err := cartStore.Generic.GetOne(Id, notFoundMsg)
	if err != nil {
		return nil, err
	}

	return &cartItem, err
}

func (cartStore *Store) ChangeCartItemQty(oldQty uint, payload *payloads.ChangeCartItemQty, cartItem *models.CartItem) (*models.CartItem, error) {
	amount := int(payload.Amount)
	if payload.Operation == "-" {
		amount = amount * -1
	}
	product, err := cartStore.GetProduct(cartItem.ProductID)
	if err != nil {
		return nil, err
	}

	newQty := int(oldQty) + amount
	if  int(product.Quantity) < newQty {
		return nil, fmt.Errorf("product has quantity: '%v' which is more than you requested",product.Quantity)
	}
	if newQty <= 0 {
		return nil, fmt.Errorf("cart item cant be 0 or minus")
	}
	err = cartStore.DB.Model(cartItem).Select("Quantity").Update("Quantity", newQty).Error
	if err != nil {
		return nil, err
	}

	return cartItem, err
}

func (cartStore *Store) GetCartByUserId(userId uint) (*respCart, error) {
	var res = make([]getCartRow, 0)
	err := cartStore.DB.Table("cart_items").
		Select(selectQ).
		Joins(joinWProducts).
		Joins(joinWImages).
		Where("user_id = ?",userId).
		Scan(&res).Error
	if err != nil {
		return nil, err
	}
	
	cartResp := convertRowsToResponse(res)

	return cartResp, err
}

func (cartStore *Store) AddToCart(payload *payloads.AddCartItem, userId uint) (*models.CartItem, error) {
	product, err := cartStore.GetProduct(payload.ProductId)
	if err != nil {
		return nil, err
	}
	if product.Quantity < payload.Quantity {
		return nil, fmt.Errorf("product has quantity: '%v' which is less than what you requested", product.Quantity)
	}
	
	var cartItem = models.CartItem{
		ProductID: payload.ProductId,
		Quantity:  payload.Quantity,
		UserID: userId,
	}

	err = cartStore.DB.Create(&cartItem).Error

	if err != nil {
		return nil, err
	}

	return &cartItem, nil
}

func (cartStore *Store) DeleteCartItem(itemId uint) error {
	err := cartStore.DB.Delete(&models.CartItem{}, itemId).Error
	if err != nil {
		return err
	}

	return nil
}

func (cartStore *Store) ClearCart(userId uint) error {
	err := cartStore.DB.First(&models.CartItem{UserID: userId}).Error
	if err != nil {
		return fmt.Errorf("cart items for the user with id: '%v' were not found", userId)
	}

	err = cartStore.DB.Where("user_id = ?", userId).Delete(&models.CartItem{}).Error

	if err != nil {
		return err
	}

	return nil
}

func (cartStore *Store) GetProduct(productId uint) (*models.Product, error) {
	var product models.Product
	err := cartStore.DB.First(&product, productId).Error
	if err != nil {
		return nil, err
	}

	return &product, err
}

//func (cartStore *Store) GetCartItemCtx(r *http.Request) (*models.CartItem, error) {
//	cartItem, err := utils.GetResourceCtx[models.CartItem](r, "cart item")
//	if err != nil {
//		return nil, err
//	}
//
//	return cartItem, nil
//}
