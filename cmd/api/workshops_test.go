package main

// package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"net/url"
// 	"strconv"
// 	"testing"
// 	"time"

// 	"github.com/Pedro-J-Kukul/police_training/internal/data"
// )

// func getSeededWorkshopData(t *testing.T) (categoryID, typeID int64) {
// 	t.Helper()
// 	t.Log("Step: Getting seeded workshop reference data")

// 	// Get first available category
// 	categories, _, err := testApp.models.TrainingCategory.GetAll("", nil, data.Filters{Page: 10, PageSize: 1, Sort: "id", SortSafelist: []string{"id"}})
// 	if err != nil || len(categories) == 0 {
// 		t.Fatal("No training categories found in seeded data")
// 	}

// 	// Get first available type
// 	types, _, err := testApp.models.TrainingType.GetAll("", data.Filters{Page: 1, PageSize: 10, Sort: "id", SortSafelist: []string{"id"}})
// 	if err != nil || len(types) == 0 {
// 		t.Fatal("No training types found in seeded data")
// 	}

// 	categoryID = categories[0].ID
// 	typeID = types[0].ID

// 	t.Logf("Step: Using category ID %d and type ID %d", categoryID, typeID)
// 	return categoryID, typeID
// }

// func getSeededWorkshop(t *testing.T) *data.Workshop {
// 	t.Helper()
// 	t.Log("Step: Getting seeded workshop")

// 	workshops, _, err := testApp.models.Workshop.GetAll("", nil, nil, nil, data.Filters{Page: 1, PageSize: 1, Sort: "id", SortSafelist: []string{"id"}})
// 	if err != nil || len(workshops) == 0 {
// 		t.Fatal("No workshops found in seeded data")
// 	}

// 	workshop := workshops[0]
// 	t.Logf("Step: Using seeded workshop ID %d: %s", workshop.ID, workshop.WorkshopName)
// 	return workshop
// }

// func createTestWorkshop(t *testing.T) *data.Workshop {
// 	t.Helper()

// 	categoryID, typeID := getSeededWorkshopData(t)

// 	workshop := &data.Workshop{
// 		WorkshopName: fmt.Sprintf("TEST_Workshop_%d", time.Now().UnixNano()),
// 		CategoryID:   categoryID,
// 		TypeID:       typeID,
// 		CreditHours:  40,
// 		Description:  stringPtr("Test workshop for API testing"),
// 		IsActive:     true,
// 	}

// 	err := testApp.models.Workshop.Insert(workshop)
// 	if err != nil {
// 		t.Fatalf("Failed to create test workshop: %v", err)
// 	}

// 	t.Logf("Step: Created test workshop ID %d: %s", workshop.ID, workshop.WorkshopName)
// 	return workshop
// }

// func stringPtr(s string) *string {
// 	return &s
// }

// func TestCreateWorkshopHandler(t *testing.T) {
// 	t.Log("=== Testing Create Workshop Handler ===")

// 	categoryID, typeID := getSeededWorkshopData(t)
// 	adminUser := getSeededUser(t, "admin1@police-training.bz")
// 	adminToken := createTokenForSeededUser(t, adminUser.ID)

// 	tests := []struct {
// 		name           string
// 		input          map[string]any
// 		expectedStatus int
// 		checkResponse  func(*testing.T, *http.Response)
// 	}{
// 		{
// 			name: "Valid workshop creation",
// 			input: map[string]any{
// 				"workshop_name": fmt.Sprintf("NEW_Workshop_%d", time.Now().UnixNano()),
// 				"category_id":   categoryID,
// 				"type_id":       typeID,
// 				"credit_hours":  60,
// 				"description":   "New workshop created via API",
// 				"is_active":     true,
// 			},
// 			expectedStatus: http.StatusCreated,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating successful workshop creation")
// 				var response map[string]any
// 				err := json.NewDecoder(res.Body).Decode(&response)
// 				if err != nil {
// 					t.Fatalf("Failed to decode response: %v", err)
// 				}

// 				if response["workshop"] == nil {
// 					t.Error("Expected workshop object in response")
// 					return
// 				}

// 				workshop := response["workshop"].(map[string]any)
// 				workshopID := int64(workshop["id"].(float64))
// 				t.Logf("Step: Created workshop ID %d: %v", workshopID, workshop["workshop_name"])

