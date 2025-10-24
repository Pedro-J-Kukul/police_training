package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Pedro-J-Kukul/police_training/internal/data"
)

// Test helper for training sessions API tests
type trainingSessionAPITestHelper struct {
	t          *testing.T
	createdIDs []int64
	db         *sql.DB
}

func newTrainingSessionAPITestHelper(t *testing.T) *trainingSessionAPITestHelper {
	return &trainingSessionAPITestHelper{
		t:          t,
		createdIDs: make([]int64, 0),
		db:         testApp.models.User.DB,
	}
}

func (h *trainingSessionAPITestHelper) addSessionID(id int64) {
	h.createdIDs = append(h.createdIDs, id)
	h.t.Logf("Tracking training session ID for cleanup: %d", id)
}

func (h *trainingSessionAPITestHelper) cleanup() {
	h.t.Helper()
	h.t.Logf("Cleaning up %d training session records", len(h.createdIDs))

	for _, id := range h.createdIDs {
		_, err := h.db.Exec("DELETE FROM training_sessions WHERE id = $1", id)
		if err != nil {
			h.t.Logf("Warning: Failed to cleanup training session ID %d: %v", id, err)
		} else {
			h.t.Logf("Successfully deleted training session ID: %d", id)
		}
	}
	h.createdIDs = nil
}

func TestCreateTrainingSessionHandler(t *testing.T) {
	helper := newTrainingSessionAPITestHelper(t)
	defer helper.cleanup()

	_, token := createTestUserWithToken(t)

	tests := []struct {
		name           string
		input          map[string]any
		token          string
		expectedStatus int
		checkResponse  func(*testing.T, *http.Response, *trainingSessionAPITestHelper)
	}{
		{
			name: "Valid training session creation",
			input: map[string]any{
				"facilitator_id":     1,
				"workshop_id":        1,
				"formation_id":       1,
				"region_id":          1,
				"session_date":       "2025-12-25",
				"start_time":         "09:00",
				"end_time":           "17:00",
				"location":           "Test Training Room",
				"max_capacity":       30,
				"training_status_id": 1,
				"notes":              "Test training session",
			},
			token:          token,
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, res *http.Response, h *trainingSessionAPITestHelper) {
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if response["training_session"] == nil {
					t.Error("Expected training_session object in response")
				} else {
					session := response["training_session"].(map[string]any)
					sessionID := int64(session["id"].(float64))
					h.addSessionID(sessionID)
					t.Logf("Created training session with ID: %d", sessionID)
				}
			},
		},
		{
			name: "Missing required fields",
			input: map[string]any{
				"facilitator_id": 1,
				"workshop_id":    1,
			},
			token:          token,
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, res *http.Response, h *trainingSessionAPITestHelper) {
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
			req := createAuthenticatedRequest(http.MethodPost, "/v1/training-sessions", body, tt.token)

			rec := httptest.NewRecorder()

			user, _ := testApp.models.User.GetForToken(data.ScopeAuthentication, tt.token)
			req = setUserContext(req, user)
			testApp.createTrainingSessionHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}

			tt.checkResponse(t, res, helper)
		})
	}
}
