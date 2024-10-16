package auth

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"main.go/config"
	"main.go/pkg/models"
	"main.go/pkg/utils"
)

func DenyPermission(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("forbidden"))
}

func Unauthorized(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
}

func CreateJWT(user models.User, w http.ResponseWriter, r *http.Request) (string, error) {
	secret := config.Envs.JWT_SECRET
	secretAsBytes := []byte(secret)
	durationInt, err := strconv.Atoi(config.Envs.JWT_EXPIRATION_IN_SECONDS)
	if err != nil {
		return "", err
	}

	expiration := time.Second * time.Duration(durationInt)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    user.ID,
		"email":     user.Email,
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secretAsBytes)
	if err != nil {
		return "", err
	}

	SetCookie(w, r, &user, tokenString)
	return tokenString, nil
}

func ValidateToken(tokenString *string) (*jwt.Token, error) {
	token, err := jwt.Parse(*tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unauthenticated")
		}

		return []byte(config.Envs.JWT_SECRET), nil
	})

	jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name})

	if err != nil {
		return nil, err
	}

	isValidToken := token.Valid
	if !isValidToken {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func GetToken(r *http.Request) (*string, error) {
	session, err := GetCookie(r)
	if err != nil {
		return nil, fmt.Errorf("unauthorized")
	}

	token := session.Values["token"]
	var tokenAsString string
	tokenAsString, ok := token.(string); 
	if !ok {
		return nil, fmt.Errorf("error during parsing token")
	}

	return &tokenAsString, nil
}

func GetUserIdFromJWT(jwtToken *jwt.Token) (*uint, jwt.MapClaims, error) {
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	var userId uint
	if !ok {
		return &userId, nil, fmt.Errorf("error has occurred during parsing token process")
	}

	userIdInt, err := GetUserIdFromClaims(claims)
	if err != nil {
		return nil, nil, err
	}
	
	return userIdInt, claims,nil
}

// this jwt library has issues with numbers assertions from token.
//
// they recommend asserting all the numbers as float64, but sometimes it throws an error as float64 and it works as string then
// when it refuses to work with float64 assertion somehow it manage to work with string assertion and vice versa, 
// therefore this type of assertion was made.
func GetUserIdFromClaims(claims jwt.MapClaims) (*uint, error) {
	var userId uint
	userIdAsStr, ok := claims["userId"].(string)
	if !ok {
		userIdAsFloat, okAsFloat := claims["userId"].(float64)
		if !okAsFloat {
			return &userId, fmt.Errorf("user id not exist in token")
		}
		userIdAsInt := uint(userIdAsFloat)
		return &userIdAsInt, nil
	}
	userIdInt, err := strconv.Atoi(userIdAsStr)
	if err != nil {
		return &userId, fmt.Errorf("invalid user id type")
	}
	userId = uint(userIdInt)

	return &userId, nil
}