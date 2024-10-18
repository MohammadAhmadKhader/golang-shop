package order

import (
	"time"
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

type GetOneOrderRow struct {
	Id                   uint
	UserId               uint
	TotalPrice           float64
	Status               string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	OrderItemId          uint
	UnitPrice            float64
	OrderItemQuantity    uint8
	AddressId            uint
	AddressFullName      string
	AddressCity          string
	AddressCountry       string
	AddressStreetAddress string
	AddressZipCode       *string
	AddressState         *string
	ProductId            uint
	ProductName          string
	ProductQuantity      uint
	ProductPrice         float64
	ProductImageUrl string
	ProductImageIsMain bool
}
type GetAllOrdersRows struct {
	Id                   uint
	UserId               uint
	TotalPrice           float64
	Status               string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	OrderItemsCount      uint
	AddressId            uint
	AddressFullName      string
	AddressCity          string
	AddressCountry       string
	AddressStreetAddress string
	AddressZipCode       *string
	AddressState         *string
}

type respAllOrders struct {
	Id              uint        `json:"id"`
	UserId          uint        `json:"userId"`
	TotalPrice      float64     `json:"totalPrice"`
	Status          string      `json:"status"`
	OrderItemsCount uint        `json:"orderItemsCount"`
	Address         respAddress `json:"address"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt"`
}

type respOneOrder struct {
	Id         uint            `json:"id"`
	UserId     uint            `json:"userId"`
	TotalPrice float64         `json:"totalPrice"`
	Status     string          `json:"status"`
	OrderItems []respOrderItem `json:"orderItems"`
	Address    respAddress     `json:"address"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
}

type respOrderItem struct {
	Id       uint        `json:"id"`
	Price    float64     `json:"price"`
	Quantity uint8       `json:"quantity"`
	Product  respProduct `json:"product"`
}

type respProduct struct {
	Id       uint    `json:"id"`
	Name     string  `json:"name"`
	Quantity uint    `json:"quantity"`
	Price    float64 `json:"price"`
	MainImage respImage `json:"mainImage"`
}

type respImage struct {
	IsMain bool `json:"isMain"`
	ImageUrl string `json:"imageUrl"`
}

type respAddress struct {
	Id            uint    `json:"id"`
	FullName      string  `json:"fullName"`
	City          string  `json:"city"`
	StreetAddress string  `json:"streetAddress"`
	State         *string `json:"state"`
	ZipCode       *string `json:"zipCode"`
	Country       string  `json:"country"`
}

func convertRowsToResp(rows []GetAllOrdersRows) []*respAllOrders {
	orderMap := make(map[uint]uint)
	ordersSlice := make([]*respAllOrders, 0)

	for _, row := range rows {
		if _, exists := orderMap[row.Id]; !exists {
			orderMap[row.Id] = row.Id
			order := &respAllOrders{
				Id:         row.Id,
				UserId:     row.UserId,
				TotalPrice: row.TotalPrice,
				Status:     row.Status,
				CreatedAt:  row.CreatedAt,

				UpdatedAt: row.UpdatedAt,
				Address: respAddress{
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

func convertRowToResp(rows []GetOneOrderRow) *respOneOrder {
	orderMap := make(map[uint]uint)
	order := new(respOneOrder)

	for _, row := range rows {
		if _, exists := orderMap[row.Id]; !exists {
			orderMap[row.Id] = row.Id
			rowOrder := &respOneOrder{
				Id:         row.Id,
				UserId:     row.UserId,
				TotalPrice: row.TotalPrice,
				Status:     row.Status,
				CreatedAt:  row.CreatedAt,
				UpdatedAt:  row.UpdatedAt,
				Address: respAddress{
					Id:            row.AddressId,
					FullName:      row.AddressFullName,
					City:          row.AddressCity,
					StreetAddress: row.AddressStreetAddress,
					Country:       row.AddressCountry,
					State:         row.AddressState,
					ZipCode:       row.AddressZipCode,
				},
				OrderItems: []respOrderItem{},
			}

			order = rowOrder
		}

		order.OrderItems = append(order.OrderItems, respOrderItem{
			Id:       row.OrderItemId,
			Price:    row.UnitPrice,
			Quantity: row.OrderItemQuantity,
			Product: respProduct{
				Id:       row.ProductId,
				Name:     row.ProductName,
				Quantity: row.ProductQuantity,
				Price:    row.ProductPrice,
				MainImage: respImage{
					IsMain: row.ProductImageIsMain,
					ImageUrl: row.ProductImageUrl,
				},
			},
		})
	}

	return order
}
