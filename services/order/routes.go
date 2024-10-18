package order

import (
	"fmt"
	"net/http"

	"main.go/constants"
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
	router.HandleFunc(utils.RoutePath("GET", "/orders/{id}"), Authenticate(h.GetOrderById))
	router.HandleFunc(utils.RoutePath("GET", "/orders"), Authenticate(h.GetAllOrders))
	router.HandleFunc(utils.RoutePath("POST", "/orders"), Authenticate(h.CreateOrder))
	router.HandleFunc(utils.RoutePath("DELETE", "/orders/{id}"), Authenticate(h.CancelOrderById))
	router.HandleFunc(utils.RoutePath("PATCH", "/orders/{id}/status"), Authenticate(h.UpdateOrderStatusById))
}

func (h *Handler) GetOrderById(w http.ResponseWriter, r *http.Request) {
	Id, err := utils.GetValidateId(r, constants.IdUrlPathKey)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}

	order, err := h.store.GetPopulatedOrderById(*Id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{"order": convertRowToResp(order)})
}

func (h *Handler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	pagination := middlewares.GetPagination(r)

	conditions := utils.GetFilterConditions(r, whiteListedParams)
	sortString := utils.GetSortQ(r, whiteListedSortParams)

	orders, count, errs := utils.GenericFilterWithJoins[models.Order, GetAllOrdersRows](&utils.GenericFilterConfigWithJoins{
		DB:                h.store.DB,
		Filters:           conditions,
		SortQ:             sortString,
		SelectQ:           selectAllOrdersQ,
		Joins:             []string{jointWAddress,joinWOrderItemsCount},
		Pagination:        pagination,
		WhiteListedParams: whiteListedParams,
	})
	if len(errs) != 0 {
		utils.WriteError(w, http.StatusBadRequest, errs[0])
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"page":   pagination.Page,
		"limit":  pagination.Limit,
		"count":  count,
		"orders": convertRowsToResp(orders),
	})
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIdCtx(r)
	if err != nil {
		auth.DenyPermission(w)
		return
	}

	coPayload, err := utils.ValidateAndParseBody[payloads.CreateOrder](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	address, err := h.store.GetAddressById(coPayload.AddressId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if address.UserID != *userId {
		auth.Unauthorized(w)
		return
	}

	cart, err := h.store.GetCart(*userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	orderItems := h.store.ConvertToOrderItems(cart)
	productsIds := h.store.ExtractProductIds(cart)
	prods, err := h.store.GetProductsByIds(productsIds)
	if err != nil {
		// most likely at this point user is hacking so error 500 is returned
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	totalPrice, err := h.store.ValidateAndCalTotalPrice(prods, orderItems)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	order:= models.Order{
		AddressID: address.ID,
		UserID: *userId,
		TotalPrice: *totalPrice,
	}

	cartItemsCount, err := h.store.GetCartItemsCount(*userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if *cartItemsCount == 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("you have no cart items"))
		return
	}

	err = h.store.CreateOrderWithItems(&order, userId, orderItems)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	order.TotalPrice = utils.TruncateToTwoDecimals(order.TotalPrice)
		
	utils.WriteJSON(w, http.StatusCreated, map[string]any{"order": order})
}

func (h *Handler) UpdateOrderStatusById(w http.ResponseWriter, r *http.Request) {
	Id, err := utils.GetValidateId(r, constants.IdUrlPathKey)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}
	userId, err := utils.GetUserIdCtx(r)
	if err != nil {
		auth.DenyPermission(w)
		return
	}
	uPayload, err := utils.ValidateAndParseBody[payloads.UpdateOrder](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid status"))
		return
	}

	err = h.store.UpdateOrderStatus(*Id, *userId, uPayload.Status)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, map[string]any{})
}

func (h *Handler) CancelOrderById(w http.ResponseWriter, r *http.Request) {
	Id, err := utils.GetValidateId(r, constants.IdUrlPathKey)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}
	userId, err := utils.GetUserIdCtx(r)
	if err != nil {
		auth.DenyPermission(w)
		return
	}

	err = h.store.CancelOrder(*Id, *userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, map[string]any{})
}
