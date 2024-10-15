package auth

import (
	"net/http"

	"github.com/gorilla/sessions"
	"main.go/config"
	_"main.go/constants"
	"main.go/pkg/models"
)

const CookieMaxAge = 1209600 // 14 days

var cookiesStore = sessions.NewCookieStore([]byte(config.Envs.JWT_SECRET))

func GetCookie(r *http.Request) (*sessions.Session, error) {
	session, err := cookiesStore.Get(r, "session_token")
	if err != nil {
		return nil, err
	}

	return session, nil
}

func SetCookie(w http.ResponseWriter, r *http.Request, user *models.User, token string) (*sessions.Session, error) {
	session, err := cookiesStore.New(r, "session_token")
	if err != nil {
		return nil, err
	}

	session.Options = &sessions.Options{
		MaxAge: CookieMaxAge,
		//Path:   constants.Prefix,
		//HttpOnly: true,
		//Secure: false,
		//SameSite: http.SameSiteStrictMode,
	}

	session.Values["userId"] = user.ID
	session.Values["email"] = user.Email
	session.Values["token"] = token

	err = session.Save(r, w)
	if err != nil {
		return nil, err
	}

	return session, nil
}
