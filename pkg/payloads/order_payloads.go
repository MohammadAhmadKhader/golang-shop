package payloads

import "main.go/pkg/models"

type OrderPayloadItem struct {
	ProductId uint `json:"productId"`
	Quantity  uint `json:"quantity"`
}

type CreateOrder struct {
	AddressId uint  `json:"addressId"`
}

type UpdateOrder struct {
	Status models.Status `json:"status"`
}

func (co *CreateOrder) GetProductsIds(orderItems []OrderPayloadItem) []uint {
	var productsIds = make([]uint, 0, len(orderItems))
	for _, orderPayloadItem := range orderItems {
		productsIds = append(productsIds, orderPayloadItem.ProductId)
	}
	
	return productsIds
}