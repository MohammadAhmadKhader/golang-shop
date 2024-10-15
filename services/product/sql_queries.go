package product

import (
	"time"

	"main.go/pkg/models"
)

// requires 2 arguments both of them are product's id
var qGetProductById = `
		SELECT 
		    p.id,
		    p.name,
		    p.quantity,
		    p.description,
		    p.price,
		    avg_rating.avg_rating,
		    i.id,
		    i.image_url,
			i.image_public_id,
		    reviews.review_id,
		    reviews.comment,
		    reviews.rate,
		    reviews.user_id
		FROM 
		    products p
		LEFT JOIN 
		    images i ON i.product_id = p.id
		LEFT JOIN 
		    (SELECT 
		         product_id, 
		         AVG(rate) AS avg_rating 
		     FROM 
		         reviews 
		     GROUP BY 
		         product_id) avg_rating ON avg_rating.product_id = p.id
		LEFT JOIN 
		    (SELECT 
		         r.id AS review_id,
		         r.comment,
		         r.rate,
		         r.user_id,
		         r.product_id
		     FROM 
		         reviews r 
		     WHERE 
		         r.product_id = ?
		     ORDER BY 
		         r.rate DESC 
		     LIMIT 9) reviews ON reviews.product_id = p.id
		WHERE 
		    p.id = ?
		GROUP BY 
		    p.id, i.id, reviews.review_id; 
`

func scanGetProductById(product *models.Product, image *models.Image, review *models.Review) (id *uint, name *string, qty *uint, desc **string, price *float64, avgRating **float64, imageId *uint, imageUrl *string, imagePublicId *string, reviewId *uint, comment *string, rate *uint8, userId *uint) {
	return &product.ID, &product.Name, &product.Quantity, &product.Description, &product.Price, &product.AvgRating, &image.ID, &image.ImageUrl, &image.ImagePublicId, &review.ID, &review.Comment, &review.Rate, &review.UserID
}

func scanProduct(product *models.Product) (id *uint, name *string, qty *uint, desc *string, price *float64, avgRating *float64) {
	return &product.ID, &product.Name, &product.Quantity, product.Description, &product.Price, product.AvgRating
}

func scanReview(review *models.Review) (reviewId *uint, comment *string, rate *uint8, userId *uint) {
	return &review.ID, &review.Comment, &review.Rate, &review.UserID
}

func scanImage(image *models.Image) (imageId *uint, imageUrl *string, imagePublicId *string) {
	return &image.ID, &image.ImageUrl, &image.ImagePublicId
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

// mapping mechanism to ensure the duplicates wont be included

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
	Id          uint       `json:"id"`
	Name        string     `json:"name"`
	Quantity    uint       `json:"quantity"`
	Description string     `json:"description"`
	CategoryId  uint       `json:"categoryId"`
	Price       float64    `json:"price"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	Images      []rowImage `json:"images"`
	AvgRating   float64    `json:"avgRating"`
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
				Images:      []rowImage{},
				AvgRating:   row.AvgRating,
			}
			productsSlice = append(productsSlice, resAllProducts)
		}

		if row.ImageId != 0 {
			productsSlice[len(productsSlice) - 1].Images = append(productsSlice[len(productsSlice) - 1].Images, rowImage{
				ImageId:       row.ImageId,
				ImageUrl:      row.ImageUrl,
				ImagePublicId: row.ImagePublicId,
			})
		}

	}

	var result []respAllProds
    for _, product := range productsSlice {
        result = append(result, *product)
    }

	return result
}
