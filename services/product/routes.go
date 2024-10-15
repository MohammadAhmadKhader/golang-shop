package product

import (
	"fmt"
	"net/http"

	"main.go/constants"
	"main.go/middlewares"
	"main.go/pkg/models"
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
var AuthorizeAdmin = middlewares.AuthorizeAdmin

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc(utils.RoutePath("GET", "/products/{id}"), Authenticate(AuthorizeAdmin(h.GetProductById)))
	router.HandleFunc(utils.RoutePath("GET", "/products"), Authenticate(AuthorizeAdmin(h.GetAllProducts)))
	router.HandleFunc(utils.RoutePath("POST", "/products"), Authenticate(AuthorizeAdmin(h.CreateProduct)))
	router.HandleFunc(utils.RoutePath("PUT", "/products/{id}"), Authenticate(AuthorizeAdmin(h.CreateProduct)))
}

func (h *Handler) GetProductById(w http.ResponseWriter, r *http.Request) {
	Id, err := utils.GetValidateId(r, constants.IdUrlPathKey)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}

	product, err := h.store.GetProductById(*Id)
	responseMap := getProductByIdMap(product)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{"product": responseMap})
}

func (h *Handler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	pagination := middlewares.GetPagination(r)

	conditions := utils.GetFilterConditions(r, whiteListedParams)
	sortString := utils.GetSortQ(r, whiteListedSortParams)
	rows, count, err := utils.GenericFilterWithJoins[models.Product, GetAllProductsRow](&utils.GenericFilterConfigWithJoins{
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
	resp, err := imgsHandler.UploadOne(&file, fileHeader, utils.ProductsFolder, r.Context())
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
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	urPayload, err := utils.ValidateAndParseFormData(r, func() (*payloads.UpdateProduct, error) {
		payload, err := payloads.NewUpdatePayload(r, utils.ConvertStrToUint, utils.ConvertStrToFloat64)
		if err != nil{
			return nil, err
		}

		return payload, nil
	})

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	prod := urPayload.TrimStrs().ToModel()
	product, err := h.store.UpdateProduct(*Id, prod, urPayload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, map[string]any{
		"product":*product,
	})
}
