package utils

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func ParseJSON(req *http.Request, payload any) error {
	if req.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(req.Body).Decode(payload)
}

func WriteJSON(res http.ResponseWriter, status int, v any) error {
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(status)
	return json.NewEncoder(res).Encode(v)
}

func WriteError(res http.ResponseWriter, status int, err error) {
	// Use the colorful error writer
	S.WriteError(res, status, err)
}

func GetTokenFromRequest(r *http.Request) string {
	// First check Authorization header
	tokenAuth := r.Header.Get("Authorization")
	if tokenAuth != "" {
		// Handle "Bearer <token>" format
		if len(tokenAuth) > 7 && tokenAuth[:7] == "Bearer " {
			return tokenAuth[7:]
		}
		return tokenAuth
	}

	// Then check token query parameter
	tokenQuery := r.URL.Query().Get("token")
	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}