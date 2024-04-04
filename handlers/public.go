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
	Item string `json:"item"` // Use struct tags to match JSON keys
}

type PublicData struct {
	Item string
	Data map[string]string
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
	if p.Item == "" {
		ErrorResponse(w, "Key needed!", http.StatusBadRequest)
		return
	}

	// Respond with the key
	w.Header().Set("Content-Type", "application/json")

	result := h.db.Collection("content").FindOne(r.Context(), bson.M{"item": p.Item})

	// Check if the key is empty
	if result.Err() != nil {
		ErrorResponse(w, "item not found!", http.StatusBadRequest)
		return
	}
	var data PublicData

	result.Decode(&data)

	json.NewEncoder(w).Encode(map[string]interface{}{"payload": data})
}

func (*PublicHandler) Pattern() string {
	return "/public"
}
