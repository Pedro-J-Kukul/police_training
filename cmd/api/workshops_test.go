package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Pedro-J-Kukul/police_training/internal/data"
)

// Test-specific cleanup helper
type workshopTestHelper struct {
	t          *testing.T
	createdIDs []int64
	db         *sql.DB
}

func newWorkshopTestHelper(t *testing.T) *workshopTestHelper {
	return &workshopTestHelper{
		t:          t,
		createdIDs: make([]int64, 0),
		db:         testApp.models.User.DB,
	}
}

func (h *workshopTestHelper) addWorkshopID(id int64) {
	h.createdIDs = append(h.createdIDs, id)
	h.t.Logf("Tracking workshop ID for cleanup: %d", id)
}

func (h *workshopTestHelper) cleanup() {
	h.t.Helper()
	h.t.Logf("Cleaning up %d workshop records", len(h.createdIDs))

	for _, id := range h.createdIDs {
		_, err := h.db.Exec("DELETE FROM workshops WHERE id = $1", id)
		if err != nil {
			h.t.Logf("Warning: Failed to cleanup workshop ID %d: %v", id, err)
		} else {
			h.t.Logf("Successfully deleted workshop ID: %d", id)
		}
	}
	h.createdIDs = nil
}

func TestCreateWorkshopHandler(t *testing.T) {
	helper := newWorkshopTestHelper(t)
	defer helper.cleanup()

	// Get existing seed data
	category, err := testApp.models.TrainingCategory.Get(1)
	if err != nil {
		t.Skipf("Skipping test - seed data not available: %v", err)
	}

	trainingType, err := testApp.models.TrainingType.Get(1)
	if err != nil {
		t.Skipf("Skipping test - seed data not available: %v", err)
	}

	_, token := createTestUserWithToken(t)

	tests := []struct {
		name           string
		input          map[string]any
		token          string
		expectedStatus int
		checkResponse  func(*testing.T, *http.Response, *workshopTestHelper)
	}{
		{
			name: "Valid workshop creation",
			input: map[string]any{
				"workshop_name": fmt.Sprintf("TEST_API_Workshop_%d", time.Now().UnixNano()),
				"category_id":   category.ID,
				"type_id":       trainingType.ID,
				"credit_hours":  40,
				"description":   "Test workshop description",
				"is_active":     true,
			},
			token:          token,
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, res *http.Response, h *workshopTestHelper) {
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if response["workshop"] == nil {
					t.Error("Expected workshop object in response")
				} else {
					workshop := response["workshop"].(map[string]any)
					workshopID := int64(workshop["id"].(float64))
					h.addWorkshopID(workshopID) // Track for cleanup
					t.Logf("Created workshop with ID: %d", workshopID)
				}
			},
		},
		{
			name: "Missing required fields",
			input: map[string]any{
				"workshop_name": "Test",
			},
			token:          token,
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, res *http.Response, h *workshopTestHelper) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := createAuthenticatedRequest(http.MethodPost, "/v1/workshops", body, tt.token)

			rec := httptest.NewRecorder()

			if tt.token != "" {
				user, _ := testApp.models.User.GetForToken(data.ScopeAuthentication, tt.token)
				req = setUserContext(req, user)
				testApp.createWorkshopHandler(rec, req)
			} else {
				handler := testApp.authenticate(http.HandlerFunc(testApp.createWorkshopHandler))
				handler.ServeHTTP(rec, req)
			}

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}

			tt.checkResponse(t, res, helper)
		})
	}
}

func TestListWorkshopsHandler(t *testing.T) {
	helper := newWorkshopTestHelper(t)
	defer helper.cleanup()

	_, token := createTestUserWithToken(t)

	// Get existing seed data
	category, err := testApp.models.TrainingCategory.Get(1)
	if err != nil {
		t.Skipf("Skipping test - seed data not available: %v", err)
	}

	trainingType, err := testApp.models.TrainingType.Get(1)
	if err != nil {
		t.Skipf("Skipping test - seed data not available: %v", err)
	}

	// Create test workshops
	numTestWorkshops := 3
	for i := 0; i < numTestWorkshops; i++ {
		workshop := &data.Workshop{
			WorkshopName: fmt.Sprintf("TEST_List_Workshop_%d_%d", i, time.Now().UnixNano()),
			CategoryID:   category.ID,
			TypeID:       trainingType.ID,
			CreditHours:  40,
			IsActive:     true,
		}
		err := testApp.models.Workshop.Insert(workshop)
		if err != nil {
			t.Fatalf("Failed to create test workshop: %v", err)
		}
		helper.addWorkshopID(workshop.ID)
		t.Logf("Created test workshop with ID: %d", workshop.ID)
	}

	tests := []struct {
		name           string
		queryParams    string
		token          string
		expectedStatus int
		minWorkshops   int // Minimum workshops expected (since seed data exists)
	}{
		{
			name:           "Valid list request",
			queryParams:    "",
			token:          token,
			expectedStatus: http.StatusOK,
			minWorkshops:   numTestWorkshops, // At least our test workshops
		},
		{
			name:           "With pagination",
			queryParams:    "?page=1&page_size=2",
			token:          token,
			expectedStatus: http.StatusOK,
			minWorkshops:   0, // Could be 0-2 depending on seed data
		},
		{
			name:           "Filter by category",
			queryParams:    fmt.Sprintf("?category_id=%d", category.ID),
			token:          token,
			expectedStatus: http.StatusOK,
			minWorkshops:   numTestWorkshops, // At least our test workshops
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/v1/workshops" + tt.queryParams
			req := createAuthenticatedRequest(http.MethodGet, path, nil, tt.token)

			rec := httptest.NewRecorder()

			user, _ := testApp.models.User.GetForToken(data.ScopeAuthentication, tt.token)
			req = setUserContext(req, user)
			testApp.listWorkshopsHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				workshops := response["workshops"].([]any)
				if len(workshops) < tt.minWorkshops {
					t.Errorf("Expected at least %d workshops; got %d", tt.minWorkshops, len(workshops))
				}
				t.Logf("Retrieved %d workshops", len(workshops))
			}
		})
	}
}

