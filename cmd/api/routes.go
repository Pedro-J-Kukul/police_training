// Filename: cmd/api/routes.go

package main

import (
	"expvar"
	"net/http"

	_ "github.com/Pedro-J-Kukul/police_training/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/julienschmidt/httprouter"
)

// routes sets up the HTTP routes for the application
func (app *appDependencies) routes() http.Handler {
	router := httprouter.New() // create a new router instance

	// Error handling for unsupported methods
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.Handler(http.MethodGet, "/swagger/:any", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), //The url pointing to API definition
	))
	// Health and observability
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheckHandler)
	router.Handler(http.MethodGet, "/v1/observability/metrics", expvar.Handler())

	// Authentication and user lifecycle
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activate", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/password-reset", app.createPasswordResetTokenHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/password-reset", app.resetPasswordHandler)

	// Authenticated user endpoints
	router.Handler(http.MethodGet, "/v1/me", app.requireActivatedUser(http.HandlerFunc(app.showCurrentUserHandler))) // this static route was conflicting with /v1/users/:id
	router.Handler(http.MethodGet, "/v1/users", app.requireActivatedUser(app.requireRole("admin", http.HandlerFunc(app.listUsersHandler))))
	router.Handler(http.MethodGet, "/v1/users/:id", app.requireActivatedUser(app.requireRole("admin", http.HandlerFunc(app.showUserHandler))))
	router.Handler(http.MethodPatch, "/v1/users/:id", app.requireActivatedUser(app.requireRole("admin", http.HandlerFunc(app.updateUserHandler))))
	router.Handler(http.MethodDelete, "/v1/users/:id", app.requireActivatedUser(app.requireRole("admin", http.HandlerFunc(app.deleteUserHandler))))

	// ------------------ Domain-specific routes (standardized) ----------------------

	// Workshop routes (Admin only for write operations)
	router.Handler(http.MethodPost, "/v1/workshops", app.requireRole("admin", http.HandlerFunc(app.createWorkshopHandler)))
	router.Handler(http.MethodGet, "/v1/workshops", app.requireActivatedUser(http.HandlerFunc(app.listWorkshopsHandler)))
	router.Handler(http.MethodGet, "/v1/workshops/:id", app.requireActivatedUser(http.HandlerFunc(app.showWorkshopHandler)))
	router.Handler(http.MethodPatch, "/v1/workshops/:id", app.requireRole("admin", http.HandlerFunc(app.updateWorkshopHandler)))

	// Training Categories routes (Admin only for write operations)
	router.Handler(http.MethodPost, "/v1/training/categories", app.requireRole("admin", http.HandlerFunc(app.createTrainingCategoryHandler)))
	router.Handler(http.MethodGet, "/v1/training/categories", app.requireActivatedUser(http.HandlerFunc(app.listTrainingCategoriesHandler)))
	router.Handler(http.MethodGet, "/v1/training/categories/:id", app.requireActivatedUser(http.HandlerFunc(app.showTrainingCategoryHandler)))
	router.Handler(http.MethodPatch, "/v1/training/categories/:id", app.requireRole("admin", http.HandlerFunc(app.updateTrainingCategoryHandler)))

	// Training Types routes (Admin only for write operations)
	router.Handler(http.MethodPost, "/v1/training/types", app.requireRole("admin", http.HandlerFunc(app.createTrainingTypeHandler)))
	router.Handler(http.MethodGet, "/v1/training/types", app.requireActivatedUser(http.HandlerFunc(app.listTrainingTypesHandler)))
	router.Handler(http.MethodGet, "/v1/training/types/:id", app.requireActivatedUser(http.HandlerFunc(app.showTrainingTypeHandler)))
	router.Handler(http.MethodPatch, "/v1/training/types/:id", app.requireRole("admin", http.HandlerFunc(app.updateTrainingTypeHandler)))

	// Training Status routes
	router.Handler(http.MethodPost, "/v1/training/status", app.requireActivatedUser(http.HandlerFunc(app.createTrainingStatusHandler)))
	router.Handler(http.MethodGet, "/v1/training/status/:id", app.requireActivatedUser(http.HandlerFunc(app.showTrainingStatusHandler)))
	router.Handler(http.MethodPatch, "/v1/training/status/:id", app.requireActivatedUser(http.HandlerFunc(app.updateTrainingStatusHandler)))
	router.Handler(http.MethodGet, "/v1/training/status", app.requireActivatedUser(http.HandlerFunc(app.getTrainingStatusesHandler)))

	// Enrollment Status routes (Admin only for write operations)
	router.Handler(http.MethodPost, "/v1/enrollment/status", app.requireRole("admin", http.HandlerFunc(app.createEnrollmentStatusHandler)))
	router.Handler(http.MethodGet, "/v1/enrollment/status", app.requireActivatedUser(http.HandlerFunc(app.listEnrollmentStatusesHandler)))
	router.Handler(http.MethodGet, "/v1/enrollment/status/:id", app.requireActivatedUser(http.HandlerFunc(app.showEnrollmentStatusHandler)))
	router.Handler(http.MethodPatch, "/v1/enrollment/status/:id", app.requireRole("admin", http.HandlerFunc(app.updateEnrollmentStatusHandler)))

	// Postings routes (Admin only for write operations)
	router.Handler(http.MethodPost, "/v1/postings", app.requireRole("admin", http.HandlerFunc(app.createPostingHandler)))
	router.Handler(http.MethodGet, "/v1/postings", app.requireActivatedUser(http.HandlerFunc(app.listPostingsHandler)))
	router.Handler(http.MethodGet, "/v1/postings/:id", app.requireActivatedUser(http.HandlerFunc(app.showPostingHandler)))
	router.Handler(http.MethodPatch, "/v1/postings/:id", app.requireRole("admin", http.HandlerFunc(app.updatePostingHandler)))

	// Ranks routes (Admin only for write operations)
	router.Handler(http.MethodPost, "/v1/ranks", app.requireRole("admin", http.HandlerFunc(app.createRankHandler)))
	router.Handler(http.MethodGet, "/v1/ranks", app.requireActivatedUser(http.HandlerFunc(app.listRanksHandler)))
	router.Handler(http.MethodGet, "/v1/ranks/:id", app.requireActivatedUser(http.HandlerFunc(app.showRankHandler)))
	router.Handler(http.MethodPatch, "/v1/ranks/:id", app.requireRole("admin", http.HandlerFunc(app.updateRankHandler)))

	// Regions routes (Admin only for write operations)
	router.Handler(http.MethodPost, "/v1/regions", app.requireRole("admin", http.HandlerFunc(app.createRegionHandler)))
	router.Handler(http.MethodGet, "/v1/regions", app.requireActivatedUser(http.HandlerFunc(app.listRegionsHandler)))
	router.Handler(http.MethodGet, "/v1/regions/:id", app.requireActivatedUser(http.HandlerFunc(app.showRegionHandler)))
	router.Handler(http.MethodPatch, "/v1/regions/:id", app.requireRole("admin", http.HandlerFunc(app.updateRegionHandler)))

	// Formations routes (Admin only for write operations)
	router.Handler(http.MethodPost, "/v1/formations", app.requireRole("admin", http.HandlerFunc(app.createFormationHandler)))
	router.Handler(http.MethodGet, "/v1/formations", app.requireActivatedUser(http.HandlerFunc(app.listFormationsHandler)))
	router.Handler(http.MethodGet, "/v1/formations/:id", app.requireActivatedUser(http.HandlerFunc(app.showFormationHandler)))
	router.Handler(http.MethodPatch, "/v1/formations/:id", app.requireRole("admin", http.HandlerFunc(app.updateFormationHandler)))

	// Officer routes
	router.Handler(http.MethodPost, "/v1/officers", app.requireRole("admin", http.HandlerFunc(app.createOfficerHandler)))
	router.Handler(http.MethodGet, "/v1/officers", app.requireActivatedUser(http.HandlerFunc(app.getAllOfficersHandler)))
	router.Handler(http.MethodGet, "/v1/officers/:id", app.requireActivatedUser(http.HandlerFunc(app.showOfficerHandler)))
	router.Handler(http.MethodPatch, "/v1/officers/:id", app.requireRole("admin", http.HandlerFunc(app.updateOfficerHandler)))
	router.Handler(http.MethodDelete, "/v1/officers/:id", app.requireRole("admin", http.HandlerFunc(app.deleteOfficerHandler)))

	// // Training Sessions routes (Admin only for write operations)
	// router.Handler(http.MethodPost, "/v1/training/sessions", app.requireRole("admin", http.HandlerFunc(app.createTrainingSessionHandler)))
	// router.Handler(http.MethodGet, "/v1/training/sessions", app.requireActivatedUser(http.HandlerFunc(app.listTrainingSessionsHandler)))
	// router.Handler(http.MethodGet, "/v1/training/sessions/:id", app.requireActivatedUser(http.HandlerFunc(app.showTrainingSessionHandler)))
	// router.Handler(http.MethodPatch, "/v1/training/sessions/:id", app.requireRole("admin", http.HandlerFunc(app.updateTrainingSessionHandler)))

	// // Training Enrollments routes
	// router.Handler(http.MethodPost, "/v1/training/enrollments", app.requireActivatedUser(http.HandlerFunc(app.createTrainingEnrollmentHandler)))
	// router.Handler(http.MethodGet, "/v1/training/enrollments", app.requireActivatedUser(http.HandlerFunc(app.listTrainingEnrollmentsHandler)))
	// router.Handler(http.MethodGet, "/v1/training/enrollments/:id", app.requireActivatedUser(http.HandlerFunc(app.showTrainingEnrollmentHandler)))
	// router.Handler(http.MethodPatch, "/v1/training/enrollments/:id", app.requireActivatedUser(http.HandlerFunc(app.updateTrainingEnrollmentHandler)))
	return app.recoverPanic(app.enableCORS(app.metrics(app.rateLimit(app.authenticate(router)))))
}
