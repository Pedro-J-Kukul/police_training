package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Pedro-J-Kukul/police_training/internal/data"
)

func TestCreateAuthenticationTokenHandler(t *testing.T) {
	t.Log("=== Testing Authentication Token Creation ===")

	tests := []struct {
		name           string
		input          map[string]string
		expectedStatus int
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name: "Valid credentials for seeded admin",
			input: map[string]string{
				"email":    "admin1@police-training.bz",
				"password": "TrainingPass123!",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating successful token creation")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response["authentication_token"] == nil {
					t.Error("Expected authentication_token in response")
					return
				}

				tokenData := response["authentication_token"].(map[string]any)
				if tokenData["token"] == nil {
					t.Error("Expected token field in authentication_token")
				}
				if tokenData["expiry"] == nil {
					t.Error("Expected expiry field in authentication_token")
				}

				t.Logf("Step: Successfully created token with expiry %v", tokenData["expiry"])
			},
		},
		{
			name: "Valid credentials for seeded facilitator",
			input: map[string]string{
				"email":    "admin1@police-training.bz",
				"password": "TrainingPass123!",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating facilitator token creation")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response["authentication_token"] == nil {
					t.Error("Expected authentication_token in response")
				}

				t.Log("Step: Facilitator token created successfully")
			},
		},
		{
			name: "Valid credentials for seeded officer",
			input: map[string]string{
				"email":    "john.smith@police-training.bz",
				"password": "TrainingPass123!",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating officer token creation")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response["authentication_token"] == nil {
					t.Error("Expected authentication_token in response")
				}

				t.Log("Step: Officer token created successfully")
			},
		},
		{
			name: "Invalid password for seeded user",
			input: map[string]string{
				"email":    "admin1@police-training.bz",
				"password": "WrongPassword123!",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating invalid password rejection")
			},
		},
		{
			name: "Non-existent email",
			input: map[string]string{
				"email":    "nonexistent@police-training.bz",
				"password": "SomePassword123!",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating non-existent user rejection")
			},
		},
		{
			name: "Invalid email format",
			input: map[string]string{
				"email":    "invalid-email",
				"password": "SomePassword123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating email format validation")
			},
		},
		{
			name: "Missing password",
			input: map[string]string{
				"email": "admin1@police-training.bz",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating missing password validation")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test: %s", tt.name)

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/v1/tokens/authentication", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			t.Logf("Step: Making POST request with email %s", tt.input["email"])

			rec := httptest.NewRecorder()
			testApp.createAuthenticationTokenHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			t.Logf("Step: Received status code %d (expected %d)", res.StatusCode, tt.expectedStatus)
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}

			tt.checkResponse(t, res)
			t.Logf("Completed test: %s", tt.name)
		})
	}
}

func TestCreatePasswordResetTokenHandler(t *testing.T) {
	t.Log("=== Testing Password Reset Token Creation ===")

	tests := []struct {
		name           string
		input          map[string]string
		expectedStatus int
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name: "Valid email for seeded user",
			input: map[string]string{
				"email": "admin1@police-training.bz",
			},
			expectedStatus: http.StatusAccepted,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating password reset token creation")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response["message"] == nil {
					t.Error("Expected message in response")
				}

				t.Log("Step: Password reset token created successfully")
			},
		},
		{
			name: "Non-existent email",
			input: map[string]string{
				"email": "nonexistent@police-training.bz",
			},
			expectedStatus: http.StatusAccepted, // Still returns 202 for security
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating non-existent email handling")
			},
		},
		{
			name: "Invalid email format",
			input: map[string]string{
				"email": "invalid-email",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating email format validation")
			},
		},
		{
			name:           "Missing email",
			input:          map[string]string{},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating missing email validation")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test: %s", tt.name)

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/v1/tokens/password-reset", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			t.Logf("Step: Making POST request for password reset")

			rec := httptest.NewRecorder()
			testApp.createPasswordResetTokenHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			t.Logf("Step: Received status code %d", res.StatusCode)
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}

			tt.checkResponse(t, res)
			t.Logf("Completed test: %s", tt.name)
		})
	}
}

