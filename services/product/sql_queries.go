package product
// * this file better to be re-factored to ensure more readable queries and better structs names.
import (
	"fmt"
	"time"
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
 
 type getOneRespShape struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Description *string `json:"description"`
	Quantity uint `json:"quantity"`
	Price float64 `json:"price"`
	CategoryName string `json:"category"`
	Category rowCategory `json:"-"`
	AvgRating float64 `json:"avgRating"`
	Images []rowImageGetOneProduct `json:"images"`
	Reviews []reviewShape `json:"reviews"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
type rowGetOneById struct {
	ID uint
	Name string
	Description *string
	Quantity uint
	Price float64
	CategoryID uint
	AvgRating float64
	CreatedAt time.Time
	UpdatedAt time.Time
	ReviewID uint
	ReviewUserID uint
	ReviewComment string
	ReviewRate uint8
	ReviewCreatedAt time.Time
	ReviewUpdatedAt time.Time
	CategoryName string
	ImageId       uint   
	ImageUrl      string 
	ImagePublicId string 
	IsMain bool
	UserName string
	UserEmail string
	UserAvatar string
}
type reviewShape struct {
	ReviewID uint `json:"id"`
	UserID uint `json:"userId"`
	User rowUser `json:"user"`
	Comment string `json:"comment"`
	Rate uint8 `json:"rate"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type rowImageGetOneProduct struct {
	ImageId       uint   `json:"id"`
	ImageUrl      string `json:"imageUrl"`
	ImagePublicId string `json:"imagePublicId"`
	IsMain bool `json:"isMain"`
}

type rowCategory struct {
	CategoryName string
}

type rowUser struct {
	UserName string `json:"name"`
	UserEmail string `json:"email"`
	UserAvatar string `json:"avatar"`
}


func convertRowsToProduct(rows []rowGetOneById) *getOneRespShape{
	var product getOneRespShape
	var images = make([]rowImageGetOneProduct, 0)
	var reviews = make([]reviewShape, 0)
	IdsMap := map[string]any{}

	for _, row := range rows {
		_, exists := IdsMap[fmt.Sprintf("productId-?",row.ID)]
		if !exists {
			product = getOneRespShape{
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
			images = append(images, rowImageGetOneProduct{
				ImageId: row.ImageId,
				ImageUrl: row.ImageUrl,
				ImagePublicId: row.ImagePublicId,
				IsMain: row.IsMain,
			})

			IdsMap[fmt.Sprintf("imageId-?",row.ImageId)] = fmt.Sprintf("imageId-?",row.ImageId) 
		}

		_, reviewExist := IdsMap[fmt.Sprintf("reviewId-?",row.ImageId)]
		if !reviewExist {
			reviews = append(reviews, reviewShape{
				ReviewID: row.ReviewID,
				Comment: row.ReviewComment,
				Rate: row.ReviewRate,
				User: rowUser{
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

type GetAllProductsRow struct {
	Id            uint
	Name          string
	Quantity      uint
	Description   string
	CategoryId    uint
	Price         float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ImageId       uint
	ImageUrl      string
	ImagePublicId string
	AvgRating     float64
}

type respAllProds struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Quantity    uint      `json:"quantity"`
	Description string    `json:"description"`
	CategoryId  uint      `json:"categoryId"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Image       rowImage  `json:"mainImage"`
	AvgRating   float64   `json:"avgRating"`
}
type rowImage struct {
	ImageId       uint   `json:"id"`
	ImageUrl      string `json:"imageUrl"`
	ImagePublicId string `json:"imagePublicId"`
}

func convertRowsToResp(rows []GetAllProductsRow) []respAllProds {
	productMap := make(map[uint]uint)
	productsSlice := make([]*respAllProds, 0)

	for _, row := range rows {
		if _, exists := productMap[row.Id]; !exists {
			productMap[row.Id] = row.Id
			resAllProducts := &respAllProds{
				Id:          row.Id,
				Name:        row.Name,
				Quantity:    row.Quantity,
				Description: row.Description,
				CategoryId:  row.CategoryId,
				Price:       row.Price,
				CreatedAt:   row.CreatedAt,
				UpdatedAt:   row.UpdatedAt,
				Image:       rowImage{},
				AvgRating:   row.AvgRating,
			}
			productsSlice = append(productsSlice, resAllProducts)
		}

		if row.ImageId != 0 {
			productsSlice[len(productsSlice)-1].Image = rowImage{
				ImageId:       row.ImageId,
				ImageUrl:      row.ImageUrl,
				ImagePublicId: row.ImagePublicId,
			}
		}

	}

	var result []respAllProds
	for _, product := range productsSlice {
		result = append(result, *product)
	}

	return result
}
