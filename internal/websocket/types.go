package websocket

import (
	"main.go/pkg/payloads"
)
// TODO: Move all payloads and their validation to here and refactor handlers accordingly
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

type MessagePayload struct {
	Id      uint   `json:"id"`
	From    uint   `json:"from"`
	To      uint   `json:"to"`
	Content string `json:"content"`
}

type DeleteMessagePayload struct {
	Id uint `json:"id"`
}

type MessageStatusPayload struct {
	Id     uint   `json:"id" validate:"min=1"`
	Status string `json:"status" validate:"oneof=Seen Delivered"`
}

func (mp *MessagePayload) ToCreatePayload() payloads.CreateMessage {
	return payloads.CreateMessage{
		From:    mp.From,
		To:      mp.To,
		Content: mp.Content,
	}
}

func (mp *MessagePayload) ToUpdatePayload() payloads.UpdateMessage {
	return payloads.UpdateMessage{
		Id:      mp.Id,
		Content: mp.Content,
	}
}

func (dmp *DeleteMessagePayload) ToDeletePayload() payloads.DeleteMessage {
	return payloads.DeleteMessage{
		Id: dmp.Id,
	}
}
