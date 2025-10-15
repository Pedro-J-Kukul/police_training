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

	// Health and observability
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheckHandler)
	router.Handler(http.MethodGet, "/v1/observability/metrics", expvar.Handler())

	// Authentication and user lifecycle
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/password-reset", app.createPasswordResetTokenHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/password-reset", app.resetPasswordHandler)

	// Authenticated user endpoints
	router.Handler(http.MethodGet, "/v1/users/me", app.requireActivatedUser(http.HandlerFunc(app.showCurrentUserHandler)))
	router.Handler(http.MethodGet, "/v1/users/:id", app.requireActivatedUser(app.requireRole("admin", http.HandlerFunc(app.showUserHandler))))
	router.Handler(http.MethodGet, "/v1/users", app.requireActivatedUser(app.requireRole("admin", http.HandlerFunc(app.listUsersHandler))))
	router.Handler(http.MethodPatch, "/v1/users/:id", app.requireActivatedUser(app.requireRole("admin", http.HandlerFunc(app.updateUserHandler))))

	return app.recoverPanic(app.enableCORS(app.metrics(app.rateLimit(app.authenticate(router)))))
}
