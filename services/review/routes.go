package review

import (
	"errors"
	"fmt"
	"net/http"

	"main.go/constants"
	appErrors "main.go/errors"
	"main.go/middlewares"
	"main.go/pkg/models"
	"main.go/pkg/payloads"
	"main.go/pkg/utils"
	"main.go/types"
)

type Handler struct {
	store types.ReviewStore
}

func NewHandler(store Store) *Handler {
	return &Handler{
		store: &store,
	}
}

const (
	reviewId = "reviewId"
)

func invalidRevIdErr(id uint) error {
	return appErrors.NewInvalidIDError("review", id)
}

var Authenticate = middlewares.Authenticate
var AuthorizeAdmin = middlewares.AuthorizeAdmin
var Pagination = middlewares.PaginationMiddleware

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc(utils.RoutePath("GET", "/reviews"), Pagination(Authenticate(AuthorizeAdmin(h.GetAllReviews))))
	router.HandleFunc(utils.RoutePath("POST", "/products/{id}/reviews"), Authenticate(h.AddReview))
	router.HandleFunc(utils.RoutePath("PUT", "/products/{id}/reviews/{reviewId}"), Authenticate(h.EditReview))
	router.HandleFunc(utils.RoutePath("DELETE", "/products/{id}/reviews/{reviewId}"), Authenticate(h.DeleteReview))
}

func (h *Handler) GetAllReviews(w http.ResponseWriter, r *http.Request) {
	pagination := middlewares.GetPagination(r)

	conditions := utils.GetFilterConditions(r, whiteListedParams)
	sortString := utils.GetSortQ(r, whiteListedSortParams)

	reviews, count, err := utils.GenericFilterWithJoins[models.Review, types.GetAllReviewsRow](
		&utils.GenericFilterConfigWithJoins{
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
	userId, err:=utils.GetUserIdCtx(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, appErrors.ErrFailedToRetrieveToken)
		return
	}
	crPayload, err := utils.ValidateAndParseBody[payloads.CreateReview](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	productId, err := utils.GetValidateId(r, constants.IdUrlPathKey)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, appErrors.NewInvalidIDError("product",*productId))
		return
	}
	product, err := h.store.GetProductById(*productId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, appErrors.NewResourceWasNotFoundError("product", *productId))
		return
	}

	review := crPayload.TrimStrs().ToModel(*userId, product.ID)
	newRev, err := h.store.CreateReview(review)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{"review": newRev})
}

func (h *Handler) EditReview(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIdCtx(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, appErrors.ErrFailedToRetrieveToken)
		return
	}
	Id, err := utils.GetValidateId(r, reviewId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, invalidRevIdErr(*Id))
		return
	}

	upPayload, err := utils.ValidateAndParseBody[payloads.UpdateReview](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if upPayload.IsEmpty() {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("at least one of (comment, rate) is required"))
		return
	}

	model := upPayload.TrimStrs().ToModel()
	review, err := h.store.UpdateReview(*Id, *userId, model, upPayload)
	if err != nil {
		if errors.Is(err, appErrors.ErrForbidden) {
			utils.WriteError(w, http.StatusForbidden, err)
			return
		}
		var notFoundErr *appErrors.ResourceWasNotFoundError
		if errors.As(err, &notFoundErr){
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, map[string]any{"review": review})
}

func (h *Handler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	Id, err := utils.GetValidateId(r, reviewId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, invalidRevIdErr(*Id))
		return
	}

	userId, err := utils.GetUserIdCtx(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, appErrors.ErrFailedToRetrieveToken)
		return
	}

	err = h.store.HardDelete(*Id, *userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, map[string]any{})
}

