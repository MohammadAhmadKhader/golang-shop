package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// read and write buffers will be adjusted lately to avoid wasting memory
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}