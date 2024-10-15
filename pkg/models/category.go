package models

type Category struct {
	ModelBasicsTrackedDel
	Name     string `json:"name" gorm:"not null;size:32"`
	Products []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
}