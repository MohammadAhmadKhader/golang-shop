package order

import (
	"time"
)

var selectQ = `orders.id as id, orders.user_id as user_id, orders.total_price as total_price, 
				orders.status as status, orders.created_at as created_at, orders.updated_at as updated_at,
				order_items.id as order_item_id, order_items.unit_price as unit_price,order_items.quantity as order_item_quantity,
				products.id as product_id, products.name as product_name, products.quantity as product_quantity, products.price as product_price
				`
var joinWOrderItems = `LEFT JOIN order_items ON order_items.order_id = orders.id`
var joinWProducts = `LEFT JOIN products ON order_items.product_id = products.id`

type GetAllOrdersRow struct {
	Id         uint
	UserId     uint
	TotalPrice float64
	Status     string
	CreatedAt  time.Time
	UpdatedAt time.Time
	OrderItemId uint
	UnitPrice float64
	OrderItemQuantity uint8
	ProductId uint
	ProductName string
	ProductQuantity uint
	ProductPrice float64
}

type respAllOrders struct {
	Id         uint `json:"id"`
	UserId     uint `json:"userId"`
	TotalPrice float64 `json:"totalPrice"`
	Status     string `json:"status"`
	OrderItems []respOrderItem `json:"orderItems"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type respOrderItem struct {
	Id uint `json:"id"`
	Price float64 `json:"price"`
	Quantity uint8 `json:"quantity"`
	Product respProduct `json:"product"`
}

type respProduct struct {
	Id uint `json:"id"`
	Name string `json:"name"`
	Quantity uint `json:"quantity"`
	Price float64 `json:"price"`
}

func convertRowsToResp(rows []GetAllOrdersRow) []*respAllOrders {
	orderMap := make(map[uint]uint)
	ordersSlice := make([]*respAllOrders, 0)

	for _, row := range rows {
		if _, exists := orderMap[row.Id]; !exists {
			orderMap[row.Id] = row.Id
			order := &respAllOrders{
				Id: row.Id,
				UserId: row.UserId,
				TotalPrice: row.TotalPrice,
				Status: row.Status,
				CreatedAt: row.CreatedAt,
				UpdatedAt: row.UpdatedAt,
				OrderItems: []respOrderItem{},
			}

			ordersSlice = append(ordersSlice, order)
		}
		

		ordersSlice[len(ordersSlice)-1].OrderItems = append(ordersSlice[len(ordersSlice)-1].OrderItems, respOrderItem{
			Id: row.OrderItemId,
			Price: row.UnitPrice,
			Quantity: row.OrderItemQuantity,
			Product: respProduct{
				Id: row.ProductId,
				Name: row.ProductName,
				Quantity: row.ProductQuantity,
				Price: row.ProductPrice,
			},
		})
	}

	return ordersSlice
}

func convertRowToResp(rows []GetAllOrdersRow) *respAllOrders {
	orderMap := make(map[uint]uint)
	order := new(respAllOrders)

	for _, row := range rows {
		if _, exists := orderMap[row.Id]; !exists {
			orderMap[row.Id] = row.Id
			rowOrder := &respAllOrders{
				Id: row.Id,
				UserId: row.UserId,
				TotalPrice: row.TotalPrice,
				Status: row.Status,
				CreatedAt: row.CreatedAt,
				UpdatedAt: row.UpdatedAt,
				OrderItems: []respOrderItem{},
			}

			order = rowOrder
		}
		

		order.OrderItems = append(order.OrderItems, respOrderItem{
			Id: row.OrderItemId,
			Price: row.UnitPrice,
			Quantity: row.OrderItemQuantity,
			Product: respProduct{
				Id: row.ProductId,
				Name: row.ProductName,
				Quantity: row.ProductQuantity,
				Price: row.ProductPrice,
			},
		})
	}

	return order
}