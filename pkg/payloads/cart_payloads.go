package payloads


type AddCartItem struct {
	ProductId uint `json:"productId" validate:"required,min=1"`
	Quantity uint `json:"quantity" validate:"required,min=1"`
}

type ChangeCartItemQty struct {
	Operation string `json:"operation" validate:"required,oneof + -"`
	Amount uint `json:"amount" validate:"required,min=1"`
}