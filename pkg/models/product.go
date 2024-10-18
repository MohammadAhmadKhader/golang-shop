package models


type Product struct {
	ModelBasicsTrackedDel
	Name         string    `json:"name" gorm:"size:32;not null"`
	Quantity     uint      `json:"quantity" gorm:"check:quantity > 0"`
	Image        *Image    `json:"mainImage,omitempty" gorm:"-"`
	Images       []Image   `json:"images,omitempty" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Reviews      []Review  `json:"reviews,omitempty" gorm:"foreignKey:ProductID;OnDelete:CASCADE"`
	Description  *string   `json:"description" gorm:"size:256"`
	Category     *Category `json:"-" gorm:"-"`
	CategoryName *string   `json:"category,omitempty" gorm:"-"`
	CategoryID   uint      `json:"categoryId,omitempty" gorm:"not null"`
	Price        float64   `json:"price" gorm:"check:price > 0;type:decimal(7,2)"`
	AvgRating    *float64  `json:"avgRating,omitempty" gorm:"-"`
}
