// Package httpresponse
package httpresponse

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteError(w http.ResponseWriter, status int, message string, field string, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(map[string]any{
		"success": false,
		"error": map[string]any{
			"message": message,
			"field":   field,
			"code":    code,
		},
	})
	if err != nil {
		log.Printf("encode error json response: %v", err)
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf(
			"encode json response: %v",
			err,
		)
	}
}
