package address

import (
	"fmt"
	"net/http"

	"main.go/middlewares"
	"main.go/pkg/payloads"
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

// userId was added to avoid URL conflict with the generic routes and also its a good way to stick to REST rules
func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	var _ = middlewares.AuthorizeUser("id", "address id is required", "address with id: '%v' was not found", h.store.GetAddressById)
	
	router.HandleFunc(utils.RoutePath("GET", "/users/{id}/addresses"), Authenticate(h.GetAllAddress))
	router.HandleFunc(utils.RoutePath("GET", "/users/{id}/addresses/{addressId}"), Authenticate(h.GetAddressById))
	router.HandleFunc(utils.RoutePath("POST", "/users/{id}/addresses"), Authenticate(h.CreateAddress))
	router.HandleFunc(utils.RoutePath("PUT", "/users/{id}/addresses/{addressId}"), Authenticate(h.UpdateAddress))
	router.HandleFunc(utils.RoutePath("DELETE", "/users/{id}/addresses/{addressId}"), Authenticate(h.DeleteAddress))
}


func (h *Handler) GetAddressById(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIdCtx(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	addressId, err := utils.GetValidateId(r, "addressId")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid address id"))
		return
	}

	address, err := h.store.GetById(*addressId, *userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid address id"))
		return
	}


	utils.WriteJSON(w, http.StatusOK, map[string]any{"address": *address})
}

func (h *Handler) GetAllAddress(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIdCtx(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	addresses, err := h.store.GetAllAddresses(*userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{"addresses":addresses})
}

func (h *Handler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIdCtx(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	addrPayload, err := utils.ValidateAndParseBody[payloads.CreateAddress](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	address := addrPayload.TrimStrs().ToModel(*userId)
	newAddress, err := h.store.CreateAddress(address)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{"address":newAddress})
}

func (h *Handler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIdCtx(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	addrPayload, err := utils.ValidateAndParseBody[payloads.UpdateAddress](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if addrPayload.IsEmpty() {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("new address is required"))
		return
	}

	addressId, err := utils.GetValidateId(r, "addressId")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid address id"))
		return
	}
	count, err := h.store.GetUndeletedAddressesCount(*userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if *count == 10 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("max addresses allowed is 10"))
		return
	}

	address := addrPayload.TrimStrs().ToModel()
	updatedAddr, err := h.store.UpdateAddress(*addressId, address, addrPayload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, map[string]any{"address":updatedAddr})
}

func (h *Handler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIdCtx(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	addressId, err := utils.GetValidateId(r, "addressId")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid address id"))
		return
	}

	err = h.store.DeleteAddress(*addressId ,*userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, map[string]any{})
}