package utils

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"gorm.io/gorm"
	"main.go/internal/database"
	"main.go/types"
)

func GenericFilter[T any](config *GenericFilterConfig) (result []T, count int64, errors []error) {
	var err1_Mu sync.Mutex
	var err2_Mu sync.Mutex
	var results []T

	page := config.Pagination.Page
	limit := config.Pagination.Limit
	query := config.DB.Model(new(T))

	for _, filter := range config.Filters {
		if config.WhiteListedParams[filter.Field] != nil {
			condition := fmt.Sprintf("%s %s ?", filter.Field, filter.Operator)
			query = query.Where(condition, filter.Value)
		}
	}

	var wg sync.WaitGroup
	wg.Add(2)
// does not applu sort
	errors = make([]error, 0)
	
	go func() {
		defer wg.Done()
		offset := CalculateOffset(page, limit)
		clonedQuery := query.Session(&gorm.Session{}).Order(config.SortQ).Offset(offset).Limit(limit)
		
		for _, preload := range config.Preloads {
			clonedQuery = clonedQuery.Preload(preload)
		}	
		if err := clonedQuery.Find(&results).Error; err != nil {
			err1_Mu.Lock()
			errors = append(errors, err)
			err1_Mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		clonedQuery := query.Session(&gorm.Session{})
		if err := clonedQuery.Count(&count).Error; err != nil {
			err2_Mu.Lock()
			errors = append(errors, err)
			err2_Mu.Unlock()
		}
	}()

	wg.Wait()

	return results, count, errors
}

func GenericFilterWithJoins[TModel any,TRow any](config *GenericFilterConfigWithJoins) (result []TRow, count int64, errors []error) {
	var err1_Mu sync.Mutex
	var err2_Mu sync.Mutex
	var results []TRow

	DB := database.DB

	page := config.Pagination.Page
	limit := config.Pagination.Limit
	query := DB.Model(new(TModel))

	for _, filter := range config.Filters {
		if config.WhiteListedParams[filter.Field] != nil {
			condition := fmt.Sprintf("%s %s ?", filter.Field, filter.Operator)
			query = query.Where(condition, filter.Value)
		}
	}
	
	var wg sync.WaitGroup
	wg.Add(2)

	errors = make([]error, 0)
	
	go func() {
		defer wg.Done()
		offset := CalculateOffset(page, limit)
		clonedQuery := query.Session(&gorm.Session{})
		clonedQuery = clonedQuery.Select(config.SelectQ)
		for _, join := range config.Joins {
			clonedQuery = clonedQuery.Joins(join)
		}
		if config.Group != nil {
			clonedQuery = clonedQuery.Group(*config.Group)
		}
		err := clonedQuery.Order(config.SortQ).Offset(offset).Limit(limit).Scan(&results).Error

		if err != nil {
			err1_Mu.Lock()
			errors = append(errors, err)
			err1_Mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		clonedQuery := query.Session(&gorm.Session{})
		if err := clonedQuery.Count(&count).Error; err != nil {
			err2_Mu.Lock()
			errors = append(errors, err)
			err2_Mu.Unlock()
		}
	}()

	wg.Wait()

	return results, count, errors
}

var whiteListedOperators = map[string]any{
	"gt":  ">",
	"gte": ">=",
	"lt":  "<",
	"lte": "<=",
}

// This function handles the given params values as long as they not received in array, if received in array only takes first value.
func GetFilterConditions(r *http.Request, whiteListedParams map[string]any) []types.FilterCondition {
	params := r.URL.Query()
	var conditions = make([]types.FilterCondition, len(params))
	for key, values := range params {
		field, op := GetFieldOperator(key, whiteListedParams)
		if field != "" {
			conditions = append(conditions, types.FilterCondition{
				Field:    field,
				Value:    values[0],
				Operator: op,
			})
		}
	}

	return conditions
}

// key could price_lte as example or price, the param example price_lte=30, which translates to price <= 30.
//
// if the key is not white listed then it returns both the field and operator as empty strings.
func GetFieldOperator(key string, whiteListedParams map[string]any) (field string, operator string) {
	var cond []string
	if strings.Contains(key, "_") && len(strings.Split(key, "_")) == 2 && whiteListedParams[key] != nil {
		cond = strings.Split(key, "_")
		if whiteListedOperators[cond[1]] != nil {
			key = cond[0]
			return key, whiteListedOperators[cond[1]].(string)
		}

		return key, "="
	}

	return "", ""
}

type GenericFilterConfig struct {
	DB *gorm.DB
	Filters []types.FilterCondition
	SortQ string
	Pagination types.Pagination
	WhiteListedParams map[string]any
	Preloads []string
}

type GenericFilterConfigWithJoins struct {
	Filters []types.FilterCondition
	SortQ string
	Pagination types.Pagination
	WhiteListedParams map[string]any
	Joins []string
	SelectQ string
	Group *string
}
