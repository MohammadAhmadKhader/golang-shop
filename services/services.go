package services

import (
	"net/http"

	"gorm.io/gorm"
	"main.go/pkg/models"
	"main.go/services/address"
	"main.go/services/cart"
	"main.go/services/category"
	"main.go/services/generic"
	"main.go/services/image"
	"main.go/services/order"
	"main.go/services/product"
	"main.go/services/review"
	"main.go/services/role"
	"main.go/services/user"
)

func SetupAllServices(DB *gorm.DB, router *http.ServeMux) {
	category.Setup(DB, router)
	generic.Setup[models.Category](DB, router, "categories", *adminRoles)
	product.Setup(DB, router)
	generic.Setup[models.Product](DB, router, "products", *adminRoles)

	image.Setup(DB, router)
	order.Setup(DB, router)
	user.Setup(DB, router)
	generic.Setup[models.User](DB, router, "users", *superAdminOpts)
	review.Setup(DB, router)
	
	cart.Setup(DB,router)
	address.Setup(DB, router)
	
	role.Setup(DB, router)
	generic.Setup[models.Role](DB, router, "roles", *superAdminOpts)
}