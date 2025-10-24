package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Pedro-J-Kukul/police_training/internal/data"
	"github.com/Pedro-J-Kukul/police_training/internal/validator"
)

func (app *appDependencies) createTrainingEnrollmentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OfficerID          int64   `json:"officer_id"`
		SessionID          int64   `json:"session_id"`
		EnrollmentStatusID int64   `json:"enrollment_status_id"`
		AttendanceStatusID *int64  `json:"attendance_status_id"`
		ProgressStatusID   int64   `json:"progress_status_id"`
		CompletionDate     *string `json:"completion_date"` // "2025-01-15"
		CertificateIssued  *bool   `json:"certificate_issued"`
		CertificateNumber  *string `json:"certificate_number"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	enrollment := &data.TrainingEnrollment{
		OfficerID:          input.OfficerID,
		SessionID:          input.SessionID,
		EnrollmentStatusID: input.EnrollmentStatusID,
		AttendanceStatusID: input.AttendanceStatusID,
		ProgressStatusID:   input.ProgressStatusID,
		CertificateIssued:  false, // default
		CertificateNumber:  input.CertificateNumber,
	}

	if input.CertificateIssued != nil {
		enrollment.CertificateIssued = *input.CertificateIssued
	}

	// Parse completion date if provided
	if input.CompletionDate != nil {
		completionDate, err := time.Parse("2006-01-02", *input.CompletionDate)
		if err != nil {
			app.badRequestResponse(w, r, errors.New("invalid completion_date format, use YYYY-MM-DD"))
			return
		}
		enrollment.CompletionDate = &completionDate
	}

	v := validator.New()
	data.ValidateTrainingEnrollment(v, enrollment)

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.TrainingEnrollment.Insert(enrollment)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			app.badRequestResponse(w, r, errors.New("officer is already enrolled in this session"))
		case errors.Is(err, data.ErrForeignKeyViolation):
			app.badRequestResponse(w, r, errors.New("invalid officer_id, session_id, enrollment_status_id, or progress_status_id"))
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/training-enrollments/%d", enrollment.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"training_enrollment": enrollment}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

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

	err = app.writeJSON(w, http.StatusOK, envelope{"training_enrollment": enrollment}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) listTrainingEnrollmentsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	v := validator.New()

	filters := app.readFilters(query, "created_at", 20, []string{"created_at", "-created_at", "id", "-id", "completion_date", "-completion_date"}, v)

	officerID := app.getOptionalInt64QueryParameter(query, "officer_id", v)
	sessionID := app.getOptionalInt64QueryParameter(query, "session_id", v)
	enrollmentStatusID := app.getOptionalInt64QueryParameter(query, "enrollment_status_id", v)
	attendanceStatusID := app.getOptionalInt64QueryParameter(query, "attendance_status_id", v)
	progressStatusID := app.getOptionalInt64QueryParameter(query, "progress_status_id", v)

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	enrollments, metadata, err := app.models.TrainingEnrollment.GetAll(officerID, sessionID, enrollmentStatusID, attendanceStatusID, progressStatusID, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"training_enrollments": enrollments, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

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

	var input struct {
		OfficerID          *int64  `json:"officer_id"`
		SessionID          *int64  `json:"session_id"`
		EnrollmentStatusID *int64  `json:"enrollment_status_id"`
		AttendanceStatusID *int64  `json:"attendance_status_id"`
		ProgressStatusID   *int64  `json:"progress_status_id"`
		CompletionDate     *string `json:"completion_date"`
		CertificateIssued  *bool   `json:"certificate_issued"`
		CertificateNumber  *string `json:"certificate_number"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.OfficerID != nil {
		enrollment.OfficerID = *input.OfficerID
	}
	if input.SessionID != nil {
		enrollment.SessionID = *input.SessionID
	}
	if input.EnrollmentStatusID != nil {
		enrollment.EnrollmentStatusID = *input.EnrollmentStatusID
	}
	if input.AttendanceStatusID != nil {
		enrollment.AttendanceStatusID = input.AttendanceStatusID
	}
	if input.ProgressStatusID != nil {
		enrollment.ProgressStatusID = *input.ProgressStatusID
	}
	if input.CompletionDate != nil {
		if *input.CompletionDate == "" {
			enrollment.CompletionDate = nil
		} else {
			completionDate, err := time.Parse("2006-01-02", *input.CompletionDate)
			if err != nil {
				app.badRequestResponse(w, r, errors.New("invalid completion_date format, use YYYY-MM-DD"))
				return
			}
			enrollment.CompletionDate = &completionDate
		}
	}
	if input.CertificateIssued != nil {
		enrollment.CertificateIssued = *input.CertificateIssued
	}
	if input.CertificateNumber != nil {
		enrollment.CertificateNumber = input.CertificateNumber
	}

	v := validator.New()
	data.ValidateTrainingEnrollment(v, enrollment)

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.TrainingEnrollment.Update(enrollment)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		case errors.Is(err, data.ErrDuplicateValue):
			app.badRequestResponse(w, r, errors.New("officer is already enrolled in this session"))
		case errors.Is(err, data.ErrForeignKeyViolation):
			app.badRequestResponse(w, r, errors.New("invalid officer_id, session_id, enrollment_status_id, or progress_status_id"))
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"training_enrollment": enrollment}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

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

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "training enrollment successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Specialized handlers
func (app *appDependencies) getOfficerEnrollmentsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	query := r.URL.Query()
	v := validator.New()
	filters := app.readFilters(query, "created_at", 20, []string{"created_at", "-created_at", "completion_date", "-completion_date"}, v)

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	enrollments, metadata, err := app.models.TrainingEnrollment.GetByOfficer(id, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"training_enrollments": enrollments, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) getSessionEnrollmentsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	query := r.URL.Query()
	v := validator.New()
	filters := app.readFilters(query, "created_at", 20, []string{"created_at", "-created_at", "completion_date", "-completion_date"}, v)

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	enrollments, metadata, err := app.models.TrainingEnrollment.GetBySession(id, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"training_enrollments": enrollments, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) issueCertificateHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		CertificateNumber string `json:"certificate_number"`
		CompletionDate    string `json:"completion_date"` // "2025-01-15"
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate inputs
	v := validator.New()
	v.Check(input.CertificateNumber != "", "certificate_number", "must be provided")
	v.Check(len(input.CertificateNumber) <= 100, "certificate_number", "must not exceed 100 characters")
	v.Check(input.CompletionDate != "", "completion_date", "must be provided")

	completionDate, err := time.Parse("2006-01-02", input.CompletionDate)
	if err != nil {
		v.AddError("completion_date", "invalid date format, use YYYY-MM-DD")
	}

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.TrainingEnrollment.IssueCertificate(id, input.CertificateNumber, completionDate)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "certificate issued successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
