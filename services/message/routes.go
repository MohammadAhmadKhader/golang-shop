package message

import (
	"fmt"
	"net/http"

	"main.go/errors"
	"main.go/middlewares"
	"main.go/pkg/utils"
	"main.go/types"
)

type Handler struct {
	store types.MessageStore
}


func NewHandler(store Store) *Handler {
	return &Handler{
		store: &store,
	}
}

var Authenticate = middlewares.Authenticate
var AuthorizeAdmin = middlewares.AuthorizeAdmin

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc(utils.RoutePath("GET", "/messages/users/{userId}"), Authenticate(AuthorizeAdmin(h.GetUserMessages)))
	router.HandleFunc(utils.RoutePath("GET", "/messages"), Authenticate(h.GetChatMessages))
}


func (h *Handler) GetUserMessages(w http.ResponseWriter, r *http.Request) {
	userId, receivedStr,err := utils.GetValidateId(r, "userId")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.NewInvalidIDError("user id", receivedStr))
		return
	}
	
	lastMessageId, limit, cursor, err := GetMessagesParamsWithoutUserId(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	messages, err := h.store.GetById(*userId, uint(lastMessageId), *cursor, limit)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"messages":messages,
	})

}

func (h *Handler) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIdCtx(r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("something went wrong during retrieving user id"))
		return
	}
	to, lastMessageId, limit, cursor, err := GetMessagesParams(r, "to")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	messages, err := h.store.GetByUsersIds(*userId, uint(to), uint(lastMessageId), *cursor ,limit)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"messages":messages,
	})
}