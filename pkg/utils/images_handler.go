package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"sync"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"main.go/config"
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
	DeletedCounts    map[string]interface{}
}

// TODO : resolving the issue with cloudinary delete response type
type DeleteResponse struct {
	result string
}

type ImagesHandler struct {
	cld *cloudinary.Cloudinary
}

func NewImagesHandler() *ImagesHandler {
	return &ImagesHandler{
		cld: config.GetCloudinary(),
	}
}

func (ih *ImagesHandler) UpdateOne(newFile *multipart.File, newFileHeader *multipart.FileHeader, folder Folder, ctx context.Context, oldImagePublicId string) (*UploadResponse, error) {
	uploadParams := NewUploadParams(string(folder), newFileHeader, true, &oldImagePublicId)
	resp, err := ih.cld.Upload.Upload(ctx, *newFile, uploadParams)
	if err != nil {
		return nil, err
	}

	return ih.getUploadResponse(resp), nil
}

func (ih *ImagesHandler) UploadOne(file *multipart.File, fileHeader *multipart.FileHeader, folder Folder, ctx context.Context) (*UploadResponse, error) {
	folderAsString := string(folder)

	uploadParams := NewUploadParams(folderAsString, fileHeader, true, nil)
	resp, err := ih.cld.Upload.Upload(ctx, *file, uploadParams)
	if err != nil {
		return nil, err
	}

	return ih.getUploadResponse(resp), nil
}

func (ih *ImagesHandler) UploadMany(r *http.Request, folder Folder, keyName string, ctx context.Context) ([]*UploadResponse, []error) {
	folderAsString := string(folder)
	files := r.MultipartForm.File[keyName]
	if len(files) == 0 {
		return nil, []error{fmt.Errorf("no images were provided")}
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	UploadResponses := make([]*UploadResponse, 0, len(files))
	errors := make([]error, 0)

	for _, fileHeader := range files {
		wg.Add(1)
		go func() {
			defer wg.Done()

			file, err := fileHeader.Open()
			if err != nil {
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
				return
			}
			defer file.Close()

			uploadParams := NewUploadParams(folderAsString, fileHeader, true, nil)
			resp, err := ih.cld.Upload.Upload(ctx, file, uploadParams)
			if err != nil {
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
				return
			}

			mu.Lock()
			UploadResponses = append(UploadResponses, ih.getUploadResponse(resp))
			mu.Unlock()
		}()
	}

	wg.Wait()

	return UploadResponses, errors
}

func (ih *ImagesHandler) DeleteOne(PublicID string, ctx context.Context) error {
	_, err := ih.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: PublicID,
	})
	if err != nil {
		return err
	}

	return nil 
}

func (ih *ImagesHandler) DeleteMany(PublicIDs []string, ctx context.Context) (*DeleteManyResponse, error) {
	delResp, err := ih.cld.Admin.DeleteAssets(ctx, admin.DeleteAssetsParams{
		AssetType: "image",
		PublicIDs: PublicIDs,
	})
	if err != nil {
		return nil, err
	}

	return ih.getDeleteManyResponse(delResp), nil
}

func (ih *ImagesHandler) getUploadResponse(uploadResult *uploader.UploadResult) *UploadResponse {
	return &UploadResponse{
		Width: uploadResult.Width, Height: uploadResult.Height,
		SecureUrl: uploadResult.SecureURL, PublicID: uploadResult.PublicID,
		URL: uploadResult.URL, Format: uploadResult.Format,
	}
}

func (ih *ImagesHandler) getDeleteManyResponse(deleteResult *admin.DeleteAssetsResult) *DeleteManyResponse {
	return &DeleteManyResponse{
		DeletedCounts: deleteResult.DeletedCounts,
	}
}

func NewUploadParams(folderName string, fileHeader *multipart.FileHeader, useFileName bool, publicId *string) uploader.UploadParams{
	var pubId string

	if publicId == nil {
		uniqueID := fmt.Sprintf("image_%s", time.Now().Format("20060102150405"))
		pubId = fileHeader.Filename+"_"+uniqueID
	} else {
		pubId = *publicId
	}
	
	return uploader.UploadParams{
		ResourceType:   "image",
		Folder:         folderName,
		PublicID:       pubId,
		UseFilename:    &useFileName,
		Transformation: "f_webp",
	}
}