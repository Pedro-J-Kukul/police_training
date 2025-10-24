package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Pedro-J-Kukul/police_training/internal/data"
	"github.com/Pedro-J-Kukul/police_training/internal/validator"
)

func (app *appDependencies) createTrainingSessionHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FacilitatorID    int64   `json:"facilitator_id"`
		WorkshopID       int64   `json:"workshop_id"`
		FormationID      int64   `json:"formation_id"`
		RegionID         int64   `json:"region_id"`
		SessionDate      string  `json:"session_date"` // "2025-01-15"
		StartTime        string  `json:"start_time"`   // "09:00"
		EndTime          string  `json:"end_time"`     // "17:00"
		Location         *string `json:"location"`
		MaxCapacity      *int    `json:"max_capacity"`
		TrainingStatusID int64   `json:"training_status_id"`
		Notes            *string `json:"notes"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Parse date and times
	sessionDate, err := time.Parse("2006-01-02", input.SessionDate)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid session_date format, use YYYY-MM-DD"))
		return
	}

	startTime, err := time.Parse("15:04", input.StartTime)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid start_time format, use HH:MM"))
		return
	}

	endTime, err := time.Parse("15:04", input.EndTime)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid end_time format, use HH:MM"))
		return
	}

	session := &data.TrainingSession{
		FacilitatorID:    input.FacilitatorID,
		WorkshopID:       input.WorkshopID,
		FormationID:      input.FormationID,
		RegionID:         input.RegionID,
		SessionDate:      sessionDate,
		StartTime:        startTime,
		EndTime:          endTime,
		Location:         input.Location,
		MaxCapacity:      input.MaxCapacity,
		TrainingStatusID: input.TrainingStatusID,
		Notes:            input.Notes,
	}

	v := validator.New()
	data.ValidateTrainingSession(v, session)

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.TrainingSession.Insert(session)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			app.badRequestResponse(w, r, errors.New("a training session with these details already exists"))
		case errors.Is(err, data.ErrForeignKeyViolation):
			app.badRequestResponse(w, r, errors.New("invalid facilitator_id, workshop_id, formation_id, region_id, or training_status_id"))
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/training-sessions/%d", session.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"training_session": session}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

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

	err = app.writeJSON(w, http.StatusOK, envelope{"training_session": session}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) listTrainingSessionsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	v := validator.New()

	filters := app.readFilters(query, "session_date", 20, []string{"session_date", "-session_date", "id", "-id", "start_time", "-start_time", "created_at", "-created_at"}, v)

	facilitatorID := app.getOptionalInt64QueryParameter(query, "facilitator_id", v)
	workshopID := app.getOptionalInt64QueryParameter(query, "workshop_id", v)
	formationID := app.getOptionalInt64QueryParameter(query, "formation_id", v)
	regionID := app.getOptionalInt64QueryParameter(query, "region_id", v)
	statusID := app.getOptionalInt64QueryParameter(query, "training_status_id", v)

	var sessionDate *time.Time
	sessionDateStr := app.getSingleQueryParameter(query, "session_date", "")
	if sessionDateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", sessionDateStr)
		if err != nil {
			v.AddError("session_date", "invalid date format, use YYYY-MM-DD")
		} else {
			sessionDate = &parsedDate
		}
	}

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	sessions, metadata, err := app.models.TrainingSession.GetAll(facilitatorID, workshopID, formationID, regionID, statusID, sessionDate, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"training_sessions": sessions, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

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

	var input struct {
		FacilitatorID    *int64  `json:"facilitator_id"`
		WorkshopID       *int64  `json:"workshop_id"`
		FormationID      *int64  `json:"formation_id"`
		RegionID         *int64  `json:"region_id"`
		SessionDate      *string `json:"session_date"`
		StartTime        *string `json:"start_time"`
		EndTime          *string `json:"end_time"`
		Location         *string `json:"location"`
		MaxCapacity      *int    `json:"max_capacity"`
		TrainingStatusID *int64  `json:"training_status_id"`
		Notes            *string `json:"notes"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.FacilitatorID != nil {
		session.FacilitatorID = *input.FacilitatorID
	}
	if input.WorkshopID != nil {
		session.WorkshopID = *input.WorkshopID
	}
	if input.FormationID != nil {
		session.FormationID = *input.FormationID
	}
	if input.RegionID != nil {
		session.RegionID = *input.RegionID
	}
	if input.SessionDate != nil {
		sessionDate, err := time.Parse("2006-01-02", *input.SessionDate)
		if err != nil {
			app.badRequestResponse(w, r, errors.New("invalid session_date format, use YYYY-MM-DD"))
			return
		}
		session.SessionDate = sessionDate
	}
	if input.StartTime != nil {
		startTime, err := time.Parse("15:04", *input.StartTime)
		if err != nil {
			app.badRequestResponse(w, r, errors.New("invalid start_time format, use HH:MM"))
			return
		}
		session.StartTime = startTime
	}
	if input.EndTime != nil {
		endTime, err := time.Parse("15:04", *input.EndTime)
		if err != nil {
			app.badRequestResponse(w, r, errors.New("invalid end_time format, use HH:MM"))
			return
		}
		session.EndTime = endTime
	}
	if input.Location != nil {
		session.Location = input.Location
	}
	if input.MaxCapacity != nil {
		session.MaxCapacity = input.MaxCapacity
	}
	if input.TrainingStatusID != nil {
		session.TrainingStatusID = *input.TrainingStatusID
	}
	if input.Notes != nil {
		session.Notes = input.Notes
	}

	v := validator.New()
	data.ValidateTrainingSession(v, session)

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.TrainingSession.Update(session)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		case errors.Is(err, data.ErrDuplicateValue):
			app.badRequestResponse(w, r, errors.New("a training session with these details already exists"))
		case errors.Is(err, data.ErrForeignKeyViolation):
			app.badRequestResponse(w, r, errors.New("invalid facilitator_id, workshop_id, formation_id, region_id, or training_status_id"))
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"training_session": session}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

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

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "training session successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
