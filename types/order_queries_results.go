package types

import "time"

// ** Get all orders types

type GetAllOrdersRows struct {
	Id                   uint
	UserId               uint
	TotalPrice           float64
	Status               string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	OrderItemsCount      uint
	AddressId            uint
	AddressFullName      string
	AddressCity          string
	AddressCountry       string
	AddressStreetAddress string
	AddressZipCode       *string
	AddressState         *string
}

type RespAllOrders struct {
	Id              uint        `json:"id"`
	UserId          uint        `json:"userId"`
	TotalPrice      float64     `json:"totalPrice"`
	Status          string      `json:"status"`
	OrderItemsCount uint        `json:"orderItemsCount"`
	Address         RespOrderAddress `json:"address"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt"`
}

// ** Get one order types
type GetOneOrderRow struct {
	Id                   uint
	UserId               uint
	TotalPrice           float64
	Status               string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	OrderItemId          uint
	UnitPrice            float64
	OrderItemQuantity    uint8
	AddressId            uint
	AddressFullName      string
	AddressCity          string
	AddressCountry       string
	AddressStreetAddress string
	AddressZipCode       *string
	AddressState         *string
	ProductId            uint
	ProductName          string
	ProductQuantity      uint
	ProductPrice         float64
	ProductImageUrl string
	ProductImageIsMain bool
}

type RespOneOrder struct {
	Id         uint            `json:"id"`
	UserId     uint            `json:"userId"`
	TotalPrice float64         `json:"totalPrice"`
	Status     string          `json:"status"`
	OrderItems []RespOrderItem `json:"orderItems"`
	Address    RespOrderAddress     `json:"address"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
}

type RespOrderItem struct {
	Id       uint        `json:"id"`
	Price    float64     `json:"price"`
	Quantity uint8       `json:"quantity"`
	Product  RespOrderItemProduct `json:"product"`
}

type RespOrderItemProduct struct {
	Id       uint    `json:"id"`
	Name     string  `json:"name"`
	Quantity uint    `json:"quantity"`
	Price    float64 `json:"price"`
	MainImage RespOrderItemImage `json:"mainImage"`
}

type RespOrderItemImage struct {
	IsMain bool `json:"isMain"`
	ImageUrl string `json:"imageUrl"`
}

// ** Shared types

type RespOrderAddress struct {
	Id            uint    `json:"id"`
	FullName      string  `json:"fullName"`
	City          string  `json:"city"`
	StreetAddress string  `json:"streetAddress"`
	State         *string `json:"state"`
	ZipCode       *string `json:"zipCode"`
	Country       string  `json:"country"`
}