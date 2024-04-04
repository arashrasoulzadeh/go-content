package handlers

import (
	"encoding/json"
	"net/http"
)

type ErrorStruct struct {
	Error string `json:"error"`
}

func ErrorResponse(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorStruct{Error: msg})
}
