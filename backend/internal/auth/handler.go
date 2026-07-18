package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joshu-sajeev/echo/internal/httpresponse"
)

type Handler struct {
	Store *sessions.CookieStore
}

func NewHandler(
	store *sessions.CookieStore,
) *Handler {
	return &Handler{
		Store: store,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(
			w,
			http.StatusBadRequest,
			"invalid request body",
			"",
			"INVALID_REQUEST_BODY",
		)
		return
	}

	if !VerifyPIN(req.PIN) {
		httpresponse.WriteError(
			w,
			http.StatusUnauthorized,
			"invalid pin",
			"pin",
			"INVALID_PIN",
		)
		return
	}

	session, err := h.Store.Get(
		r,
		SessionName,
	)
	if err != nil {
		httpresponse.WriteError(
			w,
			http.StatusInternalServerError,
			"session error",
			"",
			"SESSION_ERROR",
		)
		return
	}

	session.Values["authenticated"] = true

	if err := session.Save(r, w); err != nil {
		httpresponse.WriteError(
			w,
			http.StatusInternalServerError,
			"failed to create session",
			"",
			"SESSION_ERROR",
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Logout(
	w http.ResponseWriter,
	r *http.Request,
) {
	session, err := h.Store.Get(r, SessionName)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	log.Printf("LOGOUT BEFORE: %#v", session.Values)

	session.Values = make(map[any]any)
	session.Options.MaxAge = -1

	log.Printf("LOGOUT AFTER: %#v", session.Values)

	if err := session.Save(r, w); err != nil {
		httpresponse.WriteError(
			w,
			http.StatusInternalServerError,
			"failed to destroy session",
			"",
			"SESSION_ERROR",
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Me(
	w http.ResponseWriter,
	r *http.Request,
) {
	if os.Getenv("DEMO_MODE") == "true" {
		httpresponse.WriteJSON(
			w,
			http.StatusOK,
			map[string]bool{
				"authenticated": true,
			},
		)
		return
	}

	session, err := h.Store.Get(r, SessionName)
	if err != nil {
		log.Printf("ME: session error: %v", err)

		httpresponse.WriteJSON(
			w,
			http.StatusOK,
			map[string]bool{
				"authenticated": false,
			},
		)
		return
	}

	authenticated, _ := session.Values["authenticated"].(bool)
	httpresponse.WriteJSON(
		w,
		http.StatusOK,
		map[string]bool{
			"authenticated": authenticated,
		},
	)
}
