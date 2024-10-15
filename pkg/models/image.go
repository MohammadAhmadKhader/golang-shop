package models

type Image struct {
	ModelBasicsTrackedDel
	ProductID uint `json:"productId" gorm:"index"`
	ImageUrl string `json:"imageUrl" gorm:"not null;size:256"`
	IsMain *bool `json:"isMain" gorm:"default:false;not null"`
	ImagePublicId string `json:"imagePublicId" gorm:"not null;size:128"`
}