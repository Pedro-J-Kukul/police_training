// Filename: cmd/api/training_sessions.go
package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Pedro-J-Kukul/police_training/internal/data"
	"github.com/Pedro-J-Kukul/police_training/internal/validator"
)

/*********************** Training Sessions ***********************/

// createTrainingSessionHandler handles the creation of a new training session.
//
//	@Summary		Create a new training session
//	@Description	Create a new training session
//	@Tags			training-sessions
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			training_session	body		CreateTrainingSessionRequest_T	true	"Training session data"
//	@Success		201					{object}	envelope
//	@Failure		400					{object}	errorResponse
//	@Failure		422					{object}	errorResponse
//	@Failure		500					{object}	errorResponse
//	@Router			/v1/training/sessions [post]
func (app *appDependencies) createTrainingSessionHandler(w http.ResponseWriter, r *http.Request) {
	if err := app.readJSON(w, r, &CreateTrainingSessionRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	session := &data.TrainingSession{
		FormationID:      CreateTrainingSessionRequest.FormationID,
		RegionID:         CreateTrainingSessionRequest.RegionID,
		FacilitatorID:    CreateTrainingSessionRequest.FacilitatorID,
		WorkshopID:       CreateTrainingSessionRequest.WorkshopID,
		SessionDate:      CreateTrainingSessionRequest.SessionDate,
		StartTime:        CreateTrainingSessionRequest.StartTime,
		EndTime:          CreateTrainingSessionRequest.EndTime,
		Location:         CreateTrainingSessionRequest.Location,
		MaxCapacity:      CreateTrainingSessionRequest.MaxCapacity,
		TrainingStatusID: CreateTrainingSessionRequest.TrainingStatusID,
		Notes:            CreateTrainingSessionRequest.Notes,
	}

	v := validator.New()
	data.ValidateTrainingSession(v, session)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.TrainingSession.Insert(session); err != nil {
		switch {
		case errors.Is(err, data.ErrForeignKeyViolation):
			v.AddError("workshop_id", "must reference valid workshop, formation, region, facilitator, and training status")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/training/sessions/%d", session.ID))

	if err := app.writeJSON(w, http.StatusCreated, envelope{"training_session": session}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showTrainingSessionHandler retrieves and returns a training session by its ID.
//
//	@Summary		Get a training session
//	@Description	Retrieve a training session by its ID
//	@Tags			training-sessions
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Training session ID"
//	@Success		200	{object}	envelope
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/v1/training/sessions/{id} [get]
func (app *appDependencies) showTrainingSessionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	session, err := app.models.TrainingSession.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"training_session": session}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listTrainingSessionsHandler returns a filtered list of training sessions.
//
//	@Summary		List training sessions
//	@Description	Retrieve a list of training sessions with optional filtering
//	@Tags			training-sessions
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			formation_id		query		int		false	"Filter by formation ID"
//	@Param			region_id			query		int		false	"Filter by region ID"
//	@Param			facilitator_id		query		int		false	"Filter by facilitator ID"
//	@Param			workshop_id			query		int		false	"Filter by workshop ID"
//	@Param			training_status_id	query		int		false	"Filter by training status ID"
//	@Param			location			query		string	false	"Filter by location"
//	@Param			notes				query		string	false	"Filter by notes"
//	@Param			session_date		query		string	false	"Filter by session date (YYYY-MM-DD)"
//	@Param			page				query		int		false	"Page number for pagination"
//	@Param			page_size			query		int		false	"Number of items per page"
//	@Param			sort				query		string	false	"Sort order"
//	@Success		200					{object}	envelope
//	@Failure		422					{object}	errorResponse
//	@Failure		500					{object}	errorResponse
//	@Router			/v1/training/sessions [get]
func (app *appDependencies) listTrainingSessionsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	v := validator.New()

	filters := app.readFilters(query, "session_date", 20, []string{
		"session_date", "-session_date", "start_time", "-start_time",
		"end_time", "-end_time", "id", "-id", "created_at", "-created_at"}, v)

	// Get optional filter parameters
	formationID := app.getOptionalInt64QueryParameter(query, "formation_id", v)
	regionID := app.getOptionalInt64QueryParameter(query, "region_id", v)
	facilitatorID := app.getOptionalInt64QueryParameter(query, "facilitator_id", v)
	workshopID := app.getOptionalInt64QueryParameter(query, "workshop_id", v)
	trainingStatusID := app.getOptionalInt64QueryParameter(query, "training_status_id", v)

	location := app.getSingleQueryParameter(query, "location", "")
	notes := app.getSingleQueryParameter(query, "notes", "")

	// Parse session date if provided
	sessionDateStr := app.getSingleQueryParameter(query, "session_date", "")
	var sessionDate time.Time
	if sessionDateStr != "" {
		var err error
		sessionDate, err = time.Parse("2006-01-02", sessionDateStr)
		if err != nil {
			v.AddError("session_date", "must be a valid date in YYYY-MM-DD format")
		}
	}

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Convert pointers to values with defaults for GetAll
	formationIDValue := int64(0)
	if formationID != nil {
		formationIDValue = *formationID
	}
	regionIDValue := int64(0)
	if regionID != nil {
		regionIDValue = *regionID
	}
	facilitatorIDValue := int64(0)
	if facilitatorID != nil {
		facilitatorIDValue = *facilitatorID
	}
	workshopIDValue := int64(0)
	if workshopID != nil {
		workshopIDValue = *workshopID
	}
	trainingStatusIDValue := int64(0)
	if trainingStatusID != nil {
		trainingStatusIDValue = *trainingStatusID
	}

	sessions, metadata, err := app.models.TrainingSession.GetAll(
		formationIDValue, regionIDValue, facilitatorIDValue,
		workshopIDValue, trainingStatusIDValue, location, notes,
		time.Time{}, time.Time{}, sessionDate, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"training_sessions": sessions, "metadata": metadata}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateTrainingSessionHandler performs a partial update on a training session record.
//
//	@Summary		Update a training session
//	@Description	Perform a partial update on a training session record
//	@Tags			training-sessions
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id				path		int							true	"Training session ID"
//	@Param			training_session	body		UpdateTrainingSessionRequest_T	true	"Training session data"
//	@Success		200				{object}	envelope
//	@Failure		400				{object}	errorResponse
//	@Failure		404				{object}	errorResponse
//	@Failure		422				{object}	errorResponse
//	@Failure		500				{object}	errorResponse
//	@Router			/v1/training/sessions/{id} [patch]
func (app *appDependencies) updateTrainingSessionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	session, err := app.models.TrainingSession.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.readJSON(w, r, &UpdateTrainingSessionRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Update fields if provided
	if UpdateTrainingSessionRequest.FormationID != nil {
		session.FormationID = *UpdateTrainingSessionRequest.FormationID
	}
	if UpdateTrainingSessionRequest.RegionID != nil {
		session.RegionID = *UpdateTrainingSessionRequest.RegionID
	}
	if UpdateTrainingSessionRequest.FacilitatorID != nil {
		session.FacilitatorID = *UpdateTrainingSessionRequest.FacilitatorID
	}
	if UpdateTrainingSessionRequest.WorkshopID != nil {
		session.WorkshopID = *UpdateTrainingSessionRequest.WorkshopID
	}
	if UpdateTrainingSessionRequest.SessionDate != nil {
		session.SessionDate = *UpdateTrainingSessionRequest.SessionDate
	}
	if UpdateTrainingSessionRequest.StartTime != nil {
		session.StartTime = *UpdateTrainingSessionRequest.StartTime
	}
	if UpdateTrainingSessionRequest.EndTime != nil {
		session.EndTime = *UpdateTrainingSessionRequest.EndTime
	}
	if UpdateTrainingSessionRequest.Location != nil {
		session.Location = UpdateTrainingSessionRequest.Location
	}
	if UpdateTrainingSessionRequest.MaxCapacity != nil {
		session.MaxCapacity = UpdateTrainingSessionRequest.MaxCapacity
	}
	if UpdateTrainingSessionRequest.TrainingStatusID != nil {
		session.TrainingStatusID = *UpdateTrainingSessionRequest.TrainingStatusID
	}
	if UpdateTrainingSessionRequest.Notes != nil {
		session.Notes = UpdateTrainingSessionRequest.Notes
	}

	v := validator.New()
	data.ValidateTrainingSession(v, session)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.TrainingSession.Update(session); err != nil {
		switch {
		case errors.Is(err, data.ErrForeignKeyViolation):
			v.AddError("foreign_key", "must reference valid related records")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"training_session": session}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteTrainingSessionHandler handles the deletion of a training session.
//
//	@Summary		Delete a training session
//	@Description	Delete a training session by its ID
//	@Tags			training-sessions
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Training session ID"
//	@Success		200	{object}	envelope{message=string}
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/v1/training/sessions/{id} [delete]
func (app *appDependencies) deleteTrainingSessionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.TrainingSession.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"message": "training session successfully deleted"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
