package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/Pedro-J-Kukul/police_training/internal/data"
)

// Helper function to create test dependencies for officers
func createOfficerTestDependencies(t *testing.T) (userID, regionID, formationID, postingID, rankID int64) {
	t.Helper()

	// Create test user
	user := &data.User{
		FirstName:   "Officer",
		LastName:    "Test",
		Email:       fmt.Sprintf("officer%d@example.com", time.Now().UnixNano()),
		Gender:      "m",
		IsActivated: true,
		IsOfficer:   true,
	}
	_ = user.Password.Set("TestPass123!")
	err := testApp.models.User.Insert(user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create test region
	region := &data.Region{Region: "Test Officer Region"}
	err = testApp.models.Region.Insert(region)
	if err != nil {
		t.Fatalf("Failed to create test region: %v", err)
	}

	// Create test formation
	formation := &data.Formation{Formation: "Test Officer Formation", RegionID: region.ID}
	err = testApp.models.Formation.Insert(formation)
	if err != nil {
		t.Fatalf("Failed to create test formation: %v", err)
	}

	// Create test posting
	posting := &data.Posting{Posting: "Test Officer Posting", Code: "TOP"}
	err = testApp.models.Posting.Insert(posting)
	if err != nil {
		t.Fatalf("Failed to create test posting: %v", err)
	}

	// Create test rank
	rank := &data.Rank{Rank: "Test Officer Rank", Code: "TOR", AnnualTrainingHoursRequired: 40}
	err = testApp.models.Rank.Insert(rank)
	if err != nil {
		t.Fatalf("Failed to create test rank: %v", err)
	}

	return user.ID, region.ID, formation.ID, posting.ID, rank.ID
}

// Helper function to create a test officer and return authentication token
func createTestOfficerWithToken(t *testing.T) (*data.Officer, string) {
	t.Helper()

	userID, regionID, formationID, postingID, rankID := createOfficerTestDependencies(t)

	officer := &data.Officer{
		UserID:           userID,
		RegulationNumber: fmt.Sprintf("OFF%d", time.Now().UnixNano()),
		RankID:           rankID,
		PostingID:        postingID,
		FormationID:      formationID,
		RegionID:         regionID,
	}

	err := testApp.models.Officer.Insert(officer)
	if err != nil {
		t.Fatalf("Failed to create test officer: %v", err)
	}

	// Create authentication token for the user
	token, err := testApp.models.Token.New(userID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	return officer, token.Plaintext
}

func TestCreateOfficerHandler(t *testing.T) {
	// Clean up before test
	cleanupAPITestData(t)
	defer cleanupAPITestData(t)

	userID, regionID, formationID, postingID, rankID := createOfficerTestDependencies(t)
	_, token := createTestUserWithToken(t)

	tests := []struct {
		name           string
		input          map[string]any
		token          string
		expectedStatus int
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name: "Valid officer creation",
			input: map[string]any{
				"user_id":           userID,
				"regulation_number": fmt.Sprintf("NEW%d", time.Now().UnixNano()),
				"rank_id":           rankID,
				"posting_id":        postingID,
				"formation_id":      formationID,
				"region_id":         regionID,
			},
			token:          token,
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, res *http.Response) {
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if response["officer"] == nil {
					t.Error("Expected officer object in response")
				}
			},
		},
		{
			name: "Missing required fields",
			input: map[string]any{
				"user_id": userID,
			},
			token:          token,
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
			name: "Invalid user_id",
			input: map[string]any{
				"user_id":           99999,
				"regulation_number": fmt.Sprintf("INV%d", time.Now().UnixNano()),
				"rank_id":           rankID,
				"posting_id":        postingID,
				"formation_id":      formationID,
				"region_id":         regionID,
			},
			token:          token,
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse:  func(t *testing.T, res *http.Response) {},
		},
		{
			name: "Missing authentication",
			input: map[string]any{
				"user_id":           userID,
				"regulation_number": fmt.Sprintf("NO_AUTH%d", time.Now().UnixNano()),
				"rank_id":           rankID,
				"posting_id":        postingID,
				"formation_id":      formationID,
				"region_id":         regionID,
			},
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			checkResponse:  func(t *testing.T, res *http.Response) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := createAuthenticatedRequest(http.MethodPost, "/v1/officers", body, tt.token)

			rec := httptest.NewRecorder()

			// Simulate authentication middleware
			if tt.token != "" {
				user, _ := testApp.models.User.GetForToken(data.ScopeAuthentication, tt.token)
				req = setUserContext(req, user)
				testApp.createOfficerHandler(rec, req)
			} else {
				handler := testApp.authenticate(http.HandlerFunc(testApp.createOfficerHandler))
				handler.ServeHTTP(rec, req)
			}

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}

			tt.checkResponse(t, res)
		})
	}
}

