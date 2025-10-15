// FileName: internal/data/context.go
package main

import (
	"context"
	"net/http"

	"github.com/Pedro-J-Kukul/police_training/internal/data"
)

type contextKey string // Define a custom type for context keys to avoid collisions

const contextKeyUser = contextKey("user") // Key for storing/retrieving user information in/from context

// contextSetUser adds the user information to the request context.
func (app *appDependencies) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), contextKeyUser, user) // Add user to context
	return r.WithContext(ctx)                                   // Return a new request with the updated context
}

// contextGetUser retrieves the user information from the request context.
func (app *appDependencies) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(contextKeyUser).(*data.User) // Retrieve user from context
	if !ok {
		panic("missing user value in context") // Panic if user is not found in context
	}
	return user // Return the retrieved user
}
