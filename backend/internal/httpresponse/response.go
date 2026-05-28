// Package httpresponse
package httpresponse

import (
	"encoding/json"
	"net/http"
)

func WriteError(w http.ResponseWriter, status int, message string, field string, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(map[string]any{
		"success": false,
		"error": map[string]any{
			"message": message,
			"field":   field,
			"code":    code,
		},
	})
}
