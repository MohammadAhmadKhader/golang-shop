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
	Generic *generic.GenericRepository[models.Cart]
}

func NewStore(DB *gorm.DB) *Store {
	return &Store{
		DB:      DB,
		Generic: &generic.GenericRepository[models.Cart]{DB: DB},
	}
}

// ! must be reworked
func (cartStore *Store) GetCartByUserId(userId uint) (*models.Cart, error) {
	var cart models.Cart
	err := cartStore.DB.Where("userId = ?", userId).First(&cart).Error
	if err != nil {
		return nil, err
	}

	return &cart, err
}

func (cartStore *Store) GetCartById(Id uint) (*models.Cart, error) {
	notFoundMsg := "cart with id: '%v' was not found"
	cart ,err := cartStore.Generic.GetOne(Id, notFoundMsg)
	if err != nil {
		return nil, err
	}

	return &cart, err
}

func (cartStore *Store) ChangeCartItemQty(cartItemId *uint, payload *payloads.ChangeCartItemQty) (*models.CartItem, error) {
	amount := int(payload.Amount)
	if payload.Operation == "-" {
		amount = amount * -1
	}
	var cartItem models.CartItem
	err := cartStore.DB.Model(&cartItem).First(&cartItem, *cartItemId).Error
	if err != nil {
		return nil, err
	}
	fmt.Println("before update", cartItem, "quantity: ", cartItem.Quantity)

	newQty := int(cartItem.Quantity) + amount
	if newQty <= 0 {
		return nil, fmt.Errorf("cart item cant be 0 or minus")
	}
	err = cartStore.DB.Model(&cartItem).Select("Quantity").Update("Quantity", newQty).Error
	if err != nil {
		return nil, err
	}
	fmt.Println("after update", cartItem, "quantity: ", cartItem.Quantity)
	return &cartItem, err
}

func (cartStore *Store) GetPopulatedCartByUserId(userId uint) (*respCart, error) {
	var res = make([]getCartRow, 0)
	err := cartStore.DB.Table("carts").
		Select(selectQ).
		Joins(joinWCartItems).
		Joins(joinWProducts).
		Joins(joinWImages).
		Where("user_id = ?",userId).
		Scan(&res).Error
	if err != nil {
		return nil, err
	}
	fmt.Println(res)
	cartResp := convertRowsToResponse(res)

	return cartResp, err
}

func (cartStore *Store) AddToCart(payload *payloads.AddCartItem, cartId uint) (*models.CartItem, error) {
	var cartItem = models.CartItem{
		ProductID: payload.ProductId,
		Quantity:  payload.ProductId,
		CartID:    cartId,
	}

	err := cartStore.DB.Create(&cartItem).Error

	if err != nil {
		return nil, err
	}

	return &cartItem, nil
}

func (cartStore *Store) DeleteCartItemById(itemId uint) error {
	err := cartStore.DB.First(&models.CartItem{}, itemId).Error
	if err != nil {
		return fmt.Errorf("cart item with id: '%v' was not found", itemId)
	}

	err = cartStore.DB.Delete(&models.CartItem{}, itemId).Error

	if err != nil {
		return err
	}

	return nil
}

func (cartStore *Store) ClearCart(Id uint) error {
	var cartModel models.Cart
	err := cartStore.DB.First(&cartModel, Id).Error
	if err != nil {
		return fmt.Errorf("cart with id: '%v' was not found", Id)
	}

	err = cartStore.DB.Model(&cartModel).Association("CartItems").Delete(&models.CartItem{})

	if err != nil {
		return err
	}

	return nil
}
