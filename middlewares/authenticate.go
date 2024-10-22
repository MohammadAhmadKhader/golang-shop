package middlewares

import (
	"context"
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

// it skips the authentication if no token found, this meant to be used with websocket handler to allow public clients and non public clients.
func AuthenticateIfCookieExist(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetToken(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		jwtToken, err := auth.ValidateToken(token)
		if err != nil {
			auth.Unauthorized(w)
			return
		}

		userId, claims, err := auth.GetUserIdFromJWT(jwtToken)
		if err != nil {
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
