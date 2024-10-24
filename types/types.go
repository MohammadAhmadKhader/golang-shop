package types

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
	"main.go/pkg/models"
	"main.go/pkg/payloads"
)

type GenericRepository[TModel any] interface {
	GetOne(Id uint, notFoundMsg string) (TModel, error)
	GetAll(page int, limit int) (TModelS []TModel, count int64, errors []error)
	Create(model *TModel, selectedFields []string) (*TModel, error)
	CreateTx(model *TModel, tx *gorm.DB) (error)
	Update(model *TModel, selectedFields []string) (*TModel, error)
	UpdateAndReturn(id uint, model *TModel, selectedFields []string) (*TModel, error)
	SoftDelete(Id uint, notFoundMsg string) error
	SoftDeleteWithUserId(id uint, userId uint, notFoundMsg string) error
	HardDelete(Id uint, notFoundMsg string) error
	Restore(id uint, notFoundMsg string) (*TModel, error)
	RestoreWithUserId(id uint, userId uint, notFoundMsg string) (*TModel, error)
	GetAllDeleted(page, limit int) ([]TModel, int64, []error)
	FindThenUpdate(id uint, changes *TModel, selectedFields []string, notFoundMsg string)
	FindThenUpdateWithAuth(id uint, changes *TModel, selectedFields []string, notFoundMsg string, userId uint) (*TModel, error)
	FindThenDeleteWithAuth(id uint, notFoundMsg string, userId uint) (*TModel, error)
}

type CategoryStore interface {
	GetCategoryById(Id uint) (*models.Category, error)
	GetAllCategories(page, limit int) ([]models.Category, int64, error)
	CreateCategory(category *models.Category) (*models.Category, error)
	UpdateCategory(id uint, category *models.Category) (*models.Category, error)
}

type ProductAmountDiscounter interface {
	GetProductId() uint
	GetAmountDiscount() uint
}

type OrderStore interface {
	GetPopulatedOrderById(Id uint) ([]GetOneOrderRow, error)
	CreateOrder(tx *gorm.DB, order *models.Order) error
	GetAllOrders(page, limit int) ([]models.Order, int64, error)
	GetProductsByIds(Ids []uint) ([]models.Product, error)
	ValidateAndCalTotalPrice(prods []models.Product, orderItems []models.OrderItem) (*float64, error)
	CreateOrderItems(tx *gorm.DB, order *models.Order, orderItems []models.OrderItem) error
	CreateOrderWithItems(order *models.Order, userId *uint, orderItems []models.OrderItem) error
	EmptyTheCartTx(tx *gorm.DB, userId uint) error
	CancelOrder(Id uint, userId uint) error
	UpdateOrderStatus(Id uint, status models.Status) error
	GetAddressById(addressId uint) (*models.Address, error)
	GetCartItemsCount(userId uint) (*int64, error)
	GetCart(userId uint) ([]models.CartItem, error)
	ConvertToOrderItems(cart []models.CartItem) []models.OrderItem
	ExtractProductIds(cart []models.CartItem) []uint
	UpdateProductQtys(tx *gorm.DB,orderId uint) ([]ProductAmountDiscounter, error)
	GetOrderItems(orderId uint) ([]models.OrderItem, error)
}

type ProductStore interface {
	GetProductById(Id uint) ([]RowGetProductById, error)
	GetAllProducts(page, limit int, filter func(db *gorm.DB, filters []FilterCondition) ([]models.Product, error)) ([]models.Product, int64, error)
	CreateProduct(product *models.Product) (*models.Product, error)
	UpdateProduct(id uint, changes *models.Product, excluder Excluder) (*models.Product, error)
	CreateImageTx(tx *gorm.DB, uploadResp *UploadResponse, productId uint, isMain bool) (*models.Image, error)
	CreateProductWithImage(product *models.Product, uploadResp *UploadResponse) (*models.Product, error)
}

type UserStore interface {
	GetUserById(Id uint) (*models.User, error)
	GetUserWithRolesById(Id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user payloads.UserSignUp) (*models.User, error)
	UpdatePassword(newHashedPassword, email string) error
	UpdateProfile(id uint, user *models.User, excluder Excluder) (*models.User, error)
	RemoveUserRole(roleId, userId uint) (error)
	AssignUserRole(roleId, userId uint) (*models.UserRoles, error)
}

type ReviewStore interface {
	GetReviewById(Id uint) (*models.Review, error)
	GetAllReviews(page, limit int) ([]models.Review, int64, error)
	UpdateReview(id,userId uint, updatePayload *models.Review,excluder Excluder) (*models.Review, error)
	CreateReview(createPayload *models.Review) (*models.Review, error)
	HardDelete(Id uint, userId uint) error
	GetProductById(productId uint) (*models.Product, error)
}

