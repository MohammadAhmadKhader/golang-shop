package websocket

import (
	"net/http"

	"main.go/constants"
)


func Setup(manager *Manager,router *http.ServeMux) {
	router.HandleFunc(constants.Prefix+"/ws", manager.serveWS)
}