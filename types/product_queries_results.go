package types

import "time"

// ** Get One Product By Id

type RespGetOneProductShape struct {
	ID           uint                     `json:"id"`
	Name         string                   `json:"name"`
	Description  *string                  `json:"description"`
	Quantity     uint                     `json:"quantity"`
	Price        float64                  `json:"price"`
	CategoryName string                   `json:"category"`
	Category     RowProductCategory       `json:"-"`
	AvgRating    float64                  `json:"avgRating"`
	Images       []RowGetOneProductImage  `json:"images"`
	Reviews      []RespProductReviewShape `json:"reviews"`
	CreatedAt    time.Time                `json:"createdAt"`
	UpdatedAt    time.Time                `json:"updatedAt"`
}
type RowGetProductById struct {
	ID              uint
	Name            string
	Description     *string
	Quantity        uint
	Price           float64
	CategoryID      uint
	AvgRating       float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
	ReviewID        uint
	ReviewUserID    uint
	ReviewComment   string
	ReviewRate      uint8
	ReviewCreatedAt time.Time
	ReviewUpdatedAt time.Time
	CategoryName    string
	ImageId         uint
	ImageUrl        string
	ImagePublicId   string
	IsMain          bool
	UserName        string
	UserEmail       string
	UserAvatar      string
}
type RespProductReviewShape struct {
	ReviewID  uint                 `json:"id"`
	UserID    uint                 `json:"userId"`
	User      RowProductReviewUser `json:"user"`
	Comment   string               `json:"comment"`
	Rate      uint8                `json:"rate"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt"`
}

type RowGetOneProductImage struct {
	ImageId       uint   `json:"id"`
	ImageUrl      string `json:"imageUrl"`
	ImagePublicId string `json:"imagePublicId"`
	IsMain        bool   `json:"isMain"`
}

type RowProductCategory struct {
	CategoryName string
}

type RowProductReviewUser struct {
	UserName   string `json:"name"`
	UserEmail  string `json:"email"`
	UserAvatar string `json:"avatar"`
}

// ** Get All Products

type GetAllProductsRow struct {
	Id            uint
	Name          string
	Quantity      uint
	Description   string
	CategoryId    uint
	Price         float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ImageId       uint
	ImageUrl      string
	ImagePublicId string
	AvgRating     float64
}

type RespGetAllProductsShape struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Quantity    uint      `json:"quantity"`
	Description string    `json:"description"`
	CategoryId  uint      `json:"categoryId"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Image       RowGetAllProductsImage  `json:"mainImage"`
	AvgRating   float64   `json:"avgRating"`
}
type RowGetAllProductsImage struct {
	ImageId       uint   `json:"id"`
	ImageUrl      string `json:"imageUrl"`
	ImagePublicId string `json:"imagePublicId"`
}
