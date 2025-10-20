package main

import (
	"errors"
	"fmt"
	"net/http"

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
// createRegionHandler handles the creation of a new region.
//	@Summary		Create a new region
//	@Description	Create a new region
//	@Tags			regions
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			region	body		CreateRegionRequest_T	true	"Region data"
//	@Success		201		{object}	envelope
//	@Failure		400		{object}	errorResponse
//	@Failure		422		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/v1/regions [post]
func (app *appDependencies) createRegionHandler(w http.ResponseWriter, r *http.Request) {

	if err := app.readJSON(w, r, &CreateRegionRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	region := &data.Region{Region: CreateRegionRequest.Region}

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

// showRegionHandler retrieves and returns a region by its ID.
//
//	@Summary		Get a region
//	@Description	Retrieve a region by its ID
//	@Tags			regions
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Region ID"
//	@Success		200	{object}	envelope
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/v1/regions/{id} [get]
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

// listRegionsHandler returns a filtered list of regions.
//
//	@Summary		List regions
//	@Description	Retrieve a list of regions with optional filtering
//	@Tags			regions
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			region		query		string	false	"Filter by region name"
//	@Param			page		query		int		false	"Page number for pagination"
//	@Param			page_size	query		int		false	"Number of items per page"
//	@Param			sort		query		string	false	"Sort order"
//	@Success		200			{object}	envelope
//	@Failure		422			{object}	errorResponse
//	@Failure		500			{object}	errorResponse
//	@Router			/v1/regions [get]
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

// updateRegionHandler performs a partial update on a region record.
//
//	@Summary		Update a region
//	@Description	Perform a partial update on a region record
//	@Tags			regions
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		int						true	"Region ID"
//	@Param			region	body		UpdateRegionRequest_T	true	"Region data"
//	@Success		200		{object}	envelope
//	@Failure		400		{object}	errorResponse
//	@Failure		404		{object}	errorResponse
//	@Failure		422		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/v1/regions/{id} [patch]
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

	if err := app.readJSON(w, r, &UpdateRegionRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if UpdateRegionRequest.Region != nil {
		region.Region = *UpdateRegionRequest.Region
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
// createFormationHandler handles the creation of a new formation.
//	@Summary		Create a new formation
//	@Description	Create a new formation
//	@Tags			formations
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			formation	body		CreateFormationRequest_T	true	"Formation data"
//	@Success		201			{object}	envelope
//	@Failure		400			{object}	errorResponse
//	@Failure		422			{object}	errorResponse
//	@Failure		500			{object}	errorResponse
//	@Router			/v1/formations [post]
func (app *appDependencies) createFormationHandler(w http.ResponseWriter, r *http.Request) {

	if err := app.readJSON(w, r, &UpdateFormationRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	formation := &data.Formation{
		Formation: *UpdateFormationRequest.Formation,
		RegionID:  *UpdateFormationRequest.RegionID,
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

// showFormationHandler retrieves a formation by id.
//
//	@Summary		Get a formation
//	@Description	Retrieve a formation by its ID
//	@Tags			formations
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Formation ID"
//	@Success		200	{object}	envelope
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/v1/formations/{id} [get]
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

// listFormationsHandler returns formations filtered by name and region.
//
//	@Summary		List formations
//	@Description	Retrieve a list of formations with optional filtering by name and region
//	@Tags			formations
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			formation	query		string	false	"Filter by formation name"
//	@Param			region_id	query		int		false	"Filter by region ID"
//	@Param			page		query		int		false	"Page number for pagination"
//	@Param			page_size	query		int		false	"Number of items per page"
//	@Param			sort		query		string	false	"Sort order"
//	@Success		200			{object}	envelope
//	@Failure		422			{object}	errorResponse
//	@Failure		500			{object}	errorResponse
//	@Router			/v1/formations [get]
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

// updateFormationHandler performs a partial update on a formation record.
//
//	@Summary		Update a formation
//	@Description	Perform a partial update on a formation record
//	@Tags			formations
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id			path		int							true	"Formation ID"
//	@Param			formation	body		UpdateFormationRequest_T	true	"Formation data"
//	@Success		200			{object}	envelope
//	@Failure		400			{object}	errorResponse
//	@Failure		404			{object}	errorResponse
//	@Failure		422			{object}	errorResponse
//	@Failure		500			{object}	errorResponse
//	@Router			/v1/formations/{id} [patch]
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

	if err := app.readJSON(w, r, &UpdateFormationRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if UpdateFormationRequest.Formation != nil {
		formation.Formation = *UpdateFormationRequest.Formation
	}
	if UpdateFormationRequest.RegionID != nil {
		formation.RegionID = *UpdateFormationRequest.RegionID
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

// createPostingHandler handles the creation of a new posting.
//
//	@Summary		Create a new posting
//	@Description	Create a new posting
//	@Tags			postings
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			posting	body		CreatePostingRequest_T	true	"Posting data"
//	@Success		201		{object}	envelope
//	@Failure		400		{object}	errorResponse
//	@Failure		422		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/v1/postings [post]
func (app *appDependencies) createPostingHandler(w http.ResponseWriter, r *http.Request) {

	if err := app.readJSON(w, r, &CreatePostingRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	posting := &data.Posting{
		Posting: CreatePostingRequest.Posting,
		Code:    CreatePostingRequest.Code,
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

// showPostingHandler retrieves and returns a posting by its ID.
//
//	@Summary		Get a posting
//	@Description	Retrieve a posting by its ID
//	@Tags			postings
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Posting ID"
//	@Success		200	{object}	envelope
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/v1/postings/{id} [get]
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

// listPostingsHandler returns a filtered list of postings.
//
//	@Summary		List postings
//	@Description	Retrieve a list of postings with optional filtering by name and code
//	@Tags			postings
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			posting		query		string	false	"Filter by posting name"
//	@Param			code		query		string	false	"Filter by posting code"
//	@Param			page		query		int		false	"Page number for pagination"
//	@Param			page_size	query		int		false	"Number of items per page"
//	@Param			sort		query		string	false	"Sort order"
//	@Success		200			{object}	envelope
//	@Failure		422			{object}	errorResponse
//	@Failure		500			{object}	errorResponse
//	@Router			/v1/postings [get]
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

// updatePostingHandler performs a partial update on a posting record.
//
//	@Summary		Update a posting
//	@Description	Perform a partial update on a posting record
//	@Tags			postings
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		int						true	"Posting ID"
//	@Param			posting	body		UpdatePostingRequest_T	true	"Posting data"
//	@Success		200		{object}	envelope
//	@Failure		400		{object}	errorResponse
//	@Failure		404		{object}	errorResponse
//	@Failure		422		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/v1/postings/{id} [patch]
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

	if err := app.readJSON(w, r, &UpdatePostingRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if UpdatePostingRequest.Posting != nil {
		posting.Posting = *UpdatePostingRequest.Posting
	}
	if UpdatePostingRequest.Code != nil {
		posting.Code = UpdatePostingRequest.Code
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

// createRankHandler handles the creation of a new rank.
//
//	@Summary		Create a new rank
//	@Description	Create a new rank
//	@Tags			ranks
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			rank	body		CreateRankRequest_T	true	"Rank data"
//	@Success		201		{object}	envelope
//	@Failure		400		{object}	errorResponse
//	@Failure		422		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/v1/ranks [post]
func (app *appDependencies) createRankHandler(w http.ResponseWriter, r *http.Request) {

	if err := app.readJSON(w, r, &CreateRankRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	rank := &data.Rank{
		Rank:                        CreateRankRequest.Rank,
		Code:                        CreateRankRequest.Code,
		AnnualTrainingHoursRequired: CreateRankRequest.AnnualTrainingHoursRequired,
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

// showRankHandler retrieves and returns a rank by its ID.
//
//	@Summary		Get a rank
//	@Description	Retrieve a rank by its ID
//	@Tags			ranks
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Rank ID"
//	@Success		200	{object}	envelope
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/v1/ranks/{id} [get]
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

// listRanksHandler returns a filtered list of ranks.
//
//	@Summary		List ranks
//	@Description	Retrieve a list of ranks with optional filtering by rank name and code
//	@Tags			ranks
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			rank		query		string	false	"Filter by rank name"
//	@Param			code		query		string	false	"Filter by rank code"
//	@Param			page		query		int		false	"Page number for pagination"
//	@Param			page_size	query		int		false	"Number of items per page"
//	@Param			sort		query		string	false	"Sort order"
//	@Success		200			{object}	envelope
//	@Failure		422			{object}	errorResponse
//	@Failure		500			{object}	errorResponse
//	@Router			/v1/ranks [get]
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

// updateRankHandler performs a partial update on a rank record.
//
//	@Summary		Update a rank
//	@Description	Perform a partial update on a rank record
//	@Tags			ranks
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		int					true	"Rank ID"
//	@Param			rank	body		UpdateRankRequest_T	true	"Rank data"
//	@Success		200		{object}	envelope
//	@Failure		400		{object}	errorResponse
//	@Failure		404		{object}	errorResponse
//	@Failure		422		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/v1/ranks/{id} [patch]
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

	if err := app.readJSON(w, r, &UpdateRankRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if UpdateRankRequest.Rank != nil {
		rank.Rank = *UpdateRankRequest.Rank
	}
	if UpdateRankRequest.Code != nil {
		rank.Code = *UpdateRankRequest.Code
	}
	if UpdateRankRequest.AnnualTrainingHoursRequired != nil {
		rank.AnnualTrainingHoursRequired = *UpdateRankRequest.AnnualTrainingHoursRequired
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

// createTrainingTypeHandler handles the creation of a new training type.
//
//	@Summary		Create a new training type
//	@Description	Create a new training type
//	@Tags			training-types
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			training_type	body		CreateTrainingTypeRequest_T	true	"Training type data"
//	@Success		201				{object}	envelope
//	@Failure		400				{object}	errorResponse
//	@Failure		422				{object}	errorResponse
//	@Failure		500				{object}	errorResponse
//	@Router			/v1/training/types [post]
func (app *appDependencies) createTrainingTypeHandler(w http.ResponseWriter, r *http.Request) {

	if err := app.readJSON(w, r, &CreateTrainingTypeRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	trainingType := &data.TrainingType{Type: CreateTrainingTypeRequest.Type}

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

// showTrainingTypeHandler retrieves and returns a training type by its ID.
//
//	@Summary		Get a training type
//	@Description	Retrieve a training type by its ID
//	@Tags			training-types
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Training type ID"
//	@Success		200	{object}	envelope
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/v1/training/types/{id} [get]
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

// listTrainingTypesHandler returns a filtered list of training types.
//
//	@Summary		List training types
//	@Description	Retrieve a list of training types with optional filtering by type name
//	@Tags			training-types
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			type		query		string	false	"Filter by training type name"
//	@Param			page		query		int		false	"Page number for pagination"
//	@Param			page_size	query		int		false	"Number of items per page"
//	@Param			sort		query		string	false	"Sort order"
//	@Success		200			{object}	envelope
//	@Failure		422			{object}	errorResponse
//	@Failure		500			{object}	errorResponse
//	@Router			/v1/training/types [get]
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

// updateTrainingTypeHandler performs a partial update on a training type record.
//
//	@Summary		Update a training type
//	@Description	Perform a partial update on a training type record
//	@Tags			training-types
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id				path		int							true	"Training type ID"
//	@Param			training_type	body		UpdateTrainingTypeRequest_T	true	"Training type data"
//	@Success		200				{object}	envelope
//	@Failure		400				{object}	errorResponse
//	@Failure		404				{object}	errorResponse
//	@Failure		422				{object}	errorResponse
//	@Failure		500				{object}	errorResponse
//	@Router			/v1/training/types/{id} [patch]
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

	if err := app.readJSON(w, r, &UpdateTrainingTypeRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if UpdateTrainingTypeRequest.Type != nil {
		typeRecord.Type = *UpdateTrainingTypeRequest.Type
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

// createTrainingCategoryHandler handles the creation of a new training category.
//
//	@Summary		Create a new training category
//	@Description	Create a new training category
//	@Tags			training-categories
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			training_category	body		CreateTrainingCategoryRequest_T	true	"Training category data"
//	@Success		201					{object}	envelope
//	@Failure		400					{object}	errorResponse
//	@Failure		422					{object}	errorResponse
//	@Failure		500					{object}	errorResponse
//	@Router			/v1/training/categories [post]
func (app *appDependencies) createTrainingCategoryHandler(w http.ResponseWriter, r *http.Request) {

	if err := app.readJSON(w, r, &CreateTrainingCategoryRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	category := &data.TrainingCategory{
		Name:     CreateTrainingCategoryRequest.Name,
		IsActive: true,
	}
	if CreateTrainingCategoryRequest.IsActive != nil {
		category.IsActive = *CreateTrainingCategoryRequest.IsActive
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
	headers.Set("Location", fmt.Sprintf("/v1/training/categories/%d", category.ID))

	if err := app.writeJSON(w, http.StatusCreated, envelope{"training_category": category}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showTrainingCategoryHandler retrieves and returns a training category by its ID.
//
//	@Summary		Get a training category
//	@Description	Retrieve a training category by its ID
//	@Tags			training-categories
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Training category ID"
//	@Success		200	{object}	envelope
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/v1/training/categories/{id} [get]
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

// listTrainingCategoriesHandler returns a filtered list of training categories.
//
//	@Summary		List training categories
//	@Description	Retrieve a list of training categories with optional filtering by name and active status
//	@Tags			training-categories
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			name		query		string	false	"Filter by category name"
//	@Param			is_active	query		bool	false	"Filter by active status"
//	@Param			page		query		int		false	"Page number for pagination"
//	@Param			page_size	query		int		false	"Number of items per page"
//	@Param			sort		query		string	false	"Sort order"
//	@Success		200			{object}	envelope
//	@Failure		422			{object}	errorResponse
//	@Failure		500			{object}	errorResponse
//	@Router			/v1/training/categories [get]
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

// updateTrainingCategoryHandler performs a partial update on a training category record.
//
//	@Summary		Update a training category
//	@Description	Perform a partial update on a training category record
//	@Tags			training-categories
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id					path		int								true	"Training category ID"
//	@Param			training_category	body		UpdateTrainingCategoryRequest_T	true	"Training category data"
//	@Success		200					{object}	envelope
//	@Failure		400					{object}	errorResponse
//	@Failure		404					{object}	errorResponse
//	@Failure		422					{object}	errorResponse
//	@Failure		500					{object}	errorResponse
//	@Router			/v1/training/categories/{id} [patch]
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

	if err := app.readJSON(w, r, &UpdateTrainingCategoryRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if UpdateTrainingCategoryRequest.Name != nil {
		category.Name = *UpdateTrainingCategoryRequest.Name
	}
	if UpdateTrainingCategoryRequest.IsActive != nil {
		category.IsActive = *UpdateTrainingCategoryRequest.IsActive
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

/*********************** Training Status ***********************/

// createTrainingStatusHandler handles the creation of a new training status.
//
//	@Summary		Create a new training status
//	@Description	Create a new training status
//	@Tags			training-statuses
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			training_status	body		CreateTrainingStatusRequest_T	true	"Training status data"
//	@Success		201				{object}	envelope
//	@Failure		400				{object}	errorResponse
//	@Failure		422				{object}	errorResponse
//	@Failure		500				{object}	errorResponse
//	@Router			/v1/training/status [post]
func (app *appDependencies) createTrainingStatusHandler(w http.ResponseWriter, r *http.Request) {

	if err := app.readJSON(w, r, &CreateTrainingStatusRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	status := &data.TrainingStatus{Status: CreateTrainingStatusRequest.Status}

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
	headers.Set("Location", fmt.Sprintf("/v1/training/status/%d", status.ID))

	if err := app.writeJSON(w, http.StatusCreated, envelope{"training_status": status}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showTrainingStatusHandler retrieves and returns a training status by its ID.
//
//	@Summary		Get a training status
//	@Description	Retrieve a training status by its ID
//	@Tags			training-statuses
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Training status ID"
//	@Success		200	{object}	envelope
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/v1/training/status/{id} [get]
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

// getTrainingStatusesHandler returns a filtered list of training statuses.
//
//	@Summary		List training statuses
//	@Description	Retrieve a list of training statuses with optional filtering by status name
//	@Tags			training-statuses
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			status		query		string	false	"Filter by status name"
//	@Param			page		query		int		false	"Page number for pagination"
//	@Param			page_size	query		int		false	"Number of items per page"
//	@Param			sort		query		string	false	"Sort order"
//	@Success		200			{object}	envelope
//	@Failure		422			{object}	errorResponse
//	@Failure		500			{object}	errorResponse
//	@Router			/v1/training/status [get]
func (app *appDependencies) getTrainingStatusesHandler(w http.ResponseWriter, r *http.Request) {
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

// updateTrainingStatusHandler performs a partial update on a training status record.
//
//	@Summary		Update a training status
//	@Description	Perform a partial update on a training status record
//	@Tags			training-statuses
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id				path		int								true	"Training status ID"
//	@Param			training_status	body		UpdateTrainingStatusRequest_T	true	"Training status data"
//	@Success		200				{object}	envelope
//	@Failure		400				{object}	errorResponse
//	@Failure		404				{object}	errorResponse
//	@Failure		422				{object}	errorResponse
//	@Failure		500				{object}	errorResponse
//	@Router			/v1/training/status/{id} [patch]
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

	if err := app.readJSON(w, r, &UpdateTrainingStatusRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if UpdateTrainingStatusRequest.Status != nil {
		status.Status = *UpdateTrainingStatusRequest.Status
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

// createEnrollmentStatusHandler handles the creation of a new enrollment status.
//
//	@Summary		Create a new enrollment status
//	@Description	Create a new enrollment status
//	@Tags			enrollment-statuses
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			enrollment_status	body		CreateEnrollmentStatusRequest_T	true	"Enrollment status data"
//	@Success		201					{object}	envelope
//	@Failure		400					{object}	errorResponse
//	@Failure		422					{object}	errorResponse
//	@Failure		500					{object}	errorResponse
//	@Router			/v1/enrollment/status [post]
func (app *appDependencies) createEnrollmentStatusHandler(w http.ResponseWriter, r *http.Request) {

	if err := app.readJSON(w, r, &CreateEnrollmentStatusRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	status := &data.EnrollmentStatus{Status: CreateEnrollmentStatusRequest.Status}

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
	headers.Set("Location", fmt.Sprintf("/v1/enrollment/status/%d", status.ID))

	if err := app.writeJSON(w, http.StatusCreated, envelope{"enrollment_status": status}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showEnrollmentStatusHandler retrieves and returns an enrollment status by its ID.
//
//	@Summary		Get an enrollment status
//	@Description	Retrieve an enrollment status by its ID
//	@Tags			enrollment-statuses
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Enrollment status ID"
//	@Success		200	{object}	envelope
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/v1/enrollment/status/{id} [get]
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

// listEnrollmentStatusesHandler returns a filtered list of enrollment statuses.
//
//	@Summary		List enrollment statuses
//	@Description	Retrieve a list of enrollment statuses with optional filtering by status name
//	@Tags			enrollment-statuses
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			status		query		string	false	"Filter by status name"
//	@Param			page		query		int		false	"Page number for pagination"
//	@Param			page_size	query		int		false	"Number of items per page"
//	@Param			sort		query		string	false	"Sort order"
//	@Success		200			{object}	envelope
//	@Failure		422			{object}	errorResponse
//	@Failure		500			{object}	errorResponse
//	@Router			/v1/enrollment/status [get]
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

// updateEnrollmentStatusHandler performs a partial update on an enrollment status record.
//
//	@Summary		Update an enrollment status
//	@Description	Perform a partial update on an enrollment status record
//	@Tags			enrollment-statuses
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id					path		int								true	"Enrollment status ID"
//	@Param			enrollment_status	body		UpdateEnrollmentStatusRequest_T	true	"Enrollment status data"
//	@Success		200					{object}	envelope
//	@Failure		400					{object}	errorResponse
//	@Failure		404					{object}	errorResponse
//	@Failure		422					{object}	errorResponse
//	@Failure		500					{object}	errorResponse
//	@Router			/v1/enrollment/status/{id} [patch]
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

	if err := app.readJSON(w, r, &UpdateEnrollmentStatusRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if UpdateEnrollmentStatusRequest.Status != nil {
		status.Status = *UpdateEnrollmentStatusRequest.Status
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

/*********************** Attendance Status ***********************/

// createAttendanceStatusHandler handles the creation of a new attendance status.
//
//	@Summary		Create a new attendance status
//	@Description	Create a new attendance status
//	@Tags			attendance-statuses
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			attendance_status	body		CreateAttendanceStatusRequest_T	true	"Attendance status data"
//	@Success		201					{object}	envelope
//	@Failure		400					{object}	errorResponse
//	@Failure		422					{object}	errorResponse
//	@Failure		500					{object}	errorResponse
//	@Router			/v1/attendance/status [post]
func (app *appDependencies) createAttendanceStatusHandler(w http.ResponseWriter, r *http.Request) {
	if err := app.readJSON(w, r, &CreateAttendanceStatusRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	status := &data.AttendanceStatus{
		Status:          CreateAttendanceStatusRequest.Status,
		CountsAsPresent: CreateAttendanceStatusRequest.CountsAsPresent,
	}

	v := validator.New()
	data.ValidateAttendanceStatus(v, status)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.AttendanceStatus.Insert(status); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("status", "an attendance status with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/attendance/status/%d", status.ID))

	if err := app.writeJSON(w, http.StatusCreated, envelope{"attendance_status": status}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showAttendanceStatusHandler retrieves and returns an attendance status by its ID.
//
//	@Summary		Get an attendance status
//	@Description	Retrieve an attendance status by its ID
//	@Tags			attendance-statuses
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Attendance status ID"
//	@Success		200	{object}	envelope
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/v1/attendance/status/{id} [get]
func (app *appDependencies) showAttendanceStatusHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	status, err := app.models.AttendanceStatus.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"attendance_status": status}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listAttendanceStatusesHandler returns a filtered list of attendance statuses.
//
//	@Summary		List attendance statuses
//	@Description	Retrieve a list of attendance statuses with optional filtering
//	@Tags			attendance-statuses
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			status				query		string	false	"Filter by status name"
//	@Param			counts_as_present	query		bool	false	"Filter by counts as present flag"
//	@Param			page				query		int		false	"Page number for pagination"
//	@Param			page_size			query		int		false	"Number of items per page"
//	@Param			sort				query		string	false	"Sort order"
//	@Success		200					{object}	envelope
//	@Failure		422					{object}	errorResponse
//	@Failure		500					{object}	errorResponse
//	@Router			/v1/attendance/status [get]
func (app *appDependencies) listAttendanceStatusesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	v := validator.New()

	filters := app.readFilters(query, "status", 20, []string{"status", "-status", "id", "-id", "created_at", "-created_at"}, v)

	countsAsPresent := app.getOptionalBoolQueryParameter(query, "counts_as_present", v)

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	name := likeSearch(app.getSingleQueryParameter(query, "status", ""))

	statuses, metadata, err := app.models.AttendanceStatus.GetAll(name, countsAsPresent, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"attendance_statuses": statuses, "metadata": metadata}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateAttendanceStatusHandler performs a partial update on an attendance status record.
//
//	@Summary		Update an attendance status
//	@Description	Perform a partial update on an attendance status record
//	@Tags			attendance-statuses
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id					path		int								true	"Attendance status ID"
//	@Param			attendance_status	body		UpdateAttendanceStatusRequest_T	true	"Attendance status data"
//	@Success		200					{object}	envelope
//	@Failure		400					{object}	errorResponse
//	@Failure		404					{object}	errorResponse
//	@Failure		422					{object}	errorResponse
//	@Failure		500					{object}	errorResponse
//	@Router			/v1/attendance/status/{id} [patch]
func (app *appDependencies) updateAttendanceStatusHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	status, err := app.models.AttendanceStatus.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.readJSON(w, r, &UpdateAttendanceStatusRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if UpdateAttendanceStatusRequest.Status != nil {
		status.Status = *UpdateAttendanceStatusRequest.Status
	}
	if UpdateAttendanceStatusRequest.CountsAsPresent != nil {
		status.CountsAsPresent = *UpdateAttendanceStatusRequest.CountsAsPresent
	}

	v := validator.New()
	data.ValidateAttendanceStatus(v, status)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.AttendanceStatus.Update(status); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("status", "an attendance status with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"attendance_status": status}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

/*********************** Progress Status ***********************/

// createProgressStatusHandler handles the creation of a new progress status.
//
//	@Summary		Create a new progress status
//	@Description	Create a new progress status
//	@Tags			progress-statuses
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			progress_status	body		CreateProgressStatusRequest_T	true	"Progress status data"
//	@Success		201				{object}	envelope
//	@Failure		400				{object}	errorResponse
//	@Failure		422				{object}	errorResponse
//	@Failure		500				{object}	errorResponse
//	@Router			/v1/progress/status [post]
func (app *appDependencies) createProgressStatusHandler(w http.ResponseWriter, r *http.Request) {
	if err := app.readJSON(w, r, &CreateProgressStatusRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	status := &data.ProgressStatus{Status: CreateProgressStatusRequest.Status}

	v := validator.New()
	data.ValidateProgressStatus(v, status)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.ProgressStatus.Insert(status); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("status", "a progress status with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/progress/status/%d", status.ID))

	if err := app.writeJSON(w, http.StatusCreated, envelope{"progress_status": status}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showProgressStatusHandler retrieves and returns a progress status by its ID.
//
//	@Summary		Get a progress status
//	@Description	Retrieve a progress status by its ID
//	@Tags			progress-statuses
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Progress status ID"
//	@Success		200	{object}	envelope
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/v1/progress/status/{id} [get]
func (app *appDependencies) showProgressStatusHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	status, err := app.models.ProgressStatus.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"progress_status": status}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listProgressStatusesHandler returns a filtered list of progress statuses.
//
//	@Summary		List progress statuses
//	@Description	Retrieve a list of progress statuses with optional filtering
//	@Tags			progress-statuses
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			status		query		string	false	"Filter by status name"
//	@Param			page		query		int		false	"Page number for pagination"
//	@Param			page_size	query		int		false	"Number of items per page"
//	@Param			sort		query		string	false	"Sort order"
//	@Success		200			{object}	envelope
//	@Failure		422			{object}	errorResponse
//	@Failure		500			{object}	errorResponse
//	@Router			/v1/progress/status [get]
func (app *appDependencies) listProgressStatusesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	v := validator.New()

	filters := app.readFilters(query, "status", 20, []string{"status", "-status", "id", "-id", "created_at", "-created_at"}, v)

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	name := likeSearch(app.getSingleQueryParameter(query, "status", ""))

	statuses, metadata, err := app.models.ProgressStatus.GetAll(name, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"progress_statuses": statuses, "metadata": metadata}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateProgressStatusHandler performs a partial update on a progress status record.
//
//	@Summary		Update a progress status
//	@Description	Perform a partial update on a progress status record
//	@Tags			progress-statuses
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id				path		int							true	"Progress status ID"
//	@Param			progress_status	body		UpdateProgressStatusRequest_T	true	"Progress status data"
//	@Success		200				{object}	envelope
//	@Failure		400				{object}	errorResponse
//	@Failure		404				{object}	errorResponse
//	@Failure		422				{object}	errorResponse
//	@Failure		500				{object}	errorResponse
//	@Router			/v1/progress/status/{id} [patch]
func (app *appDependencies) updateProgressStatusHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	status, err := app.models.ProgressStatus.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.readJSON(w, r, &UpdateProgressStatusRequest); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if UpdateProgressStatusRequest.Status != nil {
		status.Status = *UpdateProgressStatusRequest.Status
	}

	v := validator.New()
	data.ValidateProgressStatus(v, status)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.ProgressStatus.Update(status); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateValue):
			v.AddError("status", "a progress status with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"progress_status": status}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
