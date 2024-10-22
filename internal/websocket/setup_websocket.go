package websocket

import (
	"net/http"

	"main.go/constants"
	"main.go/middlewares"
)


func Setup(manager *Manager,router *http.ServeMux) {
	router.HandleFunc(constants.Prefix+"/ws", middlewares.AuthenticateIfCookieExist(manager.serveWS))
}