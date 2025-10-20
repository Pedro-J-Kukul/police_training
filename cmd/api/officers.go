// Filename: cmd/api/officers.go

package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Pedro-J-Kukul/police_training/internal/data"
	"github.com/Pedro-J-Kukul/police_training/internal/validator"
)

/*********************** Officers ***********************/

// createOfficerHandler handles the creation of a new officer.
//
//	@Summary		Create a new officer
//	@Description	Create a new officer
//	@Tags			officers
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			officer	body		CreateOfficerRequest_T	true	"Officer data"
//	@Success		201		{object}	envelope
//	@Failure		400		{object}	errorResponse
//	@Failure		422		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/v1/officers [post]
func (app *appDependencies) createOfficerHandler(w http.ResponseWriter, r *http.Request) {

	if err := app.readJSON(w, r, &CreateOfficerRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	officer := &data.Officer{
		RegulationNumber: CreateOfficerRequest.RegulationNumber,
		PostingID:        CreateOfficerRequest.PostingID,
		RankID:           CreateOfficerRequest.RankID,
		FormationID:      CreateOfficerRequest.FormationID,
		RegionID:         CreateOfficerRequest.RegionID,
		UserId:           CreateOfficerRequest.UserID,
	}

	v := validator.New()
	data.ValidateOfficer(v, officer)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Officer.Insert(officer); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("regulation_number", "an officer with this regulation number already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrForeignKeyViolation):
			v.AddError("user_id", "must reference an existing activated user and related records")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/officers/%d", officer.ID))

	if err := app.writeJSON(w, http.StatusCreated, envelope{"officer": officer}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showOfficerHandler retrieves and returns an officer by its ID.
//
//	@Summary		Get an officer
//	@Description	Retrieve an officer by its ID
//	@Tags			officers
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Officer ID"
//	@Success		200	{object}	envelope
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/v1/officers/{id} [get]
func (app *appDependencies) showOfficerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	officer, err := app.models.Officer.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"officer": officer}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// getAllOfficersHandler returns a filtered list of officers.
//
//	@Summary		List officers
//	@Description	Retrieve a list of officers with optional filtering by regulation number, posting, rank, formation, and region
//	@Tags			officers
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			regulation_number	query		string	false	"Filter by regulation number"
//	@Param			posting_id			query		int		false	"Filter by posting ID"
//	@Param			rank_id				query		int		false	"Filter by rank ID"
//	@Param			formation_id		query		int		false	"Filter by formation ID"
//	@Param			region_id			query		int		false	"Filter by region ID"
//	@Param			page				query		int		false	"Page number for pagination"
//	@Param			page_size			query		int		false	"Number of items per page"
//	@Param			sort				query		string	false	"Sort order"
//	@Success		200					{object}	envelope
//	@Failure		422					{object}	errorResponse
//	@Failure		500					{object}	errorResponse
//	@Router			/v1/officers [get]
func (app *appDependencies) getAllOfficersHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	v := validator.New()

	filters := app.readFilters(query, "regulation_number", 20, []string{"regulation_number", "-regulation_number", "created_at", "-created_at", "updated_at", "-updated_at", "id", "-id"}, v)

	postingID := app.getOptionalInt64QueryParameter(query, "posting_id", v)
	rankID := app.getOptionalInt64QueryParameter(query, "rank_id", v)
	formationID := app.getOptionalInt64QueryParameter(query, "formation_id", v)
	regionID := app.getOptionalInt64QueryParameter(query, "region_id", v)

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	regulation := likeSearch(app.getSingleQueryParameter(query, "regulation_number", ""))

	officers, metadata, err := app.models.Officer.GetAll(regulation, postingID, rankID, formationID, regionID, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"officers": officers, "metadata": metadata}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateOfficerHandler performs a partial update on an officer record.
//
//	@Summary		Update an officer
//	@Description	Perform a partial update on an officer record
//	@Tags			officers
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		int						true	"Officer ID"
//	@Param			officer	body		UpdateOfficerRequest_T	true	"Officer data"
//	@Success		200		{object}	envelope
//	@Failure		400		{object}	errorResponse
//	@Failure		404		{object}	errorResponse
//	@Failure		409		{object}	errorResponse
//	@Failure		422		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/v1/officers/{id} [patch]
func (app *appDependencies) updateOfficerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	officer, err := app.models.Officer.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.readJSON(w, r, &UpdateOfficerRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if UpdateOfficerRequest.UpdatedAt == nil || UpdateOfficerRequest.UpdatedAt.IsZero() {
		app.editConflictResponse(w, r)
		return
	}

	originalUpdatedAt := *UpdateOfficerRequest.UpdatedAt

	if UpdateOfficerRequest.RegulationNumber != nil {
		officer.RegulationNumber = *UpdateOfficerRequest.RegulationNumber
	}
	if UpdateOfficerRequest.PostingID != nil {
		officer.PostingID = *UpdateOfficerRequest.PostingID
	}
	if UpdateOfficerRequest.RankID != nil {
		officer.RankID = *UpdateOfficerRequest.RankID
	}
	if UpdateOfficerRequest.FormationID != nil {
		officer.FormationID = *UpdateOfficerRequest.FormationID
	}
	if UpdateOfficerRequest.RegionID != nil {
		officer.RegionID = *UpdateOfficerRequest.RegionID
	}

	v := validator.New()
	data.ValidateOfficer(v, officer)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Officer.Update(officer, originalUpdatedAt); err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		case errors.Is(err, data.ErrForeignKeyViolation):
			v.AddError("posting_id", "must reference valid related records")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("regulation_number", "an officer with this regulation number already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"officer": officer}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteOfficerHandler handles the deletion of an officer.
//
//	@Summary		Delete an officer
//	@Description	Delete an officer by its ID
//	@Tags			officers
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Officer ID"
//	@Success		200	{object}	envelope{message=string}
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/v1/officers/{id} [delete]
func (app *appDependencies) deleteOfficerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Officer.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"message": "officer successfully deleted"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
