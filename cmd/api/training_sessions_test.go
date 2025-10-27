package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/Pedro-J-Kukul/police_training/internal/data"
)

func getSeededSessionData(t *testing.T) (facilitatorID, workshopID, formationID, regionID, statusID int64) {
	t.Helper()
	t.Log("Step: Getting seeded session reference data")

	// Get facilitator user
	facilitator := getSeededUser(t, "maria.rodriguez@police-training.bz")
	facilitatorID = facilitator.ID

	// Get workshop
	workshops, _, err := testApp.models.Workshop.GetAll("", nil, nil, nil, data.Filters{Page: 1, PageSize: 1, Sort: "id", SortSafelist: []string{"id"}})
	if err != nil || len(workshops) == 0 {
		t.Fatal("No workshops found in seeded data")
	}
	workshopID = workshops[0].ID

	// Get formation
	formation, err := testApp.models.Formation.Get(1)
	if err != nil {
		t.Fatal("Formation ID 1 not found in seeded data")
	}
	formationID = formation.ID

	// Get region
	region, err := testApp.models.Region.GetByName("Northern Region")
	if err != nil {
		t.Fatal("Northern Region not found in seeded data")
	}
	regionID = region.ID

	// Get training status
	status, err := testApp.models.TrainingStatus.GetByName("Scheduled")
	if err != nil {
		t.Fatal("Scheduled status not found in seeded data")
	}
	statusID = status.ID

	t.Logf("Step: Using facilitator ID %d, workshop ID %d, formation ID %d, region ID %d, status ID %d",
		facilitatorID, workshopID, formationID, regionID, statusID)
	return
}

func getSeededSession(t *testing.T) *data.TrainingSession {
	t.Helper()
	t.Log("Step: Getting seeded training session")

	sessions, _, err := testApp.models.TrainingSession.GetAll(nil, nil, nil, nil, nil, nil, data.Filters{
		Page: 1, PageSize: 1, Sort: "id", SortSafelist: []string{"id"},
	})
	if err != nil || len(sessions) == 0 {
		t.Fatal("No training sessions found in seeded data")
	}

	session := sessions[0]
	t.Logf("Step: Using seeded session ID %d on %s", session.ID, session.SessionDate.Format("2006-01-02"))
	return session
}

func createTestSession(t *testing.T) *data.TrainingSession {
	t.Helper()

	facilitatorID, workshopID, formationID, regionID, statusID := getSeededSessionData(t)

	session := &data.TrainingSession{
		FacilitatorID:    facilitatorID,
		WorkshopID:       workshopID,
		FormationID:      formationID,
		RegionID:         regionID,
		TrainingStatusID: statusID,
		SessionDate:      time.Now().AddDate(0, 0, 30), // 30 days from now
		StartTime:        time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
		EndTime:          time.Date(0, 1, 1, 17, 0, 0, 0, time.UTC),
		Location:         stringPtr("Test Training Room"),
		MaxCapacity:      intPtr(25),
		Notes:            stringPtr("Test session for API testing"),
	}

	err := testApp.models.TrainingSession.Insert(session)
	if err != nil {
		t.Fatalf("Failed to create test session: %v", err)
	}

	t.Logf("Step: Created test session ID %d", session.ID)
	return session
}

func intPtr(i int) *int {
	return &i
}

func TestCreateTrainingSessionHandler(t *testing.T) {
	t.Log("=== Testing Create Training Session Handler ===")

	facilitatorID, workshopID, formationID, regionID, statusID := getSeededSessionData(t)
	facilitator := getSeededUser(t, "maria.rodriguez@police-training.bz")
	facilitatorToken := createTokenForSeededUser(t, facilitator.ID)

	tests := []struct {
		name           string
		input          map[string]any
		expectedStatus int
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name: "Valid session creation",
			input: map[string]any{
				"facilitator_id":     facilitatorID,
				"workshop_id":        workshopID,
				"formation_id":       formationID,
				"region_id":          regionID,
				"training_status_id": statusID,
				"session_date":       time.Now().AddDate(0, 0, 45).Format("2006-01-02"),
				"start_time":         "09:00",
				"end_time":           "17:00",
				"location":           "API Test Room",
				"max_capacity":       30,
				"notes":              "Session created via API test",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating successful session creation")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response["training_session"] == nil {
					t.Error("Expected training_session object in response")
					return
				}

				session := response["training_session"].(map[string]any)
				sessionID := int64(session["id"].(float64))
				t.Logf("Step: Created session ID %d", sessionID)

				// Cleanup
				testApp.models.TrainingSession.Delete(sessionID)
			},
		},
		{
			name: "Missing required fields",
			input: map[string]any{
				"facilitator_id": facilitatorID,
				"workshop_id":    workshopID,
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating missing fields validation")
			},
		},
		{
			name: "Invalid date format",
			input: map[string]any{
				"facilitator_id":     facilitatorID,
				"workshop_id":        workshopID,
				"formation_id":       formationID,
				"region_id":          regionID,
				"training_status_id": statusID,
				"session_date":       "invalid-date",
				"start_time":         "09:00",
				"end_time":           "17:00",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating invalid date format handling")
			},
		},
		{
			name: "Invalid time format",
			input: map[string]any{
				"facilitator_id":     facilitatorID,
				"workshop_id":        workshopID,
				"formation_id":       formationID,
				"region_id":          regionID,
				"training_status_id": statusID,
				"session_date":       time.Now().AddDate(0, 0, 50).Format("2006-01-02"),
				"start_time":         "invalid-time",
				"end_time":           "17:00",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating invalid time format handling")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test: %s", tt.name)

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/v1/training-sessions", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", facilitatorToken))

			req = setUserContext(req, facilitator)

			rec := httptest.NewRecorder()
			testApp.createTrainingSessionHandler(rec, req)

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

