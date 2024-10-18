package models
// the productId userId index must be search on by productId first then userId to use the index efficiently
type CartItem struct {
	ModelBasics
	Product   *Product `json:"product,omitempty" gorm:"foreignkey:ProductID;constraint:OnDelete:CASCADE"`
	ProductID uint `json:"productId" gorm:"uniqueIndex:idx_user_product,not null"`
	Quantity uint `json:"quantity" gorm:"not null;check:quantity > 0"`
	UserID uint `json:"userId" gorm:"uniqueIndex:idx_user_product,not null"`
}