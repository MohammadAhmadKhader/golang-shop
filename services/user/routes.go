package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"main.go/constants"
	appErrors "main.go/errors"
	"main.go/internal/websocket"
	"main.go/middlewares"
	"main.go/pkg/payloads"
	"main.go/pkg/utils"
	"main.go/services/auth"
	"main.go/types"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store Store) *Handler {
	return &Handler{
		store: &store,
	}
}

var Authenticate = middlewares.Authenticate
var AuthorizeSuperAdmin = middlewares.AuthorizeSuperAdmin

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc(utils.RoutePath("GET", "/users"), Authenticate(h.GetUserByToken))
	router.HandleFunc(utils.RoutePath("POST", "/users/login"), h.Login)
	router.HandleFunc(utils.RoutePath("POST", "/users/sign-up"), h.SignUp)
	router.HandleFunc(utils.RoutePath("PATCH", "/users/{id}/reset-password"), Authenticate(h.ResetPassword))
	router.HandleFunc(utils.RoutePath("PUT", "/users/{id}/profile"), Authenticate(h.UpdateProfile))
	router.HandleFunc(utils.RoutePath("POST", "/users/{id}/roles"), Authenticate(AuthorizeSuperAdmin(h.AssignUserRole)))
	router.HandleFunc(utils.RoutePath("DELETE", "/users/{id}/roles/{roleId}"), Authenticate(AuthorizeSuperAdmin(h.RemoveUserRole)))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	loginPayload, err := utils.ValidateAndParseBody[payloads.UserLogin](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, appErrors.ErrWrongPWOrEmail)
		return
	}
	loginPayload.TrimStrs()

	user, err := h.store.GetUserByEmail(loginPayload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, appErrors.ErrWrongPWOrEmail)
		return
	}
	isEqual := auth.ComparePassword(user.Password, []byte(loginPayload.Password))
	if !isEqual {
		utils.WriteError(w, http.StatusBadRequest, appErrors.ErrWrongPWOrEmail)
		return
	}

	token, err := auth.CreateJWT(*user, w, r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	otp := websocket.GlobalManager.Otps.NewOTP()

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"user":  user,
		"token": token,
		"otp": otp.Key,
	})
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	signUpPayload, err := utils.ValidateAndParseBody[payloads.UserSignUp](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	signUpPayload.TrimStrs()
	signUpPayload.Email = strings.ToLower(signUpPayload.Email)
	// hash password and create user
	user, err := h.store.CreateUser(*signUpPayload)
	if err != nil {
		if utils.IsDuplicateKeyErr(err) {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with this email already existing"))
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, appErrors.ErrGenericMessage)
		return
	}

	token, err := auth.CreateJWT(*user, w, r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, appErrors.ErrGenericMessage)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"user":  user,
		"token": token,
	})
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	rpPayload, err := utils.ValidateAndParseBody[payloads.ResetPassword](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	rpPayload.TrimStrs()
	if rpPayload.NewPassword == rpPayload.OldPassword {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("old and new password are same"))
		return
	}

	email, err := utils.GetEmailFromToken(r)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, appErrors.ErrUnauthorized)
		return
	}

	user, err := h.store.GetUserByEmail(*email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	isEqual := auth.ComparePassword(user.Password, []byte(rpPayload.OldPassword))
	if !isEqual {
		utils.WriteError(w, http.StatusBadRequest, appErrors.ErrPasswordsNotMatching)
		return
	}

	hashedPW, err := auth.HashPassword(rpPayload.NewPassword)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.UpdatePassword(hashedPW, *email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	newToken, err := auth.CreateJWT(*user, w, r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, map[string]any{
		"token": newToken,
	})
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIdCtx(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, appErrors.ErrGenericMessage)
		return
	}
	upPayload, err := utils.ValidateAndParseFormData(r, func() (*payloads.UpdateProfile, error) {
		payload := &payloads.UpdateProfile{
			Name:         r.FormValue("name"),
			Email:        strings.ToLower(r.FormValue("email")),
			MobileNumber: r.FormValue("mobileNumber"),
		}
		return payload, nil
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	size1MB := 1
	file, fileHeader, err := utils.HandleOneFileUpload(r, int64(size1MB), "avatar")
	if err != nil && !errors.Is(err, appErrors.ErrNoFileFound) {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	isEmpty := upPayload.TrimStrs().IsEmpty()
	if isEmpty {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("at least one of (name, email, mobile number) is required"))
		return
	}

	model := upPayload.ToModel()
	if file != nil {
		imgHandler := utils.NewImagesHandler()
		upResult, err := imgHandler.UploadOne(&file, fileHeader, types.UsersFolder, context.Background())
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, appErrors.ErrUnexpectedDuringImageUpload)
			return
		}
		// both them must have the avatar to ensure the field will not be excluded during the update
		model.Avatar = &upResult.URL
		upPayload.Avatar = upResult.URL
	}

	user, err := h.store.UpdateProfile(*userId, model, upPayload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	user.ID = *userId

	utils.WriteJSON(w, http.StatusAccepted, map[string]any{
		"user": user,
	})
}

func (h *Handler) AssignUserRole(w http.ResponseWriter, r *http.Request) {
	Id, err := utils.GetValidateId(r, constants.IdUrlPathKey)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, appErrors.NewInvalidIDError("user", *Id))
		return
	}

	rPayload, err := utils.ValidateAndParseBody[payloads.AssignRolePayload](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userRole, err := h.store.AssignUserRole(rPayload.RoleId, *Id)
	if err != nil {
		if utils.IsDuplicateKeyErr(err) {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user already has this role"))
			return
		}
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, map[string]any{"userRole": userRole})
}

func (h *Handler) RemoveUserRole(w http.ResponseWriter, r *http.Request) {
	Id, err := utils.GetValidateId(r, constants.IdUrlPathKey)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, appErrors.NewInvalidIDError("user", *Id))
		return
	}

	roleId, err := utils.GetValidateId(r, "roleId")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.store.RemoveUserRole(*roleId, *Id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, map[string]any{})
}

func (h *Handler) GetUserByToken(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUserCtx(r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, appErrors.ErrGenericMessage)
		return
	}
	
	otp := websocket.GlobalManager.Otps.NewOTP()

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"user":user,
		"otp":otp.Key,
	})
}