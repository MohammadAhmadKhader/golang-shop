package payloads

import "main.go/pkg/models"

type OrderPayloadItem struct {
	ProductId uint `json:"productId" validate:"required,min=1"`
	Quantity  uint `json:"quantity" validate:"required,min=1"`
}

type CreateOrder struct {
	AddressId uint  `json:"addressId" validate:"required,min=1"`
}

type UpdateOrder struct {
	Status models.Status `json:"status" validate:"required"`
}

func (co *CreateOrder) GetProductsIds(orderItems []OrderPayloadItem) []uint {
	var productsIds = make([]uint, 0, len(orderItems))
	for _, orderPayloadItem := range orderItems {
		productsIds = append(productsIds, orderPayloadItem.ProductId)
	}
	
	return productsIds
}