// 				// Cleanup created workshop
// 				testApp.models.Workshop.Delete(workshopID)
// 				t.Logf("Step: Cleaned up test workshop ID %d", workshopID)
// 			},
// 		},
// 		{
// 			name: "Missing required fields",
// 			input: map[string]any{
// 				"workshop_name": "Incomplete Workshop",
// 			},
// 			expectedStatus: http.StatusUnprocessableEntity,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating missing fields validation")
// 			},
// 		},
// 		{
// 			name: "Invalid category_id",
// 			input: map[string]any{
// 				"workshop_name": fmt.Sprintf("Invalid_Category_%d", time.Now().UnixNano()),
// 				"category_id":   999999,
// 				"type_id":       typeID,
// 				"credit_hours":  40,
// 			},
// 			expectedStatus: http.StatusBadRequest,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating invalid category_id handling")
// 			},
// 		},
// 		{
// 			name: "Duplicate workshop name",
// 			input: map[string]any{
// 				"workshop_name": getSeededWorkshop(t).WorkshopName,
// 				"category_id":   categoryID,
// 				"type_id":       typeID,
// 				"credit_hours":  40,
// 			},
// 			expectedStatus: http.StatusUnprocessableEntity,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating duplicate workshop name handling")
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Logf("Starting test: %s", tt.name)

// 			body, _ := json.Marshal(tt.input)
// 			req := httptest.NewRequest(http.MethodPost, "/v1/workshops", bytes.NewReader(body))
// 			req.Header.Set("Content-Type", "application/json")
// 			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))

// 			req = setUserContext(req, adminUser)

// 			rec := httptest.NewRecorder()
// 			testApp.createWorkshopHandler(rec, req)

// 			res := rec.Result()
// 			defer res.Body.Close()

// 			t.Logf("Step: Received status code %d", res.StatusCode)
// 			if res.StatusCode != tt.expectedStatus {
// 				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
// 			}

// 			tt.checkResponse(t, res)
// 			t.Logf("Completed test: %s", tt.name)
// 		})
// 	}
// }

// func TestShowWorkshopHandler(t *testing.T) {
// 	t.Log("=== Testing Show Workshop Handler ===")

// 	seededWorkshop := getSeededWorkshop(t)
// 	facilitatorUser := getSeededUser(t, "maria.rodriguez@police-training.bz")
// 	facilitatorToken := createTokenForSeededUser(t, facilitatorUser.ID)

// 	tests := []struct {
// 		name           string
// 		workshopID     string
// 		expectedStatus int
// 		checkResponse  func(*testing.T, *http.Response)
// 	}{
// 		{
// 			name:           "Valid seeded workshop ID",
// 			workshopID:     strconv.FormatInt(seededWorkshop.ID, 10),
// 			expectedStatus: http.StatusOK,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating workshop details response")
// 				var response map[string]any
// 				err := json.NewDecoder(res.Body).Decode(&response)
// 				if err != nil {
// 					t.Fatalf("Failed to decode response: %v", err)
// 				}

// 				if response["workshop"] == nil {
// 					t.Error("Expected workshop object in response")
// 					return
// 				}

// 				workshop := response["workshop"].(map[string]any)
// 				if workshop["id"] != float64(seededWorkshop.ID) {
// 					t.Errorf("Expected workshop ID %d, got %v", seededWorkshop.ID, workshop["id"])
// 				}

// 				t.Logf("Step: Retrieved workshop: %v", workshop["workshop_name"])

// 				// Verify expected fields
// 				requiredFields := []string{"id", "workshop_name", "category_id", "type_id", "credit_hours", "is_active"}
// 				for _, field := range requiredFields {
// 					if workshop[field] == nil {
// 						t.Errorf("Expected field %s in workshop response", field)
// 					}
// 				}
// 			},
// 		},
// 		{
// 			name:           "Non-existent workshop ID",
// 			workshopID:     "999999",
// 			expectedStatus: http.StatusNotFound,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating non-existent workshop response")
// 			},
// 		},
// 		{
// 			name:           "Invalid workshop ID format",
// 			workshopID:     "abc",
// 			expectedStatus: http.StatusNotFound,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating invalid ID format response")
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Logf("Starting test: %s with workshop ID %s", tt.name, tt.workshopID)

