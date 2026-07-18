// Package auth handles authentication
package auth

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joshu-sajeev/echo/internal/httpresponse"
)

func RequireAuth(
	store *sessions.CookieStore,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				if os.Getenv("DEMO_MODE") == "true" {
					next.ServeHTTP(w, r)
					return
				}

				session, err := store.Get(
					r,
					SessionName,
				)
				if err != nil {
					httpresponse.WriteError(
						w,
						http.StatusUnauthorized,
						"unauthorized",
						"",
						"UNAUTHORIZED",
					)
					return
				}

				authenticated, ok := session.Values["authenticated"].(bool)

				if !ok || !authenticated {
					httpresponse.WriteError(
						w,
						http.StatusUnauthorized,
						"unauthorized",
						"",
						"UNAUTHORIZED",
					)
					return
				}

				next.ServeHTTP(w, r)
			},
		)
	}
}
