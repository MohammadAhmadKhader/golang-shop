package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"main.go/pkg/models"
	"main.go/types"
)

// TODO: there functions must be moved the manager, so they used the manager locks, or these locks will not function properly
// read and write buffers will be adjusted lately to avoid wasting memory
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var muRegistedClients sync.RWMutex
// we will use userId => connection, so when we want to send a message to a specific user we do not need to loop all over users.
var muClients sync.RWMutex
func BroadcastProductQtyChange(products []types.ProductAmountDiscounter) {
	productsMsg, err := json.Marshal(products)
	if err != nil {
		log.Fatal(err)
		return
	}

	muClients.RLock()
	defer muClients.RUnlock()
	for client := range GlobalManager.clients {
		err := client.conn.WriteJSON(NewProductStockUpdateEvent(productsMsg))
		if err != nil {
			log.Fatal(err)
			GlobalManager.deleteClient(client)
			return
		}
	}
}

// broadcasts Create and Update for messages
func BroadcastCUMessage(message models.Message, userIds []uint, eventType EventType) {
	productMsg, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
		return
	}

	muRegistedClients.RLock()
	defer muRegistedClients.RUnlock()
	for client := range GlobalManager.clients {
		err := client.conn.WriteJSON(NewEvent(eventType, productMsg))
		if err != nil {
			log.Println(err)
			GlobalManager.deleteClient(client)
			return
		}
	}
}

// broadcasts Delete for messages
func BroadcastDMessage(payload DeleteMessagePayload, userIds []uint) {
	deletePayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
		return
	}

	muRegistedClients.RLock()
	defer muRegistedClients.RUnlock()
	for _, userId := range userIds {
		clients, isOk := GlobalManager.registedClients[userId]
		if isOk {
			for _, client := range clients {
				err := client.conn.WriteJSON(NewEvent(MessageDeleted,deletePayload))
				if err != nil {
					log.Println(err)
					GlobalManager.deleteClient(client)
					return
				}	
			}
		}
	}
}