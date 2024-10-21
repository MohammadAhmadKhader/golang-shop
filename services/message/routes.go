package message

import (
	"fmt"
	"net/http"

	"main.go/middlewares"
	"main.go/pkg/utils"
)

type Handler struct {
	store Store
}


func NewHandler(store Store) *Handler {
	return &Handler{
		store: store,
	}
}

var Authenticate = middlewares.Authenticate
var AuthorizeAdmin = middlewares.AuthorizeAdmin

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc(utils.RoutePath("GET", "/messages/{id}"), Authenticate(AuthorizeAdmin(h.GetUserMessages)))
	router.HandleFunc(utils.RoutePath("GET", "/messages"), Authenticate(h.GetChatMessages))
}


func (h *Handler) GetUserMessages(w http.ResponseWriter, r *http.Request) {
	from, lastMessageId,limit,cursor, err := GetMessagesParams(r, "from")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	messages, err := h.store.GetById(uint(from), uint(lastMessageId), *cursor, limit)
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
	to, lastMessageId,limit,cursor, err := GetMessagesParams(r, "to")
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