package generic

import (
	"fmt"
	"net/http"

	"gorm.io/gorm"
	"main.go/constants"
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
	"categories":"category",
	"products":"product",
	"roles":"role",
	"reviews":"review",
}

func (h *Handler[TModel]) RegisterRoutesGeneric(router *http.ServeMux, modelName string, options Options) {
	
	if options.SoftDeleteRoutes.IsEnabled {
		_ = options.SoftDeleteRoutes.AuthenticateMiddleware
		if options.HardDelete.AuthorizeMiddleware != nil {
			AuthorizeMW := *options.HardDelete.AuthorizeMiddleware
			router.HandleFunc(utils.RoutePath("GET","/"+modelName+"/deleted"), middlewares.Authenticate(AuthorizeMW(h.GenerateGetAllDeleted(modelName))))
			router.HandleFunc(utils.RoutePath("PATCH","/"+modelName+"/{id}/restore"), middlewares.Authenticate(AuthorizeMW(h.GenerateRestoreRoute(modelName))))
			router.HandleFunc(utils.RoutePath("DELETE","/"+modelName+"/{id}/soft-delete"), middlewares.Authenticate(AuthorizeMW(h.GenerateSoftDeleteRoute(modelName))))
		} else {
			router.HandleFunc(utils.RoutePath("GET","/"+modelName+"/deleted"), middlewares.Authenticate(h.GenerateGetAllDeleted(modelName)))
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
		Id, err := utils.GetValidateId(r, constants.IdUrlPathKey)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest , fmt.Errorf("invalid id"))
			return
		}

		notFoundErr := fmt.Errorf("%v with id: '%v' was not found",ModelNameMapper[modelName],*Id)

		err = h.store.Generic.SoftDelete(*Id, notFoundErr)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest , err)
			return
		}
	
		utils.WriteJSON(w, http.StatusNoContent , map[string]any{})
	}
}

func (h *Handler[TModel]) GenerateRestoreRoute(modelName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Id, err := utils.GetValidateId(r, constants.IdUrlPathKey)
		if err != nil {
			utils.WriteJSON(w, http.StatusBadRequest , fmt.Errorf("invalid id"))
			return
		}

		notFoundErr := fmt.Errorf("%v with id: '%v' was not found",ModelNameMapper[modelName],*Id)

		item, err := h.store.Generic.Restore(*Id, notFoundErr)
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
		Id, err := utils.GetValidateId(r, constants.IdUrlPathKey)
		if err != nil {
			utils.WriteJSON(w, http.StatusBadRequest , fmt.Errorf("invalid id"))
			return
		}

		notFoundErr := fmt.Errorf("%v with id: '%v' was not found",ModelNameMapper[modelName],*Id)

		err = h.store.Generic.HardDelete(*Id, notFoundErr)
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

//for _, model := range models {
//	var typed = interface{}(model).(types.SoftDeletable)
//	bytes, err := typed.ShowDeletedAt()
//	if err != nil {
//		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to parse deleted_at"))
//		return
//	}
//	fmt.Println("DeletedAt JSON:", string(bytes))
//}