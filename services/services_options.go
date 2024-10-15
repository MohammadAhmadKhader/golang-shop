package services

import (
	"main.go/middlewares"
	"main.go/services/generic"
)

var superAdminOnlyRO = generic.NewRouteOptions(&generic.RouteOptions{
	IsEnabled:      true,
	AuthenticateMiddleware: middlewares.Authenticate,
	AuthorizeMiddleware: &middlewares.AuthorizeSuperAdmin,
})

var adminRO = generic.NewRouteOptions(&generic.RouteOptions{
	IsEnabled:      true,
	AuthenticateMiddleware: middlewares.Authenticate,
	AuthorizeMiddleware: &middlewares.AuthorizeAdmin,
}) 

var adminRoles = generic.NewOptions(&generic.Options{
	SoftDeleteRoutes: adminRO,
	HardDelete:       adminRO,
})

var superAdminOpts = generic.NewOptions(&generic.Options{
	SoftDeleteRoutes: superAdminOnlyRO,
	HardDelete:       superAdminOnlyRO,
})
