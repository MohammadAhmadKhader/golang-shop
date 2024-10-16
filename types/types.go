package types

import (
	"strings"

	"main.go/pkg/models"
	"main.go/pkg/payloads"
)

type IGenericRepository[TModel any] interface {
	GetOne(Id uint, notFoundErr error) (TModel, error)
	GetAll(page int, limit int) (TModelS []TModel, count int64, errors []error)
	Create(model *TModel, selectedFields []string) *TModel
	Update(model *TModel, selectedFields []string) *TModel
	SoftDelete(Id uint, notFoundErr error) error
	HardDelete(Id uint, notFoundErr error) error
}

type ICategoryStore interface {
	GetCategoryById(Id uint) (*models.Category, error)
	GetAllCategories(page, limit int) ([]models.Category, int64, error)
	CreateCategory(category *models.Category) (*models.Category, error)
	UpdateCategory(category *models.Category) (*models.Category, error)
	SoftDeleteCategory(Id uint) error
	HardDeleteCategory(Id uint) error
}

type IOrderStore interface {
	GetOrderById(Id uint) (*models.Order, error)
	GetAllOrders(page, limit int) ([]models.Order, int64, error)
	SoftDeleteOrder(Id uint) error
	HardDeleteOrder(Id uint) error
}

type IProductStore interface {
	GetProductById(Id uint) (*models.Product, error)
	GetAllProducts(page, limit int) ([]models.Product, int64, error)
	CreateProduct(product *models.Product) (*models.Product, error)
	SoftDeleteProduct(Id uint) error
	HardDeleteProduct(Id uint) error
}

type IUserStore interface {
	GetUserById(Id uint) (*models.User, error)
	GetUserWithRolesById(Id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user payloads.UserSignUp) (*models.User, error)
	UpdatePassword(newHashedPassword, email string) error
}

type IReviewStore interface {
	GetReviewById(Id uint) (*models.Review, error)
	GetAllReviews(page, limit int) ([]models.Review, int64, error)
	SoftDeleteReview(Id uint) error
	HardDeleteReview(Id uint) error
}

type IImageStore interface {
	CreateImage(image *models.Image) (*models.Image, error)
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