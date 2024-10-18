package models

type Review struct {
	ModelBasicsTrackedDel
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;OnDelete:CASCADE"`
	UserID uint `json:"userId" gorm:"uniqueIndex:idx_user_product,not null"`
	Product *Product `json:"product,omitempty"`
	ProductID uint `json:"productId" gorm:"uniqueIndex:idx_user_product,not null"`
	Comment string `json:"comment" gorm:"not null;size:256"`
	Rate uint8 `json:"rate" gorm:"not null;type:TINYINT;check:rate >=1 AND rate <= 5"`
}