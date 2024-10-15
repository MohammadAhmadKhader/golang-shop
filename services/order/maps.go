package order

var whiteListedParams = map[string]any{
	"user_id":          1,
	"status":     1,
	"total_price":     1,

}

var whiteListedSortParams = map[string]any{
	"created_at": 1,
	"updated_at": 1,
	"total_price":  1,
	"status":   1,
}