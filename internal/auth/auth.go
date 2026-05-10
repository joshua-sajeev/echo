// Package auth provides simple PIN-based authentication with
// a signed HttpOnly session cookie. No external deps needed.
package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"time"
)

const (
	cookieName = "echo_session"
	cookieTTL  = 30 * 24 * time.Hour // 30 days
)

// sessionSecret is derived from SESSION_SECRET env var.
// Falls back to a random value (sessions won't survive restarts without the env var).
var sessionSecret = func() []byte {
	s := os.Getenv("SESSION_SECRET")
	if s != "" {
		return []byte(s)
	}
	b := make([]byte, 32)
	rand.Read(b)
	return b
}()

// sign creates an HMAC-SHA256 signature for the given value.
func sign(value string) string {
	mac := hmac.New(sha256.New, sessionSecret)
	mac.Write([]byte(value))
	return hex.EncodeToString(mac.Sum(nil))
}

// SetSession writes the authenticated session cookie.
func SetSession(w http.ResponseWriter) {
	// token = timestamp.signature
	ts := time.Now().Format(time.RFC3339)
	sig := sign(ts)
	token := ts + "." + sig

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(cookieTTL.Seconds()),
		// Set Secure: true if you're on HTTPS (recommended in production)
		// Secure: true,
	})
}

// ClearSession deletes the session cookie.
func ClearSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    cookieName,
		Value:   "",
		Path:    "/",
		MaxAge:  -1,
		Expires: time.Unix(0, 0),
	})
}

// IsAuthenticated checks whether the request carries a valid session cookie.
func IsAuthenticated(r *http.Request) bool {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return false
	}

	// split token into ts.sig
	v := c.Value
	dot := len(v) - 65 // sha256 hex = 64 chars + 1 dot
	if dot <= 0 {
		return false
	}
	ts := v[:dot]
	sig := v[dot+1:]

	// verify signature
	if !hmac.Equal([]byte(sign(ts)), []byte(sig)) {
		return false
	}

	// verify not expired
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return false
	}
	return time.Since(t) < cookieTTL
}

// Middleware redirects unauthenticated requests to /login.
// HTMX requests get a special header so the client can redirect.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// always allow login page and its POST
		if r.URL.Path == "/login" || r.URL.Path == "/logout" {
			next.ServeHTTP(w, r)
			return
		}

		if !IsAuthenticated(r) {
			// HTMX requests need HX-Redirect, not a 302
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Redirect", "/login")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
