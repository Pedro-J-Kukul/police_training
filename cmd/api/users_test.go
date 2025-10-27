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
	dbDSN := "postgres://police:police@localhost/police_training_testing?sslmode=disable"

	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		panic(fmt.Sprintf("Failed to open database connection: %v", err))
	}

	if err = db.Ping(); err != nil {
		panic(fmt.Sprintf("Could not connect to test database: %v\nUsing DSN: %s", err, dbDSN))
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
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

func getSeededUser(t *testing.T, email string) *data.User {
	t.Helper()
	t.Logf("Step: Fetching seeded user with email: %s", email)

	user, err := testApp.models.User.GetByEmail(email)
	if err != nil {
		t.Fatalf("Failed to get seeded user %s: %v", email, err)
	}

	t.Logf("Step: Successfully retrieved user ID %d (%s %s)", user.ID, user.FirstName, user.LastName)
	return user
}

func createTokenForSeededUser(t *testing.T, userID int64) string {
	t.Helper()
	t.Logf("Step: Creating authentication token for user ID %d", userID)

	token, err := testApp.models.Token.New(userID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		t.Fatalf("Failed to create token for user %d: %v", userID, err)
	}

	t.Logf("Step: Successfully created token for user ID %d", userID)
	return token.Plaintext
}

func TestRegisterUserHandler(t *testing.T) {
	t.Log("=== Testing User Registration Handler ===")

	tests := []struct {
		name           string
		input          map[string]any
		expectedStatus int
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name: "Valid new user registration",
			input: map[string]any{
				"first_name":     "TestUser",
				"last_name":      "Registration",
				"email":          fmt.Sprintf("newuser%d@test-police-training.bz", time.Now().UnixNano()),
				"gender":         "m",
				"password":       "NewUserPass123!",
				"is_facilitator": false,
				"is_officer":     true,
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating successful registration response")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response["user"] == nil {
					t.Error("Expected user object in response")
					return
				}

				userResponse := response["user"].(map[string]any)
				t.Logf("Step: Created user ID %v with email %v", userResponse["id"], userResponse["email"])

				// Cleanup created user
				if userID, ok := userResponse["id"].(float64); ok {
					testApp.models.User.HardDelete(int64(userID))
					t.Logf("Step: Cleaned up test user ID %d", int64(userID))
				}
			},
		},
		{
			name: "Duplicate email with seeded user",
			input: map[string]any{
				"first_name": "Duplicate",
				"last_name":  "User",
				"email":      "admin1@police-training.bz",
				"gender":     "m",
				"password":   "DuplicatePass123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating duplicate email rejection")
			},
		},
		{
			name: "Invalid email format",
			input: map[string]any{
				"first_name": "Invalid",
				"last_name":  "Email",
				"email":      "not-an-email",
				"gender":     "f",
				"password":   "ValidPass123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating email format validation")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test: %s", tt.name)

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			testApp.registerUserHandler(rec, req)

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

func TestActivateUserHandler(t *testing.T) {
	t.Log("=== Testing Activate User Handler ===")

	// Create test user first
	testEmail := fmt.Sprintf("activate%d@test-police-training.bz", time.Now().UnixNano())
	userInput := map[string]any{
		"first_name": "Activate",
		"last_name":  "Test",
		"email":      testEmail,
		"gender":     "f",
		"password":   "ActivatePass123!",
		"is_officer": true,
	}

	body, _ := json.Marshal(userInput)
	req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	testApp.registerUserHandler(rec, req)

	var createResponse map[string]any
	_ = json.NewDecoder(rec.Result().Body).Decode(&createResponse)
	userResponse := createResponse["user"].(map[string]any)
	userID := int64(userResponse["id"].(float64))

	// Get activation token
	token, err := testApp.models.Token.New(userID, 72*time.Hour, data.ScopeActivation)
	if err != nil {
		t.Fatalf("Failed to create activation token: %v", err)
	}

	defer testApp.models.User.HardDelete(userID) // Cleanup

	tests := []struct {
		name           string
		input          map[string]string
		expectedStatus int
	}{
		{
			name: "Valid activation token",
			input: map[string]string{
				"token": token.Plaintext,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid activation token",
			input: map[string]string{
				"token": "invalid-token-123",
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Missing token",
			input:          map[string]string{},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test: %s", tt.name)

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPut, "/v1/users/activate", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			testApp.activateUserHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			t.Logf("Step: Received status code %d", res.StatusCode)
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}

			t.Logf("Completed test: %s", tt.name)
		})
	}
}

func TestShowCurrentUserHandler(t *testing.T) {
	t.Log("=== Testing Show Current User Handler ===")

	adminUser := getSeededUser(t, "admin1@police-training.bz")
	adminToken := createTokenForSeededUser(t, adminUser.ID)

	tests := []struct {
		name           string
		user           *data.User
		token          string
		expectedStatus int
		setupUser      bool
	}{
		{
			name:           "Authenticated seeded admin user",
			user:           adminUser,
			token:          adminToken,
			expectedStatus: http.StatusOK,
			setupUser:      true,
		},
		{
			name:           "Anonymous user context",
			user:           data.AnonymousUser,
			token:          "",
			expectedStatus: http.StatusOK,
			setupUser:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test: %s", tt.name)

			req := httptest.NewRequest(http.MethodGet, "/v1/me", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))
				t.Logf("Step: Set authorization header with token")
			}

			req = setUserContext(req, tt.user)
			t.Logf("Step: Set user context for user ID %d", tt.user.ID)

			rec := httptest.NewRecorder()
			testApp.showCurrentUserHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			t.Logf("Step: Received status code %d", res.StatusCode)
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
				t.Logf("Step: Response contains user email %v", userResponse["email"])

				if userResponse["email"] != tt.user.Email {
					t.Errorf("Expected email %s; got %s", tt.user.Email, userResponse["email"])
				}
			}

			t.Logf("Completed test: %s", tt.name)
		})
	}
}

