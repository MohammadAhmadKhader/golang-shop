package auth

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"main.go/config"
	"main.go/errors"
	"main.go/pkg/models"
	"main.go/pkg/utils"
)

func DenyPermission(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, errors.ErrForbidden)
}

func Unauthorized(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusUnauthorized, errors.ErrUnauthorized)
}

func ValidateToken(tokenString *string) (*jwt.Token, error) {
	token, err := jwt.Parse(*tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrInvalidToken
		}

		return []byte(config.Envs.JWT_SECRET), nil
	})

	jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name})

	if err != nil {
		return nil, err
	}

	isValidToken := token.Valid
	if !isValidToken {
		return nil, errors.ErrInvalidToken
	}

	return token, nil
}

func GetUserIdFromJWT(jwtToken *jwt.Token) (*uint, jwt.MapClaims, error) {
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	var userId uint
	if !ok {
		return &userId, nil, errors.ErrGenericMessage
	}

	userIdInt, err := GetUserIdFromClaims(claims)
	if err != nil {
		return nil, nil, err
	}

	return userIdInt, claims,nil
}

func GetUserIdFromClaims(claims jwt.MapClaims) (*uint, error) {
	var userId uint
	userIdAsFloat, ok := claims["userId"].(float64)
	if !ok {
		return &userId, fmt.Errorf("user id not exist in token")
	}
	userId = uint(userIdAsFloat)

	return &userId, nil
}

// Generate access_token and refresh_token and set them to the cookie
func GenerateAndSetTokens(user models.User, w http.ResponseWriter, r *http.Request) (access_token string, refresh_token string, err error) {
	secret := config.Envs.JWT_SECRET
	secretAsBytes := []byte(secret)
	accessToken, err := createAccessToken(&user, secretAsBytes)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := createRefreshToken(&user, secretAsBytes)
	if err != nil {
		return "", "", err
	}

	SetCookie(w, r, &user, accessToken, refreshToken)
	return accessToken, refreshToken, nil
}

func GetRefreshToken(r *http.Request)(*string, error) {
	session, err := GetCookie(r)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}

	token := session.Values["refresh_token"]
	var tokenAsString string
	tokenAsString, ok := token.(string); 
	if !ok {
		return nil, errors.ErrGenericMessage
	}

	return &tokenAsString, nil
}

func GetAccessToken(r *http.Request) (*string, error) {
	session, err := GetCookie(r)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}

	token := session.Values["access_token"]
	var tokenAsString string
	tokenAsString, ok := token.(string); 
	if !ok {
		return nil, errors.ErrGenericMessage
	}

	return &tokenAsString, nil
}

// generates access token and sets the cookie with the refresh token and the new generated access token.
func GenerateAccessTokenAndSetTokens(refreshToken string ,user *models.User, w http.ResponseWriter, r *http.Request) (string, error) {
	secret := config.Envs.JWT_SECRET
	secretAsBytes := []byte(secret)
	accessToken, err := createAccessToken(user, secretAsBytes)
	if err != nil {
		return "", err
	}

	SetCookie(w, r, user, accessToken, refreshToken)
	return accessToken, nil
}

func createAccessToken(user *models.User, secret []byte) (string, error){
	accessDurationInt, err := strconv.Atoi(config.Envs.ACCESS_JWT_EXPIRATION_IN_SECONDS)
	if err != nil {
		return "", err
	}
	accessExpiration := time.Second * time.Duration(accessDurationInt)
	accessClaims := jwt.MapClaims{
		"userId":    user.ID,
		"email":     user.Email,
		"expiredAt": time.Now().Add(accessExpiration).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	stringAccessToken, err := accessToken.SignedString(secret)
	if err != nil {
		return "", err
	}

	return stringAccessToken, nil
}

func createRefreshToken(user *models.User, secret []byte) (string, error){
	refreshDurationInt, err := strconv.Atoi(config.Envs.REFRESH_JWT_EXPIRATION_IN_SECONDS)
	if err != nil {
		return "",err
	}
	refreshExpiration := time.Second * time.Duration(refreshDurationInt)
	refreshClaims := jwt.MapClaims{
		"userId":    user.ID,
		"expiredAt": time.Now().Add(refreshExpiration).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(secret)
	if err != nil {
		return "", err
	}

	return refreshTokenString, nil
}