func TestUpdateWorkshopHandler(t *testing.T) {
	helper := newWorkshopTestHelper(t)
	defer helper.cleanup()

	_, token := createTestUserWithToken(t)

	// Get existing seed data and create test workshop
	category, err := testApp.models.TrainingCategory.Get(1)
	if err != nil {
		t.Skipf("Skipping test - seed data not available: %v", err)
	}

	trainingType, err := testApp.models.TrainingType.Get(1)
	if err != nil {
		t.Skipf("Skipping test - seed data not available: %v", err)
	}

	workshop := &data.Workshop{
		WorkshopName: fmt.Sprintf("TEST_Update_Workshop_%d", time.Now().UnixNano()),
		CategoryID:   category.ID,
		TypeID:       trainingType.ID,
		CreditHours:  40,
		IsActive:     true,
	}
	err = testApp.models.Workshop.Insert(workshop)
	if err != nil {
		t.Fatalf("Failed to create test workshop: %v", err)
	}
	helper.addWorkshopID(workshop.ID)

	tests := []struct {
		name           string
		workshopID     string
		input          map[string]any
		token          string
		expectedStatus int
	}{
		{
			name:       "Valid workshop update",
			workshopID: fmt.Sprintf("%d", workshop.ID),
			input: map[string]any{
				"workshop_name": fmt.Sprintf("TEST_Updated_Workshop_%d", time.Now().UnixNano()),
				"credit_hours":  60,
			},
			token:          token,
			expectedStatus: http.StatusOK,
		},
		{
			name:       "Non-existent workshop",
			workshopID: "999999",
			input: map[string]any{
				"workshop_name": "Non-existent Workshop",
			},
			token:          token,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			path := fmt.Sprintf("/v1/workshops/%s", tt.workshopID)
			req := createAuthenticatedRequest(http.MethodPatch, path, body, tt.token)

			req = setURLParam(req, "id", tt.workshopID)

			rec := httptest.NewRecorder()

			user, _ := testApp.models.User.GetForToken(data.ScopeAuthentication, tt.token)
			req = setUserContext(req, user)
			testApp.updateWorkshopHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}
		})
	}
}

func TestDeleteWorkshopHandler(t *testing.T) {
	helper := newWorkshopTestHelper(t)
	defer helper.cleanup()

	_, token := createTestUserWithToken(t)

	// Get existing seed data and create test workshop
	category, err := testApp.models.TrainingCategory.Get(1)
	if err != nil {
		t.Skipf("Skipping test - seed data not available: %v", err)
	}

	trainingType, err := testApp.models.TrainingType.Get(1)
	if err != nil {
		t.Skipf("Skipping test - seed data not available: %v", err)
	}

	workshop := &data.Workshop{
		WorkshopName: fmt.Sprintf("TEST_Delete_Workshop_%d", time.Now().UnixNano()),
		CategoryID:   category.ID,
		TypeID:       trainingType.ID,
		CreditHours:  40,
		IsActive:     true,
	}
	err = testApp.models.Workshop.Insert(workshop)
	if err != nil {
		t.Fatalf("Failed to create test workshop: %v", err)
	}
	// Note: We don't add to helper here since the test itself will delete it

	tests := []struct {
		name           string
		workshopID     string
		token          string
		expectedStatus int
	}{
		{
			name:           "Valid workshop deletion",
			workshopID:     fmt.Sprintf("%d", workshop.ID),
			token:          token,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Non-existent workshop",
			workshopID:     "999999",
			token:          token,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := fmt.Sprintf("/v1/workshops/%s", tt.workshopID)
			req := createAuthenticatedRequest(http.MethodDelete, path, nil, tt.token)

			req = setURLParam(req, "id", tt.workshopID)

			rec := httptest.NewRecorder()

			user, _ := testApp.models.User.GetForToken(data.ScopeAuthentication, tt.token)
			req = setUserContext(req, user)
			testApp.deleteWorkshopHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}
		})
	}
}
