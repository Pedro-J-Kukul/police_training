package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Pedro-J-Kukul/police_training/internal/data"
	"github.com/Pedro-J-Kukul/police_training/internal/validator"
)

// Mocks and helpers

func newTestApp() *appDependencies {
	// Setup a test appDependencies struct.
	// Replace with real test DB and config as needed.
	return &appDependencies{
		models: data.NewModels(nil), // nil DB for validation-only tests
	}
}

// --- Handler-Level Tests ---

func TestRegisterUserHandler_Validation(t *testing.T) {
	app := newTestApp()

	// Missing required fields
	payload := `{"first_name": "", "last_name": "", "email": "", "password": "", "gender": ""}`
	req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	app.registerUserHandler(w, req)

	if w.Code != http.StatusUnprocessableEntity && w.Code != http.StatusBadRequest {
		t.Fatalf("expected 422/400, got %d", w.Code)
	}
	// Optionally check error response body
}

// func TestRegisterUserHandler_Success(t *testing.T) {
// 	app := newTestApp()
// 	// You'd need a real or mock DB to actually persist

// 	payload := `{
// 		"first_name": "Jane",
// 		"last_name": "Doe",
// 		"email": "jane.doe@example.com",
// 		"password": "StrongPass123!",
// 		"gender": "F"
// 	}`
// 	req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewBufferString(payload))
// 	req.Header.Set("Content-Type", "application/json")
// 	w := httptest.NewRecorder()

// 	app.registerUserHandler(w, req)

// 	if w.Code != http.StatusCreated {
// 		t.Fatalf("expected 201 Created, got %d", w.Code)
// 	}
// }

func TestActivateUserHandler_BadToken(t *testing.T) {
	app := newTestApp()
	payload := `{"token": "invalidtoken"}`
	req := httptest.NewRequest(http.MethodPut, "/v1/users/activate", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	app.activateUserHandler(w, req)
	// Should fail with 400 or 422 if token is bad
	if w.Code != http.StatusUnprocessableEntity && w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400/422 for bad token, got %d", w.Code)
	}
}

func TestCreateAuthenticationTokenHandler_InvalidUser(t *testing.T) {
	app := newTestApp()
	payload := `{"email": "fake@notreal.com", "password": "wrongpass"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/tokens/authentication", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	app.createAuthenticationTokenHandler(w, req)
	if w.Code != http.StatusUnprocessableEntity && w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400/422 for wrong credentials, got %d", w.Code)
	}
}

func TestShowCurrentUserHandler_Unauthorized(t *testing.T) {
	app := newTestApp()
	req := httptest.NewRequest(http.MethodGet, "/v1/me", nil)
	w := httptest.NewRecorder()
	// No user in context

	app.showCurrentUserHandler(w, req)
	if w.Code != http.StatusUnauthorized && w.Code != http.StatusForbidden {
		t.Fatalf("expected 401/403, got %d", w.Code)
	}
}

func TestShowUserHandler_NotFound(t *testing.T) {
	app := newTestApp()
	req := httptest.NewRequest(http.MethodGet, "/v1/users/99999", nil)
	w := httptest.NewRecorder()
	// Simulate user not found

	app.showUserHandler(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateUserHandler_InvalidPayload(t *testing.T) {
	app := newTestApp()
	req := httptest.NewRequest(http.MethodPatch, "/v1/users/1", bytes.NewBufferString(`{bad}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	app.updateUserHandler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", w.Code)
	}
}

func TestDeleteUserHandler_NotFound(t *testing.T) {
	app := newTestApp()
	req := httptest.NewRequest(http.MethodDelete, "/v1/users/99999", nil)
	w := httptest.NewRecorder()

	app.deleteUserHandler(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// --- Add more tests for success cases with a test DB or a mock ---

// Table-driven example for validation:
func TestValidateUser_EmailFormats(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"john@example.com", true},
		{"bad-email", false},
		{"", false},
	}

	for _, tc := range tests {
		user := &data.User{Email: tc.email}
		v := validator.New()
		data.ValidateUser(v, user)
		if (len(v.Errors) == 0) != tc.valid {
			t.Errorf("email %q, expected valid=%v, got errors=%v", tc.email, tc.valid, v.Errors)
		}
	}
}

func TestRegisterUserHandler_Success(t *testing.T) {
	app := newTestApp()
	payload := `{"first_name":"Test","last_name":"User","email":"testuser@example.com","password":"StrongPass123!","gender":"M"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	app.registerUserHandler(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", w.Code)
	}
}