// 			path := fmt.Sprintf("/v1/workshops/%s", tt.workshopID)
// 			req := httptest.NewRequest(http.MethodGet, path, nil)
// 			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", facilitatorToken))

// 			req = setURLParam(req, "id", tt.workshopID)
// 			req = setUserContext(req, facilitatorUser)

// 			t.Logf("Step: Set URL param id=%s", tt.workshopID)

// 			rec := httptest.NewRecorder()
// 			testApp.showWorkshopHandler(rec, req)

// 			res := rec.Result()
// 			defer res.Body.Close()

// 			t.Logf("Step: Received status code %d", res.StatusCode)
// 			if res.StatusCode != tt.expectedStatus {
// 				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
// 			}

// 			tt.checkResponse(t, res)
// 			t.Logf("Completed test: %s", tt.name)
// 		})
// 	}
// }

// func TestListWorkshopsHandler(t *testing.T) {
// 	t.Log("=== Testing List Workshops Handler ===")

// 	categoryID, _ := getSeededWorkshopData(t)
// 	adminUser := getSeededUser(t, "admin1@police-training.bz")
// 	adminToken := createTokenForSeededUser(t, adminUser.ID)

// 	tests := []struct {
// 		name           string
// 		queryParams    string
// 		expectedStatus int
// 		checkResponse  func(*testing.T, *http.Response)
// 	}{
// 		{
// 			name:           "List all workshops",
// 			queryParams:    "",
// 			expectedStatus: http.StatusOK,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating workshops list response")
// 				var response map[string]any
// 				err := json.NewDecoder(res.Body).Decode(&response)
// 				if err != nil {
// 					t.Fatalf("Failed to decode response: %v", err)
// 				}

// 				if response["workshops"] == nil {
// 					t.Error("Expected workshops array in response")
// 					return
// 				}

// 				workshops := response["workshops"].([]any)
// 				t.Logf("Step: Retrieved %d workshops", len(workshops))

// 				// Should have seeded workshops
// 				if len(workshops) == 0 {
// 					t.Error("Expected at least some seeded workshops")
// 				}

// 				// Verify first workshop structure
// 				if len(workshops) > 0 {
// 					workshop := workshops[0].(map[string]any)
// 					requiredFields := []string{"id", "workshop_name", "category_id", "type_id", "credit_hours"}
// 					for _, field := range requiredFields {
// 						if workshop[field] == nil {
// 							t.Errorf("Expected field %s in workshop response", field)
// 						}
// 					}

// 					t.Logf("Step: First workshop ID %v, Name: %v", workshop["id"], workshop["workshop_name"])
// 				}

// 				if response["metadata"] == nil {
// 					t.Error("Expected metadata in response")
// 				}
// 			},
// 		},
// 		{
// 			name:           "List workshops with pagination",
// 			queryParams:    "?page=1&page_size=2",
// 			expectedStatus: http.StatusOK,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating paginated workshops response")
// 				var response map[string]any
// 				err := json.NewDecoder(res.Body).Decode(&response)
// 				if err != nil {
// 					t.Fatalf("Failed to decode response: %v", err)
// 				}

// 				workshops := response["workshops"].([]any)
// 				t.Logf("Step: Retrieved %d workshops with page_size=2", len(workshops))

// 				if len(workshops) > 2 {
// 					t.Errorf("Expected max 2 workshops, got %d", len(workshops))
// 				}
// 			},
// 		},
// 		{
// 			name:           "Filter workshops by category",
// 			queryParams:    fmt.Sprintf("?category_id=%d", categoryID),
// 			expectedStatus: http.StatusOK,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating category filter response")
// 				var response map[string]any
// 				err := json.NewDecoder(res.Body).Decode(&response)
// 				if err != nil {
// 					t.Fatalf("Failed to decode response: %v", err)
// 				}

// 				workshops := response["workshops"].([]any)
// 				t.Logf("Step: Retrieved %d workshops for category %d", len(workshops), categoryID)

