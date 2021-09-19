package v1

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func WithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	return json.NewEncoder(w).Encode(payload)
}

func WithErrWrapJSON(w http.ResponseWriter, code int, in error) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := Response{
		Code:  code,
		Error: in.Error(),
	}

	return json.NewEncoder(w).Encode(response)
}
