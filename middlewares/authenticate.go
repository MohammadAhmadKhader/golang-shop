package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"main.go/constants"
	"main.go/services/auth"
)


var userLookup = NewUserLookup()

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetToken(r)
		if err != nil {
			auth.Unauthorized(w)
			return
		}

		jwtToken, err := auth.ValidateToken(token)
		if err != nil {
			auth.Unauthorized(w)
			return
		}

		userId, claims, err := auth.GetUserIdFromJWT(jwtToken)
		if err != nil {
			fmt.Println(err)
			auth.Unauthorized(w)
			return
		}
		
		user, err := userLookup.GetUserById(*userId)
		if err != nil {
			auth.Unauthorized(w)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, constants.UserKey, user)
		ctx = context.WithValue(ctx, constants.TokenPayload, claims)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})

}
