package models


type User struct {
	ModelBasicsTrackedDel
	Name         string  `json:"name" gorm:"not null;size:32"`
	Avatar       *string `json:"avatar" gorm:"default:NULL;size:256"`
	Email        string  `json:"email" gorm:"uniqueIndex;not null;size:64"`
	Password     string  `json:"-" gorm:"size:128;not null"`
	MobileNumber *string `json:"mobileNumber" gorm:"default:NULL;size:32"`
	Roles        []Role  `json:"roles,omitempty" gorm:"many2many:user_roles;foreignKey:ID;joinForeignKey:UserID;References:ID;joinReferences:RoleID;constraint:OnDelete:CASCADE;"`
	CartItems []CartItem `json:"cart,omitempty" gorm:"foreignkey:UserID;constraint:OnDelete:CASCADE"`
}