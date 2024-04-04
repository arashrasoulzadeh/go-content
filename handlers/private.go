package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	// Split the URL path by "/"
	pathParts := strings.Split(r.URL.Path, "/")

	// Get the last part of the path, which should be the language code
	languageCode := pathParts[len(pathParts)-1]

	var p map[string]string

	// Decode the request body into the struct
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		ErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Respond with the data
	w.Header().Set("Content-Type", "application/json")

	filter := bson.M{"language_code": languageCode}
	update := bson.M{
		"$set": p,
	}
	opts := options.Update().SetUpsert(true)
	_, err = h.db.Collection("content").UpdateOne(r.Context(), filter, update, opts)
	if err != nil {
		ErrorResponse(w, "cant update", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"status": "done"})
}

func (*PrivateHandler) Pattern() string {
	return "/private/{languageCode}"
}
