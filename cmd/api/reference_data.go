package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Pedro-J-Kukul/police_training/internal/data"
	"github.com/Pedro-J-Kukul/police_training/internal/validator"
)

func likeSearch(value string) string {
	if value == "" {
		return ""
	}
	return fmt.Sprintf("%%%s%%", value)
}

/*********************** Regions ***********************/

func (app *appDependencies) createRegionHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Region string `json:"region"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	region := &data.Region{Region: input.Region}

	v := validator.New()
	data.ValidateRegion(v, region)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Region.Insert(region); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("region", "a region with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/regions/%d", region.ID))

	if err := app.writeJSON(w, http.StatusCreated, envelope{"region": region}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) showRegionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	region, err := app.models.Region.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"region": region}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) listRegionsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	v := validator.New()

	filters := app.readFilters(query, "region", 20, []string{"region", "-region", "id", "-id"}, v)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	search := likeSearch(app.getSingleQueryParameter(query, "region", ""))

	regions, metadata, err := app.models.Region.GetAll(search, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"regions": regions, "metadata": metadata}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) updateRegionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	region, err := app.models.Region.Get(id)
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
		Region *string `json:"region"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Region != nil {
		region.Region = *input.Region
	}

	v := validator.New()
	data.ValidateRegion(v, region)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Region.Update(region); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("region", "a region with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"region": region}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

/*********************** Formations ***********************/

func (app *appDependencies) createFormationHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Formation string `json:"formation"`
		RegionID  int64  `json:"region_id"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	formation := &data.Formation{
		Formation: input.Formation,
		RegionID:  input.RegionID,
	}

	v := validator.New()
	data.ValidateFormation(v, formation)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Formation.Insert(formation); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("formation", "a formation with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrForeignKeyViolation):
			v.AddError("region_id", "must reference an existing region")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/formations/%d", formation.ID))

	if err := app.writeJSON(w, http.StatusCreated, envelope{"formation": formation}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) showFormationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	formation, err := app.models.Formation.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"formation": formation}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) listFormationsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	v := validator.New()

	filters := app.readFilters(query, "formation", 20, []string{"formation", "-formation", "id", "-id", "region_id", "-region_id", "created_at", "-created_at"}, v)

	regionID := app.getOptionalInt64QueryParameter(query, "region_id", v)

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	name := likeSearch(app.getSingleQueryParameter(query, "formation", ""))

	formations, metadata, err := app.models.Formation.GetAll(name, regionID, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"formations": formations, "metadata": metadata}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) updateFormationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	formation, err := app.models.Formation.Get(id)
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
		Formation *string `json:"formation"`
		RegionID  *int64  `json:"region_id"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Formation != nil {
		formation.Formation = *input.Formation
	}
	if input.RegionID != nil {
		formation.RegionID = *input.RegionID
	}

	v := validator.New()
	data.ValidateFormation(v, formation)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Formation.Update(formation); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("formation", "a formation with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrForeignKeyViolation):
			v.AddError("region_id", "must reference an existing region")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"formation": formation}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

/*********************** Postings ***********************/

func (app *appDependencies) createPostingHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Posting string  `json:"posting"`
		Code    *string `json:"code"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	posting := &data.Posting{
		Posting: input.Posting,
		Code:    input.Code,
	}

	v := validator.New()
	data.ValidatePosting(v, posting)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Posting.Insert(posting); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("posting", "a posting with these details already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/postings/%d", posting.ID))

	if err := app.writeJSON(w, http.StatusCreated, envelope{"posting": posting}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) showPostingHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	posting, err := app.models.Posting.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"posting": posting}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) listPostingsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	v := validator.New()

	filters := app.readFilters(query, "posting", 20, []string{"posting", "-posting", "id", "-id", "code", "-code", "created_at", "-created_at"}, v)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	name := likeSearch(app.getSingleQueryParameter(query, "posting", ""))
	code := likeSearch(app.getSingleQueryParameter(query, "code", ""))

	postings, metadata, err := app.models.Posting.GetAll(name, code, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"postings": postings, "metadata": metadata}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) updatePostingHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	posting, err := app.models.Posting.Get(id)
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
		Posting *string  `json:"posting"`
		Code    **string `json:"code"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Posting != nil {
		posting.Posting = *input.Posting
	}
	if input.Code != nil {
		posting.Code = *input.Code
	}

	v := validator.New()
	data.ValidatePosting(v, posting)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Posting.Update(posting); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("posting", "a posting with these details already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"posting": posting}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

/*********************** Ranks ***********************/

func (app *appDependencies) createRankHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Rank                        string `json:"rank"`
		Code                        string `json:"code"`
		AnnualTrainingHoursRequired int    `json:"annual_training_hours_required"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	rank := &data.Rank{
		Rank:                        input.Rank,
		Code:                        input.Code,
		AnnualTrainingHoursRequired: input.AnnualTrainingHoursRequired,
	}

	v := validator.New()
	data.ValidateRank(v, rank)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Rank.Insert(rank); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("rank", "a rank with these details already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/ranks/%d", rank.ID))

	if err := app.writeJSON(w, http.StatusCreated, envelope{"rank": rank}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) showRankHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	rank, err := app.models.Rank.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"rank": rank}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) listRanksHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	v := validator.New()

	filters := app.readFilters(query, "rank", 20, []string{"rank", "-rank", "code", "-code", "annual_training_hours_required", "-annual_training_hours_required", "id", "-id"}, v)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	rankFilter := likeSearch(app.getSingleQueryParameter(query, "rank", ""))
	codeFilter := likeSearch(app.getSingleQueryParameter(query, "code", ""))

	ranks, metadata, err := app.models.Rank.GetAll(rankFilter, codeFilter, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"ranks": ranks, "metadata": metadata}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) updateRankHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	rank, err := app.models.Rank.Get(id)
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
		Rank                        *string `json:"rank"`
		Code                        *string `json:"code"`
		AnnualTrainingHoursRequired *int    `json:"annual_training_hours_required"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Rank != nil {
		rank.Rank = *input.Rank
	}
	if input.Code != nil {
		rank.Code = *input.Code
	}
	if input.AnnualTrainingHoursRequired != nil {
		rank.AnnualTrainingHoursRequired = *input.AnnualTrainingHoursRequired
	}

	v := validator.New()
	data.ValidateRank(v, rank)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Rank.Update(rank); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("rank", "a rank with these details already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"rank": rank}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

/*********************** Training Types ***********************/

func (app *appDependencies) createTrainingTypeHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Type string `json:"type"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	trainingType := &data.TrainingType{Type: input.Type}

	v := validator.New()
	data.ValidateTrainingType(v, trainingType)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.TrainingType.Insert(trainingType); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("type", "a training type with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/training-types/%d", trainingType.ID))

	if err := app.writeJSON(w, http.StatusCreated, envelope{"training_type": trainingType}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) showTrainingTypeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	typeRecord, err := app.models.TrainingType.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"training_type": typeRecord}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) listTrainingTypesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	v := validator.New()

	filters := app.readFilters(query, "type", 20, []string{"type", "-type", "id", "-id"}, v)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	name := likeSearch(app.getSingleQueryParameter(query, "type", ""))

	types, metadata, err := app.models.TrainingType.GetAll(name, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"training_types": types, "metadata": metadata}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) updateTrainingTypeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	typeRecord, err := app.models.TrainingType.Get(id)
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
		Type *string `json:"type"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Type != nil {
		typeRecord.Type = *input.Type
	}

	v := validator.New()
	data.ValidateTrainingType(v, typeRecord)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.TrainingType.Update(typeRecord); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("type", "a training type with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"training_type": typeRecord}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

/*********************** Training Categories ***********************/

func (app *appDependencies) createTrainingCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		IsActive *bool  `json:"is_active"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	category := &data.TrainingCategory{
		Name:     input.Name,
		IsActive: true,
	}
	if input.IsActive != nil {
		category.IsActive = *input.IsActive
	}

	v := validator.New()
	data.ValidateTrainingCategory(v, category)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.TrainingCategory.Insert(category); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("name", "a training category with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/training-categories/%d", category.ID))

	if err := app.writeJSON(w, http.StatusCreated, envelope{"training_category": category}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) showTrainingCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	category, err := app.models.TrainingCategory.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"training_category": category}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) listTrainingCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	v := validator.New()

	filters := app.readFilters(query, "name", 20, []string{"name", "-name", "id", "-id", "created_at", "-created_at"}, v)

	isActive := app.getOptionalBoolQueryParameter(query, "is_active", v)

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	name := likeSearch(app.getSingleQueryParameter(query, "name", ""))

	categories, metadata, err := app.models.TrainingCategory.GetAll(name, isActive, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"training_categories": categories, "metadata": metadata}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) updateTrainingCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	category, err := app.models.TrainingCategory.Get(id)
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
		Name     *string `json:"name"`
		IsActive *bool   `json:"is_active"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		category.Name = *input.Name
	}
	if input.IsActive != nil {
		category.IsActive = *input.IsActive
	}

	v := validator.New()
	data.ValidateTrainingCategory(v, category)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.TrainingCategory.Update(category); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("name", "a training category with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"training_category": category}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

/*********************** Workshops ***********************/

func (app *appDependencies) createWorkshopHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		WorkshopName   string  `json:"workshop_name"`
		CategoryID     int64   `json:"category_id"`
		TrainingTypeID int64   `json:"training_type_id"`
		CreditHours    int     `json:"credit_hours"`
		Description    *string `json:"description"`
		Objectives     *string `json:"objectives"`
		IsActive       *bool   `json:"is_active"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	workshop := &data.Workshop{
		WorkshopName:   input.WorkshopName,
		CategoryID:     input.CategoryID,
		TrainingTypeID: input.TrainingTypeID,
		CreditHours:    input.CreditHours,
		Description:    input.Description,
		Objectives:     input.Objectives,
		IsActive:       true,
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

func (app *appDependencies) listWorkshopsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	v := validator.New()

	filters := app.readFilters(query, "workshop_name", 20, []string{"workshop_name", "-workshop_name", "id", "-id", "created_at", "-created_at", "updated_at", "-updated_at"}, v)

	categoryID := app.getOptionalInt64QueryParameter(query, "category_id", v)
	trainingTypeID := app.getOptionalInt64QueryParameter(query, "training_type_id", v)
	isActive := app.getOptionalBoolQueryParameter(query, "is_active", v)

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	name := likeSearch(app.getSingleQueryParameter(query, "workshop_name", ""))

	workshops, metadata, err := app.models.Workshop.GetAll(name, categoryID, trainingTypeID, isActive, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"workshops": workshops, "metadata": metadata}, nil); err != nil {
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
		WorkshopName   *string    `json:"workshop_name"`
		CategoryID     *int64     `json:"category_id"`
		TrainingTypeID *int64     `json:"training_type_id"`
		CreditHours    *int       `json:"credit_hours"`
		Description    **string   `json:"description"`
		Objectives     **string   `json:"objectives"`
		IsActive       *bool      `json:"is_active"`
		UpdatedAt      *time.Time `json:"updated_at"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.UpdatedAt == nil || input.UpdatedAt.IsZero() {
		app.editConflictResponse(w, r)
		return
	}

	originalUpdatedAt := *input.UpdatedAt

	if input.WorkshopName != nil {
		workshop.WorkshopName = *input.WorkshopName
	}
	if input.CategoryID != nil {
		workshop.CategoryID = *input.CategoryID
	}
	if input.TrainingTypeID != nil {
		workshop.TrainingTypeID = *input.TrainingTypeID
	}
	if input.CreditHours != nil {
		workshop.CreditHours = *input.CreditHours
	}
	if input.Description != nil {
		workshop.Description = *input.Description
	}
	if input.Objectives != nil {
		workshop.Objectives = *input.Objectives
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

/*********************** Training Status ***********************/

func (app *appDependencies) createTrainingStatusHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Status string `json:"status"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	status := &data.TrainingStatus{Status: input.Status}

	v := validator.New()
	data.ValidateTrainingStatus(v, status)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.TrainingStatus.Insert(status); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("status", "a training status with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/training-statuses/%d", status.ID))

	if err := app.writeJSON(w, http.StatusCreated, envelope{"training_status": status}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) showTrainingStatusHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	status, err := app.models.TrainingStatus.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"training_status": status}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) listTrainingStatusesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	v := validator.New()

	filters := app.readFilters(query, "status", 20, []string{"status", "-status", "id", "-id", "created_at", "-created_at"}, v)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	name := likeSearch(app.getSingleQueryParameter(query, "status", ""))

	statuses, metadata, err := app.models.TrainingStatus.GetAll(name, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"training_statuses": statuses, "metadata": metadata}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) updateTrainingStatusHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	status, err := app.models.TrainingStatus.Get(id)
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
		Status *string `json:"status"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Status != nil {
		status.Status = *input.Status
	}

	v := validator.New()
	data.ValidateTrainingStatus(v, status)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.TrainingStatus.Update(status); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("status", "a training status with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"training_status": status}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

/*********************** Enrollment Status ***********************/

func (app *appDependencies) createEnrollmentStatusHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Status string `json:"status"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	status := &data.EnrollmentStatus{Status: input.Status}

	v := validator.New()
	data.ValidateEnrollmentStatus(v, status)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.EnrollmentStatus.Insert(status); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("status", "an enrollment status with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/enrollment-statuses/%d", status.ID))

	if err := app.writeJSON(w, http.StatusCreated, envelope{"enrollment_status": status}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) showEnrollmentStatusHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	status, err := app.models.EnrollmentStatus.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"enrollment_status": status}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) listEnrollmentStatusesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	v := validator.New()

	filters := app.readFilters(query, "status", 20, []string{"status", "-status", "id", "-id", "created_at", "-created_at"}, v)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	name := likeSearch(app.getSingleQueryParameter(query, "status", ""))

	statuses, metadata, err := app.models.EnrollmentStatus.GetAll(name, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"enrollment_statuses": statuses, "metadata": metadata}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) updateEnrollmentStatusHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	status, err := app.models.EnrollmentStatus.Get(id)
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
		Status *string `json:"status"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Status != nil {
		status.Status = *input.Status
	}

	v := validator.New()
	data.ValidateEnrollmentStatus(v, status)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.EnrollmentStatus.Update(status); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("status", "an enrollment status with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"enrollment_status": status}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

/*********************** Officers ***********************/

func (app *appDependencies) createOfficerHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID           int64  `json:"user_id"`
		RegulationNumber string `json:"regulation_number"`
		PostingID        int64  `json:"posting_id"`
		RankID           int64  `json:"rank_id"`
		FormationID      int64  `json:"formation_id"`
		RegionID         int64  `json:"region_id"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	officer := &data.Officer{
		RegulationNumber: input.RegulationNumber,
		PostingID:        input.PostingID,
		RankID:           input.RankID,
		FormationID:      input.FormationID,
		RegionID:         input.RegionID,
		UserId:           input.UserID,
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
		RegulationNumber *string    `json:"regulation_number"`
		PostingID        *int64     `json:"posting_id"`
		RankID           *int64     `json:"rank_id"`
		FormationID      *int64     `json:"formation_id"`
		RegionID         *int64     `json:"region_id"`
		UpdatedAt        *time.Time `json:"updated_at"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.UpdatedAt == nil || input.UpdatedAt.IsZero() {
		app.editConflictResponse(w, r)
		return
	}

	originalUpdatedAt := *input.UpdatedAt

	if input.RegulationNumber != nil {
		officer.RegulationNumber = *input.RegulationNumber
	}
	if input.PostingID != nil {
		officer.PostingID = *input.PostingID
	}
	if input.RankID != nil {
		officer.RankID = *input.RankID
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