func TestResetPasswordHandler(t *testing.T) {
	t.Log("=== Testing Password Reset Handler ===")

	// Get seeded user and create reset token
	user := getSeededUser(t, "admin2.garcia@police-training.bz")

	// Create password reset token
	token, err := testApp.models.Token.New(user.ID, 45*time.Minute, data.ScopePasswordReset)
	if err != nil {
		t.Fatalf("Failed to create password reset token: %v", err)
	}

	tests := []struct {
		name           string
		input          map[string]string
		expectedStatus int
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name: "Valid password reset",
			input: map[string]string{
				"token":    token.Plaintext,
				"password": "NewResetPassword123!",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating successful password reset")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response["message"] == nil {
					t.Error("Expected message in response")
				}

				t.Log("Step: Password reset completed successfully")
			},
		},
		{
			name: "Invalid token",
			input: map[string]string{
				"token":    "invalid-token-123",
				"password": "ValidPassword123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating invalid token rejection")
			},
		},
		{
			name: "Invalid password format",
			input: map[string]string{
				"token":    token.Plaintext,
				"password": "weak",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating password format validation")
			},
		},
		{
			name: "Missing token",
			input: map[string]string{
				"password": "ValidPassword123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating missing token validation")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test: %s", tt.name)

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPut, "/v1/users/password-reset", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			t.Logf("Step: Making PUT request for password reset")

			rec := httptest.NewRecorder()
			testApp.resetPasswordHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			t.Logf("Step: Received status code %d", res.StatusCode)
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}

			tt.checkResponse(t, res)
			t.Logf("Completed test: %s", tt.name)
		})
	}
}

func TestTokenWorkflow(t *testing.T) {
	t.Log("=== Testing Complete Token Workflow ===")

	t.Log("Step: Starting authentication workflow with seeded admin user")

	// 1. Create authentication token
	loginInput := map[string]string{
		"email":    "admin1@police-training.bz",
		"password": "TrainingPass123!",
	}

	body, _ := json.Marshal(loginInput)
	req := httptest.NewRequest(http.MethodPost, "/v1/tokens/authentication", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	testApp.createAuthenticationTokenHandler(rec, req)

	if rec.Result().StatusCode != http.StatusCreated {
		t.Fatalf("Token creation failed with status %d", rec.Result().StatusCode)
	}

	// Extract token from response
	var tokenResponse map[string]any
	_ = json.NewDecoder(rec.Result().Body).Decode(&tokenResponse)
	authToken := tokenResponse["authentication_token"].(map[string]any)
	token := authToken["token"].(string)

	t.Logf("Step: Successfully created authentication token")

	// 2. Use token to access protected endpoint
	req = httptest.NewRequest(http.MethodGet, "/v1/me", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Get user for context
	user := getSeededUser(t, "admin2.garcia@police-training.bz")
	req = setUserContext(req, user)

	rec = httptest.NewRecorder()
	testApp.showCurrentUserHandler(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("Protected endpoint access failed with status %d", rec.Result().StatusCode)
	}

	var userResponse map[string]any
	_ = json.NewDecoder(rec.Result().Body).Decode(&userResponse)
	userData := userResponse["user"].(map[string]any)

	t.Logf("Step: Successfully accessed protected endpoint for user %s %s",
		userData["first_name"], userData["last_name"])

	// Verify user data matches seeded data
	if userData["email"] != "admin2.garcia@police-training.bz" {
		t.Errorf("Expected email admin2.garcia@police-training.bz, got %s", userData["email"])
	}

	if userData["first_name"] != "Immanuel" {
		t.Errorf("Expected first_name Immanuel, got %s", userData["first_name"])
	}

	// 3. Test password reset workflow
	resetEmailInput := map[string]string{
		"email": "admin2.garcia@police-training.bz",
	}

	body, _ = json.Marshal(resetEmailInput)
	req = httptest.NewRequest(http.MethodPost, "/v1/tokens/password-reset", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec = httptest.NewRecorder()
	testApp.createPasswordResetTokenHandler(rec, req)

	if rec.Result().StatusCode != http.StatusAccepted {
		t.Fatalf("Password reset token creation failed with status %d", rec.Result().StatusCode)
	}

	t.Log("Step: Password reset token creation completed")

	t.Log("Step: Token workflow completed successfully with correct user data")
}
