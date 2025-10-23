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
	// Clean up before test
	cleanupAPITestData(t)
	defer cleanupAPITestData(t)

	// Create a test user
	user := &data.User{
		FirstName:   "Auth",
		LastName:    "Test",
		Email:       fmt.Sprintf("auth%d@example.com", time.Now().UnixNano()),
		Gender:      "m",
		IsActivated: true,
		IsOfficer:   true,
	}
	_ = user.Password.Set("AuthPass123!")
	err := testApp.models.User.Insert(user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name           string
		input          map[string]string
		expectedStatus int
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name: "Valid credentials",
			input: map[string]string{
				"email":    user.Email,
				"password": "AuthPass123!",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, res *http.Response) {
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
			},
		},
		{
			name: "Invalid password",
			input: map[string]string{
				"email":    user.Email,
				"password": "WrongPassword", // This might be failing password validation rules
			},
			expectedStatus: http.StatusUnauthorized, // Change this based on your validation rules
			checkResponse: func(t *testing.T, res *http.Response) {
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if response["error"] == nil {
					t.Error("Expected error message in response")
				}
			},
		},
		{
			name: "Non-existent email",
			input: map[string]string{
				"email":    "nonexistent@example.com",
				"password": "SomePassword123!",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse:  func(t *testing.T, res *http.Response) {},
		},
		{
			name: "Invalid email format",
			input: map[string]string{
				"email":    "invalid-email",
				"password": "SomePassword123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse:  func(t *testing.T, res *http.Response) {},
		},
		{
			name: "Missing password",
			input: map[string]string{
				"email": user.Email,
			},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse:  func(t *testing.T, res *http.Response) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/v1/tokens/authentication", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			testApp.createAuthenticationTokenHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}

			tt.checkResponse(t, res)
		})
	}
}
