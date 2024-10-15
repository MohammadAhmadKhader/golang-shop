package order

import (
	"net/http"

	"gorm.io/gorm"
)

func Setup(DB *gorm.DB, router *http.ServeMux) {
	store := NewStore(DB)
	handler := NewHandler(*store)
	handler.RegisterRoutes(router)
}
