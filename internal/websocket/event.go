package websocket

import "encoding/json"

type EventType string

const (
	MessageCreate         EventType = "message_create"
	MessageUpdate         EventType = "message_update"
	MessageDelete         EventType = "message_delete"
	ProductsStockUpdate EventType = "products_stock_update"
)

type Event struct {
	Type    EventType       `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, client *Client) error

func NewProductStockUpdateEvent(payload []byte) Event{
	return Event{
		Type: ProductsStockUpdate,
		Payload: payload,
	}
}