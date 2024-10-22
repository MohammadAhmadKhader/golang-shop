package websocket

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"main.go/internal/database"
	"main.go/pkg/utils"
)

type Manager struct {
	sync.RWMutex
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

	client := NewClient(ws, m)
	userId, err := utils.GetUserIdCtx(r)
	if err != nil || !isCorrectOTP {
		log.Println("Guest client")
	} else {
		log.Println("RegularUser client")
	}
	m.addClient(client, userId)

	
	// this only available for validated OTP (session Id)
	// un-registed clients can listen to the product stocks changes only so we run the heart beats for them

	if isCorrectOTP {
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
	m.Lock()
	if _, ok := m.clients[client]; ok {
		client.conn.Close()
		delete(m.clients, client)
	}
	m.Unlock()
}

func (m *Manager) addClient(client *Client, userId *uint) {
	m.Lock()
	m.clients[client] = true
	if userId != nil {
		m.registedClients[*userId] = append(m.registedClients[*userId], client)
	}

	m.Unlock()
}
