package review

var whiteListedParams = map[string]any{
	"userId":    1,
	"rate_lte":  1,
	"rate_gte":  1,
	"rate":      1,
	"comment":   1,
	"productId": 1,
}

var whiteListedSortParams = map[string]any{
	"comment":    1,
	"rate":       1,
	"created_at": 1,
}