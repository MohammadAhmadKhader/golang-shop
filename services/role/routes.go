package role

import (
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
var AuthorizeSuperAdmin = middlewares.AuthorizeSuperAdmin
var Pagination = middlewares.PaginationMiddleware

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc(utils.RoutePath("GET","/roles"), Pagination(Authenticate(AuthorizeSuperAdmin(h.GetAllRoles))))
	router.HandleFunc(utils.RoutePath("PUT","/roles/{id}"), Authenticate(AuthorizeSuperAdmin(h.UpdateRole)))
	router.HandleFunc(utils.RoutePath("POST","/roles"), Authenticate(AuthorizeSuperAdmin(h.CreateRole)))
	router.HandleFunc(utils.RoutePath("DELETE","/roles/{id}"), Authenticate(AuthorizeSuperAdmin(h.DeleteRole)))
}

func (h *Handler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	pagination := middlewares.GetPagination(r)

	orders, count, err := h.store.GetAllRoles(pagination.Page, pagination.Limit)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"page":pagination.Page,
		"limit":pagination.Limit,
		"count": count,
		"orders": orders,
	})
}

func (h *Handler) CreateRole(w http.ResponseWriter, r *http.Request) {
	rcPayload , err := utils.ValidateAndParseBody[payloads.CreateRole](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	role, err := h.store.CreateRole(rcPayload.TrimStrs().ToModel())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"role": role,
	})
}

func (h *Handler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	ruPayload , err := utils.ValidateAndParseBody[payloads.CreateRole](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	Id, err := utils.GetValidateId(r, "id")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	role, err := h.store.UpdateRole(*Id, ruPayload.TrimStrs().ToModel())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, map[string]any{
		"role":role,
	})
}

func (h *Handler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	Id, err := utils.GetValidateId(r, "id")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	_,err = h.store.GetRole(*Id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.store.DeleteRole(*Id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, map[string]any{})
}