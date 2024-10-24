package order

import (
	"fmt"
	"sync"

	"gorm.io/gorm"
	"main.go/internal/websocket"
	"main.go/pkg/models"
	"main.go/services/generic"
	"main.go/types"
)

type Store struct {
	DB      *gorm.DB
	Generic *generic.GenericRepository[models.Order]
}

func NewStore(DB *gorm.DB) *Store {
	return &Store{
		DB:      DB,
		Generic: &generic.GenericRepository[models.Order]{DB: DB},
	}
}

var (
	notFoundMsg = "order with id: '%v' was not found"
)

func (orderStore *Store) GetPopulatedOrderById(Id uint) ([]types.GetOneOrderRow, error) {
	var order []types.GetOneOrderRow
	err := orderStore.DB.Model(&models.Order{}).Select(selectOneOrderQ).Where("orders.id = ?", Id).
		Joins(jointWAddress).Joins(joinWOrderItems).Joins(joinWProducts).Joins(jointWProductImages).
		Scan(&order).Error
	if err != nil || order == nil {
		notFoundErr := fmt.Errorf(notFoundMsg, Id)
		return nil, notFoundErr
	}

	return order, err
}

func (orderStore *Store) GetAllOrders(page, limit int) ([]models.Order, int64, error) {
	orders, count, errs := orderStore.Generic.GetAll(page, limit)
	if len(errs) != 0 {
		return nil, 0, errs[0]
	}

	return orders, count, nil
}

func (orderStore *Store) CreateOrder(tx *gorm.DB, order *models.Order) error {
	order.Status = models.Pending
	return tx.Create(order).Error
}

func (orderStore *Store) GetProductsByIds(Ids []uint) ([]models.Product, error) {
	var products = make([]models.Product, 0, len(Ids))
	err := orderStore.DB.Where(Ids).Find(&products).Error

	return products, err
}

func (orderStore *Store) ValidateAndCalTotalPrice(prods []models.Product, orderItems []models.OrderItem) (*float64, error) {
	var totalPrice = 0.0

	for _, prod := range prods {
		for _, orderItem := range orderItems {
			if orderItem.ProductID == prod.ID {
				if orderItem.Quantity > prod.Quantity {
					return nil, fmt.Errorf("product with name: '%v' has '%v' quantity which is less than its cart item", prod.Name, prod.Quantity)
				}

				totalPrice = totalPrice + (prod.Price * float64(orderItem.Quantity))
			}
		}
	}

	return &totalPrice, nil
}

func (orderStore *Store) CreateOrderItems(tx *gorm.DB, order *models.Order, orderItems []models.OrderItem) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(orderItems))
	for _, orderItem := range orderItems {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := tx.Create(&models.OrderItem{
				OrderID:   order.ID,
				ProductID: orderItem.ProductID,
				UnitPrice: orderItem.Product.Price,
				Quantity:  orderItem.Quantity,
			}).Error
			if err != nil {
				errChan <- err
			}
		}()
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

func (orderStore *Store) CreateOrderWithItems(order *models.Order, userId *uint, orderItems []models.OrderItem) error {
	return orderStore.DB.Transaction(func(tx *gorm.DB) error {
		err := orderStore.CreateOrder(tx, order)
		if err != nil {
			return err
		}

		err = orderStore.CreateOrderItems(tx, order, orderItems)
		if err != nil {
			return err
		}

		err = orderStore.EmptyTheCartTx(tx, *userId)
		if err != nil {
			return err
		}

		return nil
	})
}

func (orderStore *Store) CancelOrder(id uint, userId uint) error {
	var order models.Order
	err := orderStore.DB.Model(&order).Where("id = ? AND user_id = ?", id, userId).First(&order).Error
	if err != nil {
		return fmt.Errorf(notFoundMsg, id)
	}
	if order.Status == models.Cancelled {
		return fmt.Errorf("order is already cancelled")
	}

	err = orderStore.DB.Model(&order).Where("id = ? AND user_id = ?", id, userId).Update("status", string(models.Cancelled)).Error
	if err != nil {
		return err
	}

	return nil
}

