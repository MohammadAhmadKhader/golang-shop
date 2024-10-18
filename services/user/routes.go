package user

import (
	"fmt"
	"net/http"
	"strings"

	"main.go/errors"
	"main.go/middlewares"
	"main.go/pkg/payloads"
	"main.go/pkg/utils"
	"main.go/services/auth"
)

type Handler struct {
	store Store
}

func NewHandler(store Store) *Handler {
	return &Handler{
		store: store,
	}
}

var Authenticate = middlewares.Authenticate
var AuthorizeSuperAdmin = middlewares.AuthorizeSuperAdmin

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc(utils.RoutePath("POST", "/users/login"), h.Login)
	router.HandleFunc(utils.RoutePath("POST", "/users/sign-up"), h.SignUp)
	router.HandleFunc(utils.RoutePath("PUT", "/users/reset-password"), Authenticate(h.ResetPassword))
	router.HandleFunc(utils.RoutePath("PUT", "/users/{id}/profile"), Authenticate(h.UpdateProfile))
	router.HandleFunc(utils.RoutePath("POST", "/users/{id}/role"), AuthorizeSuperAdmin(Authenticate(h.AssignUserRole)))
	router.HandleFunc(utils.RoutePath("DELETE", "/users/{id}/role"), AuthorizeSuperAdmin(Authenticate(h.RemoveUserRole)))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	loginPayload, err := utils.ValidateAndParseBody[payloads.UserLogin](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.ErrWrongPWOrEmail)
		return
	}
	loginPayload.TrimStrs()
	
	user, err := h.store.GetUserByEmail(loginPayload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errors.ErrWrongPWOrEmail)
		return
	}
	isEqual := auth.ComparePassword(user.Password, []byte(loginPayload.Password))
	if !isEqual {
		utils.WriteError(w, http.StatusBadRequest, errors.ErrWrongPWOrEmail)
		return
	}

	token, err := auth.CreateJWT(*user, w, r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"user":  user,
		"token": token,
	})
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	signUpPayload, err := utils.ValidateAndParseBody[payloads.UserSignUp](r)
	if err != nil {
		fmt.Println(signUpPayload)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	signUpPayload.TrimStrs()
	signUpPayload.Email = strings.ToLower(signUpPayload.Email)
	// hash password and create user
	user, err := h.store.CreateUser(*signUpPayload)

	if err != nil {
		if utils.IsDuplicateKeyErr(err) {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("user with this email already existing"))
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, errors.ErrGenericMessage)
		return
	}

	token, err := auth.CreateJWT(*user, w, r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errors.ErrGenericMessage)
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

	email, err := utils.GetEmailFromToken(r)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.ErrUnauthorized)
		return
	}

	user, err := h.store.GetUserByEmail(*email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	isEqual := auth.ComparePassword(rpPayload.Password, []byte(user.Password))
	if !isEqual {
		utils.WriteError(w, http.StatusBadRequest, errors.ErrWrongPWOrEmail)
		return
	}

	hashedPW, err := auth.HashPassword(rpPayload.Password)
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
	Id, err := utils.GetValidateId(r, "id")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}
	upPayload, err := utils.ValidateAndParseFormData[payloads.UpdateProfile](r, func() (*payloads.UpdateProfile, error) {
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
	isEmpty := upPayload.TrimStrs().IsEmpty()
	if isEmpty {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("at least one of (name, email, mobile number) is required"))
		return
	}

	userId, err := utils.GetUserIdCtx(r)
	if err != nil || *userId != *Id {
		auth.DenyPermission(w)
		return
	}

	user, err := h.store.UpdateProfile(*Id, upPayload.ToModel(), upPayload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	user.ID = *Id

	utils.WriteJSON(w, http.StatusAccepted, map[string]any{
		"user":    user,
	})
}

func (h *Handler) AssignUserRole(w http.ResponseWriter, r *http.Request) {
	Id, err := utils.GetValidateId(r, "id")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}

	rPayload ,err:= utils.ValidateAndParseBody[payloads.AssignRolePayload](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.store.AssignUserRole(rPayload.RoleId, *Id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, map[string]any{})
}

func (h *Handler) RemoveUserRole(w http.ResponseWriter, r *http.Request) {
	Id, err := utils.GetValidateId(r, "id")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}

	rPayload ,err:= utils.ValidateAndParseBody[payloads.RemoveRolePayload](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.store.RemoveUserRole(rPayload.RoleId, *Id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, map[string]any{})
}