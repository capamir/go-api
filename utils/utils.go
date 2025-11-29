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
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")
	
	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}