func (orderStore *Store) UpdateOrderStatus(id uint, status models.Status) error {
	var order models.Order
	return orderStore.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&order).First(&order, id).Error
		if err != nil {
			return fmt.Errorf(notFoundMsg, id)
		}

		err = tx.Model(&order).First(&order, id).Update("status", string(status)).Error
		if err != nil {
			return err
		}

		if status == models.Delivered {
			WSProduct,err := orderStore.UpdateProductQtys(tx, id)
			if err != nil {
				return err
			}

			go websocket.GlobalManager.BroadcastProductQtyChange(WSProduct)
		}
		return nil
	})
}

func (orderStore *Store) GetAddressById(addressId uint) (*models.Address, error) {
	var address models.Address
	err := orderStore.DB.Model(&address).First(&address, addressId).Error
	if err != nil {
		return nil, fmt.Errorf("address with id: '%v' was not found", addressId)
	}

	return &address, nil
}

func (orderStore *Store) EmptyTheCartTx(tx *gorm.DB, userId uint) error {
	return tx.Model(&models.CartItem{}).Where("user_id = ?", userId).Delete(&models.CartItem{}).Error
}

func (orderStore *Store) GetCart(userId uint) ([]models.CartItem, error) {
	var cart []models.CartItem
	err := orderStore.DB.Where("user_id = ?", userId).Preload("Product").Find(&cart).Error
	if err != nil {
		return nil, err
	}
	return cart, err
}

func (orderStore *Store) GetCartItemsCount(userId uint) (*int64, error) {
	var count int64
	err := orderStore.DB.Model(&models.CartItem{}).Where("user_id = ?", userId).Count(&count).Error
	if err != nil {
		return nil, err
	}
	return &count, err
}

func (orderStore *Store) ConvertToOrderItems(cart []models.CartItem) []models.OrderItem {
	var orderItems = make([]models.OrderItem, 0, len(cart))
	for _, cartItem := range cart {
		orderItems = append(orderItems, models.OrderItem{
			Quantity:  cartItem.Quantity,
			ProductID: cartItem.ProductID,
			Product:   cartItem.Product,
		})
	}

	return orderItems
}

func (orderStore *Store) ExtractProductIds(cart []models.CartItem) []uint {
	var productsIds = make([]uint, 0, len(cart))
	for _, cartItem := range cart {
		productsIds = append(productsIds, cartItem.ProductID)
	}

	return productsIds
}

func (orderStore *Store) UpdateProductQtys(tx *gorm.DB,orderId uint) ([]types.ProductAmountDiscounter, error) {
	orderItems, err := orderStore.GetOrderItems(orderId)
	if err != nil {
		return nil, err
	}
	productsIds := make([]uint, 0, len(orderItems))
	for _, orderItem := range orderItems {
		productsIds = append(productsIds, orderItem.ProductID)
	}

	var productQtyChange = make([]types.ProductAmountDiscounter, 0, len(orderItems))

	caseStatement := "CASE "
	for _, orderItem := range orderItems {
		caseStatement += fmt.Sprintf("WHEN id = %d THEN quantity - %d ", orderItem.ProductID, orderItem.Quantity)
		productQtyChange = append(productQtyChange, websocket.WSProduct{
			ID: orderItem.ProductID,
			DiscountAmount: orderItem.Quantity,
		})
		
	}
	caseStatement += " END"

	err = orderStore.DB.Model(&models.Product{}).Unscoped().Where("id IN ?", productsIds).
	UpdateColumn("quantity", gorm.Expr(caseStatement)).Error
	if err != nil {
		return nil, err
	}

	return productQtyChange, nil
}

func (orderStore *Store) GetOrderItems(orderId uint) ([]models.OrderItem, error) {
	var orderItems []models.OrderItem
	err := orderStore.DB.Model(&orderItems).Where(models.OrderItem{OrderID: orderId}).Find(&orderItems).Error
	if err != nil {
		return nil, err
	}

	return orderItems, nil
}
