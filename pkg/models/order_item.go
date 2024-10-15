package models

type OrderItem struct {
	Identifier
	OrderID uint `json:"orderId" gorm:"index;not null"`
	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	ProductID uint `json:"productId" gorm:"not null"`
	UnitPrice float64 `json:"unitPrice" gorm:"not null;check: quantity >= 0;type:decimal(7,2)"`
	Quantity uint `json:"quantity" gorm:"not null;type:TINYINT"`
}