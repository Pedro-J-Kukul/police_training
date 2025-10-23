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

	// Authentication and user lifecycle (no permissions required)
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activate", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/password-reset", app.createPasswordResetTokenHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/password-reset", app.resetPasswordHandler)

	// Authenticated user endpoints
	router.Handler(http.MethodGet, "/v1/me", app.requireActivatedUser(http.HandlerFunc(app.showCurrentUserHandler)))
	router.Handler(http.MethodGet, "/v1/users", app.requirePermissions("users:view")(http.HandlerFunc(app.listUsersHandler)))
	router.Handler(http.MethodGet, "/v1/users/:id", app.requirePermissions("users:view")(http.HandlerFunc(app.showUserHandler)))
	router.Handler(http.MethodPatch, "/v1/users/:id", app.requirePermissions("users:edit")(http.HandlerFunc(app.updateUserHandler)))
	router.Handler(http.MethodDelete, "/v1/users/:id", app.requirePermissions("users:delete")(http.HandlerFunc(app.deleteUserHandler)))

	// ------------------ Domain-specific routes (standardized) ----------------------

	// Workshop routes
	router.Handler(http.MethodPost, "/v1/workshops", app.requirePermissions("workshops:create")(http.HandlerFunc(app.createWorkshopHandler)))
	router.Handler(http.MethodGet, "/v1/workshops", app.requirePermissions("workshops:view")(http.HandlerFunc(app.listWorkshopsHandler)))
	router.Handler(http.MethodGet, "/v1/workshops/:id", app.requirePermissions("workshops:view")(http.HandlerFunc(app.showWorkshopHandler)))
	router.Handler(http.MethodPatch, "/v1/workshops/:id", app.requirePermissions("workshops:edit")(http.HandlerFunc(app.updateWorkshopHandler)))

	// Training Categories routes
	router.Handler(http.MethodPost, "/v1/training/categories", app.requirePermissions("training:categories:create")(http.HandlerFunc(app.createTrainingCategoryHandler)))
	router.Handler(http.MethodGet, "/v1/training/categories", app.requirePermissions("training:categories:view")(http.HandlerFunc(app.listTrainingCategoriesHandler)))
	router.Handler(http.MethodGet, "/v1/training/categories/:id", app.requirePermissions("training:categories:view")(http.HandlerFunc(app.showTrainingCategoryHandler)))
	router.Handler(http.MethodPatch, "/v1/training/categories/:id", app.requirePermissions("training:categories:edit")(http.HandlerFunc(app.updateTrainingCategoryHandler)))

	// Training Types routes
	router.Handler(http.MethodPost, "/v1/training/types", app.requirePermissions("training:types:create")(http.HandlerFunc(app.createTrainingTypeHandler)))
	router.Handler(http.MethodGet, "/v1/training/types", app.requirePermissions("training:types:view")(http.HandlerFunc(app.listTrainingTypesHandler)))
	router.Handler(http.MethodGet, "/v1/training/types/:id", app.requirePermissions("training:types:view")(http.HandlerFunc(app.showTrainingTypeHandler)))
	router.Handler(http.MethodPatch, "/v1/training/types/:id", app.requirePermissions("training:types:edit")(http.HandlerFunc(app.updateTrainingTypeHandler)))

	// Training Status routes
	router.Handler(http.MethodPost, "/v1/training/status", app.requirePermissions("training:status:create")(http.HandlerFunc(app.createTrainingStatusHandler)))
	router.Handler(http.MethodGet, "/v1/training/status/:id", app.requirePermissions("training:status:view")(http.HandlerFunc(app.showTrainingStatusHandler)))
	router.Handler(http.MethodPatch, "/v1/training/status/:id", app.requirePermissions("training:status:edit")(http.HandlerFunc(app.updateTrainingStatusHandler)))
	router.Handler(http.MethodGet, "/v1/training/status", app.requirePermissions("training:status:view")(http.HandlerFunc(app.getTrainingStatusesHandler)))

	// Enrollment Status routes
	router.Handler(http.MethodPost, "/v1/enrollment/status", app.requirePermissions("enrollment:status:create")(http.HandlerFunc(app.createEnrollmentStatusHandler)))
	router.Handler(http.MethodGet, "/v1/enrollment/status", app.requirePermissions("enrollment:status:view")(http.HandlerFunc(app.listEnrollmentStatusesHandler)))
	router.Handler(http.MethodGet, "/v1/enrollment/status/:id", app.requirePermissions("enrollment:status:view")(http.HandlerFunc(app.showEnrollmentStatusHandler)))
	router.Handler(http.MethodPatch, "/v1/enrollment/status/:id", app.requirePermissions("enrollment:status:edit")(http.HandlerFunc(app.updateEnrollmentStatusHandler)))

	// Attendance Status routes
	router.Handler(http.MethodPost, "/v1/attendance/status", app.requirePermissions("attendance:status:create")(http.HandlerFunc(app.createAttendanceStatusHandler)))
	router.Handler(http.MethodGet, "/v1/attendance/status", app.requirePermissions("attendance:status:view")(http.HandlerFunc(app.listAttendanceStatusesHandler)))
	router.Handler(http.MethodGet, "/v1/attendance/status/:id", app.requirePermissions("attendance:status:view")(http.HandlerFunc(app.showAttendanceStatusHandler)))
	router.Handler(http.MethodPatch, "/v1/attendance/status/:id", app.requirePermissions("attendance:status:edit")(http.HandlerFunc(app.updateAttendanceStatusHandler)))

	// Progress Status routes
	router.Handler(http.MethodPost, "/v1/progress/status", app.requirePermissions("progress:status:create")(http.HandlerFunc(app.createProgressStatusHandler)))
	router.Handler(http.MethodGet, "/v1/progress/status", app.requirePermissions("progress:status:view")(http.HandlerFunc(app.listProgressStatusesHandler)))
	router.Handler(http.MethodGet, "/v1/progress/status/:id", app.requirePermissions("progress:status:view")(http.HandlerFunc(app.showProgressStatusHandler)))
	router.Handler(http.MethodPatch, "/v1/progress/status/:id", app.requirePermissions("progress:status:edit")(http.HandlerFunc(app.updateProgressStatusHandler)))

	// Postings routes
	router.Handler(http.MethodPost, "/v1/postings", app.requirePermissions("postings:create")(http.HandlerFunc(app.createPostingHandler)))
	router.Handler(http.MethodGet, "/v1/postings", app.requirePermissions("postings:view")(http.HandlerFunc(app.listPostingsHandler)))
	router.Handler(http.MethodGet, "/v1/postings/:id", app.requirePermissions("postings:view")(http.HandlerFunc(app.showPostingHandler)))
	router.Handler(http.MethodPatch, "/v1/postings/:id", app.requirePermissions("postings:edit")(http.HandlerFunc(app.updatePostingHandler)))

	// Ranks routes
	router.Handler(http.MethodPost, "/v1/ranks", app.requirePermissions("ranks:create")(http.HandlerFunc(app.createRankHandler)))
	router.Handler(http.MethodGet, "/v1/ranks", app.requirePermissions("ranks:view")(http.HandlerFunc(app.listRanksHandler)))
	router.Handler(http.MethodGet, "/v1/ranks/:id", app.requirePermissions("ranks:view")(http.HandlerFunc(app.showRankHandler)))
	router.Handler(http.MethodPatch, "/v1/ranks/:id", app.requirePermissions("ranks:edit")(http.HandlerFunc(app.updateRankHandler)))

	// Regions routes
	router.Handler(http.MethodPost, "/v1/regions", app.requirePermissions("regions:create")(http.HandlerFunc(app.createRegionHandler)))
	router.Handler(http.MethodGet, "/v1/regions", app.requirePermissions("regions:view")(http.HandlerFunc(app.listRegionsHandler)))
	router.Handler(http.MethodGet, "/v1/regions/:id", app.requirePermissions("regions:view")(http.HandlerFunc(app.showRegionHandler)))
	router.Handler(http.MethodPatch, "/v1/regions/:id", app.requirePermissions("regions:edit")(http.HandlerFunc(app.updateRegionHandler)))

	// Formations routes
	router.Handler(http.MethodPost, "/v1/formations", app.requirePermissions("formations:create")(http.HandlerFunc(app.createFormationHandler)))
	router.Handler(http.MethodGet, "/v1/formations", app.requirePermissions("formations:view")(http.HandlerFunc(app.listFormationsHandler)))
	router.Handler(http.MethodGet, "/v1/formations/:id", app.requirePermissions("formations:view")(http.HandlerFunc(app.showFormationHandler)))
	router.Handler(http.MethodPatch, "/v1/formations/:id", app.requirePermissions("formations:edit")(http.HandlerFunc(app.updateFormationHandler)))

	// Officer routes
	router.Handler(http.MethodPost, "/v1/officers", app.requirePermissions("officers:create")(http.HandlerFunc(app.createOfficerHandler)))
	router.Handler(http.MethodGet, "/v1/officers", app.requirePermissions("officers:view")(http.HandlerFunc(app.getAllOfficersHandler)))
	router.Handler(http.MethodGet, "/v1/officers/:id", app.requirePermissions("officers:view")(http.HandlerFunc(app.showOfficerHandler)))
	router.Handler(http.MethodPatch, "/v1/officers/:id", app.requirePermissions("officers:edit")(http.HandlerFunc(app.updateOfficerHandler)))
	router.Handler(http.MethodDelete, "/v1/officers/:id", app.requirePermissions("officers:delete")(http.HandlerFunc(app.deleteOfficerHandler)))

	// Training sessions routes
	router.Handler(http.MethodPost, "/v1/training/sessions", app.requirePermissions("training:sessions:create")(http.HandlerFunc(app.createTrainingSessionHandler)))
	router.Handler(http.MethodGet, "/v1/training/sessions", app.requirePermissions("training:sessions:view")(http.HandlerFunc(app.listTrainingSessionsHandler)))
	router.Handler(http.MethodGet, "/v1/training/sessions/:id", app.requirePermissions("training:sessions:view")(http.HandlerFunc(app.showTrainingSessionHandler)))
	router.Handler(http.MethodPatch, "/v1/training/sessions/:id", app.requirePermissions("training:sessions:edit")(http.HandlerFunc(app.updateTrainingSessionHandler)))
	router.Handler(http.MethodDelete, "/v1/training/sessions/:id", app.requirePermissions("training:sessions:delete")(http.HandlerFunc(app.deleteTrainingSessionHandler)))

	// Training enrollments routes
	router.Handler(http.MethodPost, "/v1/training/enrollments", app.requirePermissions("training:enrollments:create")(http.HandlerFunc(app.createTrainingEnrollmentHandler)))
	router.Handler(http.MethodGet, "/v1/training/enrollments", app.requirePermissions("training:enrollments:view")(http.HandlerFunc(app.listTrainingEnrollmentsHandler)))
	router.Handler(http.MethodGet, "/v1/training/enrollments/:id", app.requirePermissions("training:enrollments:view")(http.HandlerFunc(app.showTrainingEnrollmentHandler)))
	router.Handler(http.MethodPatch, "/v1/training/enrollments/:id", app.requirePermissions("training:enrollments:edit")(http.HandlerFunc(app.updateTrainingEnrollmentHandler)))
	router.Handler(http.MethodDelete, "/v1/training/enrollments/:id", app.requirePermissions("training:enrollments:delete")(http.HandlerFunc(app.deleteTrainingEnrollmentHandler)))

	return app.recoverPanic(app.enableCORS(app.metrics(app.rateLimit(app.authenticate(router)))))
}
