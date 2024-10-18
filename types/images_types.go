package types

import (
	"context"
	"mime/multipart"
	"net/http"

)

type Folder string

const (
	ProductsFolder Folder = "golang-shop/products"
	UsersFolder    Folder = "golang-shop/users"
)


type IImagesHandler interface {
	UploadOne(file *multipart.File, folder Folder) (*UploadResponse, error)
	UploadMany(r *http.Request, folder Folder, keyName string, ctx context.Context) ([]*UploadResponse, []error)
	DeleteOne(file *multipart.File, folder Folder) error
	DeleteMany(PublicIDs []string, ctx context.Context) (*DeleteResponse, error)
}

type UploadResponse struct {
	Width     int
	Height    int
	SecureUrl string
	PublicID  string
	URL       string
	Format    string
}

type DeleteManyResponse struct {
	DeletedCounts map[string]interface{}
}

// TODO : resolving the issue with cloudinary delete response type
type DeleteResponse struct {
	result string
}