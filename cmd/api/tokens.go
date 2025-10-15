// FileName: cmd/api/tokens.go
package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/Pedro-J-Kukul/police_training/internal/data"
	"github.com/Pedro-J-Kukul/police_training/internal/validator"
)

// createAuthenticationTokenHandler handles the creation of authentication tokens.
func (app *appDependencies) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	// incomingData struct to hold the incoming JSON data
	var incomingData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Read and parse the JSON request body
	if err := app.readJSON(w, r, &incomingData); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate the incoming data
	v := validator.New() // declare a new validator instance

	data.ValidateEmail(v, incomingData.Email)                // validate the email field
	data.ValidatePasswordPlaintext(v, incomingData.Password) // validate the password field

	// If any validation errors, send a 422 Unprocessable Entity response
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Authenticate the user
	user, err := app.models.User.GetByEmail(incomingData.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r) // invalid credentials
		default:
			app.serverErrorResponse(w, r, err) // other errors
		}
		return
	}

	// Check if the provided password matches the stored password
	match, err := user.Password.Matches(incomingData.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// If the passwords don't match, send an invalid credentials response
	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	// Generate and return the authentication token
	token, err := app.models.Token.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{"authentication_token": token}       // wrap the token in an envelope
	err = app.writeJSON(w, http.StatusCreated, data, nil) // send a 201 Created response
	if err != nil {
		app.serverErrorResponse(w, r, err) // handle any errors
	}
}
