package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Pedro-J-Kukul/police_training/internal/data"
	"github.com/Pedro-J-Kukul/police_training/internal/validator"
)

// @Summary Create a workshop
// @Description Create a new workshop
// @Tags workshops
// @Accept json
// @Produce json
// @Param workshop body map[string]interface{} true "Workshop data (workshop_name, category_id, type_id, credit_hours, description, is_active)"
// @Success 201 {object} map[string]interface{} "Created workshop envelope {\"workshop\": {...}}"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 422 {object} map[string]interface{} "Validation errors"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/workshops [post]
func (app *appDependencies) createWorkshopHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		WorkshopName string  `json:"workshop_name"`
		CategoryID   int64   `json:"category_id"`
		TypeID       int64   `json:"type_id"`
		CreditHours  int     `json:"credit_hours"`
		Description  *string `json:"description"`
		IsActive     *bool   `json:"is_active"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	workshop := &data.Workshop{
		WorkshopName: input.WorkshopName,
		CategoryID:   input.CategoryID,
		TypeID:       input.TypeID,
		CreditHours:  input.CreditHours,
		Description:  input.Description,
		IsActive:     true, // default
	}

	if input.IsActive != nil {
		workshop.IsActive = *input.IsActive
	}

	v := validator.New()
	data.ValidateWorkshop(v, workshop)

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Workshop.Insert(workshop)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("workshop_name", "a workshop with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrForeignKeyViolation):
			app.badRequestResponse(w, r, errors.New("invalid category_id or type_id"))
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/workshops/%d", workshop.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"workshop": workshop}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Get a workshop
// @Description Retrieve a workshop by ID
// @Tags workshops
// @Accept json
// @Produce json
// @Param id path int true "Workshop ID"
// @Success 200 {object} map[string]interface{} "Workshop envelope {\"workshop\": {...}}"
// @Failure 404 {object} map[string]string "Not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/workshops/{id} [get]
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

	err = app.writeJSON(w, http.StatusOK, envelope{"workshop": workshop}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) listWorkshopsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	v := validator.New()

	filters := app.readFilters(query, "workshop_name", 20, []string{"workshop_name", "-workshop_name", "id", "-id", "credit_hours", "-credit_hours", "created_at", "-created_at"}, v)

	name := likeSearch(app.getSingleQueryParameter(query, "workshop_name", ""))
	categoryID := app.getOptionalInt64QueryParameter(query, "category_id", v)
	typeID := app.getOptionalInt64QueryParameter(query, "type_id", v)
	isActive := app.getOptionalBoolQueryParameter(query, "is_active", v)

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	workshops, metadata, err := app.models.Workshop.GetAll(name, categoryID, typeID, isActive, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"workshops": workshops, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

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

	var input struct {
		WorkshopName *string `json:"workshop_name"`
		CategoryID   *int64  `json:"category_id"`
		TypeID       *int64  `json:"type_id"`
		CreditHours  *int    `json:"credit_hours"`
		Description  *string `json:"description"`
		IsActive     *bool   `json:"is_active"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.WorkshopName != nil {
		workshop.WorkshopName = *input.WorkshopName
	}
	if input.CategoryID != nil {
		workshop.CategoryID = *input.CategoryID
	}
	if input.TypeID != nil {
		workshop.TypeID = *input.TypeID
	}
	if input.CreditHours != nil {
		workshop.CreditHours = *input.CreditHours
	}
	if input.Description != nil {
		workshop.Description = input.Description
	}
	if input.IsActive != nil {
		workshop.IsActive = *input.IsActive
	}

	v := validator.New()
	data.ValidateWorkshop(v, workshop)

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Workshop.Update(workshop)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("workshop_name", "a workshop with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrForeignKeyViolation):
			app.badRequestResponse(w, r, errors.New("invalid category_id or type_id"))
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"workshop": workshop}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) deleteWorkshopHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Workshop.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "workshop successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
