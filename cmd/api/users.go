package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Pedro-J-Kukul/police_training/internal/data"
	"github.com/Pedro-J-Kukul/police_training/internal/validator"
)

// registerUserHandler creates a new user account and sends an activation email.
//
//	@Summary		Register a new user
//	@Description	Create a new user account and send an activation email
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		registerUserRequest	true	"User registration data"
//	@Success		201		{object}	envelope
//	@Failure		400		{object}	errorResponse
//	@Failure		422		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/v1/users [post]
func (app *appDependencies) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName     string `json:"first_name"`
		LastName      string `json:"last_name"`
		Email         string `json:"email"`
		Gender        string `json:"gender"`
		Password      string `json:"password"`
		IsFacilitator bool   `json:"is_facilitator"`
		IsOfficer     bool   `json:"is_officer"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		FirstName:     input.FirstName,
		LastName:      input.LastName,
		Email:         input.Email,
		Gender:        input.Gender,
		IsActivated:   false,
		IsFacilitator: input.IsFacilitator,
		IsOfficer:     input.IsOfficer,
	}

	if err := user.Password.Set(input.Password); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateUser(v, user)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.User.Insert(user); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Always clear existing activation tokens and send a new one.
	_ = app.models.Token.DeleteAllForUser(data.ScopeActivation, user.ID)
	activationToken, err := app.models.Token.New(user.ID, 72*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err) // Log the error and return a server error response
		return
	}

	if app.mailer != nil {
		app.background(func() {
			data := map[string]any{
				"userID":          user.ID,
				"firstName":       user.FirstName,
				"lastName":        user.LastName,
				"email":           user.Email,
				"password":        input.Password,
				"activationToken": activationToken.Plaintext,
			}
			if err := app.mailer.Send(user.Email, "user_welcome.tmpl", data); err != nil {
				app.logger.Error("failed to send activation email", "error", err)
			}
		})
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/users/%d", user.ID))

	if err := app.writeJSON(w, http.StatusCreated, envelope{"user": user}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// activateUserHandler uses an activation token to activate a pending account.
//
//	@Summary		Activate a user account
//	@Description	Activate a user account using an activation token
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			token	body		activateUserRequest	true	"Activation token"
//	@Success		200		{object}	envelope
//	@Failure		400		{object}	errorResponse
//	@Failure		422		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/v1/users/activate [put]
func (app *appDependencies) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Token string `json:"token"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateTokenPlaintext(v, input.Token)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.User.GetForToken(data.ScopeActivation, input.Token)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user.IsActivated = true

	if err := app.models.User.Update(user); err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.models.Token.DeleteAllForUser(data.ScopeActivation, user.ID); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showCurrentUserHandler returns the user record associated with the request context.
//
//	@Summary		Get current user
//	@Description	Retrieve the user record associated with the current request context
//	@Tags			users
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	envelope
//	@Failure		500	{object}	errorResponse
//	@Router			/v1/me [get]
func (app *appDependencies) showCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if err := app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showUserHandler returns a user by id.
func (app *appDependencies) showUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user, err := app.models.User.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listUsersHandler returns a filtered list of users.
//
//	@Summary		List users
//	@Description	Retrieve a list of users with optional filters and pagination
//	@Tags			users
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	envelope
//	@Failure		500	{object}	errorResponse
//	@Router			/v1/users [get]
func (app *appDependencies) listUsersHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	v := validator.New()

	filters := data.Filters{
		Page:     app.getSingleIntQueryParameter(query, "page", 1, v),
		PageSize: app.getSingleIntQueryParameter(query, "page_size", 20, v),
		Sort:     app.getSingleQueryParameter(query, "sort", "last_name"),
		SortSafelist: []string{
			"id", "first_name", "last_name", "email", "created_at",
			"-id", "-first_name", "-last_name", "-email", "-created_at",
		},
	}

	data.ValidateFilters(v, filters)

	isActivated := app.getOptionalBoolQueryParameter(query, "is_activated", v)
	isFacilitator := app.getOptionalBoolQueryParameter(query, "is_facilitator", v)
	isOfficer := app.getOptionalBoolQueryParameter(query, "is_officer", v)
	isDeleted := app.getOptionalBoolQueryParameter(query, "is_deleted", v)

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	users, metadata, err := app.models.User.GetAll(
		app.getSingleQueryParameter(query, "first_name", ""),
		app.getSingleQueryParameter(query, "last_name", ""),
		app.getSingleQueryParameter(query, "email", ""),
		app.getSingleQueryParameter(query, "gender", ""),
		isActivated, isFacilitator, isOfficer, isDeleted,
		filters,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	payload := envelope{
		"users":    users,
		"metadata": metadata,
	}

	if err := app.writeJSON(w, http.StatusOK, payload, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateUserHandler performs a partial update on a user record.
//
//	@Summary		Update a user
//	@Description	Perform a partial update on a user record
//	@Tags			users
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	envelope
//	@Failure		400	{object}	errorResponse
//	@Failure		422	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/v1/users/{id} [patch]
func (app *appDependencies) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user, err := app.models.User.Get(id)
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
		FirstName     *string `json:"first_name"`
		LastName      *string `json:"last_name"`
		Email         *string `json:"email"`
		Gender        *string `json:"gender"`
		Password      *string `json:"password"`
		IsFacilitator *bool   `json:"is_facilitator"`
		IsActivated   *bool   `json:"is_activated"`
		IsOfficer     *bool   `json:"is_officer"`
		Version       int     `json:"version"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Version == 0 || input.Version != user.Version {
		app.editConflictResponse(w, r)
		return
	}

	if input.FirstName != nil {
		user.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		user.LastName = *input.LastName
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Gender != nil {
		user.Gender = *input.Gender
	}
	if input.Password != nil {
		if err := user.Password.Set(*input.Password); err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	if input.IsFacilitator != nil {
		user.IsFacilitator = *input.IsFacilitator
	}
	if input.IsActivated != nil {
		user.IsActivated = *input.IsActivated
	}
	if input.IsOfficer != nil {
		user.IsOfficer = *input.IsOfficer
	}

	v := validator.New()
	data.ValidateUser(v, user)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.User.Update(user); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	payload := envelope{"user": user}

	if err := app.writeJSON(w, http.StatusOK, payload, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updatePasswordHandler updates a user's password.
func (app *appDependencies) updatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user, err := app.models.User.Get(id)
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
		Password string `json:"password"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := user.Password.Set(input.Password); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateUser(v, user)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.User.Update(user); err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	payload := envelope{"user": user}

	if err := app.writeJSON(w, http.StatusOK, payload, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteUserHandler deletes a user record.
//
//	@Summary		Delete a user
//	@Description	Delete a user record by ID
//	@Tags			users
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	envelope{message=string}
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/v1/users/{id} [delete]
func (app *appDependencies) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.User.SoftDelete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "user successfully deactivated"}, nil)
}

// restoreUserHandler restores a soft-deleted user record.
func (app *appDependencies) restoreUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.User.Restore(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "user successfully restored"}, nil)
}

// Permanently delete a user record.
func (app *appDependencies) hardDeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.User.HardDelete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "user successfully deleted"}, nil)
}
