package image

import (
	"context"
	"fmt"
	"net/http"

	"main.go/constants"
	"main.go/middlewares"
	"main.go/pkg/utils"
	"main.go/types"
)

type Handler struct {
	store types.ImageStore
}

func NewHandler(store Store) *Handler {
	return &Handler{
		store: &store,
	}
}
const (
	maxImagesCapacity = 10
	maxSizeInMBForMultipleUploads int64 = 10
	maxSizeInMBForOneUpload int64 = 2
)

var Authenticate = middlewares.Authenticate
var AuthorizeAdmin = middlewares.AuthorizeAdmin

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc(utils.RoutePath("DELETE", "/products/{id}/images/{imageId}"), Authenticate(AuthorizeAdmin(h.DeleteProductImageById)))
	router.HandleFunc(utils.RoutePath("PUT", "/products/{id}/images/{imageId}"), Authenticate(AuthorizeAdmin(h.UpdateImageById)))
	router.HandleFunc(utils.RoutePath("PATCH", "/products/{id}/images/{imageId}"), Authenticate(AuthorizeAdmin(h.SetImageProductAsMain)))
	router.HandleFunc(utils.RoutePath("POST", "/products/{id}/images"), Authenticate(AuthorizeAdmin(h.CreateImagesForProduct)))
}

func (h *Handler) DeleteProductImageById(w http.ResponseWriter, r *http.Request) {
	imageId, err := utils.GetValidateId(r, "imageId")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}

	err = h.store.DeleteImageById(*imageId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, map[string]any{})
}

func (h *Handler) SetImageProductAsMain(w http.ResponseWriter, r *http.Request) {
	imageId, err := utils.GetValidateId(r, "imageId")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}

	image, err := h.store.GetImageById(*imageId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if *image.IsMain {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("image is already a main image"))
		return
	}

	err = h.store.SwapMainStatus(*imageId, image.ProductID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("an unexpected error has occurred during transaction"))
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, map[string]any{})
}

func (h *Handler) CreateImagesForProduct(w http.ResponseWriter, r *http.Request) {
	productId, err := utils.GetValidateId(r, constants.IdUrlPathKey)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}

	_, err = h.store.GetProductById(*productId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("product with id: '%v' was not found", *productId))
		return
	}

	count ,err := h.store.GetCountOfProductImages(*productId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	
	filesCount ,err := utils.GetFilesCount(r, "images")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	sizeInMB := int64(10)
	files, err:= utils.HandleMultipleFilesUpload(r, sizeInMB, "images")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	
	if int(*count) + filesCount > maxImagesCapacity {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("you can not upload more than 10 images for a single product, you have '%v' and you trying to upload more '%v'", *count, len(files)))
		return
	}

	imgHandler := utils.NewImagesHandler()
	responses, errs := imgHandler.UploadMany(r, types.ProductsFolder, files, context.Background())
	if len(errs) != 0 {
		utils.WriteError(w, http.StatusInternalServerError, errs[0])
		return
	}
	newImages, err := h.store.CreateManyImages(responses, productId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"uploaded-data":responses,
		"images":newImages,
	})
}

func (h *Handler) UpdateImageById(w http.ResponseWriter, r *http.Request) {
	imageId, err := utils.GetValidateId(r, "imageId")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}

	image, err := h.store.GetImageById(*imageId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("image with id: '%v' was not found", imageId))
		return
	}

	file,fileHeader, err := utils.HandleOneFileUpload(r, maxSizeInMBForOneUpload, "image")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	imgHandler := utils.NewImagesHandler()
	upResult, err := imgHandler.UploadOne(&file ,fileHeader, types.ProductsFolder, context.Background())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.UpdateImageUrl(*imageId, upResult.URL)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	image.ImageUrl = upResult.URL
	utils.WriteJSON(w, http.StatusAccepted, map[string]any{
		"image":image,
	})
}