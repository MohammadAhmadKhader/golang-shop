package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"main.go/config"
	"main.go/internal/database"
	"main.go/middlewares"
	"main.go/services"
)

// issues 1- cookie not  in post man
func main() {
	port := config.Envs.Port
	server := http.NewServeMux()
	loggedServer := middlewares.Logger(server)

	corsServer := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PATCH", "PUT"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)(loggedServer)

	DB := database.DB

	services.SetupAllServices(DB, server)

	log.Printf("Listening to Port: %s", port[1:])
	err := http.ListenAndServe(port, corsServer)
	if err != nil {
		log.Fatal(err)
	}
}
