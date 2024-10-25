package category

import (
	"net/http"

	"main.go/constants"
	"main.go/errors"
	"main.go/middlewares"
	"main.go/pkg/payloads"
	"main.go/pkg/utils"
	"main.go/types"
)

type Handler struct {
	store types.CategoryStore
}

func NewHandler(store Store) *Handler {
	return &Handler{
		store: &store,
	}
}

var Authenticate = middlewares.Authenticate
var AuthorizeAdmin = middlewares.AuthorizeAdmin
var Pagination = middlewares.PaginationMiddleware

func invalidCategoryIdErr(id uint) error {
	return errors.NewInvalidIDError("category", id)
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc(utils.RoutePath("GET","/categories"), Pagination(h.GetAllCategories))
	router.HandleFunc(utils.RoutePath("GET","/categories/{id}"), h.GetCategoryById)
	router.HandleFunc(utils.RoutePath("POST","/categories"), Authenticate(AuthorizeAdmin(h.CreateCategory)))
	router.HandleFunc(utils.RoutePath("PUT","/categories/{id}"), Authenticate(AuthorizeAdmin(h.UpdateCategory)))
}

func (h *Handler) GetCategoryById(w http.ResponseWriter, r *http.Request){
	Id, err := utils.GetValidateId(r, constants.IdUrlPathKey)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, invalidCategoryIdErr(*Id))
		return
	}

	category, err := h.store.GetCategoryById(*Id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest , err)
		return
	}

	utils.WriteJSON(w, http.StatusOK , map[string]any{"category":category})
}

func (h *Handler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	pagination := middlewares.GetPagination(r)

	categories, count, err := h.store.GetAllCategories(pagination.Page, pagination.Limit)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"page":pagination.Page,
		"limit":pagination.Limit,
		"count": count,
		"categories": categories,
	})
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	catePayload, err := utils.ValidateAndParseBody[payloads.CreateCategory](r)
	catePayload.TrimStrs()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	model := catePayload.ToModel()

	category, err := h.store.CreateCategory(model)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"category": *category,
	})
}

func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	catePayload, err := utils.ValidateAndParseBody[payloads.UpdateCategory](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	catePayload.TrimStrs()
	model := catePayload.ToModel()
	Id, err := utils.GetValidateId(r, constants.IdUrlPathKey)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, invalidCategoryIdErr(*Id))
		return
	}

	category, err := h.store.UpdateCategory(*Id, model)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, map[string]any{
		"category": *category,
	})
}