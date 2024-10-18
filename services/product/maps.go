package product

var whiteListedParams = map[string]any{
	"name":          1,
	"price_lte":     1,
	"price_gte":     1,
	//"avg_rating_lte": 1,
	//"avg_rating_gte": 1,
	//"avgRating":     1,
	"price":         1,
	"quantity":      1,
	"categoryId":    1,
}

var whiteListedSortParams = map[string]any{
	"created_at": 1,
	"updated_at": 1,
	"avg_rating":  1, // only works with small eltters
	"price":      1,
	"quantity":   1,
}
