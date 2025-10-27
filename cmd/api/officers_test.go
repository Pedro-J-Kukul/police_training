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

func getSeededOfficer(t *testing.T) (*data.Officer, *data.User) {
	t.Helper()
	t.Log("Step: Getting seeded officer and user data")

	// Get seeded officer user
	user := getSeededUser(t, "john.smith@police-training.bz")

	// Get officer record for this user
	officer, err := testApp.models.Officer.GetByUserID(user.ID)
	if err != nil {
		t.Fatalf("Failed to get seeded officer for user %s: %v", user.Email, err)
	}

	t.Logf("Step: Retrieved officer ID %d for user %s", officer.ID, user.Email)
	return officer, user
}

func createTestOfficer(t *testing.T) (*data.Officer, *data.User, string) {
	t.Helper()

	// Create test user
	testEmail := fmt.Sprintf("testofficer%d@test-police-training.bz", time.Now().UnixNano())
	userInput := map[string]any{
		"first_name": "Test",
		"last_name":  "Officer",
		"email":      testEmail,
		"gender":     "m",
		"password":   "TestOfficerPass123!",
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

	user, _ := testApp.models.User.Get(userID)

	// Get reference data
	region, _ := testApp.models.Region.GetByName("Northern Region")
	formation, _ := testApp.models.Formation.Get(1)
	posting, _ := testApp.models.Posting.GetByName("Relief")
	rank, _ := testApp.models.Rank.GetByName("Constable")

	// Create officer
	officer := &data.Officer{
		UserID:           userID,
		RegulationNumber: fmt.Sprintf("TEST%d", time.Now().UnixNano()),
		RankID:           rank.ID,
		PostingID:        posting.ID,
		FormationID:      formation.ID,
		RegionID:         region.ID,
	}

	err := testApp.models.Officer.Insert(officer)
	if err != nil {
		t.Fatalf("Failed to create test officer: %v", err)
	}

	token := createTokenForSeededUser(t, userID)

	t.Logf("Step: Created test officer ID %d for user %s", officer.ID, user.Email)
	return officer, user, token
}

func TestCreateOfficerHandler(t *testing.T) {
	t.Log("=== Testing Create Officer Handler ===")

	// Create test user for officer
	testEmail := fmt.Sprintf("newofficer%d@test-police-training.bz", time.Now().UnixNano())
	userInput := map[string]any{
		"first_name": "New",
		"last_name":  "Officer",
		"email":      testEmail,
		"gender":     "f",
		"password":   "NewOfficerPass123!",
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

	defer testApp.models.User.HardDelete(userID) // Cleanup

	// Get reference data
	region, _ := testApp.models.Region.GetByName("Northern Region")
	formation, _ := testApp.models.Formation.Get(1)
	posting, _ := testApp.models.Posting.GetByName("Relief")
	rank, _ := testApp.models.Rank.GetByName("Constable")

	adminUser := getSeededUser(t, "admin1@police-training.bz")
	adminToken := createTokenForSeededUser(t, adminUser.ID)

	tests := []struct {
		name           string
		input          map[string]any
		expectedStatus int
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name: "Valid officer creation",
			input: map[string]any{
				"user_id":           userID,
				"regulation_number": fmt.Sprintf("NEW%d", time.Now().UnixNano()),
				"rank_id":           rank.ID,
				"posting_id":        posting.ID,
				"formation_id":      formation.ID,
				"region_id":         region.ID,
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating successful officer creation")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response["officer"] == nil {
					t.Error("Expected officer object in response")
					return
				}

				officer := response["officer"].(map[string]any)
				officerID := int64(officer["id"].(float64))
				t.Logf("Step: Created officer ID %d", officerID)

				// Cleanup created officer
				testApp.models.Officer.Delete(officerID)
			},
		},
		{
			name: "Missing required fields",
			input: map[string]any{
				"user_id": userID,
			},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating missing fields validation")
			},
		},
		{
			name: "Invalid user_id",
			input: map[string]any{
				"user_id":           99999,
				"regulation_number": fmt.Sprintf("INV%d", time.Now().UnixNano()),
				"rank_id":           rank.ID,
				"posting_id":        posting.ID,
				"formation_id":      formation.ID,
				"region_id":         region.ID,
			},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating invalid user_id handling")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test: %s", tt.name)

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/v1/officers", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))

			req = setUserContext(req, adminUser)

			rec := httptest.NewRecorder()
			testApp.createOfficerHandler(rec, req)

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

func TestShowOfficerHandler(t *testing.T) {
	t.Log("=== Testing Show Officer Handler ===")

	seededOfficer, seededUser := getSeededOfficer(t)
	seededToken := createTokenForSeededUser(t, seededUser.ID)

	tests := []struct {
		name           string
		officerID      string
		expectedStatus int
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name:           "Valid seeded officer ID",
			officerID:      strconv.FormatInt(seededOfficer.ID, 10),
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating officer details response")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response["officer"] == nil {
					t.Error("Expected officer object in response")
					return
				}

				officer := response["officer"].(map[string]any)
				t.Logf("Step: Retrieved officer with regulation number %v", officer["regulation_number"])

				requiredFields := []string{"id", "user_id", "regulation_number", "rank_id", "posting_id", "formation_id", "region_id"}
				for _, field := range requiredFields {
					if officer[field] == nil {
						t.Errorf("Expected field %s in officer response", field)
					}
				}
			},
		},
		{
			name:           "Non-existent officer ID",
			officerID:      "999999",
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating non-existent officer response")
			},
		},
		{
			name:           "Invalid officer ID format",
			officerID:      "abc",
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating invalid ID format response")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test: %s with officer ID %s", tt.name, tt.officerID)

			path := fmt.Sprintf("/v1/officers/%s", tt.officerID)
			req := httptest.NewRequest(http.MethodGet, path, nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", seededToken))

			req = setURLParam(req, "id", tt.officerID)
			req = setUserContext(req, seededUser)

			t.Logf("Step: Set URL param id=%s", tt.officerID)

			rec := httptest.NewRecorder()
			testApp.showOfficerHandler(rec, req)

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

func TestShowOfficerWithDetailsHandler(t *testing.T) {
	t.Log("=== Testing Show Officer With Details Handler ===")

	seededOfficer, seededUser := getSeededOfficer(t)
	seededToken := createTokenForSeededUser(t, seededUser.ID)

	tests := []struct {
		name           string
		officerID      string
		expectedStatus int
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name:           "Valid seeded officer ID with details",
			officerID:      strconv.FormatInt(seededOfficer.ID, 10),
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating officer with details response")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response["officer"] == nil {
					t.Error("Expected officer object in response")
					return
				}

				officer := response["officer"].(map[string]any)
				t.Logf("Step: Retrieved detailed officer with regulation number %v", officer["regulation_number"])

				// Check for detailed fields that should be included
				detailFields := []string{"user", "rank", "posting", "formation", "region"}
				for _, field := range detailFields {
					if officer[field] == nil {
						t.Errorf("Expected detailed field %s in officer response", field)
					} else {
						t.Logf("Step: Found detailed field %s", field)
					}
				}
			},
		},
		{
			name:           "Non-existent officer ID",
			officerID:      "999999",
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating non-existent officer details response")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test: %s", tt.name)

			path := fmt.Sprintf("/v1/officers/%s/details", tt.officerID)
			req := httptest.NewRequest(http.MethodGet, path, nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", seededToken))

			req = setURLParam(req, "id", tt.officerID)
			req = setUserContext(req, seededUser)

			rec := httptest.NewRecorder()
			testApp.showOfficerWithDetailsHandler(rec, req)

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

func TestListOfficersHandler(t *testing.T) {
	t.Log("=== Testing List Officers Handler ===")

	adminUser := getSeededUser(t, "admin1@police-training.bz")
	adminToken := createTokenForSeededUser(t, adminUser.ID)

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name:           "List all officers",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating officers list response")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response["officers"] == nil {
					t.Error("Expected officers array in response")
					return
				}

				officers := response["officers"].([]any)
				t.Logf("Step: Retrieved %d officers", len(officers))

				// Should have seeded officers (6 from populate script)
				if len(officers) < 3 {
					t.Errorf("Expected at least 3 officers, got %d", len(officers))
				}

				// Verify first officer structure
				if len(officers) > 0 {
					officer := officers[0].(map[string]any)
					requiredFields := []string{"id", "user_id", "regulation_number", "rank_id", "posting_id", "formation_id", "region_id"}
					for _, field := range requiredFields {
						if officer[field] == nil {
							t.Errorf("Expected field %s in officer response", field)
						}
					}

					t.Logf("Step: First officer ID %v, Reg Number: %v", officer["id"], officer["regulation_number"])
				}

				if response["metadata"] == nil {
					t.Error("Expected metadata in response")
				}
			},
		},
		{
			name:           "List officers with pagination",
			queryParams:    "?page=1&page_size=3",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating paginated officers response")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				officers := response["officers"].([]any)
				t.Logf("Step: Retrieved %d officers with page_size=3", len(officers))

				if len(officers) > 3 {
					t.Errorf("Expected max 3 officers, got %d", len(officers))
				}
			},
		},
		{
			name:           "Filter officers by rank",
			queryParams:    "?rank_id=1", // Assuming rank ID 1 exists
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating rank filter response")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				officers := response["officers"].([]any)
				t.Logf("Step: Retrieved %d officers for rank 1", len(officers))

				// Verify all officers have the specified rank
				for i, officerInterface := range officers {
					officer := officerInterface.(map[string]any)
					if rankID, ok := officer["rank_id"].(float64); !ok || rankID != 1 {
						t.Errorf("Officer at index %d does not have rank_id 1", i)
					}
				}
			},
		},
		{
			name:           "Search officers by regulation number",
			queryParams:    "?regulation_number=PC001", // From seeded data
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating regulation number search")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				officers := response["officers"].([]any)
				t.Logf("Step: Retrieved %d officers with regulation number PC001", len(officers))

				// Should find exactly one officer with PC001
				if len(officers) == 1 {
					officer := officers[0].(map[string]any)
					if officer["regulation_number"] != "PC001" {
						t.Errorf("Expected regulation_number PC001, got %v", officer["regulation_number"])
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test: %s", tt.name)

			path := "/v1/officers" + tt.queryParams
			req := httptest.NewRequest(http.MethodGet, path, nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))
			req = setUserContext(req, adminUser)

			t.Logf("Step: Making request to %s", path)

			rec := httptest.NewRecorder()
			testApp.listOfficersHandler(rec, req)

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

func TestUpdateOfficerHandler(t *testing.T) {
	t.Log("=== Testing Update Officer Handler ===")

	// Create test officer for updating
	testOfficer, testUser, testToken := createTestOfficer(t)
	defer testApp.models.User.HardDelete(testUser.ID) // Cleanup

	// Get different rank for update
	newRank, _ := testApp.models.Rank.GetByName("Corporal")
	newPosting, _ := testApp.models.Posting.GetByName("Station Manager")

	tests := []struct {
		name           string
		officerID      string
		input          map[string]any
		expectedStatus int
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name:      "Valid officer update",
			officerID: strconv.FormatInt(testOfficer.ID, 10),
			input: map[string]any{
				"regulation_number": fmt.Sprintf("UPD%d", time.Now().UnixNano()),
				"rank_id":           newRank.ID,
				"posting_id":        newPosting.ID,
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating successful officer update")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response["officer"] == nil {
					t.Error("Expected officer object in response")
					return
				}

				officer := response["officer"].(map[string]any)
				t.Logf("Step: Updated officer regulation number to %v", officer["regulation_number"])
			},
		},
		{
			name:      "Non-existent officer",
			officerID: "999999",
			input: map[string]any{
				"regulation_number": fmt.Sprintf("NONE%d", time.Now().UnixNano()),
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating non-existent officer update")
			},
		},
		{
			name:      "Invalid rank_id",
			officerID: strconv.FormatInt(testOfficer.ID, 10),
			input: map[string]any{
				"rank_id": 999999,
			},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating invalid rank_id handling")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test: %s", tt.name)

			body, _ := json.Marshal(tt.input)
			path := fmt.Sprintf("/v1/officers/%s", tt.officerID)
			req := httptest.NewRequest(http.MethodPatch, path, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testToken))

			req = setURLParam(req, "id", tt.officerID)
			req = setUserContext(req, testUser)

			rec := httptest.NewRecorder()
			testApp.updateOfficerHandler(rec, req)

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

func TestDeleteOfficerHandler(t *testing.T) {
	t.Log("=== Testing Delete Officer Handler ===")

	// Create test officer for deletion
	testOfficer, testUser, testToken := createTestOfficer(t)
	defer testApp.models.User.HardDelete(testUser.ID) // Cleanup user

	tests := []struct {
		name           string
		officerID      string
		expectedStatus int
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name:           "Valid officer deletion",
			officerID:      strconv.FormatInt(testOfficer.ID, 10),
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating successful officer deletion")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response["message"] == nil {
					t.Error("Expected message in response")
				}

				t.Logf("Step: Officer deleted with message: %v", response["message"])
			},
		},
		{
			name:           "Non-existent officer",
			officerID:      "999999",
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating non-existent officer deletion")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test: %s", tt.name)

			path := fmt.Sprintf("/v1/officers/%s", tt.officerID)
			req := httptest.NewRequest(http.MethodDelete, path, nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testToken))

			req = setURLParam(req, "id", tt.officerID)
			req = setUserContext(req, testUser)

			rec := httptest.NewRecorder()
			testApp.deleteOfficerHandler(rec, req)

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