func TestShowUserHandler(t *testing.T) {
	t.Log("=== Testing Show User Handler ===")

	facilitatorUser := getSeededUser(t, "maria.rodriguez@police-training.bz")
	facilitatorToken := createTokenForSeededUser(t, facilitatorUser.ID)

	tests := []struct {
		name           string
		userID         string
		expectedStatus int
	}{
		{
			name:           "Valid seeded user ID",
			userID:         strconv.FormatInt(facilitatorUser.ID, 10),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Non-existent user ID",
			userID:         "999999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid user ID format",
			userID:         "abc",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test: %s with user ID %s", tt.name, tt.userID)

			path := fmt.Sprintf("/v1/users/%s", tt.userID)
			req := httptest.NewRequest(http.MethodGet, path, nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", facilitatorToken))

			req = setURLParam(req, "id", tt.userID)
			req = setUserContext(req, facilitatorUser)

			t.Logf("Step: Set URL param id=%s and user context", tt.userID)

			rec := httptest.NewRecorder()
			testApp.showUserHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			t.Logf("Step: Received status code %d", res.StatusCode)
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				userResponse := response["user"].(map[string]any)
				t.Logf("Step: Successfully retrieved user %v (%v %v)",
					userResponse["id"], userResponse["first_name"], userResponse["last_name"])
			}

			t.Logf("Completed test: %s", tt.name)
		})
	}
}

func TestListUsersHandler(t *testing.T) {
	t.Log("=== Testing List Users Handler ===")

	adminUser := getSeededUser(t, "admin1@police-training.bz")
	adminToken := createTokenForSeededUser(t, adminUser.ID)

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
	}{
		{
			name:           "List all users with default pagination",
			queryParams:    "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "List users with pagination",
			queryParams:    "?page=1&page_size=3",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Filter users by officer status",
			queryParams:    "?is_officer=true",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Filter users by facilitator status",
			queryParams:    "?is_facilitator=true",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Search by first name",
			queryParams:    "?first_name=Pedro",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test: %s", tt.name)

			path := "/v1/users" + tt.queryParams
			req := httptest.NewRequest(http.MethodGet, path, nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))
			req = setUserContext(req, adminUser)

			t.Logf("Step: Making request to %s", path)

			rec := httptest.NewRecorder()
			testApp.listUsersHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			t.Logf("Step: Received status code %d", res.StatusCode)
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				users := response["users"].([]any)
				t.Logf("Step: Retrieved %d users", len(users))

				if response["metadata"] == nil {
					t.Error("Expected metadata in response")
				}
			}

			t.Logf("Completed test: %s", tt.name)
		})
	}
}

