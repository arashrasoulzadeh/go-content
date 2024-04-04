package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

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

func (h *PublicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Ensure that the request method is POST
	if r.Method != http.MethodPost {
		ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Split the URL path by "/"
	pathParts := strings.Split(r.URL.Path, "/")

	// Get the last part of the path, which should be the language code
	languageCode := pathParts[len(pathParts)-1]

	// Respond with the key
	w.Header().Set("Content-Type", "application/json")

	filter := bson.M{"language_code": languageCode}
	content := make(map[string]string)
	err := h.db.Collection("content").FindOne(r.Context(), filter, nil).Decode(content)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			ErrorResponse(w, "not found", http.StatusNotFound)
			return
		}
		ErrorResponse(w, "internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(content)
}

func (*PublicHandler) Pattern() string {
	return "/public/{languageCode}"
}
