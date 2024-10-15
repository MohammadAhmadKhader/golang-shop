package generic

import (
	"net/http"
)

// Soft Delete routes includes: Get soft deleted, Soft delete by id, Restore soft deleted by id
type Options struct {
	SoftDeleteRoutes RouteOptions
	HardDelete RouteOptions
}

type RouteOptions struct {
	IsEnabled      bool
	AuthenticateMiddleware func(next http.HandlerFunc) http.HandlerFunc
	AuthorizeMiddleware  *func(next http.HandlerFunc) http.HandlerFunc
}

func ifNilUseDefault[TOption any](input *TOption, defaultValue TOption) TOption {
	if input == nil {
		return defaultValue
	}
	return *input
}

func NewRouteOptions(options *RouteOptions) RouteOptions {
	defaultRouteOptions := RouteOptions{
		IsEnabled:      false,
		AuthenticateMiddleware: nil,
		AuthorizeMiddleware: options.AuthorizeMiddleware,
	}
	
	return RouteOptions{
		IsEnabled: ifNilUseDefault(&options.IsEnabled, defaultRouteOptions.IsEnabled),
		AuthenticateMiddleware: ifNilUseDefault(&options.AuthenticateMiddleware, defaultRouteOptions.AuthenticateMiddleware),
		AuthorizeMiddleware: options.AuthorizeMiddleware,
	}
}

func NewOptions(options *Options) *Options {
	return &Options{
		SoftDeleteRoutes: options.SoftDeleteRoutes,
		HardDelete: options.SoftDeleteRoutes,
	}
}