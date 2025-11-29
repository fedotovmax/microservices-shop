package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func DecodeJSON(r io.Reader, v any) error {
	defer io.Copy(io.Discard, r)
	return json.NewDecoder(r).Decode(v)
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, `{"error": "failed to encode json"}`, http.StatusInternalServerError)
	}
}
