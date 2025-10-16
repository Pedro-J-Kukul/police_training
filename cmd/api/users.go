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
func (app *appDependencies) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		Gender      string `json:"gender"`
		Password    string `json:"password"`
		Facilitator bool   `json:"facilitator"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Email:       input.Email,
		Gender:      input.Gender,
		Activated:   false,
		Facilitator: input.Facilitator,
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
		app.serverErrorResponse(w, r, err)
		return
	}

	if app.mailer != nil {
		app.background(func() {
			data := map[string]any{
				"userID":          user.ID,
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

	user.Activated = true

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

	activated := app.getOptionalBoolQueryParameter(query, "activated", v)
	facilitator := app.getOptionalBoolQueryParameter(query, "facilitator", v)

	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	users, metadata, err := app.models.User.GetAll(
		app.getSingleQueryParameter(query, "first_name", ""),
		app.getSingleQueryParameter(query, "last_name", ""),
		app.getSingleQueryParameter(query, "email", ""),
		app.getSingleQueryParameter(query, "gender", ""),
		activated,
		facilitator,
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
		FirstName   *string `json:"first_name"`
		LastName    *string `json:"last_name"`
		Email       *string `json:"email"`
		Gender      *string `json:"gender"`
		Password    *string `json:"password"`
		Facilitator *bool   `json:"facilitator"`
		Activated   *bool   `json:"activated"`
		Version     int     `json:"version"`
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
	if input.Facilitator != nil {
		user.Facilitator = *input.Facilitator
	}
	if input.Activated != nil {
		user.Activated = *input.Activated
	}

	v := validator.New()
	data.ValidateUser(v, user)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.User.Update(user); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Delete the officer as well
	err = app.models.Officer.Delete(id)
	e := app.deleteUserHandlerErrorHandler(err, w, r)
	if e == true {
		return
	}

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Delete the user
	err = app.models.User.Delete(id)

	e = app.deleteUserHandlerErrorHandler(err, w, r)
	if e == true {
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"message": "user successfully deleted"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *appDependencies) deleteUserHandlerErrorHandler(err error, w http.ResponseWriter, r *http.Request) bool {
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return true
	}

	return false
}
