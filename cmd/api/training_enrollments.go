// Filename: cmd/api/training_enrollments.go
package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Pedro-J-Kukul/police_training/internal/data"
	"github.com/Pedro-J-Kukul/police_training/internal/validator"
)

/*********************** Training Enrollments ***********************/

// createTrainingEnrollmentHandler handles the creation of a new training enrollment.
//
//	@Summary		Create a new training enrollment
//	@Description	Create a new training enrollment
//	@Tags			training-enrollments
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			training_enrollment	body		CreateTrainingEnrollmentRequest_T	true	"Training enrollment data"
//	@Success		201					{object}	envelope
//	@Failure		400					{object}	errorResponse
//	@Failure		422					{object}	errorResponse
//	@Failure		500					{object}	errorResponse
//	@Router			/v1/training/enrollments [post]
func (app *appDependencies) createTrainingEnrollmentHandler(w http.ResponseWriter, r *http.Request) {
	if err := app.readJSON(w, r, &CreateTrainingEnrollmentRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	enrollment := &data.TrainingEnrollment{
		OfficerID:          CreateTrainingEnrollmentRequest.OfficerID,
		SessionID:          CreateTrainingEnrollmentRequest.SessionID,
		EnrollmentStatusID: CreateTrainingEnrollmentRequest.EnrollmentStatusID,
		AttendanceStatusID: CreateTrainingEnrollmentRequest.AttendanceStatusID,
		ProgressStatusID:   CreateTrainingEnrollmentRequest.ProgressStatusID,
		CompletionDate:     CreateTrainingEnrollmentRequest.CompletionDate,
		CertificateIssued:  CreateTrainingEnrollmentRequest.CertificateIssued,
		CertificateNumber:  CreateTrainingEnrollmentRequest.CertificateNumber,
	}

	v := validator.New()
	data.ValidateTrainingEnrollment(v, enrollment)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.TrainingEnrollment.Insert(enrollment); err != nil {
		switch {
		case errors.Is(err, data.ErrForeignKeyViolation):
			v.AddError("foreign_key", "must reference valid officer, session, enrollment status, and progress status")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("enrollment", "officer is already enrolled in this session")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/training/enrollments/%d", enrollment.ID))

	if err := app.writeJSON(w, http.StatusCreated, envelope{"training_enrollment": enrollment}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showTrainingEnrollmentHandler retrieves and returns a training enrollment by its ID.
//
//	@Summary		Get a training enrollment
//	@Description	Retrieve a training enrollment by its ID
//	@Tags			training-enrollments
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Training enrollment ID"
//	@Success		200	{object}	envelope
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/v1/training/enrollments/{id} [get]
func (app *appDependencies) showTrainingEnrollmentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	enrollment, err := app.models.TrainingEnrollment.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"training_enrollment": enrollment}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listTrainingEnrollmentsHandler returns a filtered list of training enrollments.
//
//	@Summary		List training enrollments
//	@Description	Retrieve a list of training enrollments with optional filtering
//	@Tags			training-enrollments
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			officer_id				query		int		false	"Filter by officer ID"
//	@Param			session_id				query		int		false	"Filter by session ID"
//	@Param			enrollment_status_id	query		int		false	"Filter by enrollment status ID"
//	@Param			progress_status_id		query		int		false	"Filter by progress status ID"
//	@Param			certificate_issued		query		bool	false	"Filter by certificate issued status"
//	@Param			page					query		int		false	"Page number for pagination"
//	@Param			page_size				query		int		false	"Number of items per page"
//	@Param			sort					query		string	false	"Sort order"
//	@Success		200						{object}	envelope
//	@Failure		422						{object}	errorResponse
//	@Failure		500						{object}	errorResponse
//	@Router			/v1/training/enrollments [get]
func (app *appDependencies) listTrainingEnrollmentsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	v := validator.New()

	filters := app.readFilters(query, "created_at", 20, []string{
		"created_at", "-created_at", "updated_at", "-updated_at",
		"completion_date", "-completion_date", "id", "-id"}, v)

	// Get optional filter parameters
	officerID := app.getOptionalInt64QueryParameter(query, "officer_id", v)
	sessionID := app.getOptionalInt64QueryParameter(query, "session_id", v)
	enrollmentStatusID := app.getOptionalInt64QueryParameter(query, "enrollment_status_id", v)
	progressStatusID := app.getOptionalInt64QueryParameter(query, "progress_status_id", v)
	certificateIssued := app.getOptionalBoolQueryParameter(query, "certificate_issued", v)

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	enrollments, metadata, err := app.models.TrainingEnrollment.GetAll(
		officerID, sessionID, enrollmentStatusID,
		progressStatusID, certificateIssued, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"training_enrollments": enrollments, "metadata": metadata}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateTrainingEnrollmentHandler performs a partial update on a training enrollment record.
//
//	@Summary		Update a training enrollment
//	@Description	Perform a partial update on a training enrollment record
//	@Tags			training-enrollments
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id					path		int								true	"Training enrollment ID"
//	@Param			training_enrollment	body		UpdateTrainingEnrollmentRequest_T	true	"Training enrollment data"
//	@Success		200					{object}	envelope
//	@Failure		400					{object}	errorResponse
//	@Failure		404					{object}	errorResponse
//	@Failure		422					{object}	errorResponse
//	@Failure		500					{object}	errorResponse
//	@Router			/v1/training/enrollments/{id} [patch]
func (app *appDependencies) updateTrainingEnrollmentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	enrollment, err := app.models.TrainingEnrollment.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.readJSON(w, r, &UpdateTrainingEnrollmentRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Update fields if provided
	if UpdateTrainingEnrollmentRequest.OfficerID != nil {
		enrollment.OfficerID = *UpdateTrainingEnrollmentRequest.OfficerID
	}
	if UpdateTrainingEnrollmentRequest.SessionID != nil {
		enrollment.SessionID = *UpdateTrainingEnrollmentRequest.SessionID
	}
	if UpdateTrainingEnrollmentRequest.EnrollmentStatusID != nil {
		enrollment.EnrollmentStatusID = *UpdateTrainingEnrollmentRequest.EnrollmentStatusID
	}
	if UpdateTrainingEnrollmentRequest.AttendanceStatusID != nil {
		enrollment.AttendanceStatusID = UpdateTrainingEnrollmentRequest.AttendanceStatusID
	}
	if UpdateTrainingEnrollmentRequest.ProgressStatusID != nil {
		enrollment.ProgressStatusID = *UpdateTrainingEnrollmentRequest.ProgressStatusID
	}
	if UpdateTrainingEnrollmentRequest.CompletionDate != nil {
		enrollment.CompletionDate = UpdateTrainingEnrollmentRequest.CompletionDate
	}
	if UpdateTrainingEnrollmentRequest.CertificateIssued != nil {
		enrollment.CertificateIssued = *UpdateTrainingEnrollmentRequest.CertificateIssued
	}
	if UpdateTrainingEnrollmentRequest.CertificateNumber != nil {
		enrollment.CertificateNumber = UpdateTrainingEnrollmentRequest.CertificateNumber
	}

	v := validator.New()
	data.ValidateTrainingEnrollment(v, enrollment)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.TrainingEnrollment.Update(enrollment); err != nil {
		switch {
		case errors.Is(err, data.ErrForeignKeyViolation):
			v.AddError("foreign_key", "must reference valid related records")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("enrollment", "officer is already enrolled in this session")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"training_enrollment": enrollment}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteTrainingEnrollmentHandler handles the deletion of a training enrollment.
//
//	@Summary		Delete a training enrollment
//	@Description	Delete a training enrollment by its ID
//	@Tags			training-enrollments
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Training enrollment ID"
//	@Success		200	{object}	envelope{message=string}
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/v1/training/enrollments/{id} [delete]
func (app *appDependencies) deleteTrainingEnrollmentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.TrainingEnrollment.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"message": "training enrollment successfully deleted"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