type ImageStore interface {
	GetImageById(id uint) (*models.Image, error)
	CreateImage(image *models.Image) (*models.Image, error)
	UpdateImageUrl(id uint, newImageUrl string) (error)
	GetCountOfProductImages(productId uint) (*int64, error)
	GetProductById(productId uint) (*models.Product, error)
	DeleteImageById(id uint) error
	SetImageAsMainTx(tx *gorm.DB, id, productId uint) error
	SetImageAsNotMainTx(tx *gorm.DB, productId uint) error
	SwapMainStatus(id, productId uint) error
	CreateManyImages(uploadResults []*UploadResponse, productId *uint) ([]models.Image, error)
}

type RolesStore interface {
	GetRole(id uint) (*models.Role, error)
	GetAllRoles(page, limit int) ([]models.Role, int64, error)
	CreateRole(role *models.Role) (*models.Role, error)
	UpdateRole(id uint, role *models.Role) (*models.Role, error)
	DeleteRole(id uint) error
}

type MessageStore interface {
	GetById(userId uint, lastMessageId uint, cursor time.Time, limit int) ([]*models.Message, error)
	GetByUsersIds(userId uint, to uint, lastMessageId uint, cursor time.Time, limit int) ([]*models.Message, error)
}

type CartStore interface {
	GetCartItemById(Id uint) (*models.CartItem, error)
	ChangeCartItemQty(oldQty uint, payload *payloads.ChangeCartItemQty, cartItem *models.CartItem) (*models.CartItem, error)
	GetCartByUserId(userId uint) (*RespCartShape, error)
	AddToCart(payload *payloads.AddCartItem, userId uint) (*models.CartItem, error)
	DeleteCartItem(itemId uint) error
	ClearCart(userId uint) error
}

type AddressStore interface {
	GetById(id uint, userId uint) (*models.Address, error)
	GetAddressById(id uint) (*models.Address, error)
	CreateAddress(address *models.Address) (*models.Address, error)
	UpdateAddress(id uint, address *models.Address, excluder Excluder) (*models.Address, error)
	DeleteAddress(id uint, userId uint) (error)
	GetAllAddresses(userId uint) ([]models.Address, error)
	GetUndeletedAddressesCount(userId uint) (*int64, error)
}

type TokenPayload struct {
	Email     string `json:"email"`
	UserId    int    `json:"userId"`
	ExpiredAt int64  `json:"expiredAt"`
}

type TokenKey string
// this type is used to avoid calling the function "utils.GetUserIdFromTokenPayload" twice, once in middle ware and one in route
//
// when its called in middleware this key gets stored inside the context as stand alone variable 
// to avoid conversion some repeated un-necessary computation.

type UserRole string
type UserKey string
type AuthorizedResource string

const (
	SuperAdmin  UserRole = "SuperAdmin"
	Admin       UserRole = "Admin"
	RegularUser UserRole = "RegularUser"
)

type Pagination struct {
	Page  int
	Limit int
}

type FilterCondition struct {
	Field    string      // The field to filter by
	Operator string      // The operator to use (e.g., '=', '>', 'LIKE')
	Value    interface{} // The value to filter by
}

// this has been created instead of the user store to solve an import cycle
type IUserFetcher interface {
	GetUserById(Id uint) (*models.User, error)
	GetUserRolesByUserId(Id uint) ([]models.UserRoles, error)
}

type SortCondition struct {
	Field   *string
	SortDir *string
}

func (s *SortCondition) Validate() *SortCondition {
	if s.Field == nil || strings.Contains(*s.Field, " ") {
		*s.Field = "createdAt"
	}

	if (s.SortDir == nil) || (*s.SortDir != "ASC" && *s.SortDir != "DESC" && *s.SortDir != "desc" && *s.SortDir != "asc") {
		*s.Field = "DESC"
	}
	return s
}

type Excluder interface {
	Exclude(selectedFields []string) []string
}

type SoftDeletable interface {
	ShowDeletedAt() ([]byte, error)
}

type AppResponse struct {
    http.ResponseWriter
    StatusCode int
}

func (ar *AppResponse) WriteHeader(statusCode int) {
    ar.StatusCode = statusCode
    ar.ResponseWriter.WriteHeader(statusCode)
}

func (w *AppResponse) Hijack() (net.Conn, *bufio.ReadWriter, error) {
    hijacker, ok := w.ResponseWriter.(http.Hijacker)
    if !ok {
        return nil, nil, fmt.Errorf("response writer does not support hijacking")
    }
    return hijacker.Hijack()
}

type WSMessageStatus string

const (
	Sent WSMessageStatus = "Sent"
	Delivered WSMessageStatus = "Delivered"
	Seen WSMessageStatus = "Seen"
)