func TestUpdateUserHandler(t *testing.T) {
	t.Log("=== Testing Update User Handler ===")

	// Create test user for updating
	testEmail := fmt.Sprintf("update%d@test-police-training.bz", time.Now().UnixNano())
	userInput := map[string]any{
		"first_name": "Update",
		"last_name":  "Test",
		"email":      testEmail,
		"gender":     "m",
		"password":   "UpdatePass123!",
		"is_officer": true,
	}

	body, _ := json.Marshal(userInput)
	req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	testApp.registerUserHandler(rec, req)

	var createResponse map[string]any
	_ = json.NewDecoder(rec.Result().Body).Decode(&createResponse)
	userResponse := createResponse["user"].(map[string]any)
	userID := int64(userResponse["id"].(float64))
	version := int(userResponse["version"].(float64))

	defer testApp.models.User.HardDelete(userID) // Cleanup

	adminUser := getSeededUser(t, "admin1@police-training.bz")
	adminToken := createTokenForSeededUser(t, adminUser.ID)

	tests := []struct {
		name           string
		userID         string
		input          map[string]any
		expectedStatus int
	}{
		{
			name:   "Valid user update",
			userID: strconv.FormatInt(userID, 10),
			input: map[string]any{
				"first_name": "Updated",
				"last_name":  "Name",
				"version":    version,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Version conflict",
			userID: strconv.FormatInt(userID, 10),
			input: map[string]any{
				"first_name": "Conflict",
				"version":    999,
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:   "Non-existent user",
			userID: "999999",
			input: map[string]any{
				"first_name": "NotFound",
				"version":    1,
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test: %s", tt.name)

			body, _ := json.Marshal(tt.input)
			path := fmt.Sprintf("/v1/users/%s", tt.userID)
			req := httptest.NewRequest(http.MethodPatch, path, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))

			req = setURLParam(req, "id", tt.userID)
			req = setUserContext(req, adminUser)

			rec := httptest.NewRecorder()
			testApp.updateUserHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			t.Logf("Step: Received status code %d", res.StatusCode)
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}

			t.Logf("Completed test: %s", tt.name)
		})
	}
}

func TestUpdatePasswordHandler(t *testing.T) {
	t.Log("=== Testing Update Password Handler ===")

	facilitatorUser := getSeededUser(t, "maria.rodriguez@police-training.bz")
	facilitatorToken := createTokenForSeededUser(t, facilitatorUser.ID)

	tests := []struct {
		name           string
		userID         string
		input          map[string]string
		expectedStatus int
	}{
		{
			name:   "Valid password update",
			userID: strconv.FormatInt(facilitatorUser.ID, 10),
			input: map[string]string{
				"password": "NewPassword123!",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Invalid password format",
			userID: strconv.FormatInt(facilitatorUser.ID, 10),
			input: map[string]string{
				"password": "weak",
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:   "Non-existent user",
			userID: "999999",
			input: map[string]string{
				"password": "ValidPassword123!",
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test: %s", tt.name)

			body, _ := json.Marshal(tt.input)
			path := fmt.Sprintf("/v1/users/%s/password", tt.userID)
			req := httptest.NewRequest(http.MethodPatch, path, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", facilitatorToken))

			req = setURLParam(req, "id", tt.userID)
			req = setUserContext(req, facilitatorUser)

			rec := httptest.NewRecorder()
			testApp.updatePasswordHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			t.Logf("Step: Received status code %d", res.StatusCode)
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}

			t.Logf("Completed test: %s", tt.name)
		})
	}
}

func TestDeleteUserHandler(t *testing.T) {
	t.Log("=== Testing Delete User Handler (Soft Delete) ===")

	// Create test user for deletion
	testEmail := fmt.Sprintf("delete%d@test-police-training.bz", time.Now().UnixNano())
	userInput := map[string]any{
		"first_name": "Delete",
		"last_name":  "Test",
		"email":      testEmail,
		"gender":     "f",
		"password":   "DeletePass123!",
		"is_officer": true,
	}

	body, _ := json.Marshal(userInput)
	req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	testApp.registerUserHandler(rec, req)

	var createResponse map[string]any
	_ = json.NewDecoder(rec.Result().Body).Decode(&createResponse)
	userResponse := createResponse["user"].(map[string]any)
	userID := int64(userResponse["id"].(float64))

	defer testApp.models.User.HardDelete(userID) // Final cleanup

	adminUser := getSeededUser(t, "admin1@police-training.bz")
	adminToken := createTokenForSeededUser(t, adminUser.ID)

	tests := []struct {
		name           string
		userID         string
		expectedStatus int
	}{
		{
			name:           "Valid user soft deletion",
			userID:         strconv.FormatInt(userID, 10),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Non-existent user",
			userID:         "999999",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test: %s", tt.name)

			path := fmt.Sprintf("/v1/users/%s", tt.userID)
			req := httptest.NewRequest(http.MethodDelete, path, nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))

			req = setURLParam(req, "id", tt.userID)
			req = setUserContext(req, adminUser)

			rec := httptest.NewRecorder()
			testApp.deleteUserHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			t.Logf("Step: Received status code %d", res.StatusCode)
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}

			t.Logf("Completed test: %s", tt.name)
		})
	}
}

