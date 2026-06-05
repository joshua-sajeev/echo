package auth

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

const SessionName = "echo_session"

func NewStore() *sessions.CookieStore {
	store := sessions.NewCookieStore(
		[]byte(os.Getenv("SESSION_SECRET")),
	)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
	}

	return store
}
