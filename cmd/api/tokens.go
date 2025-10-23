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
//
//	@Summary		Create Authentication Token
//	@Description	Generates an authentication token for a user based on provided email and password.
//	@Tags			Tokens
//	@Accept			json
//	@Produce		json
//	@Param			input	body		CreateAuthenticationTokenRequest_T	true	"User credentials"
//	@Success		201		{object}	envelope
//	@Failure		400		{object}	errorResponse
//	@Failure		422		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/v1/tokens/authentication [post]
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

	// Check if the user is activated and isn't trying to bypass the flow
	if !user.IsActivated {
		v.AddError("email", "account must be activated to login")
		app.failedValidationResponse(w, r, v.Errors)
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

// createPasswordResetTokenHandler sends a password reset token to the user's email.
//
//	@Summary		Create Password Reset Token
//	@Description	Generates a password reset token and sends it to the user's email address.
//	@Tags			Tokens
//	@Accept			json
//	@Produce		json
//	@Param			input	body		CreatePasswordResetTokenRequest_T	true	"User email"
//	@Success		202		{object}	envelope
//	@Failure		400		{object}	errorResponse
//	@Failure		422		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/v1/tokens/password-reset [post]
func (app *appDependencies) createPasswordResetTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateEmail(v, input.Email)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.User.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			if err := app.writeJSON(w, http.StatusAccepted, envelope{"message": "if that account exists, a password reset email has been sent"}, nil); err != nil {
				app.serverErrorResponse(w, r, err)
			}
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Check if the user is activated and isn't trying to bypass the flow
	if !user.IsActivated {
		v.AddError("email", "account must be activated to reset password")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	_ = app.models.Token.DeleteAllForUser(data.ScopePasswordReset, user.ID)

	token, err := app.models.Token.New(user.ID, 45*time.Minute, data.ScopePasswordReset)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if app.mailer != nil {
		app.background(func() {
			payload := map[string]any{
				"resetToken": token.Plaintext,
				"userID":     user.ID,
			}
			if err := app.mailer.Send(user.Email, "password_reset.tmpl", payload); err != nil {
				app.logger.Error("failed to send password reset email", "error", err)
			}
		})
	}

	if err := app.writeJSON(w, http.StatusAccepted, envelope{"message": "if that account exists, a password reset email has been sent"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// resetPasswordHandler updates a user's password using a reset token.
//
//	@Summary		Reset Password
//	@Description	Resets a user's password using a valid password reset token.
//	@Tags			Tokens
//	@Accept			json
//	@Produce		json
//	@Param			input	body		ResetPasswordRequest_T	true	"Password reset data"
//	@Success		200		{object}	envelope
//	@Failure		400		{object}	errorResponse
//	@Failure		422		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/v1/users/password-reset [put]
func (app *appDependencies) resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateTokenPlaintext(v, input.Token)
	data.ValidatePasswordPlaintext(v, input.Password)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.User.GetForToken(data.ScopePasswordReset, input.Token)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired password reset token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := user.Password.Set(input.Password); err != nil {
		app.serverErrorResponse(w, r, err)
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

	if err := app.models.Token.DeleteAllForUser(data.ScopePasswordReset, user.ID); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.models.Token.DeleteAllForUser(data.ScopeAuthentication, user.ID); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"message": "password updated"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
