package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"main.go/constants"
	"main.go/pkg/utils"
	"main.go/services/auth"
	"main.go/types"
)

var AuthorizeAdmin = CreateAuthorizationMiddleware([]types.UserRole{types.SuperAdmin, types.Admin}, userLookup)
var AuthorizeSuperAdmin = CreateAuthorizationMiddleware([]types.UserRole{types.SuperAdmin}, userLookup)

func CreateAuthorizationMiddleware(allowedRoles []types.UserRole, userFetcher types.IUserFetcher) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userId, err := utils.GetUserIdFromTokenPayload(r)
			if err != nil {
				auth.DenyPermission(w)
				return
			}

			userRoles, err := userFetcher.GetUserRolesByUserId(*userId)
			if err != nil {
				auth.Unauthorized(w)
				return
			}

			for _, userRole := range userRoles {
				isContaining := slices.Contains(allowedRoles, types.UserRole(userRole.Role.Name))
				if !isContaining {
					auth.Unauthorized(w)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

type UserIdGetter interface {
	GetUserId() uint
}


// not found error expected to have a place holder ("%v") inside the message
func AuthorizeOwnerShipMW[TModel UserIdGetter](param, idIsRequiredErrMsg, notFoundErr string, modelGetter func(id uint) (TModel,error)) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
			resourceId, err := utils.GetValidateId(r, param)
			if err != nil {
				utils.WriteError(w, http.StatusBadRequest, errors.New(idIsRequiredErrMsg))
				return
			}

			userIdToken, err := utils.GetUserIdFromTokenPayload(r)
			if err != nil {
				auth.Unauthorized(w)
				return
			}

			model, err := modelGetter(*resourceId)
			if err != nil {
				utils.WriteError(w, http.StatusBadRequest, fmt.Errorf(notFoundErr, *resourceId))
				return
			}
			if model.GetUserId() != *userIdToken {
				auth.DenyPermission(w)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, constants.UserKey, model)
			ctx = context.WithValue(ctx, constants.ResourceKey, model)
			
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// * Deprecated
func AuthorizeUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := utils.GetValidateId(r, constants.IdUrlPathKey)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user id"))
			return
		}
	
		userIdToken, err := utils.GetUserIdFromTokenPayload(r)
		if err != nil {
			auth.Unauthorized(w)
			return
		}

		if *userId != *userIdToken {
			auth.DenyPermission(w)
			return
		}

		//ctx := context.WithValue(r.Context(), constants.UserKey , *userIdToken)
		//r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}