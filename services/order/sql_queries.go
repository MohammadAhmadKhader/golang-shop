package order

import (

	"main.go/types"
)

var selectOneOrderQ = `orders.id as id, orders.user_id as user_id, orders.total_price as total_price, 
				orders.status as status, orders.created_at as created_at, orders.updated_at as updated_at,
				order_items.id as order_item_id, order_items.unit_price as unit_price,order_items.quantity as order_item_quantity,
				products.id as product_id, products.name as product_name, products.quantity as product_quantity, products.price as product_price, 
				addresses.id as address_id, addresses.full_name as address_full_name, addresses.city as address_city, addresses.street_address as address_street_address,
				addresses.zip_code as address_zip_code, addresses.state as address_state, addresses.country as address_country,
				images.image_url as product_image_url, images.is_main as product_image_is_main`
var joinWOrderItems = `LEFT JOIN order_items ON order_items.order_id = orders.id`
var joinWProducts = `LEFT JOIN products ON order_items.product_id = products.id`
var jointWAddress = `LEFT JOIN addresses ON orders.address_id = addresses.id`
var jointWProductImages = `LEFT JOIN images ON images.product_id = products.id AND images.is_main = 1`

var selectAllOrdersQ = `orders.id as id, orders.user_id as user_id, orders.total_price as total_price, 
				orders.status as status, orders.created_at as created_at, orders.updated_at as updated_at,
				addresses.id as address_id, addresses.full_name as address_full_name, addresses.city as address_city, addresses.street_address as address_street_address,
				addresses.zip_code as address_zip_code, addresses.state as address_state, addresses.country as address_country, order_items_count.order_items_count`
var joinWOrderItemsCount = `LEFT JOIN (
	SELECT order_id, COUNT(*) AS order_items_count 
	FROM order_items 
	GROUP BY order_id
) order_items_count ON order_items_count.order_id = orders.id `

func convertRowsToResp(rows []types.GetAllOrdersRows) []*types.RespAllOrders {
	orderMap := make(map[uint]uint)
	ordersSlice := make([]*types.RespAllOrders, 0)

	for _, row := range rows {
		if _, exists := orderMap[row.Id]; !exists {
			orderMap[row.Id] = row.Id
			order := &types.RespAllOrders{
				Id:         row.Id,
				UserId:     row.UserId,
				TotalPrice: row.TotalPrice,
				Status:     row.Status,
				CreatedAt:  row.CreatedAt,

				UpdatedAt: row.UpdatedAt,
				Address: types.RespOrderAddress{
					Id:            row.AddressId,
					FullName:      row.AddressFullName,
					City:          row.AddressCity,
					StreetAddress: row.AddressStreetAddress,
					Country:       row.AddressCountry,
					State:         row.AddressState,
					ZipCode:       row.AddressZipCode,
				},
				OrderItemsCount: row.OrderItemsCount,
			}

			ordersSlice = append(ordersSlice, order)
		}
	}

	return ordersSlice
}

func convertRowToResp(rows []types.GetOneOrderRow) *types.RespOneOrder {
	orderMap := make(map[uint]uint)
	order := new(types.RespOneOrder)

	for _, row := range rows {
		if _, exists := orderMap[row.Id]; !exists {
			orderMap[row.Id] = row.Id
			rowOrder := &types.RespOneOrder{
				Id:         row.Id,
				UserId:     row.UserId,
				TotalPrice: row.TotalPrice,
				Status:     row.Status,
				CreatedAt:  row.CreatedAt,
				UpdatedAt:  row.UpdatedAt,
				Address: types.RespOrderAddress{
					Id:            row.AddressId,
					FullName:      row.AddressFullName,
					City:          row.AddressCity,
					StreetAddress: row.AddressStreetAddress,
					Country:       row.AddressCountry,
					State:         row.AddressState,
					ZipCode:       row.AddressZipCode,
				},
				OrderItems: []types.RespOrderItem{},
			}

			order = rowOrder
		}

		order.OrderItems = append(order.OrderItems, types.RespOrderItem{
			Id:       row.OrderItemId,
			Price:    row.UnitPrice,
			Quantity: row.OrderItemQuantity,
			Product: types.RespOrderItemProduct{
				Id:       row.ProductId,
				Name:     row.ProductName,
				Quantity: row.ProductQuantity,
				Price:    row.ProductPrice,
				MainImage: types.RespOrderItemImage{
					IsMain: row.ProductImageIsMain,
					ImageUrl: row.ProductImageUrl,
				},
			},
		})
	}

	return order
}
