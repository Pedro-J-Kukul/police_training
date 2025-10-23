package main

import (
	"context"
	"net/http"

	"github.com/Pedro-J-Kukul/police_training/internal/data"
	"github.com/julienschmidt/httprouter"
)

// Test helper to simulate httprouter parameter extraction
func setURLParam(r *http.Request, key, value string) *http.Request {
	ctx := r.Context()
	params := httprouter.Params{httprouter.Param{Key: key, Value: value}}
	ctx = context.WithValue(ctx, httprouter.ParamsKey, params)
	return r.WithContext(ctx)
}

// Test helper to simulate user context
func setUserContext(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), contextKeyUser, user)
	return r.WithContext(ctx)
}
