package cart

import (
	"fmt"
	"net/http"

	"main.go/middlewares"
	"main.go/pkg/models"
	"main.go/pkg/payloads"
	"main.go/pkg/utils"
	"main.go/services/auth"
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

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	AuthTestMW := middlewares.AuthorizeOwnerShipMW[*models.Cart]("id", "cart id is required", "cart with user id: %v was not found", h.store.GetCartById)
	router.HandleFunc(utils.RoutePath("GET", "/carts/{userId}"), Authenticate(h.GetCartByUserId))
	router.HandleFunc(utils.RoutePath("POST", "/carts"), Authenticate(h.AddToCart))
	router.HandleFunc(utils.RoutePath("PATCH", "/carts/{id}/{itemId}"), Authenticate(AuthTestMW(h.ChangeCartItemQty)))
	router.HandleFunc(utils.RoutePath("DELETE", "/carts/{userId}/{itemId}"), Authenticate(h.DeleteCartItemById))
	router.HandleFunc(utils.RoutePath("DELETE", "/carts/{id}"), Authenticate(h.ClearCart))
}

func (h *Handler) GetCartByUserId(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetValidateId(r, "userId")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}

	userIdToken, err := utils.GetUserIdFromTokenPayload(r)
	if err != nil {
		auth.DenyPermission(w)
		return
	}

	cart, err := h.store.GetPopulatedCartByUserId(*userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	
	if cart.UserID != *userIdToken {
		auth.DenyPermission(w)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{"cart": cart})
}

func (h *Handler) AddToCart(w http.ResponseWriter, r *http.Request) {
	userId, _ := utils.GetUserIdCtx(r)
	cartPayload, err := utils.ValidateAndParseBody[payloads.AddCartItem](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	cart, err := h.store.GetCartByUserId(*userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	cartItem, err := h.store.AddToCart(cartPayload, cart.ID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"cartItem": *cartItem,
	})
}


func (h *Handler) ChangeCartItemQty(w http.ResponseWriter, r *http.Request) {
	//cartt, isOk := r.Context().Value(constants.ResourceKey).(*models.Cart)
	payload, err := utils.ValidateAndParseBody[payloads.ChangeCartItemQty](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	Id, err := utils.GetValidateId(r, "id")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid cart id"))
		return
	}
	cartItemId, err := utils.GetValidateId(r, "itemId")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid cart item id"))
		return
	}

	cart, err := h.store.GetCartById(*Id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userIdToken, err := utils.GetUserIdFromTokenPayload(r)
	if err != nil || cart.UserID != *userIdToken {
		auth.DenyPermission(w)
		return
	}

	updatedCartItem, err := h.store.ChangeCartItemQty(cartItemId, payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, map[string]any{"cartItem":updatedCartItem})
}

func (h *Handler) DeleteCartItemById(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetValidateId(r, "userId")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}
	cartItemId, err := utils.GetValidateId(r, "itemId")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}

	userIdToken, err := utils.GetUserIdFromTokenPayload(r)
	if err != nil || *userId != *userIdToken {
		auth.DenyPermission(w)
		return
	}

	err = h.store.DeleteCartItemById(*cartItemId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, map[string]any{})
}


func (h *Handler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userIdToken, err := utils.GetUserIdFromTokenPayload(r)
	if err != nil {
		auth.DenyPermission(w)
		return
	}

	cart, err := h.store.GetCartByUserId(*userIdToken)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	
	if cart.UserID != *userIdToken {
		auth.DenyPermission(w)
		return
	}

	err = h.store.ClearCart(cart.ID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, map[string]any{})
}