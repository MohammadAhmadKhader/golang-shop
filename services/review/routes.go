package review

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
	router.HandleFunc(utils.RoutePath("GET", "/reviews"), Authenticate(h.GetAllReviews))
	router.HandleFunc(utils.RoutePath("POST", "/reviews"), Authenticate(h.AddReview))
	router.HandleFunc(utils.RoutePath("PUT", "/reviews/{id}"), Authenticate(h.EditReview))
	router.HandleFunc(utils.RoutePath("DELETE", "/reviews/{id}"), Authenticate(h.DeleteReview))
}

func (h *Handler) GetAllReviews(w http.ResponseWriter, r *http.Request) {
	pagination := middlewares.GetPagination(r)

	conditions := utils.GetFilterConditions(r, whiteListedParams)
	sortString := utils.GetSortQ(r, whiteListedSortParams)

	reviews, count, err := utils.GenericFilterWithJoins[models.Review, GetAllReviewsRow](
		&utils.GenericFilterConfigWithJoins{
			DB:                h.store.DB,
			Filters:           conditions,
			SortQ:             sortString,
			Pagination:        pagination,
			WhiteListedParams: whiteListedParams,
			SelectQ:           reviewsSelectCols,
			Joins:             []string{reviewsJoin},
		})

	if len(err) != 0 {
		utils.WriteError(w, http.StatusBadRequest, err[0])
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"page":    pagination.Page,
		"limit":   pagination.Limit,
		"count":   count,
		"reviews": *convertRowsToResp(reviews),
	})
}

func (h *Handler) AddReview(w http.ResponseWriter, r *http.Request) {
	crPayload, err := utils.ValidateAndParseBody[payloads.CreateReview](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}

	userId, err:=utils.GetUserIdFromTokenPayload(r)
	if err != nil {
		auth.DenyPermission(w)
		return
	}
	crPayload.UserId = *userId

	review, err := h.store.CreateReview(crPayload.TrimStrs().ToModel())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{"review": review})
}

func (h *Handler) EditReview(w http.ResponseWriter, r *http.Request) {
	Id, err := utils.GetValidateId(r, constants.IdUrlPathKey)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}

	upPayload, err := utils.ValidateAndParseBody[payloads.UpdateReview](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	oldRev, err := h.store.GetReviewById(*Id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	
	userId, err:=utils.GetUserIdFromTokenPayload(r)
	if err != nil {
		auth.DenyPermission(w)
		return
	}
	
	if oldRev.UserID != *userId {
		auth.DenyPermission(w)
		return
	}

	model := upPayload.TrimStrs().ToModel()
	review, err := h.store.UpdateReview(*Id, model, upPayload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, map[string]any{"review": review})
}

func (h *Handler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	Id, err := utils.GetValidateId(r, constants.IdUrlPathKey)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}
	
	userId, err:=utils.GetUserIdFromTokenPayload(r)
	if err != nil {
		auth.DenyPermission(w)
		return
	}

	err = h.store.HardDeleteReview(*Id, *userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, map[string]any{})
}

