package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"main.go/internal/database"
	"main.go/pkg/models"
	"main.go/pkg/utils"
	"main.go/types"
)
// TODO: the locking mechanism must be re-looked at, as there are chances to have a deadlock.
type Manager struct {
	clientsLock sync.RWMutex
	registedClientsLock sync.RWMutex
	clients         Clients
	registedClients RegistedClients
	Otps            RetentionMap
	handlers        map[string]EventHandler
	store   		*Store
}

var GlobalManager *Manager
var retentionPeriod = 5 * time.Second

func NewManager(ctx context.Context) *Manager {
	if GlobalManager == nil {
		GlobalManager = &Manager{
			clients:  make(Clients),
			handlers: make(map[string]EventHandler),
			registedClients: make(map[uint][]*Client),
			Otps:     NewRetentionMap(ctx, retentionPeriod),
			store: NewStore(database.DB),
		}

		GlobalManager.setupEventHandlers()
		return GlobalManager
	}
	return GlobalManager
}

func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request) {
	otp := r.URL.Query().Get("otp")
	fmt.Println("otp key: ", otp)
	isCorrectOTP := otp != "" && m.Otps.ValidateOTP(otp)
	if !isCorrectOTP {
		log.Println("invalid OTP")
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	var client *Client
	userId, err := utils.GetUserIdCtx(r)
	if err != nil || !isCorrectOTP {
		log.Println("Guest client")
		client = NewClient(ws, m, nil)
	} else {
		log.Println("RegularUser client")
		client = NewClient(ws, m, userId)
	}
	m.addClient(client)

	// this only available for validated OTP (session Id)
	// un-registed clients can listen to the product stocks changes only so we run the heart beats for them

	if isCorrectOTP && client.id != nil {
		go client.readMessage()
		go client.writeMessage()
	} else {
		go client.runHeartBeat()
	}
}

func (m *Manager) setupEventHandlers() {
	m.handlers[string(MessageCreate)] = HandleMessageCreate
	m.handlers[string(MessageUpdate)] = HandleMessageUpdate
	m.handlers[string(MessageDelete)] = HandleMessageDelete
	m.handlers[string(MessageStatusUpdate)] = HandleMessageUpdateStatus
}

func (m *Manager) routeHandler(event Event, client *Client) error {
	handler, ok := m.handlers[string(event.Type)]
	if ok {
		if err := handler(event, client); err != nil {
			return err
		}

		return nil
	} else {
		return fmt.Errorf("received unknown event type: %v", event.Type)
	}
}

func (m *Manager) deleteClient(client *Client) {
	m.clientsLock.Lock()
	if _, ok := m.clients[client]; ok {
		client.conn.Close()
		delete(m.clients, client)
	}
	m.clientsLock.Unlock()

	if client.id != nil {
		m.registedClientsLock.Lock()
		m.registedClients[*client.id] = append(m.registedClients[*client.id], client)
		m.registedClientsLock.Unlock()
	}
}

func (m *Manager) addClient(client *Client) {
	m.clientsLock.Lock()
	m.clients[client] = true
	m.clientsLock.Unlock()

	if client.id != nil {
		m.registedClientsLock.Lock()
		m.registedClients[*client.id] = append(m.registedClients[*client.id], client)
		m.registedClientsLock.Unlock()
	}
}


// broadcasts Create and Update for messages
func (m *Manager) BroadcastCUMessage(message models.Message, userIds []uint, eventType EventType) {
	productMsg, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
		return
	}

	m.registedClientsLock.RLock()
	defer m.registedClientsLock.RUnlock()
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
func (m *Manager) BroadcastDMessage(payload DeleteMessagePayload, userIds []uint) {
	deletePayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
		return
	}

	m.registedClientsLock.RLock()
	defer m.registedClientsLock.RUnlock()
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

func (m *Manager) BroadcastProductQtyChange(products []types.ProductAmountDiscounter) {
	productsMsg, err := json.Marshal(products)
	if err != nil {
		log.Fatal(err)
		return
	}

	GlobalManager.clientsLock.RLock()
	defer GlobalManager.clientsLock.RUnlock()
	for client := range GlobalManager.clients {
		err := client.conn.WriteJSON(NewProductStockUpdateEvent(productsMsg))
		if err != nil {
			log.Fatal(err)
			GlobalManager.deleteClient(client)
			return
		}
	}
}
