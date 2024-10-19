package product

// * this file better to be re-factored to ensure more readable queries and better structs names.
import (
	"fmt"

	"main.go/types"
)

var getProductByIdQ = ` 
		products.id as id,
		products.name as name,
		products.quantity as quantity,
		products.description as description,
		products.price as price,
		products.category_id as category_id,
		products.created_at as created_at,
		products.updated_at as updated_at,
		products.deleted_at,
		avg_rating.avg_rating,
		images.id as image_id,
		images.image_url as image_url,
		images.image_public_id as image_public_id,
		images.is_main as is_main,
		categories.name as category_name,
		reviews.review_id as review_id,
		reviews.comment as review_comment,
		reviews.rate as review_rate,
		reviews.user_id as review_user_id,
		reviews.created_at as review_created_at,
		reviews.updated_at as review_updated_at,
		users.name as user_name,
		users.email as user_email,
		users.avatar as user_avatar
	`
var getProductByIdJoins = `
LEFT JOIN images ON images.product_id = products.id
LEFT JOIN categories ON categories.id = products.category_id
LEFT JOIN (
	SELECT product_id, AVG(rate) AS avg_rating 
	FROM reviews 
	GROUP BY product_id)
 	avg_rating ON avg_rating.product_id = products.id
LEFT JOIN (SELECT
	 reviews.id as review_id,
	 reviews.comment,
	 reviews.rate,
	 reviews.user_id,
	 reviews.product_id,
	 reviews.created_at,
	 reviews.updated_at
 FROM
	 reviews
 WHERE
 	reviews.product_id = ?
 ORDER BY
 	reviews.rate DESC
 LIMIT 9) reviews ON reviews.product_id = products.id
 LEFT JOIN users ON reviews.user_id = users.id
 `
var groupByGetProductById = `products.id, images.id, review_id, users.id`

func convertRowsToProduct(rows []types.RowGetProductById) *types.RespGetOneProductShape{
	var product types.RespGetOneProductShape
	var images = make([]types.RowGetOneProductImage, 0)
	var reviews = make([]types.RespProductReviewShape, 0)
	IdsMap := map[string]any{}

	for _, row := range rows {
		_, exists := IdsMap[fmt.Sprintf("productId-?",row.ID)]
		if !exists {
			product = types.RespGetOneProductShape{
				ID: row.ID,
				Name: row.Name,
				Description: row.Description,
				Quantity: row.Quantity,
				Price: row.Price,
				CategoryName: row.CategoryName,
				AvgRating: row.AvgRating,
				CreatedAt: row.CreatedAt,
				UpdatedAt: row.UpdatedAt,
			}

			IdsMap[fmt.Sprintf("productId-?",row.ID)] = fmt.Sprintf("productId-?",row.ID)
		}

		_, imageExists := IdsMap[fmt.Sprintf("imageId-?",row.ImageId)]
		if !imageExists {
			images = append(images, types.RowGetOneProductImage{
				ImageId: row.ImageId,
				ImageUrl: row.ImageUrl,
				ImagePublicId: row.ImagePublicId,
				IsMain: row.IsMain,
			})

			IdsMap[fmt.Sprintf("imageId-?",row.ImageId)] = fmt.Sprintf("imageId-?",row.ImageId) 
		}

		_, reviewExist := IdsMap[fmt.Sprintf("reviewId-?",row.ImageId)]
		if !reviewExist {
			reviews = append(reviews, types.RespProductReviewShape{
				ReviewID: row.ReviewID,
				Comment: row.ReviewComment,
				Rate: row.ReviewRate,
				User: types.RowProductReviewUser{
					UserName: row.UserName,
					UserEmail: row.UserEmail,
					UserAvatar: row.UserAvatar,
				},
				UserID: row.ReviewUserID,
				CreatedAt: row.CreatedAt,
				UpdatedAt: row.UpdatedAt,
			})

			IdsMap[fmt.Sprintf("reviewId-?",row.ReviewID)] = fmt.Sprintf("reviewId-?",row.ReviewID) 
		}
		
	}
	product.Images = images
	product.Reviews = reviews
	return &product
}

var prodsSelectCols = `	products.id as id, products.name as name, products.quantity as quantity,
				products.description as description, products.category_id as category_id,
				products.price as price,products.created_at as created_at,
				products.updated_at as updated_at,
 				images.id as image_id ,images.image_url as image_url, images.image_public_id as image_public_id,
				AVG(reviews.rate) AS avg_rating
			`
var imagesJoin = "LEFT JOIN images on products.id = images.product_id AND images.is_main = 1"
var reviewsJoin = "LEFT JOIN reviews on products.id = reviews.product_id"
var prodsGroupBy = "products.id, images.id"

func convertRowsToResp(rows []types.GetAllProductsRow) []types.RespGetAllProductsShape {
	productMap := make(map[uint]uint)
	productsSlice := make([]*types.RespGetAllProductsShape, 0)

	for _, row := range rows {
		if _, exists := productMap[row.Id]; !exists {
			productMap[row.Id] = row.Id
			resAllProducts := &types.RespGetAllProductsShape{
				Id:          row.Id,
				Name:        row.Name,
				Quantity:    row.Quantity,
				Description: row.Description,
				CategoryId:  row.CategoryId,
				Price:       row.Price,
				CreatedAt:   row.CreatedAt,
				UpdatedAt:   row.UpdatedAt,
				Image:       types.RowGetAllProductsImage{},
				AvgRating:   row.AvgRating,
			}
			productsSlice = append(productsSlice, resAllProducts)
		}

		if row.ImageId != 0 {
			productsSlice[len(productsSlice)-1].Image = types.RowGetAllProductsImage{
				ImageId:       row.ImageId,
				ImageUrl:      row.ImageUrl,
				ImagePublicId: row.ImagePublicId,
			}
		}

	}

	var result []types.RespGetAllProductsShape
	for _, product := range productsSlice {
		result = append(result, *product)
	}

	return result
}
