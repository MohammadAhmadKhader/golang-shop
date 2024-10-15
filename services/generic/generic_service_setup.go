package generic

import (
	"net/http"

	"gorm.io/gorm"
)

func Setup[TModel any](DB *gorm.DB, router *http.ServeMux, modelName string, options Options) {
	store := NewStore[TModel](DB)
	handler := NewHandler[TModel](*store)
	handler.RegisterRoutesGeneric(router, modelName, options)
}