func TestGetUserOfficerHandler(t *testing.T) {
	t.Log("=== Testing Get User Officer Handler ===")

	seededOfficer, seededUser := getSeededOfficer(t)
	seededToken := createTokenForSeededUser(t, seededUser.ID)

	tests := []struct {
		name           string
		userID         string
		expectedStatus int
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name:           "Valid user ID with officer",
			userID:         strconv.FormatInt(seededUser.ID, 10),
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating user officer response")
				var response map[string]any
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response["officer"] == nil {
					t.Error("Expected officer object in response")
					return
				}

				officer := response["officer"].(map[string]any)
				if officer["regulation_number"] != seededOfficer.RegulationNumber {
					t.Errorf("Expected regulation number %s; got %s",
						seededOfficer.RegulationNumber, officer["regulation_number"])
				}
				t.Logf("Step: Set user context for user ID %d", seededUser.ID)
			},
		},
		{
			name:           "User without officer record",
			userID:         "999999",
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating user without officer record")
			},
		},
		{
			name:           "Invalid user ID format",
			userID:         "abc",
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, res *http.Response) {
				t.Log("Step: Validating invalid user ID format")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test: %s", tt.name)

			path := fmt.Sprintf("/v1/users/%s/officer", tt.userID)
			req := httptest.NewRequest(http.MethodGet, path, nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", seededToken))

			req = setURLParam(req, "id", tt.userID)
			req = setUserContext(req, seededUser)

			rec := httptest.NewRecorder()
			testApp.getUserOfficerHandler(rec, req)

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

func TestOfficerWorkflow(t *testing.T) {
	t.Log("=== Testing Complete Officer Workflow ===")

	// Create test user
	testEmail := fmt.Sprintf("workflow%d@test-police-training.bz", time.Now().UnixNano())
	userInput := map[string]any{
		"first_name": "Workflow",
		"last_name":  "Officer",
		"email":      testEmail,
		"gender":     "m",
		"password":   "WorkflowPass123!",
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

	defer testApp.models.User.HardDelete(userID) // Cleanup

	user, _ := testApp.models.User.Get(userID)
	userToken := createTokenForSeededUser(t, userID)

	// Get reference data
	region, _ := testApp.models.Region.GetByName("Northern Region")
	formation, _ := testApp.models.Formation.Get(1)
	posting, _ := testApp.models.Posting.GetByName("Relief")
	rank, _ := testApp.models.Rank.GetByName("Constable")

	// 1. Create officer
	officerInput := map[string]any{
		"user_id":           userID,
		"regulation_number": fmt.Sprintf("WF%d", time.Now().UnixNano()),
		"rank_id":           rank.ID,
		"posting_id":        posting.ID,
		"formation_id":      formation.ID,
		"region_id":         region.ID,
	}

	body, _ = json.Marshal(officerInput)
	req = httptest.NewRequest(http.MethodPost, "/v1/officers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userToken))
	req = setUserContext(req, user)

	rec = httptest.NewRecorder()
	testApp.createOfficerHandler(rec, req)

	if rec.Result().StatusCode != http.StatusCreated {
		t.Fatalf("Officer creation failed with status %d", rec.Result().StatusCode)
	}

	var officerCreateResponse map[string]any
	_ = json.NewDecoder(rec.Result().Body).Decode(&officerCreateResponse)
	officerResponse := officerCreateResponse["officer"].(map[string]any)
	officerID := int64(officerResponse["id"].(float64))

	t.Logf("Step: Created officer ID %d", officerID)

	// 2. Get the officer
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/officers/%d", officerID), nil)
	req = setURLParam(req, "id", strconv.FormatInt(officerID, 10))

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userToken))
	req = setURLParam(req, "id", strconv.FormatInt(userID, 10)) // Use userID not officerID
	req = setUserContext(req, user)

	rec = httptest.NewRecorder()
	testApp.showOfficerHandler(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("Officer retrieval failed with status %d", rec.Result().StatusCode)
	}

	t.Logf("Step: Successfully retrieved officer")

	// 3. Get officer with details
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/officers/%d/details", officerID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userToken))
	req = setURLParam(req, "id", strconv.FormatInt(userID, 10)) // Use userID not officerID
	req = setUserContext(req, user)

	rec = httptest.NewRecorder()
	testApp.showOfficerWithDetailsHandler(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("Officer details retrieval failed with status %d", rec.Result().StatusCode)
	}

	t.Logf("Step: Successfully retrieved officer details")

	// 4. Update the officer
	updateRank, _ := testApp.models.Rank.GetByName("Corporal")
	updateInput := map[string]any{
		"regulation_number": fmt.Sprintf("UPD_WF%d", time.Now().UnixNano()),
		"rank_id":           updateRank.ID,
	}

	body, _ = json.Marshal(updateInput)
	req = httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/officers/%d", officerID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userToken))
	req = setURLParam(req, "id", strconv.FormatInt(userID, 10))
	req = setUserContext(req, user)

	rec = httptest.NewRecorder()
	testApp.updateOfficerHandler(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("Officer update failed with status %d", rec.Result().StatusCode)
	}

	t.Logf("Step: Successfully updated officer")

	// 5. Get officer by user ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/users/%d/officer", userID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userToken))
	req = setURLParam(req, "user_id", strconv.FormatInt(userID, 10))
	req = setUserContext(req, user)

	rec = httptest.NewRecorder()
	testApp.getUserOfficerHandler(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("Get user officer failed with status %d", rec.Result().StatusCode)
	}

	t.Logf("Step: Successfully retrieved officer by user ID")

	// 6. Delete the officer
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/officers/%d", officerID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userToken))
	req = setURLParam(req, "id", strconv.FormatInt(userID, 10))
	req = setUserContext(req, user)

	rec = httptest.NewRecorder()
	testApp.deleteOfficerHandler(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("Officer deletion failed with status %d", rec.Result().StatusCode)
	}

	t.Logf("Step: Successfully deleted officer")

	t.Log("Step: Complete officer workflow test passed successfully!")
}
