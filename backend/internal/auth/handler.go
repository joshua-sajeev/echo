package auth

import (
	"encoding/json"
	"net/http"

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
	session, err := h.Store.Get(
		r,
		SessionName,
	)
	if err == nil {
		session.Options.MaxAge = -1
		_ = session.Save(r, w)
	}

	w.WriteHeader(http.StatusNoContent)
}
