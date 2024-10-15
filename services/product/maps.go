package product

import (
	"main.go/pkg/models"
)

// this function create the response shape for product route - Get product by id
func getProductByIdMap(product *models.Product) map[string]any {

	images := make([]map[string]any, len(product.Images))
	for i, img := range product.Images {
		images[i] = map[string]any{
			"id":            img.ID,
			"imageUrl":      img.ImageUrl,
			"isMain":        img.IsMain,
			"imagePublicId": img.ImagePublicId,
		}
	}

	reviews := make([]map[string]any, len(product.Reviews))
	for i, review := range product.Reviews {
		reviews[i] = map[string]any{
			"id":      review.ID,
			"rating":  review.Rate,
			"comment": review.Comment,
			"userId":  review.UserID,
		}
	}

	return map[string]any{
		"id":          product.ID,
		"createdAt":   product.CreatedAt,
		"updatedAt":   product.UpdatedAt,
		"name":        product.Name,
		"description": product.Description,
		"quantity":    product.Quantity,
		"price":       product.Price,
		"avgRating":   product.AvgRating,
		"images":      images,
		"reviews":     reviews,
	}
}

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
