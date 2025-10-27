// Filename: cmd/api/officers.go
package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Pedro-J-Kukul/police_training/internal/data"
	"github.com/Pedro-J-Kukul/police_training/internal/validator"
)

// createOfficerHandler creates a new officer record
//
//	@Summary		Create a new officer
//	@Description	Create a new officer record linked to a user
//	@Tags			officers
//	@Accept			json
//	@Produce		json
//	@Param			officer	body		CreateOfficerRequest_T	true	"Officer data"
//	@Success		201		{object}	envelope
//	@Failure		400		{object}	errorResponse
//	@Failure		422		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/v1/officers [post]
func (app *appDependencies) createOfficerHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID           int64  `json:"user_id"`
		RegulationNumber string `json:"regulation_number"`
		RankID           int64  `json:"rank_id"`
		PostingID        int64  `json:"posting_id"`
		FormationID      int64  `json:"formation_id"`
		RegionID         int64  `json:"region_id"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	officer := &data.Officer{
		UserID:           input.UserID,
		RegulationNumber: input.RegulationNumber,
		RankID:           input.RankID,
		PostingID:        input.PostingID,
		FormationID:      input.FormationID,
		RegionID:         input.RegionID,
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
			v.AddError("references", "one or more referenced records do not exist")
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

// showOfficerHandler retrieves an officer by ID
//
//	@Summary		Get an officer
//	@Description	Retrieve an officer by their ID
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

// showOfficerWithDetailsHandler retrieves an officer with all related information
//
//	@Summary		Get an officer with details
//	@Description	Retrieve an officer with all related information (user, rank, posting, etc.)
//	@Tags			officers
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Officer ID"
//	@Success		200	{object}	envelope
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/v1/officers/{id}/details [get]
func (app *appDependencies) showOfficerWithDetailsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	officer, err := app.models.Officer.GetWithDetails(id)
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

// listOfficersHandler returns a filtered list of officers
//
//	@Summary		List officers
//	@Description	Retrieve a list of officers with optional filtering and pagination
//	@Tags			officers
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			regulation_number	query		string	false	"Filter by regulation number"
//	@Param			rank_id				query		int		false	"Filter by rank ID"
//	@Param			posting_id			query		int		false	"Filter by posting ID"
//	@Param			formation_id		query		int		false	"Filter by formation ID"
//	@Param			region_id			query		int		false	"Filter by region ID"
//	@Success		200					{object}	envelope
//	@Failure		500					{object}	errorResponse
//	@Router			/v1/officers [get]
func (app *appDependencies) listOfficersHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	v := validator.New()

	filters := data.Filters{
		Page:     app.getSingleIntQueryParameter(query, "page", 1, v),
		PageSize: app.getSingleIntQueryParameter(query, "page_size", 20, v),
		Sort:     app.getSingleQueryParameter(query, "sort", "id"),
		SortSafelist: []string{
			"id", "regulation_number", "created_at",
			"-id", "-regulation_number", "-created_at",
		},
	}

	data.ValidateFilters(v, filters)

	regulationNumber := app.getSingleQueryParameter(query, "regulation_number", "")
	rankID := app.getOptionalInt64QueryParameter(query, "rank_id", v)
	postingID := app.getOptionalInt64QueryParameter(query, "posting_id", v)
	formationID := app.getOptionalInt64QueryParameter(query, "formation_id", v)
	regionID := app.getOptionalInt64QueryParameter(query, "region_id", v)

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	officers, metadata, err := app.models.Officer.GetAll(regulationNumber, rankID, postingID, formationID, regionID, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	payload := envelope{
		"officers": officers,
		"metadata": metadata,
	}

	if err := app.writeJSON(w, http.StatusOK, payload, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateOfficerHandler updates an existing officer
//
//	@Summary		Update an officer
//	@Description	Update an existing officer's information
//	@Tags			officers
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		int					true	"Officer ID"
//	@Param			officer	body		UpdateOfficerRequest_T	true	"Officer update data"
//	@Success		200		{object}	envelope
//	@Failure		400		{object}	errorResponse
//	@Failure		404		{object}	errorResponse
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

	var input struct {
		RegulationNumber *string `json:"regulation_number"`
		RankID           *int64  `json:"rank_id"`
		PostingID        *int64  `json:"posting_id"`
		FormationID      *int64  `json:"formation_id"`
		RegionID         *int64  `json:"region_id"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.RegulationNumber != nil {
		officer.RegulationNumber = *input.RegulationNumber
	}
	if input.RankID != nil {
		officer.RankID = *input.RankID
	}
	if input.PostingID != nil {
		officer.PostingID = *input.PostingID
	}
	if input.FormationID != nil {
		officer.FormationID = *input.FormationID
	}
	if input.RegionID != nil {
		officer.RegionID = *input.RegionID
	}

	v := validator.New()
	data.ValidateOfficer(v, officer)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Officer.Update(officer); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("regulation_number", "an officer with this regulation number already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		case errors.Is(err, data.ErrForeignKeyViolation):
			v.AddError("references", "one or more referenced records do not exist")
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

// deleteOfficerHandler removes an officer
//
//	@Summary		Delete an officer
//	@Description	Remove an officer record
//	@Tags			officers
//	@Produce		json
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

// getUserOfficerHandler retrieves an officer by user ID
//
//	@Summary		Get officer by user ID
//	@Description	Retrieve an officer record by their associated user ID
//	@Tags			officers
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"User ID"
//	@Success		200		{object}	envelope
//	@Failure		404		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/v1/users/{user_id}/officer [get]
func (app *appDependencies) getUserOfficerHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	officer, err := app.models.Officer.GetByUserID(userID)
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
