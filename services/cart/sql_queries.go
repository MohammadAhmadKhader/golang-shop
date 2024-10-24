package cart

import "main.go/types"

var selectQ = ` cart_items.id as id, cart_items.quantity as quantity,
				products.id as product_id, products.name as product_name, products.price as product_price,
				images.image_url as product_image`

var joinWProducts = `LEFT JOIN products ON cart_items.product_id = products.id`
var joinWImages = `LEFT JOIN images ON images.product_id = products.id AND images.is_main = 1`

func convertRowsToResponse(rows []types.GetCartRow) *types.RespCartShape {
	var respShape types.RespCartShape
	cartItems := make([]types.RespCartItem, 0, len(rows))
	for _, row := range rows {
		cartItems = append(cartItems, types.RespCartItem{
			ID:       row.ID,
			Quantity: row.Quantity,
			Product: types.RespCartItemProduct{
				ID:    row.ProductID,
				Name:  row.ProductName,
				Image: row.ProductImage,
				Price: row.ProductPrice,
			},
		})
	}

	respShape.CartItems = cartItems

	if len(rows) != 0 {
		respShape.Id = rows[0].ID
	}

	return &respShape
}