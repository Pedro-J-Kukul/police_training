package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/Pedro-J-Kukul/police_training/internal/data"
	_ "github.com/lib/pq"
)

var testApp *appDependencies

func TestMain(m *testing.M) {
	// Try multiple environment variable names for flexibility
	dbDSN := os.Getenv("TEST_DATABASE_DSN")
	if dbDSN == "" {
		dbDSN = os.Getenv("TEST_DB_DSN")
	}
	if dbDSN == "" {
		// Fallback to match your .envrc
		dbDSN = "postgres://police:police@localhost/police_training_testing?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		panic(fmt.Sprintf("Failed to open database connection: %v", err))
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		panic(fmt.Sprintf("Could not connect to test database: %v\nUsing DSN: %s", err, dbDSN))
	}

	// Initialize test app with proper logger and configuration
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError, // Only show errors during testing
	}))

	testApp = &appDependencies{
		models: data.NewModels(db),
		logger: logger,
		config: serverConfig{
			env: "testing",
		},
	}

	code := m.Run()
	db.Close()
	os.Exit(code)
}

// Helper function to clean up test data for API tests
func cleanupUsersAPITestData(t *testing.T) {
	t.Helper()

	// Get the database from testApp
	db := testApp.models.User.DB

	// Clean up in reverse order of dependencies
	_, err := db.Exec("TRUNCATE TABLE tokens CASCADE")
	if err != nil {
		t.Logf("Warning: Failed to cleanup tokens: %v", err)
	}

	_, err = db.Exec("TRUNCATE TABLE roles_users CASCADE")
	if err != nil {
		t.Logf("Warning: Failed to cleanup roles_users: %v", err)
	}

	_, err = db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	if err != nil {
		t.Logf("Warning: Failed to cleanup users: %v", err)
	}
}

// Helper function to create a test user and return authentication token
func createTestUserWithToken(t *testing.T) (*data.User, string) {
	t.Helper()

	user := &data.User{
		FirstName:   "Test",
		LastName:    "User",
		Email:       fmt.Sprintf("test%d@example.com", time.Now().UnixNano()),
		Gender:      "m",
		IsActivated: true,
		IsOfficer:   true,
	}
	_ = user.Password.Set("TestPass123!")

	err := testApp.models.User.Insert(user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create authentication token
	token, err := testApp.models.Token.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	return user, token.Plaintext
}

// Helper function to create test request with authentication
func createAuthenticatedRequest(method, path string, body []byte, token string) *http.Request {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	return req
}

func TestRegisterUserHandler(t *testing.T) {
	// Clean up before test
	cleanupUsersAPITestData(t)
	defer cleanupUsersAPITestData(t)

	tests := []struct {
		name           string
		input          map[string]any
		expectedStatus int
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name: "Valid user registration",
			input: map[string]any{
				"first_name":     "John",
				"last_name":      "Doe",
				"email":          fmt.Sprintf("john%d@example.com", time.Now().UnixNano()),
				"gender":         "m",
				"password":       "StrongPass123!",
				"is_facilitator": false,
				"is_officer":     true,
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, res *http.Response) {
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if response["user"] == nil {
					t.Error("Expected user object in response")
				}
			},
		},
		{
			name: "Invalid email format",
			input: map[string]any{
				"first_name": "John",
				"last_name":  "Doe",
				"email":      "invalid-email",
				"gender":     "m",
				"password":   "StrongPass123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, res *http.Response) {
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if response["error"] == nil {
					t.Error("Expected error in response")
				}
			},
		},
		{
			name: "Missing required fields",
			input: map[string]any{
				"email": fmt.Sprintf("missing%d@example.com", time.Now().UnixNano()),
			},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse:  func(t *testing.T, res *http.Response) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			testApp.registerUserHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}

			tt.checkResponse(t, res)
		})
	}
}