// 				// Verify all workshops belong to the specified category
// 				for i, workshopInterface := range workshops {
// 					workshop := workshopInterface.(map[string]any)
// 					if catID, ok := workshop["category_id"].(float64); !ok || catID != float64(categoryID) {
// 						t.Errorf("Workshop at index %d does not belong to category %d", i, categoryID)
// 					}
// 				}
// 			},
// 		},
// 		{
// 			name:           "Filter workshops by active status",
// 			queryParams:    "?is_active=true",
// 			expectedStatus: http.StatusOK,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating active status filter response")
// 				var response map[string]any
// 				err := json.NewDecoder(res.Body).Decode(&response)
// 				if err != nil {
// 					t.Fatalf("Failed to decode response: %v", err)
// 				}

// 				workshops := response["workshops"].([]any)
// 				t.Logf("Step: Retrieved %d active workshops", len(workshops))

// 				// Verify all workshops are active
// 				for i, workshopInterface := range workshops {
// 					workshop := workshopInterface.(map[string]any)
// 					if isActive, ok := workshop["is_active"].(bool); !ok || !isActive {
// 						t.Errorf("Workshop at index %d is not active", i)
// 					}
// 				}
// 			},
// 		},
// 		{
// 			name:           "Search workshops by name",
// 			queryParams:    "?workshop_name=" + getSeededWorkshop(t).WorkshopName,
// 			expectedStatus: http.StatusOK,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating workshop name search")
// 				var response map[string]any
// 				err := json.NewDecoder(res.Body).Decode(&response)
// 				if err != nil {
// 					t.Fatalf("Failed to decode response: %v", err)
// 				}

// 				workshops := response["workshops"].([]any)
// 				t.Logf("Step: Retrieved %d workshops for name search", len(workshops))
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Logf("Starting test: %s", tt.name)

// 			path := "/v1/workshops?workshop_name=" + url.QueryEscape(getSeededWorkshop(t).WorkshopName)
// 			req := httptest.NewRequest(http.MethodGet, path, nil)
// 			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))
// 			req = setUserContext(req, adminUser)

// 			t.Logf("Step: Making request to %s", path)

// 			rec := httptest.NewRecorder()
// 			testApp.listWorkshopsHandler(rec, req)

// 			res := rec.Result()
// 			defer res.Body.Close()

// 			t.Logf("Step: Received status code %d", res.StatusCode)
// 			if res.StatusCode != tt.expectedStatus {
// 				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
// 			}

// 			tt.checkResponse(t, res)
// 			t.Logf("Completed test: %s", tt.name)
// 		})
// 	}
// }

// func TestUpdateWorkshopHandler(t *testing.T) {
// 	t.Log("=== Testing Update Workshop Handler ===")

// 	// Create test workshop for updating
// 	testWorkshop := createTestWorkshop(t)
// 	defer testApp.models.Workshop.Delete(testWorkshop.ID) // Cleanup

// 	adminUser := getSeededUser(t, "admin1@police-training.bz")
// 	adminToken := createTokenForSeededUser(t, adminUser.ID)

// 	tests := []struct {
// 		name           string
// 		workshopID     string
// 		input          map[string]any
// 		expectedStatus int
// 		checkResponse  func(*testing.T, *http.Response)
// 	}{
// 		{
// 			name:       "Valid workshop update",
// 			workshopID: strconv.FormatInt(testWorkshop.ID, 10),
// 			input: map[string]any{
// 				"workshop_name": fmt.Sprintf("UPDATED_Workshop_%d", time.Now().UnixNano()),
// 				"credit_hours":  80,
// 				"description":   "Updated workshop description",
// 			},
// 			expectedStatus: http.StatusOK,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating successful workshop update")
// 				var response map[string]any
// 				err := json.NewDecoder(res.Body).Decode(&response)
// 				if err != nil {
// 					t.Fatalf("Failed to decode response: %v", err)
// 				}

// 				if response["workshop"] == nil {
// 					t.Error("Expected workshop object in response")
// 					return
// 				}

// 				workshop := response["workshop"].(map[string]any)
// 				if workshop["credit_hours"] != float64(80) {
// 					t.Errorf("Expected credit_hours 80, got %v", workshop["credit_hours"])
// 				}

