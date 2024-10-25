package utils

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
)

// handles sort queries, they expected to be received as sort=name,dir, default sort is descending by created_at, it only accepts one sort.
func GetSortQ(r *http.Request, whiteListedParams map[string]any) string {
	defaultSort := "created_at DESC"

	values := r.URL.Query()
	sortStr := values.Get("sort")
	sortArr := strings.Split(sortStr, ",")
	
	if len(sortArr) != 2 {
		return defaultSort
	}

	sortField := sortArr[0]
	sortDir := sortArr[1]

	if sortDir != "desc" && sortDir != "asc" {
		return defaultSort
	}
	if whiteListedParams[sortField] == nil {
		return defaultSort
	}

	return fmt.Sprintf(sortField + " " + sortDir)
}

// handles sort queries, they expected to be received as sort=name,dir, default sort is descending by created_at, it only accepts more than one sort.
// 
// this function is meant to be used with dashboards.
func GetSortQArr(r *http.Request, whiteListedParams map[string]any) string {
	defaultSort := "created_at DESC"

	sortsArr := r.URL.Query()["sort"]
	if len(sortsArr) == 0 {
		return defaultSort
	}
	
	sorts := make([]string, 0)
	sortedFields := make([]string, 0)
	for _, singleSort := range sortsArr {
		if len(singleSort) != 2 {
			return defaultSort
		}

		singleSortArr := strings.Split(singleSort, ",")
		sortField := singleSortArr[0]
		sortDir := singleSortArr[1]

		if sortDir != "desc" && sortDir != "asc" {
			return defaultSort
		}

		
		if whiteListedParams[sortField] == nil && !slices.Contains(sortedFields, sortField) {
			break
		}
		
		sorts = append(sorts, fmt.Sprintf("%s %s", sortField, sortDir))
		sortedFields = append(sortedFields, sortField)
	}

	if len(sorts) > 0 {
		return strings.Join(sorts, ",")
	}

	return defaultSort
}
