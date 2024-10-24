package cart

import (
	"net/http"

	"main.go/errors"
	"main.go/middlewares"
	"main.go/pkg/payloads"
	"main.go/pkg/utils"
	"main.go/services/auth"
	"main.go/types"
)

type Handler struct {
	store types.CartStore
}

func NewHandler(store Store) *Handler {
	return &Handler{
		store: &store,
	}
}

var itemId = "itemId"
var Authenticate = middlewares.Authenticate
func invalidCartItemIdErr(id uint) error {
	return errors.NewInvalidIDError("cart item", id)
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc(utils.RoutePath("GET", "/carts"), Authenticate(h.GetUserCart))
	router.HandleFunc(utils.RoutePath("POST", "/carts"), Authenticate(h.AddToCart))
	router.HandleFunc(utils.RoutePath("PATCH", "/carts/{itemId}"), Authenticate(h.ChangeCartItemQty))
	router.HandleFunc(utils.RoutePath("DELETE", "/carts/{itemId}"), Authenticate(h.DeleteCartItem))
	router.HandleFunc(utils.RoutePath("DELETE", "/carts"), Authenticate(h.ClearCart))
}

func (h *Handler) GetUserCart(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIdCtx(r)
	if err != nil {
		auth.Unauthorized(w)
		return
	}

	cart, err := h.store.GetCartByUserId(*userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{"cart": cart})
}

func (h *Handler) AddToCart(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIdCtx(r)
	if err != nil {
		auth.Unauthorized(w)
		return
	}

	cartPayload, err := utils.ValidateAndParseBody[payloads.AddCartItem](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	cartItem, err := h.store.AddToCart(cartPayload, *userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"cartItem": *cartItem,
	})
}


func (h *Handler) ChangeCartItemQty(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.ValidateAndParseBody[payloads.ChangeCartItemQty](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	cartItemId, err := utils.GetValidateId(r, itemId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, invalidCartItemIdErr(*cartItemId))
		return
	}

	cartItem, err := h.store.GetCartItemById(*cartItemId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userId, err := utils.GetUserIdCtx(r)
	if err != nil {
		auth.Unauthorized(w)
		return
	} else if cartItem.UserID != *userId {
		auth.DenyPermission(w)
		return
	}

	updatedCartItem, err := h.store.ChangeCartItemQty(cartItem.Quantity, payload, cartItem)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, map[string]any{"cartItem":updatedCartItem})
}

func (h *Handler) DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIdCtx(r)
	if err != nil {
		auth.Unauthorized(w)
		return
	}
	cartItemId, err := utils.GetValidateId(r, itemId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, invalidCartItemIdErr(*cartItemId))
		return
	}

	carItem, err := h.store.GetCartItemById(*cartItemId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if carItem.UserID != *userId {
		auth.DenyPermission(w)
		return
	}

	err = h.store.DeleteCartItem(*cartItemId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, map[string]any{})
}


func (h *Handler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userIdToken, err := utils.GetUserIdCtx(r)
	if err != nil {
		auth.DenyPermission(w)
		return
	}
	
	err = h.store.ClearCart(*userIdToken)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, map[string]any{})
}