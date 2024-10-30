package generic

import (
	"net/http"

	"gorm.io/gorm"
	"main.go/constants"
	"main.go/errors"
	"main.go/middlewares"
	"main.go/pkg/utils"
)

type Store[TModel any] struct {
	DB      *gorm.DB
	Generic *GenericRepository[TModel]
}

func NewStore[TModel any](DB *gorm.DB) *Store[TModel] {
	return &Store[TModel]{
		DB:      DB,
		Generic: &GenericRepository[TModel]{DB: DB},
	}
}

type Handler[TModel any] struct {
	store Store[TModel]
}

func NewHandler[TModel any](store Store[TModel]) *Handler[TModel] {
	return &Handler[TModel]{
		store: store,
	}
}

var ModelNameMapper = map[string]string{
	"users":"user",
	"categories":"category",
	"products":"product",
	"roles":"role",
	"reviews":"review",
}

func invalidModelIdErr(modelName string,id string) error {
	return errors.NewInvalidIDError(ModelNameMapper[modelName], id)
}
var Pagination = middlewares.PaginationMiddleware
func (h *Handler[TModel]) RegisterRoutesGeneric(router *http.ServeMux, modelName string, options Options) {
	
	if options.SoftDeleteRoutes.IsEnabled {
		_ = options.SoftDeleteRoutes.AuthenticateMiddleware
		if options.HardDelete.AuthorizeMiddleware != nil {
			AuthorizeMW := *options.HardDelete.AuthorizeMiddleware
			router.HandleFunc(utils.RoutePath("GET","/"+modelName+"/deleted"), Pagination(middlewares.Authenticate(AuthorizeMW(h.GenerateGetAllDeleted(modelName)))))
			router.HandleFunc(utils.RoutePath("PATCH","/"+modelName+"/{id}/restore"), middlewares.Authenticate(AuthorizeMW(h.GenerateRestoreRoute(modelName))))
			router.HandleFunc(utils.RoutePath("DELETE","/"+modelName+"/{id}/soft-delete"), middlewares.Authenticate(AuthorizeMW(h.GenerateSoftDeleteRoute(modelName))))
		} else {
			router.HandleFunc(utils.RoutePath("GET","/"+modelName+"/deleted"), Pagination(middlewares.Authenticate(h.GenerateGetAllDeleted(modelName))))
			router.HandleFunc(utils.RoutePath("PATCH","/"+modelName+"/{id}/restore"), middlewares.Authenticate(h.GenerateRestoreRoute(modelName)))
			router.HandleFunc(utils.RoutePath("DELETE","/"+modelName+"/{id}/soft-delete"), middlewares.Authenticate(h.GenerateSoftDeleteRoute(modelName)))
		}
	}
	
	if options.HardDelete.IsEnabled  {
		AuthenticateMW := options.SoftDeleteRoutes.AuthenticateMiddleware
		if options.HardDelete.AuthorizeMiddleware != nil {
			router.HandleFunc(utils.RoutePath("DELETE","/"+modelName+"/{id}"), AuthenticateMW(h.GenerateHardDeleteRoute(modelName)))
		} else {
			AuthorizeMW := *options.HardDelete.AuthorizeMiddleware
			router.HandleFunc(utils.RoutePath("DELETE","/"+modelName+"/{id}"), AuthenticateMW(AuthorizeMW(h.GenerateHardDeleteRoute(modelName))))
		}
	}
}

func (h *Handler[TModel]) GenerateSoftDeleteRoute(modelName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Id, receivedStr, err := utils.GetValidateId(r, constants.IdUrlPathKey)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest , invalidModelIdErr(modelName ,receivedStr))
			return
		}

		notFoundMsg := ModelNameMapper[modelName]+" "+"with id: '%v' was not found"

		err = h.store.Generic.SoftDelete(*Id, notFoundMsg)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest , err)
			return
		}
	
		utils.WriteJSON(w, http.StatusNoContent , map[string]any{})
	}
}

func (h *Handler[TModel]) GenerateRestoreRoute(modelName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Id, receivedStr, err := utils.GetValidateId(r, constants.IdUrlPathKey)
		if err != nil {
			utils.WriteJSON(w, http.StatusBadRequest ,  invalidModelIdErr(modelName ,receivedStr))
			return
		}

		notFoundMsg := ModelNameMapper[modelName]+" deleted "+"with id: '%v' was not found"

		item, err := h.store.Generic.Restore(*Id, notFoundMsg)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest , err)
			return
		}

		utils.WriteJSON(w, http.StatusOK , map[string]any{
			"message":"item was restored successfully",
			ModelNameMapper[modelName]:item,
		})
	}
}

func (h *Handler[TModel]) GenerateHardDeleteRoute(modelName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Id, receivedStr, err := utils.GetValidateId(r, constants.IdUrlPathKey)
		if err != nil {
			utils.WriteJSON(w, http.StatusBadRequest , invalidModelIdErr(modelName ,receivedStr))
			return
		}

		notFoundMsg := ModelNameMapper[modelName]+" "+"with id: '%v' was not found"

		err = h.store.Generic.HardDelete(*Id, notFoundMsg)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest , err)
			return
		}

		utils.WriteJSON(w, http.StatusNoContent , map[string]any{})
	}
}

func (h *Handler[TModel]) GenerateGetAllDeleted(modelName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination := middlewares.GetPagination(r)

		models, count, errors := h.store.Generic.GetAllDeleted(pagination.Page, pagination.Limit)
		if len(errors) != 0 {
			utils.WriteError(w, http.StatusBadRequest, errors[0])
			return
		}
	
		utils.WriteJSON(w, http.StatusOK, map[string]any{
			"page":     pagination.Page,
			"limit":    pagination.Limit,
			"count":    count,
			modelName: models,
		})
	}
}