func TestShowOfficerHandler(t *testing.T) {
	cleanupAPITestData(t)
	defer cleanupAPITestData(t)

	officer, token := createTestOfficerWithToken(t)

	tests := []struct {
		name           string
		officerID      string
		token          string
		expectedStatus int
	}{
		{
			name:           "Valid officer ID with authentication",
			officerID:      strconv.FormatInt(officer.ID, 10),
			token:          token,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid officer ID",
			officerID:      "999999",
			token:          token,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Non-numeric officer ID",
			officerID:      "abc",
			token:          token,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Missing authentication",
			officerID:      strconv.FormatInt(officer.ID, 10),
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := fmt.Sprintf("/v1/officers/%s", tt.officerID)
			req := createAuthenticatedRequest(http.MethodGet, path, nil, tt.token)

			// Set URL parameter
			req = setURLParam(req, "id", tt.officerID)

			rec := httptest.NewRecorder()

			if tt.token != "" {
				user, _ := testApp.models.User.GetForToken(data.ScopeAuthentication, tt.token)
				req = setUserContext(req, user)
				testApp.showOfficerHandler(rec, req)
			} else {
				handler := testApp.authenticate(http.HandlerFunc(testApp.showOfficerHandler))
				handler.ServeHTTP(rec, req)
			}

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

				officerResponse := response["officer"].(map[string]any)
				if officerResponse["regulation_number"] != officer.RegulationNumber {
					t.Errorf("Expected regulation number %s; got %s",
						officer.RegulationNumber, officerResponse["regulation_number"])
				}
			}
		})
	}
}

func TestShowOfficerWithDetailsHandler(t *testing.T) {
	cleanupAPITestData(t)
	defer cleanupAPITestData(t)

	officer, token := createTestOfficerWithToken(t)

	path := fmt.Sprintf("/v1/officers/%d/details", officer.ID)
	req := createAuthenticatedRequest(http.MethodGet, path, nil, token)

	// Set URL parameter and user context
	req = setURLParam(req, "id", strconv.FormatInt(officer.ID, 10))
	user, _ := testApp.models.User.GetForToken(data.ScopeAuthentication, token)
	req = setUserContext(req, user)

	rec := httptest.NewRecorder()
	testApp.showOfficerWithDetailsHandler(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, res.StatusCode)
	}

	var response map[string]any
	err := json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	officerResponse := response["officer"].(map[string]any)
	if officerResponse["user"] == nil {
		t.Error("Expected user details to be included")
	}
	if officerResponse["rank"] == nil {
		t.Error("Expected rank details to be included")
	}
	if officerResponse["posting"] == nil {
		t.Error("Expected posting details to be included")
	}
	if officerResponse["formation"] == nil {
		t.Error("Expected formation details to be included")
	}
	if officerResponse["region"] == nil {
		t.Error("Expected region details to be included")
	}
}

func TestListOfficersHandler(t *testing.T) {
	cleanupAPITestData(t)
	defer cleanupAPITestData(t)

	// Create multiple test officers
	_, token := createTestUserWithToken(t)

	userID, regionID, formationID, postingID, rankID := createOfficerTestDependencies(t)

	// Create 3 test officers
	for i := 0; i < 3; i++ {
		// Create unique user for each officer
		testUser := &data.User{
			FirstName:   fmt.Sprintf("Officer%d", i),
			LastName:    "Test",
			Email:       fmt.Sprintf("officer%d_%d@example.com", i, time.Now().UnixNano()),
			Gender:      "m",
			IsActivated: true,
			IsOfficer:   true,
		}
		_ = userID
		_ = testUser.Password.Set("TestPass123!")
		_ = testApp.models.User.Insert(testUser)

		officer := &data.Officer{
			UserID:           testUser.ID,
			RegulationNumber: fmt.Sprintf("LIST%d_%d", i, time.Now().UnixNano()),
			RankID:           rankID,
			PostingID:        postingID,
			FormationID:      formationID,
			RegionID:         regionID,
		}
		_ = testApp.models.Officer.Insert(officer)
	}

	tests := []struct {
		name           string
		queryParams    string
		token          string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "Valid list request",
			queryParams:    "",
			token:          token,
			expectedStatus: http.StatusOK,
			expectedCount:  3,
		},
		{
			name:           "With pagination",
			queryParams:    "?page=1&page_size=2",
			token:          token,
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "Filter by rank",
			queryParams:    fmt.Sprintf("?rank_id=%d", rankID),
			token:          token,
			expectedStatus: http.StatusOK,
			expectedCount:  3,
		},
		{
			name:           "Missing authentication",
			queryParams:    "",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/v1/officers" + tt.queryParams
			req := createAuthenticatedRequest(http.MethodGet, path, nil, tt.token)

			rec := httptest.NewRecorder()

			if tt.token != "" {
				user, _ := testApp.models.User.GetForToken(data.ScopeAuthentication, tt.token)
				req = setUserContext(req, user)
				testApp.listOfficersHandler(rec, req)
			} else {
				handler := testApp.authenticate(http.HandlerFunc(testApp.listOfficersHandler))
				handler.ServeHTTP(rec, req)
			}

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

				officers := response["officers"].([]any)
				if len(officers) != tt.expectedCount {
					t.Errorf("Expected %d officers; got %d", tt.expectedCount, len(officers))
				}

				if response["metadata"] == nil {
					t.Error("Expected metadata object in response")
				}
			}
		})
	}
}