func TestShowTrainingSessionHandler(t *testing.T) {
	t.Log("=== Testing Show Training Session Handler ===")

	seededSession := getSeededSession(t)
	facilitator := getSeededUser(t, "maria.rodriguez@police-training.bz")
	facilitatorToken := createTokenForSeededUser(t, facilitator.ID)

	tests := []struct {
		name           string
		sessionID      string
		expectedStatus int
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name:           "Valid seeded session ID",
			sessionID:      strconv.FormatInt(seededSession.ID, 10),
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating session details response")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response["training_session"] == nil {
					t.Error("Expected training_session object in response")
					return
				}

				session := response["training_session"].(map[string]any)
				if session["id"] != float64(seededSession.ID) {
					t.Errorf("Expected session ID %d, got %v", seededSession.ID, session["id"])
				}

				t.Logf("Step: Retrieved session on %v", session["session_date"])

				// Verify expected fields
				requiredFields := []string{"id", "facilitator_id", "workshop_id", "formation_id", "region_id", "training_status_id", "session_date", "start_time", "end_time"}
				for _, field := range requiredFields {
					if session[field] == nil {
						t.Errorf("Expected field %s in session response", field)
					}
				}
			},
		},
		{
			name:           "Non-existent session ID",
			sessionID:      "999999",
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating non-existent session response")
			},
		},
		{
			name:           "Invalid session ID format",
			sessionID:      "abc",
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating invalid ID format response")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test: %s with session ID %s", tt.name, tt.sessionID)

			path := fmt.Sprintf("/v1/training-sessions/%s", tt.sessionID)
			req := httptest.NewRequest(http.MethodGet, path, nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", facilitatorToken))

			req = setURLParam(req, "id", tt.sessionID)
			req = setUserContext(req, facilitator)

			rec := httptest.NewRecorder()
			testApp.showTrainingSessionHandler(rec, req)

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

func TestListTrainingSessionsHandler(t *testing.T) {
	t.Log("=== Testing List Training Sessions Handler ===")

	_, workshopID, formationID, _, _ := getSeededSessionData(t)
	facilitator := getSeededUser(t, "maria.rodriguez@police-training.bz")
	facilitatorToken := createTokenForSeededUser(t, facilitator.ID)

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name:           "List all sessions",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating sessions list response")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response["training_sessions"] == nil {
					t.Error("Expected training_sessions array in response")
					return
				}

				sessions := response["training_sessions"].([]any)
				t.Logf("Step: Retrieved %d sessions", len(sessions))

				if len(sessions) > 0 {
					session := sessions[0].(map[string]any)
					requiredFields := []string{"id", "facilitator_id", "workshop_id", "formation_id", "region_id", "session_date"}
					for _, field := range requiredFields {
						if session[field] == nil {
							t.Errorf("Expected field %s in session response", field)
						}
					}

					t.Logf("Step: First session ID %v on %v", session["id"], session["session_date"])
				}

				if response["metadata"] == nil {
					t.Error("Expected metadata in response")
				}
			},
		},
		{
			name:           "List sessions with pagination",
			queryParams:    "?page=1&page_size=3",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating paginated sessions response")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				sessions := response["training_sessions"].([]any)
				t.Logf("Step: Retrieved %d sessions with page_size=3", len(sessions))

				if len(sessions) > 3 {
					t.Errorf("Expected max 3 sessions, got %d", len(sessions))
				}
			},
		},
		{
			name:           "Filter sessions by facilitator",
			queryParams:    fmt.Sprintf("?facilitator_id=%d", facilitator.ID),
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating facilitator filter response")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				sessions := response["training_sessions"].([]any)
				t.Logf("Step: Retrieved %d sessions for facilitator %d", len(sessions), facilitator.ID)

				// Verify all sessions belong to the facilitator
				for i, sessionInterface := range sessions {
					session := sessionInterface.(map[string]any)
					if facilID, ok := session["facilitator_id"].(float64); !ok || facilID != float64(facilitator.ID) {
						t.Errorf("Session at index %d does not belong to facilitator %d", i, facilitator.ID)
					}
				}
			},
		},
		{
			name:           "Filter sessions by workshop",
			queryParams:    fmt.Sprintf("?workshop_id=%d", workshopID),
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating workshop filter response")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				sessions := response["training_sessions"].([]any)
				t.Logf("Step: Retrieved %d sessions for workshop %d", len(sessions), workshopID)
			},
		},
		{
			name:           "Filter sessions by formation",
			queryParams:    fmt.Sprintf("?formation_id=%d", formationID),
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating formation filter response")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				sessions := response["training_sessions"].([]any)
				t.Logf("Step: Retrieved %d sessions for formation %d", len(sessions), formationID)
			},
		},
		{
			name:           "Filter sessions by date",
			queryParams:    "?session_date=" + time.Now().Format("2006-01-02"),
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating date filter response")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				sessionsInterface := response["training_sessions"]
				if sessionsInterface == nil {
					t.Logf("Step: Retrieved 0 sessions for today (no sessions found)")
					return
				}
				sessions := sessionsInterface.([]any)
				t.Logf("Step: Retrieved %d sessions for today", len(sessions))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test: %s", tt.name)

			path := "/v1/training-sessions" + tt.queryParams
			req := httptest.NewRequest(http.MethodGet, path, nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", facilitatorToken))
			req = setUserContext(req, facilitator)

			rec := httptest.NewRecorder()
			testApp.listTrainingSessionsHandler(rec, req)

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

func TestUpdateTrainingSessionHandler(t *testing.T) {
	t.Log("=== Testing Update Training Session Handler ===")

	// Create test session for updating
	testSession := createTestSession(t)
	defer testApp.models.TrainingSession.Delete(testSession.ID) // Cleanup

	facilitator := getSeededUser(t, "maria.rodriguez@police-training.bz")
	facilitatorToken := createTokenForSeededUser(t, facilitator.ID)

	tests := []struct {
		name           string
		sessionID      string
		input          map[string]any
		expectedStatus int
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name:      "Valid session update",
			sessionID: strconv.FormatInt(testSession.ID, 10),
			input: map[string]any{
				"location":     "Updated Training Room",
				"max_capacity": 40,
				"notes":        "Updated session notes",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating successful session update")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response["training_session"] == nil {
					t.Error("Expected training_session object in response")
					return
				}

				session := response["training_session"].(map[string]any)
				if session["max_capacity"] != float64(40) {
					t.Errorf("Expected max_capacity 40, got %v", session["max_capacity"])
				}

				t.Logf("Step: Updated session location to %v", session["location"])
			},
		},
		{
			name:      "Non-existent session",
			sessionID: "999999",
			input: map[string]any{
				"location": "Non-existent Session Room",
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating non-existent session update")
			},
		},
		{
			name:      "Invalid date format",
			sessionID: strconv.FormatInt(testSession.ID, 10),
			input: map[string]any{
				"session_date": "invalid-date",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating invalid date format handling")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test: %s", tt.name)

			body, _ := json.Marshal(tt.input)
			path := fmt.Sprintf("/v1/training-sessions/%s", tt.sessionID)
			req := httptest.NewRequest(http.MethodPatch, path, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", facilitatorToken))

			req = setURLParam(req, "id", tt.sessionID)
			req = setUserContext(req, facilitator)

			rec := httptest.NewRecorder()
			testApp.updateTrainingSessionHandler(rec, req)

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

func TestDeleteTrainingSessionHandler(t *testing.T) {
	t.Log("=== Testing Delete Training Session Handler ===")

	// Create test session for deletion
	testSession := createTestSession(t)
	// No defer cleanup since the test will delete it

	facilitator := getSeededUser(t, "maria.rodriguez@police-training.bz")
	facilitatorToken := createTokenForSeededUser(t, facilitator.ID)

	tests := []struct {
		name           string
		sessionID      string
		expectedStatus int
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name:           "Valid session deletion",
			sessionID:      strconv.FormatInt(testSession.ID, 10),
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating successful session deletion")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response["message"] == nil {
					t.Error("Expected message in response")
				}

				t.Logf("Step: Session deleted with message: %v", response["message"])
			},
		},
		{
			name:           "Non-existent session",
			sessionID:      "999999",
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating non-existent session deletion")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test: %s", tt.name)

			path := fmt.Sprintf("/v1/training-sessions/%s", tt.sessionID)
			req := httptest.NewRequest(http.MethodDelete, path, nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", facilitatorToken))

			req = setURLParam(req, "id", tt.sessionID)
			req = setUserContext(req, facilitator)

			rec := httptest.NewRecorder()
			testApp.deleteTrainingSessionHandler(rec, req)

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

func TestSessionWorkflow(t *testing.T) {
	t.Log("=== Testing Complete Session Workflow ===")

	facilitatorID, workshopID, formationID, regionID, statusID := getSeededSessionData(t)
	facilitator := getSeededUser(t, "maria.rodriguez@police-training.bz")
	facilitatorToken := createTokenForSeededUser(t, facilitator.ID)

	// 1. Create session
	sessionInput := map[string]any{
		"facilitator_id":     facilitatorID,
		"workshop_id":        workshopID,
		"formation_id":       formationID,
		"region_id":          regionID,
		"training_status_id": statusID,
		"session_date":       time.Now().AddDate(0, 0, 60).Format("2006-01-02"),
		"start_time":         "10:00",
		"end_time":           "16:00",
		"location":           "Workflow Test Room",
		"max_capacity":       20,
		"notes":              "Complete workflow test session",
	}

	body, _ := json.Marshal(sessionInput)
	req := httptest.NewRequest(http.MethodPost, "/v1/training-sessions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", facilitatorToken))
	req = setUserContext(req, facilitator)

	rec := httptest.NewRecorder()
	testApp.createTrainingSessionHandler(rec, req)

	if rec.Result().StatusCode != http.StatusCreated {
		t.Fatalf("Session creation failed with status %d", rec.Result().StatusCode)
	}

	var createResponse map[string]any
	_ = json.NewDecoder(rec.Result().Body).Decode(&createResponse)
	sessionResponse := createResponse["training_session"].(map[string]any)
	sessionID := int64(sessionResponse["id"].(float64))

	defer testApp.models.TrainingSession.Delete(sessionID) // Cleanup

	t.Logf("Step: Created session ID %d", sessionID)

	// 2. Get the session
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/training-sessions/%d", sessionID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", facilitatorToken))
	req = setURLParam(req, "id", strconv.FormatInt(sessionID, 10))
	req = setUserContext(req, facilitator)

	rec = httptest.NewRecorder()
	testApp.showTrainingSessionHandler(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("Session retrieval failed with status %d", rec.Result().StatusCode)
	}

	t.Logf("Step: Successfully retrieved session")

	// 3. Update the session
	updateInput := map[string]any{
		"location":     "Updated Workflow Room",
		"max_capacity": 35,
		"notes":        "Updated workflow session notes",
	}

	body, _ = json.Marshal(updateInput)
	req = httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/training-sessions/%d", sessionID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", facilitatorToken))
	req = setURLParam(req, "id", strconv.FormatInt(sessionID, 10))
	req = setUserContext(req, facilitator)

	rec = httptest.NewRecorder()
	testApp.updateTrainingSessionHandler(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("Session update failed with status %d", rec.Result().StatusCode)
	}

	t.Logf("Step: Successfully updated session")

	// 4. List sessions (should include our session)
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/training-sessions?facilitator_id=%d", facilitator.ID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", facilitatorToken))
	req = setUserContext(req, facilitator)

	rec = httptest.NewRecorder()
	testApp.listTrainingSessionsHandler(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("Session listing failed with status %d", rec.Result().StatusCode)
	}

	t.Logf("Step: Successfully listed sessions")

	// 5. Delete the session
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/training-sessions/%d", sessionID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", facilitatorToken))
	req = setURLParam(req, "id", strconv.FormatInt(sessionID, 10))
	req = setUserContext(req, facilitator)

	rec = httptest.NewRecorder()
	testApp.deleteTrainingSessionHandler(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("Session deletion failed with status %d", rec.Result().StatusCode)
	}

	t.Logf("Step: Successfully deleted session")

	t.Log("Step: Complete session workflow test passed successfully!")
}
