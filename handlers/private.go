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
type PrivateHandler struct {
	log *zap.Logger
	db  *mongo.Database
}

// NewPrivateHandler builds a new PrivateHandler.
func NewPrivateHandler(log *zap.Logger, db *mongo.Database) *PrivateHandler {
	return &PrivateHandler{
		log: log,
		db:  db,
	}
}

type PrivateRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Item  string `json:"item"`
}

func (h *PrivateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Ensure that the request method is POST
	if r.Method != http.MethodPost {
		ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var p PrivateRequest

	// Decode the request body into the struct
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		ErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Check if the data is empty
	if p.Item == "" {
		ErrorResponse(w, "item needed!", http.StatusBadRequest)
		return
	}
	// Check if the key is empty
	if p.Key == "" {
		ErrorResponse(w, "key needed!", http.StatusBadRequest)
		return
	}
	// Check if the value is empty
	if p.Key == "" {
		ErrorResponse(w, "value needed!", http.StatusBadRequest)
		return
	}

	// Respond with the data
	w.Header().Set("Content-Type", "application/json")

	result := h.db.Collection("content").FindOne(r.Context(), bson.M{"item": p.Item})

	// Check if the data is empty
	if result.Err() != nil {

		document := bson.M{"item": p.Item, "data": bson.M{p.Key: p.Value}}

		// Insert the document into the collection
		doc, err := h.db.Collection("content").InsertOne(r.Context(), document)
		if err != nil {
			ErrorResponse(w, "Failed to create document:"+err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "created new item", "id": doc.InsertedID})
		return
	}

	var data interface{}
	result.Decode(&data)
	update := bson.M{"$set": bson.M{"data." + p.Key: p.Value}}

	updateResult, err := h.db.Collection("content").UpdateOne(r.Context(), bson.M{"item": p.Item}, update)
	if err != nil {
		ErrorResponse(w, "Failed to update document:"+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "updated old item!", "data": updateResult})
}

func (*PrivateHandler) Pattern() string {
	return "/private"
}
