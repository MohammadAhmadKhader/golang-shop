package websocket

import (
	"time"
)

// this is the default and lowest value
const (
	defaultLimit = 10
)



type Cursor struct {
	CreatedAt     time.Time `json:"createdAt"`
    Limit         int       `json:"limit"`   
}

func (c *Cursor) ValidateCursor() *Cursor {
	if c != nil {
		if c.Limit < defaultLimit {
			c.Limit = defaultLimit
		}
	}

	return c
}