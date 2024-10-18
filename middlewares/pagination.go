package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"main.go/pkg/utils"
	"main.go/types"
)

type contextKey string

const (
	minLimit                 = 3
	maxLimit                 = 30
	defaultLimit             = 9
	paginationKey contextKey = "pagination"
)

func PaginationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")
		
		page := pageHandler(pageStr)
		limit := limitHandler(limitStr)
		fmt.Println(page,limit)
		
		ctx := context.WithValue(r.Context(), paginationKey, &types.Pagination{Page: page, Limit: limit})
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func limitHandler(limitAsString string) int {
	limit, err := strconv.Atoi(limitAsString)
	if err != nil {
		return defaultLimit
	}
	if limit < minLimit {
		return minLimit
	}
	if limit > maxLimit {
		return maxLimit
	}

	return limit
}

func pageHandler(pageAsString string) int {
	page, err := strconv.Atoi(pageAsString)
	
	if err != nil || page < 1 {
		return 1
	}
	
	return page
}

func GetPagination(r *http.Request) types.Pagination {
	pagination, ok := r.Context().Value(paginationKey).(*types.Pagination)

	if !ok {
		return types.Pagination{
			Page: 1,
			Limit: defaultLimit,
		}
	}

	return *pagination
}

func CalculateOffset(pagination types.Pagination) int {
	return utils.CalculateOffset(pagination.Page, pagination.Limit)
}