func TestUpdateOfficerHandler(t *testing.T) {
	cleanupAPITestData(t)
	defer cleanupAPITestData(t)

	officer, token := createTestOfficerWithToken(t)

	tests := []struct {
		name           string
		officerID      string
		input          map[string]any
		token          string
		expectedStatus int
	}{
		{
			name:      "Valid officer update",
			officerID: strconv.FormatInt(officer.ID, 10),
			input: map[string]any{
				"regulation_number": fmt.Sprintf("UPD%d", time.Now().UnixNano()),
			},
			token:          token,
			expectedStatus: http.StatusOK,
		},
		{
			name:      "Invalid regulation number (too long)",
			officerID: strconv.FormatInt(officer.ID, 10),
			input: map[string]any{
				"regulation_number": "This is a very long regulation number that exceeds the fifty character limit",
			},
			token:          token,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:      "Non-existent officer",
			officerID: "999999",
			input: map[string]any{
				"regulation_number": fmt.Sprintf("NONE%d", time.Now().UnixNano()),
			},
			token:          token,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:      "Missing authentication",
			officerID: strconv.FormatInt(officer.ID, 10),
			input: map[string]any{
				"regulation_number": fmt.Sprintf("NOAUTH%d", time.Now().UnixNano()),
			},
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			path := fmt.Sprintf("/v1/officers/%s", tt.officerID)
			req := createAuthenticatedRequest(http.MethodPatch, path, body, tt.token)

			// Set URL parameter
			req = setURLParam(req, "id", tt.officerID)

			rec := httptest.NewRecorder()

			if tt.token != "" {
				user, _ := testApp.models.User.GetForToken(data.ScopeAuthentication, tt.token)
				req = setUserContext(req, user)
				testApp.updateOfficerHandler(rec, req)
			} else {
				handler := testApp.authenticate(http.HandlerFunc(testApp.updateOfficerHandler))
				handler.ServeHTTP(rec, req)
			}

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}
		})
	}
}

func TestDeleteOfficerHandler(t *testing.T) {
	cleanupAPITestData(t)
	defer cleanupAPITestData(t)

	officer, token := createTestOfficerWithToken(t)

	tests := []struct {
		name           string
		officerID      string
		token          string
		expectedStatus int
	}{
		{
			name:           "Valid officer deletion",
			officerID:      strconv.FormatInt(officer.ID, 10),
			token:          token,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Non-existent officer",
			officerID:      "999999",
			token:          token,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Missing authentication",
			officerID:      strconv.FormatInt(officer.ID, 10),
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := fmt.Sprintf("/v1/officers/%s", tt.officerID)
			req := createAuthenticatedRequest(http.MethodDelete, path, nil, tt.token)

			// Set URL parameter
			req = setURLParam(req, "id", tt.officerID)

			rec := httptest.NewRecorder()

			if tt.token != "" {
				user, _ := testApp.models.User.GetForToken(data.ScopeAuthentication, tt.token)
				req = setUserContext(req, user)
				testApp.deleteOfficerHandler(rec, req)
			} else {
				handler := testApp.authenticate(http.HandlerFunc(testApp.deleteOfficerHandler))
				handler.ServeHTTP(rec, req)
			}

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}
		})
	}
}

