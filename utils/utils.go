package utils

import (
	"fmt"
	"net/http"
	"encoding/json"
)

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
	WriteJSON(res, status, map[string]string{"error": err.Error()})
}