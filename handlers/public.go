package handlers

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// EchoHandler is an http.Handler that copies its request body
// back to the response.
type PublicHandler struct {
	log *zap.Logger
	db  *mongo.Database
}

// NewPublicHandler builds a new PublicHandler.
func NewPublicHandler(log *zap.Logger, db *mongo.Database) *PublicHandler {
	return &PublicHandler{
		log: log,
		db:  db,
	}
}

type PublicRequest struct {
	Key string `json:"key"` // Use struct tags to match JSON keys
}
type ErrorStruct struct {
	Error string `json:"error"`
}

func (h *PublicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Ensure that the request method is POST
	if r.Method != http.MethodPost {
		ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var p PublicRequest

	// Decode the request body into the struct
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		ErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Check if the key is empty
	if p.Key == "" {
		ErrorResponse(w, "Key needed!", http.StatusBadRequest)
		return
	}

	// Respond with the key
	w.Header().Set("Content-Type", "application/json")

	cursor, err := h.db.Collection("conent").Find(r.Context(), bson.M{"key": p.Key})

	// Check if the key is empty
	if err != nil {
		ErrorResponse(w, "Key not found!", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"key": cursor.Current.String()})
}

func ErrorResponse(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorStruct{Error: msg})
}

func (*PublicHandler) Pattern() string {
	return "/public"
}
