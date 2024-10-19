package product

import (
	"fmt"
	"net/http"

	"main.go/constants"
	"main.go/errors"
	"main.go/middlewares"
	"main.go/pkg/models"
	"main.go/pkg/payloads"
	"main.go/pkg/utils"
	"main.go/types"
)

type Handler struct {
	store Store
}

func NewHandler(store Store) *Handler {
	return &Handler{
		store: store,
	}
}

func invalidProdIdErr(id uint) error {
	return errors.NewInvalidIDError("product", id)
}

var Authenticate = middlewares.Authenticate
var AuthorizeAdmin = middlewares.AuthorizeAdmin
var Pagination = middlewares.PaginationMiddleware

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc(utils.RoutePath("GET", "/products/{id}"), Authenticate(AuthorizeAdmin(h.GetProductById)))
	router.HandleFunc(utils.RoutePath("GET", "/products"),Pagination(Authenticate(AuthorizeAdmin(h.GetAllProducts))))
	router.HandleFunc(utils.RoutePath("POST", "/products"), Authenticate(AuthorizeAdmin(h.CreateProduct)))
	router.HandleFunc(utils.RoutePath("PUT", "/products/{id}"), Authenticate(AuthorizeAdmin(h.UpdateProduct)))
}

func (h *Handler) GetProductById(w http.ResponseWriter, r *http.Request) {
	Id, err := utils.GetValidateId(r, constants.IdUrlPathKey)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, invalidProdIdErr(*Id))
		return
	}

	productRows, err := h.store.GetProductById(*Id)
	if err != nil || len(productRows) == 0 {
		utils.WriteError(w, http.StatusBadRequest, errors.NewResourceWasNotFoundError("product", *Id))
		return
	}
	product := convertRowsToProduct(productRows)

	utils.WriteJSON(w, http.StatusOK, map[string]any{"product": product})
}

func (h *Handler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	pagination := middlewares.GetPagination(r)

	conditions := utils.GetFilterConditions(r, whiteListedParams)
	sortString := utils.GetSortQ(r, whiteListedSortParams)
	rows, count, err := utils.GenericFilterWithJoins[models.Product, types.GetAllProductsRow](&utils.GenericFilterConfigWithJoins{
		DB:                h.store.DB,
		Filters:           conditions,
		SortQ:             sortString,
		Pagination:        pagination,
		WhiteListedParams: whiteListedParams,
		SelectQ:           prodsSelectCols,
		Joins:             []string{imagesJoin, reviewsJoin},
		Group:             &prodsGroupBy,
	})

	if len(err) != 0 {
		utils.WriteError(w, http.StatusBadRequest, err[0])
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"count":    count,
		"page":     pagination.Page,
		"limit":    pagination.Limit,
		"products": convertRowsToResp(rows),
	})
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	crPayload, err := utils.ValidateAndParseFormData(r, func() (*payloads.CreateProduct, error) {
		payload ,err := payloads.NewCreatePayload(r, utils.ConvertStrToUint,utils.ConvertStrToFloat64)
		if err != nil {
			return nil, err
		}
		
		return payload, nil
	})

	
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	imgsHandler := utils.NewImagesHandler()
	resp, err := imgsHandler.UploadOne(&file, fileHeader, types.ProductsFolder, r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	prod := crPayload.TrimStrs().ToModelWithImage(resp.SecureUrl)
	product, err := h.store.CreateProductWithImage(prod, resp)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"product": product,
	})
}

func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	Id, err := utils.GetValidateId(r, constants.IdUrlPathKey)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, invalidProdIdErr(*Id))
		return
	}

	upPayload, err := utils.ValidateAndParseBody[payloads.UpdateProduct](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if upPayload.IsEmpty() {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("at least one of (name, quantity, description, categoryId, price) is required"))
		return
	}

	prod := upPayload.TrimStrs().ToModel()
	product, err := h.store.UpdateProduct(*Id, prod, upPayload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, map[string]any{
		"product":*product,
	})
}
