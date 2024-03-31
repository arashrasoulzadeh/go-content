package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"go.uber.org/zap"
)

// EchoHandler is an http.Handler that copies its request body
// back to the response.
type PublicHandler struct {
	log *zap.Logger
}

// NewPublicHandler builds a new PublicHandler.
func NewPublicHandler(log *zap.Logger) *PublicHandler {
	return &PublicHandler{
		log: log,
	}
}

// ServeHTTP handles an HTTP request to the /public endpoint.
func (*PublicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := io.Copy(w, r.Body); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to handle request:", err)
	}
}

func (*PublicHandler) Pattern() string {
	return "/public"
}
