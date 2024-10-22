package websocket

import "encoding/json"

type EventType string

const (
	MessageCreate EventType = "message_create"
	MessageCreated EventType = "message_created"
	MessageUpdate EventType = "message_update"
	MessageUpdated EventType = "message_updated"
	MessageDelete EventType = "message_delete"
	MessageDeleted EventType = "message_deleted"
	MessageStatusUpdate EventType = "message_status_update"
	MessageStatusUpdated EventType = "message_status_updated"
	ProductsStockUpdate EventType = "products_stock_update"
)

type Event struct {
	Type    EventType       `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, client *Client) error

func NewProductStockUpdateEvent(payload []byte) Event {
	return Event{
		Type:    ProductsStockUpdate,
		Payload: payload,
	}
}

func NewEvent(eventType EventType, payload []byte) Event {
	return Event{
		Type:    eventType,
		Payload: payload,
	}
}
