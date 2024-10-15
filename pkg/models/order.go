package models

type Status string

const (
	Pending   Status = "Pending"
	Delivered Status = "Delivered"
	Cancelled  Status = "Cancelled"
)

type Order struct {
	ModelBasics
	User       *User   `json:"user,omitempty" gorm:"foreignkey:UserID;constraint:OnDelete:CASCADE"`
	UserID     uint    `json:"userId" gorm:"index;not null"`
	TotalPrice float64 `json:"totalPrice" gorm:"check: total_price > 0;type:decimal(7,2)"`
	Status     Status `json:"status" gorm:"default:Pending;size:16;not null;index"`
	OrderItems []OrderItem `json:"orderItems" gorm:"foreignkey:OrderID;constraint:OnDelete:CASCADE"`
	AddressID uint `json:"addressId" gorm:"not null"`
	Address *Address `json:"address,omitempty" gorm:"foreignkey:AddressID"`
}