func TestGetUserOfficerHandler(t *testing.T) {
	cleanupAPITestData(t)
	defer cleanupAPITestData(t)

	officer, token := createTestOfficerWithToken(t)

	tests := []struct {
		name           string
		userID         string
		token          string
		expectedStatus int
	}{
		{
			name:           "Valid user ID with officer",
			userID:         strconv.FormatInt(officer.UserID, 10),
			token:          token,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "User without officer record",
			userID:         "999999",
			token:          token,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Missing authentication",
			userID:         strconv.FormatInt(officer.UserID, 10),
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := fmt.Sprintf("/v1/users/%s/officer", tt.userID)
			req := createAuthenticatedRequest(http.MethodGet, path, nil, tt.token)

			// Set URL parameter
			req = setURLParam(req, "id", tt.userID)

			rec := httptest.NewRecorder()

			if tt.token != "" {
				user, _ := testApp.models.User.GetForToken(data.ScopeAuthentication, tt.token)
				req = setUserContext(req, user)
				testApp.getUserOfficerHandler(rec, req)
			} else {
				handler := testApp.authenticate(http.HandlerFunc(testApp.getUserOfficerHandler))
				handler.ServeHTTP(rec, req)
			}

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

				officerResponse := response["officer"].(map[string]any)
				if officerResponse["regulation_number"] != officer.RegulationNumber {
					t.Errorf("Expected regulation number %s; got %s",
						officer.RegulationNumber, officerResponse["regulation_number"])
				}
			}
		})
	}
}

// Integration test for complete officer workflow
func TestOfficerWorkflow(t *testing.T) {
	cleanupAPITestData(t)
	defer cleanupAPITestData(t)

	userID, regionID, formationID, postingID, rankID := createOfficerTestDependencies(t)
	_, token := createTestUserWithToken(t)

	// 1. Create an officer
	officerInput := map[string]any{
		"user_id":           userID,
		"regulation_number": fmt.Sprintf("WF%d", time.Now().UnixNano()),
		"rank_id":           rankID,
		"posting_id":        postingID,
		"formation_id":      formationID,
		"region_id":         regionID,
	}

	body, _ := json.Marshal(officerInput)
	req := createAuthenticatedRequest(http.MethodPost, "/v1/officers", body, token)
	user, _ := testApp.models.User.GetForToken(data.ScopeAuthentication, token)
	req = setUserContext(req, user)

	rec := httptest.NewRecorder()
	testApp.createOfficerHandler(rec, req)

	if rec.Result().StatusCode != http.StatusCreated {
		t.Fatalf("Officer creation failed with status %d", rec.Result().StatusCode)
	}

	// Extract officer from response
	var createResponse map[string]any
	_ = json.NewDecoder(rec.Result().Body).Decode(&createResponse)
	officerResponse := createResponse["officer"].(map[string]any)
	officerID := int64(officerResponse["id"].(float64))

	// 2. Get the officer
	req = createAuthenticatedRequest(http.MethodGet, fmt.Sprintf("/v1/officers/%d", officerID), nil, token)
	req = setURLParam(req, "id", strconv.FormatInt(officerID, 10))
	req = setUserContext(req, user)

	rec = httptest.NewRecorder()
	testApp.showOfficerHandler(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("Officer retrieval failed with status %d", rec.Result().StatusCode)
	}

	// 3. Update the officer
	updateInput := map[string]any{
		"regulation_number": fmt.Sprintf("UPD_WF%d", time.Now().UnixNano()),
	}

	body, _ = json.Marshal(updateInput)
	req = createAuthenticatedRequest(http.MethodPatch, fmt.Sprintf("/v1/officers/%d", officerID), body, token)
	req = setURLParam(req, "id", strconv.FormatInt(officerID, 10))
	req = setUserContext(req, user)

	rec = httptest.NewRecorder()
	testApp.updateOfficerHandler(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("Officer update failed with status %d", rec.Result().StatusCode)
	}

	// 4. Get officer with details
	req = createAuthenticatedRequest(http.MethodGet, fmt.Sprintf("/v1/officers/%d/details", officerID), nil, token)
	req = setURLParam(req, "id", strconv.FormatInt(officerID, 10))
	req = setUserContext(req, user)

	rec = httptest.NewRecorder()
	testApp.showOfficerWithDetailsHandler(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("Officer details retrieval failed with status %d", rec.Result().StatusCode)
	}

	// 5. Delete the officer
	req = createAuthenticatedRequest(http.MethodDelete, fmt.Sprintf("/v1/officers/%d", officerID), nil, token)
	req = setURLParam(req, "id", strconv.FormatInt(officerID, 10))
	req = setUserContext(req, user)

	rec = httptest.NewRecorder()
	testApp.deleteOfficerHandler(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("Officer deletion failed with status %d", rec.Result().StatusCode)
	}

	t.Log("Complete officer workflow test passed successfully!")
}