// 				t.Logf("Step: Updated workshop name to %v", workshop["workshop_name"])
// 			},
// 		},
// 		{
// 			name:       "Non-existent workshop",
// 			workshopID: "999999",
// 			input: map[string]any{
// 				"workshop_name": "Non-existent Workshop",
// 			},
// 			expectedStatus: http.StatusNotFound,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating non-existent workshop update")
// 			},
// 		},
// 		{
// 			name:       "Invalid category_id",
// 			workshopID: strconv.FormatInt(testWorkshop.ID, 10),
// 			input: map[string]any{
// 				"category_id": 999999,
// 			},
// 			expectedStatus: http.StatusBadRequest,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating invalid category_id handling")
// 			},
// 		},
// 		{
// 			name:       "Update to duplicate name",
// 			workshopID: strconv.FormatInt(testWorkshop.ID, 10),
// 			input: map[string]any{
// 				"workshop_name": getSeededWorkshop(t).WorkshopName,
// 			},
// 			expectedStatus: http.StatusUnprocessableEntity,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating duplicate name handling")
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Logf("Starting test: %s", tt.name)

// 			body, _ := json.Marshal(tt.input)
// 			path := fmt.Sprintf("/v1/workshops/%s", tt.workshopID)
// 			req := httptest.NewRequest(http.MethodPatch, path, bytes.NewReader(body))
// 			req.Header.Set("Content-Type", "application/json")
// 			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))

// 			req = setURLParam(req, "id", tt.workshopID)
// 			req = setUserContext(req, adminUser)

// 			rec := httptest.NewRecorder()
// 			testApp.updateWorkshopHandler(rec, req)

// 			res := rec.Result()
// 			defer res.Body.Close()

// 			t.Logf("Step: Received status code %d", res.StatusCode)
// 			if res.StatusCode != tt.expectedStatus {
// 				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
// 			}

// 			tt.checkResponse(t, res)
// 			t.Logf("Completed test: %s", tt.name)
// 		})
// 	}
// }

// func TestDeleteWorkshopHandler(t *testing.T) {
// 	t.Log("=== Testing Delete Workshop Handler ===")

// 	// Create test workshop for deletion
// 	testWorkshop := createTestWorkshop(t)
// 	// No defer cleanup since the test will delete it

// 	adminUser := getSeededUser(t, "admin1@police-training.bz")
// 	adminToken := createTokenForSeededUser(t, adminUser.ID)

// 	tests := []struct {
// 		name           string
// 		workshopID     string
// 		expectedStatus int
// 		checkResponse  func(*testing.T, *http.Response)
// 	}{
// 		{
// 			name:           "Valid workshop deletion",
// 			workshopID:     strconv.FormatInt(testWorkshop.ID, 10),
// 			expectedStatus: http.StatusOK,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating successful workshop deletion")
// 				var response map[string]any
// 				err := json.NewDecoder(res.Body).Decode(&response)
// 				if err != nil {
// 					t.Fatalf("Failed to decode response: %v", err)
// 				}

// 				if response["message"] == nil {
// 					t.Error("Expected message in response")
// 				}

// 				t.Logf("Step: Workshop deleted with message: %v", response["message"])
// 			},
// 		},
// 		{
// 			name:           "Non-existent workshop",
// 			workshopID:     "999999",
// 			expectedStatus: http.StatusNotFound,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating non-existent workshop deletion")
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Logf("Starting test: %s", tt.name)

// 			path := fmt.Sprintf("/v1/workshops/%s", tt.workshopID)
// 			req := httptest.NewRequest(http.MethodDelete, path, nil)
// 			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))

// 			req = setURLParam(req, "id", tt.workshopID)
// 			req = setUserContext(req, adminUser)

// 			rec := httptest.NewRecorder()
// 			testApp.deleteWorkshopHandler(rec, req)

// 			res := rec.Result()
// 			defer res.Body.Close()

// 			t.Logf("Step: Received status code %d", res.StatusCode)
// 			if res.StatusCode != tt.expectedStatus {
// 				t.Errorf("Expected status %d; got %d", tt.expectedStatus, res.StatusCode)
// 			}

// 			tt.checkResponse(t, res)
// 			t.Logf("Completed test: %s", tt.name)
// 		})
// 	}
// }

