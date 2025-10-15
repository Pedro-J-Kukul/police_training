// Filename: cmd/api/routes.go

package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// routes sets up the HTTP routes for the application
func (app *appDependencies) routes() http.Handler {
	router := httprouter.New() // create a new router instance

	// Error handling for unsupported methods
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Metrics endpoint
	router.Handler(http.MethodGet, "/v1/observability/quotes/metrics", expvar.Handler())

	// Return the router as an http.Handler
	return nil // temp: replace with 'router' when ready
}
