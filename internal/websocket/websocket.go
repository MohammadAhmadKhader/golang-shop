package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

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

// we will use userId => connection, so when we want to send a message to a specific user we do not need to loop all over users.
var muProduct sync.RWMutex
func BroadcastProductQtyChange(products []WSProduct) {
	productsMsg, err := json.Marshal(products)
	if err != nil {
		log.Fatal(err)
		return
	}

	muProduct.Lock()
	defer muProduct.Unlock()
	for client := range GlobalManager.clients {
		err := client.conn.WriteJSON(NewProductStockUpdateEvent(productsMsg))
		if err != nil {
			log.Fatal(err)
			GlobalManager.deleteClient(client)
			return
		}
	}
}
