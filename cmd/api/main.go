package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"main.go/config"
	"main.go/internal/database"
	"main.go/internal/websocket"
	"main.go/middlewares"
	"main.go/services"
)

func main() {
	port := config.Envs.Port
	server := http.NewServeMux()

	DB := database.DB
	wsManager := websocket.NewManager(context.Background())
	websocket.Setup(wsManager, server)

	services.SetupAllServices(DB, server)
	loggedServer := middlewares.Logger(server)

	corsServer := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PATCH", "PUT"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)(loggedServer)

	log.Printf("Listening to Port: %s", port[1:])
	err := http.ListenAndServe(port, corsServer)
	if err != nil {
		log.Fatal(err)
	}
}
