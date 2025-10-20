// Filename: cmd/api/workshops.go

package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Pedro-J-Kukul/police_training/internal/data"
	"github.com/Pedro-J-Kukul/police_training/internal/validator"
)

/*********************** Workshops ***********************/

// createWorkshopHandler handles the creation of a new workshop.
//
//	@Summary		Create a new workshop
//	@Description	Create a new workshop
//	@Tags			workshops
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			workshop	body		CreateWorkshopRequest_T	true	"Workshop data"
//	@Success		201			{object}	envelope
//	@Failure		400			{object}	errorResponse
//	@Failure		422			{object}	errorResponse
//	@Failure		500			{object}	errorResponse
//	@Router			/v1/workshops [post]
func (app *appDependencies) createWorkshopHandler(w http.ResponseWriter, r *http.Request) {

	if err := app.readJSON(w, r, &CreateWorkshopRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	workshop := &data.Workshop{
		WorkshopName:   CreateWorkshopRequest.WorkshopName,
		CategoryID:     CreateWorkshopRequest.CategoryID,
		TrainingTypeID: CreateWorkshopRequest.TrainingTypeID,
		CreditHours:    CreateWorkshopRequest.CreditHours,
		Description:    CreateWorkshopRequest.Description,
		Objectives:     CreateWorkshopRequest.Objectives,
		IsActive:       true,
	}
	if CreateWorkshopRequest.IsActive != nil {
		workshop.IsActive = *CreateWorkshopRequest.IsActive
	}

	v := validator.New()
	data.ValidateWorkshop(v, workshop)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Workshop.Insert(workshop); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("workshop_name", "a workshop with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrForeignKeyViolation):
			v.AddError("category_id", "must reference valid category and training type")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/workshops/%d", workshop.ID))

	if err := app.writeJSON(w, http.StatusCreated, envelope{"workshop": workshop}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showWorkshopHandler retrieves and returns a workshop by its ID.
//
//	@Summary		Get a workshop
//	@Description	Retrieve a workshop by its ID
//	@Tags			workshops
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Workshop ID"
//	@Success		200	{object}	envelope
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/v1/workshops/{id} [get]
func (app *appDependencies) showWorkshopHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	workshop, err := app.models.Workshop.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"workshop": workshop}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listWorkshopsHandler returns a filtered list of workshops.
//
//	@Summary		List workshops
//	@Description	Retrieve a list of workshops with optional filtering by name, category, training type, and active status
//	@Tags			workshops
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			workshop_name		query		string	false	"Filter by workshop name"
//	@Param			category_id			query		int		false	"Filter by category ID"
//	@Param			training_type_id	query		int		false	"Filter by training type ID"
//	@Param			is_active			query		bool	false	"Filter by active status"
//	@Param			page				query		int		false	"Page number for pagination"
//	@Param			page_size			query		int		false	"Number of items per page"
//	@Param			sort				query		string	false	"Sort order"
//	@Success		200					{object}	envelope
//	@Failure		422					{object}	errorResponse
//	@Failure		500					{object}	errorResponse
//	@Router			/v1/workshops [get]
func (app *appDependencies) listWorkshopsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	v := validator.New()

	filters := app.readFilters(query, "workshop_name", 20, []string{"workshop_name", "-workshop_name", "id", "-id", "created_at", "-created_at", "updated_at", "-updated_at"}, v)

	categoryID := app.getOptionalInt64QueryParameter(query, "category_id", v)
	trainingTypeID := app.getOptionalInt64QueryParameter(query, "training_type_id", v)
	creditHours := app.getOptionalInt64QueryParameter(query, "credit_hours", v)
	isActive := app.getOptionalBoolQueryParameter(query, "is_active", v)

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	name := likeSearch(app.getSingleQueryParameter(query, "workshop_name", ""))

	workshops, metadata, err := app.models.Workshop.GetAll(name, categoryID, trainingTypeID, creditHours, isActive, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"workshops": workshops, "metadata": metadata}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateWorkshopHandler performs a partial update on a workshop record.
//
//	@Summary		Update a workshop
//	@Description	Perform a partial update on a workshop record
//	@Tags			workshops
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id			path		int						true	"Workshop ID"
//	@Param			workshop	body		UpdateWorkshopRequest_T	true	"Workshop data"
//	@Success		200			{object}	envelope
//	@Failure		400			{object}	errorResponse
//	@Failure		404			{object}	errorResponse
//	@Failure		409			{object}	errorResponse
//	@Failure		422			{object}	errorResponse
//	@Failure		500			{object}	errorResponse
//	@Router			/v1/workshops/{id} [patch]
func (app *appDependencies) updateWorkshopHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	workshop, err := app.models.Workshop.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.readJSON(w, r, &UpdateWorkshopRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if UpdateWorkshopRequest.UpdatedAt == nil || UpdateWorkshopRequest.UpdatedAt.IsZero() {
		app.editConflictResponse(w, r)
		return
	}

	originalUpdatedAt := *UpdateWorkshopRequest.UpdatedAt

	if UpdateWorkshopRequest.WorkshopName != nil {
		workshop.WorkshopName = *UpdateWorkshopRequest.WorkshopName
	}
	if UpdateWorkshopRequest.CategoryID != nil {
		workshop.CategoryID = *UpdateWorkshopRequest.CategoryID
	}
	if UpdateWorkshopRequest.TrainingTypeID != nil {
		workshop.TrainingTypeID = *UpdateWorkshopRequest.TrainingTypeID
	}
	if UpdateWorkshopRequest.CreditHours != nil {
		workshop.CreditHours = *UpdateWorkshopRequest.CreditHours
	}
	if UpdateWorkshopRequest.Description != nil {
		workshop.Description = *UpdateWorkshopRequest.Description
	}
	if UpdateWorkshopRequest.Objectives != nil {
		workshop.Objectives = *UpdateWorkshopRequest.Objectives
	}
	if UpdateWorkshopRequest.IsActive != nil {
		workshop.IsActive = *UpdateWorkshopRequest.IsActive
	}

	v := validator.New()
	data.ValidateWorkshop(v, workshop)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Workshop.Update(workshop, originalUpdatedAt); err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		case errors.Is(err, data.ErrForeignKeyViolation):
			v.AddError("category_id", "must reference valid category and training type")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("workshop_name", "a workshop with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"workshop": workshop}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
