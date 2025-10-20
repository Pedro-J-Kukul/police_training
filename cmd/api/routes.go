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
	router.Handler(http.MethodGet, "/v1/users", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ"}, http.HandlerFunc(app.listUsersHandler))))
	router.Handler(http.MethodGet, "/v1/users/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ"}, http.HandlerFunc(app.showUserHandler))))
	router.Handler(http.MethodPatch, "/v1/users/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_MODIFY", "CAN_MODIFY_SELF", "CAN_MODIFY_USER"}, http.HandlerFunc(app.updateUserHandler))))
	router.Handler(http.MethodDelete, "/v1/users/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_DELETE", "CAN_DELETE_SELF", "CAN_DELETE_USER"}, http.HandlerFunc(app.deleteUserHandler))))

	// ------------------ Domain-specific routes (standardized) ----------------------

	// Workshop routes (Admin only for write operations)
	router.Handler(http.MethodPost, "/v1/workshops", app.requireActivatedUser(app.requirePermissions([]string{"CAN_CREATE", "CAN_CREATE_WORKSHOPS"}, http.HandlerFunc(app.createWorkshopHandler))))
	router.Handler(http.MethodGet, "/v1/workshops", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_WORKSHOPS"}, http.HandlerFunc(app.listWorkshopsHandler))))
	router.Handler(http.MethodGet, "/v1/workshops/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_WORKSHOPS"}, http.HandlerFunc(app.showWorkshopHandler))))
	router.Handler(http.MethodPatch, "/v1/workshops/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_MODIFY", "CAN_MODIFY_WORKSHOPS"}, http.HandlerFunc(app.updateWorkshopHandler))))

	// Training Categories routes (Admin only for write operations)
	router.Handler(http.MethodPost, "/v1/training/categories", app.requireActivatedUser(app.requirePermissions([]string{"CAN_CREATE", "CAN_CREATE_TRAINING_CATEGORIES"}, http.HandlerFunc(app.createTrainingCategoryHandler))))
	router.Handler(http.MethodGet, "/v1/training/categories", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_TRAINING_CATEGORIES"}, http.HandlerFunc(app.listTrainingCategoriesHandler))))
	router.Handler(http.MethodGet, "/v1/training/categories/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_TRAINING_CATEGORIES"}, http.HandlerFunc(app.showTrainingCategoryHandler))))
	router.Handler(http.MethodPatch, "/v1/training/categories/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_MODIFY", "CAN_MODIFY_TRAINING_CATEGORIES"}, http.HandlerFunc(app.updateTrainingCategoryHandler))))

	// Training Types routes (Admin only for write operations)
	router.Handler(http.MethodPost, "/v1/training/types", app.requireActivatedUser(app.requirePermissions([]string{"CAN_CREATE", "CAN_CREATE_TRAININGS"}, http.HandlerFunc(app.createTrainingTypeHandler))))
	router.Handler(http.MethodGet, "/v1/training/types", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_TRAININGS"}, http.HandlerFunc(app.listTrainingTypesHandler))))
	router.Handler(http.MethodGet, "/v1/training/types/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_TRAININGS"}, http.HandlerFunc(app.showTrainingTypeHandler))))
	router.Handler(http.MethodPatch, "/v1/training/types/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_MODIFY", "CAN_MODIFY_TRAININGS"}, http.HandlerFunc(app.updateTrainingTypeHandler))))

	// Training Status routes
	router.Handler(http.MethodPost, "/v1/training/status", app.requireActivatedUser(app.requirePermissions([]string{"CAN_MODIFY", "CAN_MODIFY_SELF"}, http.HandlerFunc(app.createTrainingStatusHandler))))
	router.Handler(http.MethodGet, "/v1/training/status/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_TRAININGS"}, http.HandlerFunc(app.showTrainingStatusHandler))))
	router.Handler(http.MethodPatch, "/v1/training/status/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_MODIFY", "CAN_MODIFY_TRAININGS"}, http.HandlerFunc(app.updateTrainingStatusHandler))))
	router.Handler(http.MethodGet, "/v1/training/status", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_TRAININGS"}, http.HandlerFunc(app.getTrainingStatusesHandler))))

	// Enrollment Status routes (Admin only for write operations)
	router.Handler(http.MethodPost, "/v1/enrollment/status", app.requireActivatedUser(app.requirePermissions([]string{"CAN_CREATE", "CAN_CREATE_ENROLLMENT_STATUSES"}, http.HandlerFunc(app.createEnrollmentStatusHandler))))
	router.Handler(http.MethodGet, "/v1/enrollment/status", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_ENROLLMENT_STATUSES"}, http.HandlerFunc(app.listEnrollmentStatusesHandler))))
	router.Handler(http.MethodGet, "/v1/enrollment/status/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_ENROLLMENT_STATUSES"}, http.HandlerFunc(app.showEnrollmentStatusHandler))))
	router.Handler(http.MethodPatch, "/v1/enrollment/status/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_MODIFY", "CAN_MODIFY_ENROLLMENT_STATUSES"}, http.HandlerFunc(app.updateEnrollmentStatusHandler))))

	// Add these new routes:
	router.Handler(http.MethodPost, "/v1/attendance/status", app.requireActivatedUser(app.requirePermissions([]string{"CAN_CREATE", "CAN_CREATE_ATTENDANCE_STATUSES"}, http.HandlerFunc(app.createAttendanceStatusHandler))))
	router.Handler(http.MethodGet, "/v1/attendance/status", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_ATTENDANCE_STATUSES"}, http.HandlerFunc(app.listAttendanceStatusesHandler))))
	router.Handler(http.MethodGet, "/v1/attendance/status/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_ATTENDANCE_STATUSES"}, http.HandlerFunc(app.showAttendanceStatusHandler))))
	router.Handler(http.MethodPatch, "/v1/attendance/status/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_MODIFY", "CAN_MODIFY_ATTENDANCE_STATUSES"}, http.HandlerFunc(app.updateAttendanceStatusHandler))))

	// Add these new routes:
	router.Handler(http.MethodPost, "/v1/progress/status", app.requireActivatedUser(app.requirePermissions([]string{"CAN_CREATE", "CAN_CREATE_PROGRESS_STATUSES"}, http.HandlerFunc(app.createProgressStatusHandler))))
	router.Handler(http.MethodGet, "/v1/progress/status", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_PROGRESS_STATUSES"}, http.HandlerFunc(app.listProgressStatusesHandler))))
	router.Handler(http.MethodGet, "/v1/progress/status/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_PROGRESS_STATUSES"}, http.HandlerFunc(app.showProgressStatusHandler))))
	router.Handler(http.MethodPatch, "/v1/progress/status/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_MODIFY", "CAN_MODIFY_PROGRESS_STATUSES"}, http.HandlerFunc(app.updateProgressStatusHandler))))
	// Postings routes (Admin only for write operations)
	router.Handler(http.MethodPost, "/v1/postings", app.requireActivatedUser(app.requirePermissions([]string{"CAN_CREATE", "CAN_CREATE_POSTINGS"}, http.HandlerFunc(app.createPostingHandler))))
	router.Handler(http.MethodGet, "/v1/postings", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_POSTINGS"}, http.HandlerFunc(app.listPostingsHandler))))
	router.Handler(http.MethodGet, "/v1/postings/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_POSTINGS"}, http.HandlerFunc(app.showPostingHandler))))
	router.Handler(http.MethodPatch, "/v1/postings/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_MODIFY", "CAN_MODIFY_POSTINGS"}, http.HandlerFunc(app.updatePostingHandler))))

	// Ranks routes (Admin only for write operations)
	router.Handler(http.MethodPost, "/v1/ranks", app.requireActivatedUser(app.requirePermissions([]string{"CAN_CREATE", "CAN_CREATE_RANKS"}, http.HandlerFunc(app.createRankHandler))))
	router.Handler(http.MethodGet, "/v1/ranks", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_RANKS"}, http.HandlerFunc(app.listRanksHandler))))
	router.Handler(http.MethodGet, "/v1/ranks/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_RANKS"}, http.HandlerFunc(app.showRankHandler))))
	router.Handler(http.MethodPatch, "/v1/ranks/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_MODIFY", "CAN_MODIFY_RANKS"}, http.HandlerFunc(app.updateRankHandler))))

	// Regions routes (Admin only for write operations)
	router.Handler(http.MethodPost, "/v1/regions", app.requireActivatedUser(app.requirePermissions([]string{"CAN_CREATE", "CAN_CREATE_REGIONS"}, http.HandlerFunc(app.createRegionHandler))))
	router.Handler(http.MethodGet, "/v1/regions", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_REGIONS"}, http.HandlerFunc(app.listRegionsHandler))))
	router.Handler(http.MethodGet, "/v1/regions/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_REGIONS"}, http.HandlerFunc(app.showRegionHandler))))
	router.Handler(http.MethodPatch, "/v1/regions/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_MODIFY", "CAN_MODIFY_REGIONS"}, http.HandlerFunc(app.updateRegionHandler))))

	// Formations routes (Admin only for write operations)
	router.Handler(http.MethodPost, "/v1/formations", app.requireActivatedUser(app.requirePermissions([]string{"CAN_CREATE", "CAN_CREATE_FORMATIONS"}, http.HandlerFunc(app.createFormationHandler))))
	router.Handler(http.MethodGet, "/v1/formations", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_FORMATIONS"}, http.HandlerFunc(app.listFormationsHandler))))
	router.Handler(http.MethodGet, "/v1/formations/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_FORMATIONS"}, http.HandlerFunc(app.showFormationHandler))))
	router.Handler(http.MethodPatch, "/v1/formations/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_MODIFY", "CAN_MODIFY_FORMATIONS"}, http.HandlerFunc(app.updateFormationHandler))))

	// Officer routes
	router.Handler(http.MethodPost, "/v1/officers", app.requireActivatedUser(app.requirePermissions([]string{"CAN_CREATE", "CAN_CREATE_OFFICERS"}, http.HandlerFunc(app.createOfficerHandler))))
	router.Handler(http.MethodGet, "/v1/officers", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_OFFICERS"}, http.HandlerFunc(app.getAllOfficersHandler))))
	router.Handler(http.MethodGet, "/v1/officers/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_OFFICERS"}, http.HandlerFunc(app.showOfficerHandler))))
	router.Handler(http.MethodPatch, "/v1/officers/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_MODIFY", "CAN_MODIFY_OFFICERS"}, http.HandlerFunc(app.updateOfficerHandler))))
	router.Handler(http.MethodDelete, "/v1/officers/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_DELETE", "CAN_DELETE_OFFICERS"}, http.HandlerFunc(app.deleteOfficerHandler))))

	// Training sessions routes
	router.Handler(http.MethodPost, "/v1/training/sessions", app.requireActivatedUser(app.requirePermissions([]string{"CAN_CREATE", "CAN_CREATE_TRAINING_SESSIONS"}, http.HandlerFunc(app.createTrainingSessionHandler))))
	router.Handler(http.MethodGet, "/v1/training/sessions", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_TRAINING_SESSIONS"}, http.HandlerFunc(app.listTrainingSessionsHandler))))
	router.Handler(http.MethodGet, "/v1/training/sessions/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_TRAINING_SESSIONS"}, http.HandlerFunc(app.showTrainingSessionHandler))))
	router.Handler(http.MethodPatch, "/v1/training/sessions/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_MODIFY", "CAN_MODIFY_TRAINING_SESSIONS"}, http.HandlerFunc(app.updateTrainingSessionHandler))))
	router.Handler(http.MethodDelete, "/v1/training/sessions/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_DELETE", "CAN_DELETE_TRAINING_SESSIONS"}, http.HandlerFunc(app.deleteTrainingSessionHandler))))

	// Training enrollments routes
	router.Handler(http.MethodPost, "/v1/training/enrollments", app.requireActivatedUser(app.requirePermissions([]string{"CAN_CREATE", "CAN_CREATE_TRAINING_ENROLLMENTS"}, http.HandlerFunc(app.createTrainingEnrollmentHandler))))
	router.Handler(http.MethodGet, "/v1/training/enrollments", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_TRAINING_ENROLLMENTS"}, http.HandlerFunc(app.listTrainingEnrollmentsHandler))))
	router.Handler(http.MethodGet, "/v1/training/enrollments/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_READ", "CAN_READ_TRAINING_ENROLLMENTS"}, http.HandlerFunc(app.showTrainingEnrollmentHandler))))
	router.Handler(http.MethodPatch, "/v1/training/enrollments/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_MODIFY", "CAN_MODIFY_TRAINING_ENROLLMENTS"}, http.HandlerFunc(app.updateTrainingEnrollmentHandler))))
	router.Handler(http.MethodDelete, "/v1/training/enrollments/:id", app.requireActivatedUser(app.requirePermissions([]string{"CAN_DELETE", "CAN_DELETE_TRAINING_ENROLLMENTS"}, http.HandlerFunc(app.deleteTrainingEnrollmentHandler))))

	return app.recoverPanic(app.enableCORS(app.metrics(app.rateLimit(app.authenticate(router)))))
}
