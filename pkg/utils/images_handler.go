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
	"main.go/types"
)

type ImagesHandler struct {
	cld *cloudinary.Cloudinary
}

func NewImagesHandler() *ImagesHandler {
	return &ImagesHandler{
		cld: config.GetCloudinary(),
	}
}

func (ih *ImagesHandler) UpdateOne(newFile *multipart.File, newFileHeader *multipart.FileHeader, folder types.Folder, ctx context.Context, oldImagePublicId string) (*types.UploadResponse, error) {
	uploadParams := NewUploadParams(string(folder), newFileHeader, true, &oldImagePublicId)
	resp, err := ih.cld.Upload.Upload(ctx, *newFile, uploadParams)
	if err != nil {
		return nil, err
	}

	return ih.getUploadResponse(resp), nil
}

func (ih *ImagesHandler) UploadOne(file *multipart.File, fileHeader *multipart.FileHeader, folder types.Folder, ctx context.Context) (*types.UploadResponse, error) {
	folderAsString := string(folder)

	uploadParams := NewUploadParams(folderAsString, fileHeader, true, nil)
	resp, err := ih.cld.Upload.Upload(ctx, *file, uploadParams)
	if err != nil {
		return nil, err
	}

	return ih.getUploadResponse(resp), nil
}

func (ih *ImagesHandler) UploadMany(r *http.Request, folder types.Folder,files []*multipart.FileHeader, ctx context.Context) ([]*types.UploadResponse, []error) {
	folderAsString := string(folder)
	var mu sync.Mutex
	var wg sync.WaitGroup
	UploadResponses := make([]*types.UploadResponse, 0, len(files))
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

func (ih *ImagesHandler) DeleteMany(PublicIDs []string, ctx context.Context) (*types.DeleteManyResponse, error) {
	delResp, err := ih.cld.Admin.DeleteAssets(ctx, admin.DeleteAssetsParams{
		AssetType: "image",
		PublicIDs: PublicIDs,
	})
	if err != nil {
		return nil, err
	}

	return ih.getDeleteManyResponse(delResp), nil
}

func (ih *ImagesHandler) getUploadResponse(uploadResult *uploader.UploadResult) *types.UploadResponse {
	return &types.UploadResponse{
		Width: uploadResult.Width, Height: uploadResult.Height,
		SecureUrl: uploadResult.SecureURL, PublicID: uploadResult.PublicID,
		URL: uploadResult.URL, Format: uploadResult.Format,
	}
}

func (ih *ImagesHandler) getDeleteManyResponse(deleteResult *admin.DeleteAssetsResult) *types.DeleteManyResponse {
	return &types.DeleteManyResponse{
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