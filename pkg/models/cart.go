package models


type Cart struct {
	Identifier
	UserID uint `json:"userId" gorm:"uniqueIndex;not null"`
	CartItems []CartItem `json:"cartItems,omitempty" gorm:"foreignkey:CartID;constraint:OnDelete:CASCADE"`
}

func (c *Cart) GetUserId() uint {
	return c.UserID
}