func TestShowCurrentUserHandler(t *testing.T) {
	cleanupUsersAPITestData(t)
	defer cleanupUsersAPITestData(t)

	user, token := createTestUserWithToken(t)

	tests := []struct {
		name           string
		token          string
		expectedStatus int
		setupUser      bool
	}{
		{
			name:           "Valid authenticated request",
			token:          token,
			expectedStatus: http.StatusOK,
			setupUser:      true,
		},
		{
			name:           "Missing authentication token",
			token:          "",
			expectedStatus: http.StatusOK, // Anonymous user context
			setupUser:      false,
		},
		{
			name:           "Invalid authentication token",
			token:          "invalid-token",
			expectedStatus: http.StatusUnauthorized,
			setupUser:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := createAuthenticatedRequest(http.MethodGet, "/v1/me", nil, tt.token)

			// Set user context directly for testing
			if tt.setupUser && tt.token != "" {
				req = setUserContext(req, user)
			} else if tt.token == "" {
				req = setUserContext(req, data.AnonymousUser)
			}

			rec := httptest.NewRecorder()

			if tt.token == "" || tt.setupUser {
				// Call handler directly for anonymous user or with user context
				testApp.showCurrentUserHandler(rec, req)
			} else {
				// Test authentication middleware
				handler := testApp.authenticate(http.HandlerFunc(testApp.showCurrentUserHandler))
				handler.ServeHTTP(rec, req)
			}

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}

			if tt.expectedStatus == http.StatusOK && tt.setupUser {
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				userResponse := response["user"].(map[string]any)
				if userResponse["email"] != user.Email {
					t.Errorf("Expected email %s; got %s", user.Email, userResponse["email"])
				}
			}
		})
	}
}

func TestShowUserHandler(t *testing.T) {
	cleanupUsersAPITestData(t)
	defer cleanupUsersAPITestData(t)

	user, token := createTestUserWithToken(t)

	tests := []struct {
		name           string
		userID         string
		token          string
		expectedStatus int
	}{
		{
			name:           "Valid user ID with authentication",
			userID:         strconv.FormatInt(user.ID, 10),
			token:          token,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid user ID",
			userID:         "999999",
			token:          token,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Non-numeric user ID",
			userID:         "abc",
			token:          token,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := fmt.Sprintf("/v1/users/%s", tt.userID)
			req := createAuthenticatedRequest(http.MethodGet, path, nil, tt.token)

			// Set URL parameter and user context
			req = setURLParam(req, "id", tt.userID)
			req = setUserContext(req, user)

			rec := httptest.NewRecorder()
			testApp.showUserHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}
		})
	}
}

func TestDeleteUserHandler(t *testing.T) {
	cleanupUsersAPITestData(t)
	defer cleanupUsersAPITestData(t)

	user, token := createTestUserWithToken(t)

	tests := []struct {
		name           string
		userID         string
		token          string
		expectedStatus int
	}{
		{
			name:           "Valid user deletion",
			userID:         strconv.FormatInt(user.ID, 10),
			token:          token,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Non-existent user",
			userID:         "999999",
			token:          token,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := fmt.Sprintf("/v1/users/%s", tt.userID)
			req := createAuthenticatedRequest(http.MethodDelete, path, nil, tt.token)

			// Set URL parameter and user context
			req = setURLParam(req, "id", tt.userID)
			req = setUserContext(req, user)

			rec := httptest.NewRecorder()
			testApp.deleteUserHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}
		})
	}
}

// Integration test for complete user workflow
func TestUserWorkflow(t *testing.T) {
	cleanupUsersAPITestData(t)
	defer cleanupUsersAPITestData(t)

	// 1. Register a new user
	userInput := map[string]any{
		"first_name":     "Workflow",
		"last_name":      "Test",
		"email":          fmt.Sprintf("workflow%d@example.com", time.Now().UnixNano()),
		"gender":         "f",
		"password":       "WorkflowPass123!",
		"is_facilitator": false,
		"is_officer":     true,
	}

	body, _ := json.Marshal(userInput)
	req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	testApp.registerUserHandler(rec, req)

	if rec.Result().StatusCode != http.StatusCreated {
		t.Fatalf("User registration failed with status %d", rec.Result().StatusCode)
	}

	// Extract user from response
	var registerResponse map[string]any
	_ = json.NewDecoder(rec.Result().Body).Decode(&registerResponse)
	userResponse := registerResponse["user"].(map[string]any)
	userID := int64(userResponse["id"].(float64))

	// 2. Create activation token and activate user
	token, err := testApp.models.Token.New(userID, 72*time.Hour, data.ScopeActivation)
	if err != nil {
		t.Fatalf("Failed to create activation token: %v", err)
	}

	activationInput := map[string]string{"token": token.Plaintext}
	body, _ = json.Marshal(activationInput)
	req = httptest.NewRequest(http.MethodPut, "/v1/users/activate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec = httptest.NewRecorder()
	testApp.activateUserHandler(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("User activation failed with status %d", rec.Result().StatusCode)
	}

	t.Log("Complete user workflow test passed successfully!")
}
