package types

type GetCartRow struct {
	ID           uint
	Quantity     uint
	ProductID    uint
	ProductName  string
	ProductPrice float64
	ProductImage string
}

type RespCartItemProduct struct {
	ID    uint    `json:"id"`
	Image string  `json:"image"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type RespCartItem struct {
	ID       uint        `json:"id"`
	Quantity uint        `json:"quantity"`
	Product  RespCartItemProduct `json:"product"`
}

type RespCartShape struct {
	CartItems []RespCartItem `json:"cartItems"`
}