// func TestWorkshopWorkflow(t *testing.T) {
// 	t.Log("=== Testing Complete Workshop Workflow ===")

// 	categoryID, typeID := getSeededWorkshopData(t)
// 	adminUser := getSeededUser(t, "admin1@police-training.bz")
// 	adminToken := createTokenForSeededUser(t, adminUser.ID)

// 	// 1. Create workshop
// 	workshopInput := map[string]any{
// 		"workshop_name": fmt.Sprintf("WORKFLOW_Workshop_%d", time.Now().UnixNano()),
// 		"category_id":   categoryID,
// 		"type_id":       typeID,
// 		"credit_hours":  50,
// 		"description":   "Complete workflow test workshop",
// 		"is_active":     true,
// 	}

// 	body, _ := json.Marshal(workshopInput)
// 	req := httptest.NewRequest(http.MethodPost, "/v1/workshops", bytes.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))
// 	req = setUserContext(req, adminUser)

// 	rec := httptest.NewRecorder()
// 	testApp.createWorkshopHandler(rec, req)

// 	if rec.Result().StatusCode != http.StatusCreated {
// 		t.Fatalf("Workshop creation failed with status %d", rec.Result().StatusCode)
// 	}

// 	var createResponse map[string]any
// 	_ = json.NewDecoder(rec.Result().Body).Decode(&createResponse)
// 	workshopResponse := createResponse["workshop"].(map[string]any)
// 	workshopID := int64(workshopResponse["id"].(float64))

// 	defer testApp.models.Workshop.Delete(workshopID) // Cleanup

// 	t.Logf("Step: Created workshop ID %d", workshopID)

// 	// 2. Get the workshop
// 	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/workshops/%d", workshopID), nil)
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))
// 	req = setURLParam(req, "id", strconv.FormatInt(workshopID, 10))
// 	req = setUserContext(req, adminUser)

// 	rec = httptest.NewRecorder()
// 	testApp.showWorkshopHandler(rec, req)

// 	if rec.Result().StatusCode != http.StatusOK {
// 		t.Fatalf("Workshop retrieval failed with status %d", rec.Result().StatusCode)
// 	}

// 	t.Logf("Step: Successfully retrieved workshop")

// 	// 3. Update the workshop
// 	updateInput := map[string]any{
// 		"workshop_name": fmt.Sprintf("UPDATED_WORKFLOW_Workshop_%d", time.Now().UnixNano()),
// 		"credit_hours":  75,
// 		"description":   "Updated workflow workshop description",
// 	}

// 	body, _ = json.Marshal(updateInput)
// 	req = httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/workshops/%d", workshopID), bytes.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))
// 	req = setURLParam(req, "id", strconv.FormatInt(workshopID, 10))
// 	req = setUserContext(req, adminUser)

// 	rec = httptest.NewRecorder()
// 	testApp.updateWorkshopHandler(rec, req)

// 	if rec.Result().StatusCode != http.StatusOK {
// 		t.Fatalf("Workshop update failed with status %d", rec.Result().StatusCode)
// 	}

// 	t.Logf("Step: Successfully updated workshop")

// 	// 4. List workshops (should include our workshop)
// 	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/workshops?category_id=%d", categoryID), nil)
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))
// 	req = setUserContext(req, adminUser)

// 	rec = httptest.NewRecorder()
// 	testApp.listWorkshopsHandler(rec, req)

// 	if rec.Result().StatusCode != http.StatusOK {
// 		t.Fatalf("Workshop listing failed with status %d", rec.Result().StatusCode)
// 	}

// 	t.Logf("Step: Successfully listed workshops")

// 	// 5. Delete the workshop
// 	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/workshops/%d", workshopID), nil)
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))
// 	req = setURLParam(req, "id", strconv.FormatInt(workshopID, 10))
// 	req = setUserContext(req, adminUser)

// 	rec = httptest.NewRecorder()
// 	testApp.deleteWorkshopHandler(rec, req)

// 	if rec.Result().StatusCode != http.StatusOK {
// 		t.Fatalf("Workshop deletion failed with status %d", rec.Result().StatusCode)
// 	}

// 	t.Logf("Step: Successfully deleted workshop")

// 	t.Log("Step: Complete workshop workflow test passed successfully!")
// }
