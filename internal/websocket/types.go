package websocket


type WSProduct struct {
	ID             uint `json:"id"`
	DiscountAmount uint `json:"discountAmount"`
}


func (p WSProduct) GetProductId() uint {
	return p.ID
}

func (p WSProduct) GetAmountDiscount() uint {
	return p.DiscountAmount
}

type UserMessagePayload struct {
	From    uint   `json:"from"`
	To      uint   `json:"to"`
	Message string `json:"message"`
}