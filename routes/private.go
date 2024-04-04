package routes

import "net/http"

// Route is an http.Handler that knows the mux pattern
// under which it will be registered.
type PrivateRoute interface {
	http.Handler

	// Pattern reports the path at which this is registered.
	Pattern() string
}
