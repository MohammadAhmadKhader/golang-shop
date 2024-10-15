package cart

import "main.go/pkg/models"

var selectQ = ` carts.id as id,carts.user_id as user_id,
				cart_items.id as cart_item_id, cart_items.quantity as quantity,
				products.id as product_id, products.name as product_name, products.price as product_price,
				images.image_url as product_image
				`
var joinWCartItems = `LEFT JOIN cart_items ON cart_items.cart_id = carts.id`
var joinWProducts = `LEFT JOIN products ON cart_items.product_id = products.id`
var joinWImages = `LEFT JOIN images ON images.product_id = products.id AND images.is_main = 1`

type getCartRow struct {
	ID           uint
	UserID       uint
	CartItemID   uint
	Quantity     uint
	ProductID    uint
	ProductName  string
	ProductPrice float64
	ProductImage string
}

type respProduct struct {
	ID    uint    `json:"id"`
	Image string  `json:"image"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type respCartItem struct {
	ID       uint        `json:"id"`
	Quantity uint        `json:"quantity"`
	Product  respProduct `json:"product"`
}

type respCart struct {
	Id        uint           `json:"id"`
	UserID    uint           `json:"userId"`
	CartItems []respCartItem `json:"cartItems"`
}

func convertRowsToResponse(rows []getCartRow) *respCart {
	var respShape respCart
	cartItems := make([]respCartItem, 0, len(rows))
	for _, row := range rows {
		if row.CartItemID != 0 {
			cartItems = append(cartItems, respCartItem{
				ID:       row.CartItemID,
				Quantity: row.Quantity,
				Product: respProduct{
					ID:    row.ProductID,
					Name:  row.ProductName,
					Image: row.ProductImage,
					Price: row.ProductPrice,
				},
			})
		}

	}

	respShape.CartItems = cartItems

	if len(rows) != 0 {
		respShape.UserID = rows[0].UserID
		respShape.Id = rows[0].ID
	}

	return &respShape
}


type getCartWithCartItemsRow struct {
	ID         uint
	UserID     uint
	CartItemID uint
	Quantity   uint
}

func convertRowsToModel(rows []getCartWithCartItemsRow) *models.Cart {
	var cart models.Cart
	cartItems := make([]models.CartItem, 0, len(rows))
	for _, row := range rows {
		if row.CartItemID != 0 {
			cartItems = append(cartItems, models.CartItem{
				ModelBasics: models.ModelBasics{ID: row.CartItemID},
				Quantity: row.Quantity,
			})
		}
	}

	cart.CartItems = cartItems
	if len(rows) != 0 {
		cart.UserID = rows[0].UserID
		cart.ID = rows[0].ID
	}

	return &cart
}