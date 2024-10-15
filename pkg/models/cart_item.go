package models

type CartItem struct {
	ModelBasics
	Product   *Product `json:"product,omitempty" gorm:"foreignkey:ProductID"`
	ProductID uint `json:"productId" gorm:"not null"`
	Quantity uint `json:"quantity" gorm:"not null;check:quantity > 0"`
	CartID uint `json:"cartId" gorm:"not null"`
}