func TestRestoreUserHandler(t *testing.T) {
	t.Log("=== Testing Restore User Handler ===")

	// Create and delete test user first
	testEmail := fmt.Sprintf("restore%d@test-police-training.bz", time.Now().UnixNano())
	userInput := map[string]any{
		"first_name": "Restore",
		"last_name":  "Test",
		"email":      testEmail,
		"gender":     "m",
		"password":   "RestorePass123!",
		"is_officer": true,
	}

	body, _ := json.Marshal(userInput)
	req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	testApp.registerUserHandler(rec, req)

	var createResponse map[string]any
	_ = json.NewDecoder(rec.Result().Body).Decode(&createResponse)
	userResponse := createResponse["user"].(map[string]any)
	userID := int64(userResponse["id"].(float64))

	defer testApp.models.User.HardDelete(userID) // Final cleanup

	// Soft delete the user first
	testApp.models.User.SoftDelete(userID)

	adminUser := getSeededUser(t, "admin1@police-training.bz")
	adminToken := createTokenForSeededUser(t, adminUser.ID)

	tests := []struct {
		name           string
		userID         string
		expectedStatus int
	}{
		{
			name:           "Valid user restore",
			userID:         strconv.FormatInt(userID, 10),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Non-existent user",
			userID:         "999999",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test: %s", tt.name)

			path := fmt.Sprintf("/v1/users/%s/restore", tt.userID)
			req := httptest.NewRequest(http.MethodPatch, path, nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))

			req = setURLParam(req, "id", tt.userID)
			req = setUserContext(req, adminUser)

			rec := httptest.NewRecorder()
			testApp.restoreUserHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			t.Logf("Step: Received status code %d", res.StatusCode)
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}

			t.Logf("Completed test: %s", tt.name)
		})
	}
}

func TestHardDeleteUserHandler(t *testing.T) {
	t.Log("=== Testing Hard Delete User Handler ===")

	// Create test user for hard deletion
	testEmail := fmt.Sprintf("harddelete%d@test-police-training.bz", time.Now().UnixNano())
	userInput := map[string]any{
		"first_name": "HardDelete",
		"last_name":  "Test",
		"email":      testEmail,
		"gender":     "f",
		"password":   "HardDeletePass123!",
		"is_officer": true,
	}

	body, _ := json.Marshal(userInput)
	req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	testApp.registerUserHandler(rec, req)

	var createResponse map[string]any
	_ = json.NewDecoder(rec.Result().Body).Decode(&createResponse)
	userResponse := createResponse["user"].(map[string]any)
	userID := int64(userResponse["id"].(float64))

	adminUser := getSeededUser(t, "admin1@police-training.bz")
	adminToken := createTokenForSeededUser(t, adminUser.ID)

	tests := []struct {
		name           string
		userID         string
		expectedStatus int
	}{
		{
			name:           "Valid user hard deletion",
			userID:         strconv.FormatInt(userID, 10),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Non-existent user",
			userID:         "999999",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test: %s", tt.name)

			path := fmt.Sprintf("/v1/users/%s/hard-delete", tt.userID)
			req := httptest.NewRequest(http.MethodDelete, path, nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))

			req = setURLParam(req, "id", tt.userID)
			req = setUserContext(req, adminUser)

			rec := httptest.NewRecorder()
			testApp.hardDeleteUserHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			t.Logf("Step: Received status code %d", res.StatusCode)
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}

			t.Logf("Completed test: %s", tt.name)
		})
	}
}

func TestUserWorkflow(t *testing.T) {
	t.Log("=== Testing Complete User Workflow ===")

	testEmail := fmt.Sprintf("workflow%d@test-police-training.bz", time.Now().UnixNano())

	// 1. Register user
	userInput := map[string]any{
		"first_name":     "Workflow",
		"last_name":      "Test",
		"email":          testEmail,
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

	var registerResponse map[string]any
	_ = json.NewDecoder(rec.Result().Body).Decode(&registerResponse)
	userResponse := registerResponse["user"].(map[string]any)
	userID := int64(userResponse["id"].(float64))

	defer testApp.models.User.HardDelete(userID) // Cleanup

	t.Logf("Step: Created user ID %d", userID)

	// 2. Create and use activation token
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

	t.Logf("Step: Successfully activated user")

	// 3. Create auth token and access user
	authToken, err := testApp.models.Token.New(userID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		t.Fatalf("Failed to create auth token: %v", err)
	}

	user, err := testApp.models.User.Get(userID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	req = httptest.NewRequest(http.MethodGet, "/v1/me", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken.Plaintext))
	req = setUserContext(req, user)

	rec = httptest.NewRecorder()
	testApp.showCurrentUserHandler(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("Current user access failed with status %d", rec.Result().StatusCode)
	}

	t.Log("Step: Complete user workflow test passed successfully!")
}
