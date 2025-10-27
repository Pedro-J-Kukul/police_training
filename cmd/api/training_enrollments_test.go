package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"strconv"
// 	"testing"
// 	"time"

// 	"github.com/Pedro-J-Kukul/police_training/internal/data"
// )

// func getSeededEnrollmentData(t *testing.T) (officerID, sessionID, enrollmentStatusID, progressStatusID int64) {
// 	t.Helper()
// 	t.Log("Step: Getting seeded enrollment reference data")

// 	// Get officer
// 	officers, _, err := testApp.models.Officer.GetAll("", nil, nil, nil, nil, data.Filters{
// 		Page: 1, PageSize: 10, Sort: "id", SortSafelist: []string{"id"}, // Start from page 2
// 	})
// 	if err != nil || len(officers) == 0 {
// 		t.Fatal("No officers found in seeded data")
// 	}
// 	officerID = officers[0].ID

// 	// Get session
// 	sessions, _, err := testApp.models.TrainingSession.GetAll(nil, nil, nil, nil, nil, nil, data.Filters{
// 		Page: 1, PageSize: 10, Sort: "id", SortSafelist: []string{"id"}, // Start from page 2
// 	})
// 	if err != nil || len(sessions) == 0 {
// 		t.Fatal("No training sessions found in seeded data")
// 	}
// 	sessionID = sessions[0].ID

// 	// Get enrollment status
// 	enrollmentStatus, err := testApp.models.EnrollmentStatus.GetByName("Enrolled")
// 	if err != nil {
// 		t.Fatal("Enrolled status not found in seeded data")
// 	}
// 	enrollmentStatusID = enrollmentStatus.ID

// 	// Get progress status
// 	progressStatus, err := testApp.models.ProgressStatus.GetByName("In Progress")
// 	if err != nil {
// 		t.Fatal("In Progress status not found in seeded data")
// 	}
// 	progressStatusID = progressStatus.ID

// 	t.Logf("Step: Using officer ID %d, session ID %d, enrollment status ID %d, progress status ID %d",
// 		officerID, sessionID, enrollmentStatusID, progressStatusID)
// 	return
// }

// func getSeededEnrollment(t *testing.T) *data.TrainingEnrollment {
// 	t.Helper()
// 	t.Log("Step: Getting seeded training enrollment")

// 	enrollments, _, err := testApp.models.TrainingEnrollment.GetAll(nil, nil, nil, nil, nil, data.Filters{
// 		Page: 1, PageSize: 10, Sort: "id", SortSafelist: []string{"id"},
// 	})
// 	if err != nil || len(enrollments) == 0 {
// 		t.Fatal("No training enrollments found in seeded data")
// 	}

// 	enrollment := enrollments[0]
// 	t.Logf("Step: Using seeded enrollment ID %d", enrollment.ID)
// 	return enrollment
// }

// func createTestEnrollment(t *testing.T) *data.TrainingEnrollment {
// 	t.Helper()

// 	// Get officers starting from a higher page to avoid conflicts
// 	officers, _, err := testApp.models.Officer.GetAll("", nil, nil, nil, nil, data.Filters{
// 		Page: 1, PageSize: 10, Sort: "id", SortSafelist: []string{"id"}, // Use page 3
// 	})
// 	if err != nil || len(officers) == 0 {
// 		t.Fatal("No officers found for enrollment testing")
// 	}

// 	// Get sessions starting from a higher page
// 	sessions, _, err := testApp.models.TrainingSession.GetAll(nil, nil, nil, nil, nil, nil, data.Filters{
// 		Page: 1, PageSize: 10, Sort: "id", SortSafelist: []string{"id"}, // Use page 3
// 	})
// 	if err != nil || len(sessions) == 0 {
// 		t.Fatal("No sessions found for enrollment testing")
// 	}

// 	// Use the first officer and session from this page (should be unique)
// 	officerID := officers[0].ID
// 	sessionID := sessions[0].ID

// 	// Get enrollment status
// 	enrollmentStatus, err := testApp.models.EnrollmentStatus.GetByName("Enrolled")
// 	if err != nil {
// 		t.Fatal("Enrolled status not found")
// 	}
// 	enrollmentStatusID := enrollmentStatus.ID

// 	// Get progress status
// 	progressStatus, err := testApp.models.ProgressStatus.GetByName("In Progress")
// 	if err != nil {
// 		t.Fatal("In Progress status not found")
// 	}
// 	progressStatusID := progressStatus.ID

// 	// Get attendance status
// 	attendanceStatus, err := testApp.models.AttendanceStatus.GetByName("Absent")
// 	if err != nil {
// 		t.Fatal("Present attendance status not found")
// 	}

// 	enrollment := &data.TrainingEnrollment{
// 		OfficerID:          officerID,
// 		SessionID:          sessionID,
// 		EnrollmentStatusID: enrollmentStatusID,
// 		AttendanceStatusID: &attendanceStatus.ID,
// 		ProgressStatusID:   progressStatusID,
// 		CertificateIssued:  false,
// 		CertificateNumber:  nil,
// 		CompletionDate:     nil,
// 	}

// 	err = testApp.models.TrainingEnrollment.Insert(enrollment)
// 	if err != nil {
// 		t.Fatalf("Failed to create test enrollment: %v", err)
// 	}

// 	t.Logf("Step: Created test enrollment ID %d", enrollment.ID)
// 	return enrollment
// }

// func TestCreateTrainingEnrollmentHandler(t *testing.T) {
// 	t.Log("=== Testing Create Training Enrollment Handler ===")

// 	// Get different officer and session to avoid duplicate enrollment
// 	officers, _, err := testApp.models.Officer.GetAll("", nil, nil, nil, nil, data.Filters{
// 		Page: 1, PageSize: 10, Sort: "id", SortSafelist: []string{"id"},
// 	})
// 	if err != nil || len(officers) < 2 {
// 		t.Skip("Need at least 2 officers for enrollment testing")
// 	}

// 	sessions, _, err := testApp.models.TrainingSession.GetAll(nil, nil, nil, nil, nil, nil, data.Filters{
// 		Page: 1, PageSize: 10, Sort: "id", SortSafelist: []string{"id"},
// 	})
// 	if err != nil || len(sessions) < 2 {
// 		t.Skip("Need at least 2 sessions for enrollment testing")
// 	}

// 	officerID := officers[1].ID // Use second officer
// 	sessionID := sessions[1].ID // Use second session

// 	enrollmentStatus, _ := testApp.models.EnrollmentStatus.GetByName("Enrolled")
// 	progressStatus, _ := testApp.models.ProgressStatus.GetByName("In Progress")
// 	attendanceStatus, _ := testApp.models.AttendanceStatus.GetByName("Present")

// 	adminUser := getSeededUser(t, "admin1@police-training.bz")
// 	adminToken := createTokenForSeededUser(t, adminUser.ID)

// 	tests := []struct {
// 		name           string
// 		input          map[string]any
// 		expectedStatus int
// 		checkResponse  func(*testing.T, *http.Response)
// 	}{
// 		{
// 			name: "Valid enrollment creation",
// 			input: map[string]any{
// 				"officer_id":           officerID,
// 				"session_id":           sessionID,
// 				"enrollment_status_id": enrollmentStatus.ID,
// 				"attendance_status_id": attendanceStatus.ID,
// 				"progress_status_id":   progressStatus.ID,
// 				"certificate_issued":   false,
// 			},
// 			expectedStatus: http.StatusCreated,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating successful enrollment creation")
// 				var response map[string]any
// 				err := json.NewDecoder(res.Body).Decode(&response)
// 				if err != nil {
// 					t.Fatalf("Failed to decode response: %v", err)
// 				}

// 				if response["training_enrollment"] == nil {
// 					t.Error("Expected training_enrollment object in response")
// 					return
// 				}

// 				enrollment := response["training_enrollment"].(map[string]any)
// 				enrollmentID := int64(enrollment["id"].(float64))
// 				t.Logf("Step: Created enrollment ID %d", enrollmentID)

// 				// Cleanup
// 				testApp.models.TrainingEnrollment.Delete(enrollmentID)
// 			},
// 		},
// 		{
// 			name: "Missing required fields",
// 			input: map[string]any{
// 				"officer_id": officerID,
// 			},
// 			expectedStatus: http.StatusUnprocessableEntity,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating missing fields validation")
// 			},
// 		},
// 		{
// 			name: "Invalid officer_id",
// 			input: map[string]any{
// 				"officer_id":           999999,
// 				"session_id":           sessionID,
// 				"enrollment_status_id": enrollmentStatus.ID,
// 				"progress_status_id":   progressStatus.ID,
// 			},
// 			expectedStatus: http.StatusBadRequest,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating invalid officer_id handling")
// 			},
// 		},
// 		{
// 			name: "Invalid completion date format",
// 			input: map[string]any{
// 				"officer_id":           officerID,
// 				"session_id":           sessionID,
// 				"enrollment_status_id": enrollmentStatus.ID,
// 				"progress_status_id":   progressStatus.ID,
// 				"completion_date":      "invalid-date",
// 			},
// 			expectedStatus: http.StatusBadRequest,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating invalid date format handling")
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Logf("Starting test: %s", tt.name)

// 			body, _ := json.Marshal(tt.input)
// 			req := httptest.NewRequest(http.MethodPost, "/v1/training-enrollments", bytes.NewReader(body))
// 			req.Header.Set("Content-Type", "application/json")
// 			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))

// 			req = setUserContext(req, adminUser)

// 			rec := httptest.NewRecorder()
// 			testApp.createTrainingEnrollmentHandler(rec, req)

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

// func TestShowTrainingEnrollmentHandler(t *testing.T) {
// 	t.Log("=== Testing Show Training Enrollment Handler ===")

// 	seededEnrollment := getSeededEnrollment(t)
// 	adminUser := getSeededUser(t, "admin1@police-training.bz")
// 	adminToken := createTokenForSeededUser(t, adminUser.ID)

// 	tests := []struct {
// 		name           string
// 		enrollmentID   string
// 		expectedStatus int
// 		checkResponse  func(*testing.T, *http.Response)
// 	}{
// 		{
// 			name:           "Valid seeded enrollment ID",
// 			enrollmentID:   strconv.FormatInt(seededEnrollment.ID, 10),
// 			expectedStatus: http.StatusOK,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating enrollment details response")
// 				var response map[string]any
// 				err := json.NewDecoder(res.Body).Decode(&response)
// 				if err != nil {
// 					t.Fatalf("Failed to decode response: %v", err)
// 				}

// 				if response["training_enrollment"] == nil {
// 					t.Error("Expected training_enrollment object in response")
// 					return
// 				}

// 				enrollment := response["training_enrollment"].(map[string]any)
// 				if enrollment["id"] != float64(seededEnrollment.ID) {
// 					t.Errorf("Expected enrollment ID %d, got %v", seededEnrollment.ID, enrollment["id"])
// 				}

// 				t.Logf("Step: Retrieved enrollment for officer %v in session %v",
// 					enrollment["officer_id"], enrollment["session_id"])

// 				// Verify expected fields
// 				requiredFields := []string{"id", "officer_id", "session_id", "enrollment_status_id", "progress_status_id", "certificate_issued"}
// 				for _, field := range requiredFields {
// 					if enrollment[field] == nil {
// 						t.Errorf("Expected field %s in enrollment response", field)
// 					}
// 				}
// 			},
// 		},
// 		{
// 			name:           "Non-existent enrollment ID",
// 			enrollmentID:   "999999",
// 			expectedStatus: http.StatusNotFound,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating non-existent enrollment response")
// 			},
// 		},
// 		{
// 			name:           "Invalid enrollment ID format",
// 			enrollmentID:   "abc",
// 			expectedStatus: http.StatusNotFound,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating invalid ID format response")
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Logf("Starting test: %s with enrollment ID %s", tt.name, tt.enrollmentID)

// 			path := fmt.Sprintf("/v1/training-enrollments/%s", tt.enrollmentID)
// 			req := httptest.NewRequest(http.MethodGet, path, nil)
// 			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))

// 			req = setURLParam(req, "id", tt.enrollmentID)
// 			req = setUserContext(req, adminUser)

// 			rec := httptest.NewRecorder()
// 			testApp.showTrainingEnrollmentHandler(rec, req)

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

// func TestListTrainingEnrollmentsHandler(t *testing.T) {
// 	t.Log("=== Testing List Training Enrollments Handler ===")

// 	adminUser := getSeededUser(t, "admin1@police-training.bz")
// 	adminToken := createTokenForSeededUser(t, adminUser.ID)

// 	tests := []struct {
// 		name           string
// 		queryParams    string
// 		expectedStatus int
// 		checkResponse  func(*testing.T, *http.Response)
// 	}{
// 		{
// 			name:           "List all enrollments",
// 			queryParams:    "",
// 			expectedStatus: http.StatusOK,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating enrollments list response")
// 				var response map[string]any
// 				err := json.NewDecoder(res.Body).Decode(&response)
// 				if err != nil {
// 					t.Fatalf("Failed to decode response: %v", err)
// 				}

// 				if response["training_enrollments"] == nil {
// 					t.Error("Expected training_enrollments array in response")
// 					return
// 				}

// 				enrollmentsInterface := response["training_enrollments"]
// 				if enrollmentsInterface == nil {
// 					t.Logf("Step: Retrieved 0 enrollments (empty database)")
// 					return
// 				}
// 				enrollments := enrollmentsInterface.([]any)
// 				t.Logf("Step: Retrieved %d enrollments", len(enrollments))

// 				if len(enrollments) > 0 {
// 					enrollment := enrollments[0].(map[string]any)
// 					requiredFields := []string{"id", "officer_id", "session_id", "enrollment_status_id", "progress_status_id"}
// 					for _, field := range requiredFields {
// 						if enrollment[field] == nil {
// 							t.Errorf("Expected field %s in enrollment response", field)
// 						}
// 					}

// 					t.Logf("Step: First enrollment ID %v, Officer: %v, Session: %v",
// 						enrollment["id"], enrollment["officer_id"], enrollment["session_id"])
// 				}

// 				if response["metadata"] == nil {
// 					t.Error("Expected metadata in response")
// 				}
// 			},
// 		},
// 		{
// 			name:           "List enrollments with pagination",
// 			queryParams:    "?page=1&page_size=3",
// 			expectedStatus: http.StatusOK,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating paginated enrollments response")
// 				var response map[string]any
// 				err := json.NewDecoder(res.Body).Decode(&response)
// 				if err != nil {
// 					t.Fatalf("Failed to decode response: %v", err)
// 				}

// 				enrollmentsInterface := response["training_enrollments"]
// 				if enrollmentsInterface == nil {
// 					t.Logf("Step: Retrieved 0 enrollments (empty database)")
// 					return
// 				}
// 				enrollments := enrollmentsInterface.([]any)
// 				t.Logf("Step: Retrieved %d enrollments with page_size=3", len(enrollments))

// 				if len(enrollments) > 3 {
// 					t.Errorf("Expected max 3 enrollments, got %d", len(enrollments))
// 				}
// 			},
// 		},
// 		{
// 			name:           "Filter enrollments by officer",
// 			queryParams:    "?officer_id=1",
// 			expectedStatus: http.StatusOK,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating officer filter response")
// 				var response map[string]any
// 				err := json.NewDecoder(res.Body).Decode(&response)
// 				if err != nil {
// 					t.Fatalf("Failed to decode response: %v", err)
// 				}

// 				enrollmentsInterface := response["training_enrollments"]
// 				if enrollmentsInterface == nil {
// 					t.Logf("Step: Retrieved 0 enrollments (empty database)")
// 					return
// 				}
// 				enrollments := enrollmentsInterface.([]any)
// 				t.Logf("Step: Retrieved %d enrollments for officer 1", len(enrollments))

// 				// Verify all enrollments belong to officer 1
// 				for i, enrollmentInterface := range enrollments {
// 					enrollment := enrollmentInterface.(map[string]any)
// 					if officerID, ok := enrollment["officer_id"].(float64); !ok || officerID != 1 {
// 						t.Errorf("Enrollment at index %d does not belong to officer 1", i)
// 					}
// 				}
// 			},
// 		},
// 		{
// 			name:           "Filter enrollments by session",
// 			queryParams:    "?session_id=1",
// 			expectedStatus: http.StatusOK,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating session filter response")
// 				var response map[string]any
// 				err := json.NewDecoder(res.Body).Decode(&response)
// 				if err != nil {
// 					t.Fatalf("Failed to decode response: %v", err)
// 				}

// 				enrollmentsInterface := response["training_enrollments"]
// 				if enrollmentsInterface == nil {
// 					t.Logf("Step: Retrieved 0 enrollments (empty database)")
// 					return
// 				}
// 				enrollments := enrollmentsInterface.([]any)
// 				t.Logf("Step: Retrieved %d enrollments for session 1", len(enrollments))
// 			},
// 		},
// 		{
// 			name:           "Filter by enrollment status",
// 			queryParams:    "?enrollment_status_id=1",
// 			expectedStatus: http.StatusOK,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating enrollment status filter response")
// 				var response map[string]any
// 				err := json.NewDecoder(res.Body).Decode(&response)
// 				if err != nil {
// 					t.Fatalf("Failed to decode response: %v", err)
// 				}

// 				enrollmentsInterface := response["training_enrollments"]
// 				if enrollmentsInterface == nil {
// 					t.Logf("Step: Retrieved 0 enrollments (empty database)")
// 					return
// 				}
// 				enrollments := enrollmentsInterface.([]any)
// 				t.Logf("Step: Retrieved %d enrollments with status 1", len(enrollments))
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Logf("Starting test: %s", tt.name)

// 			path := "/v1/training-enrollments" + tt.queryParams
// 			req := httptest.NewRequest(http.MethodGet, path, nil)
// 			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))
// 			req = setUserContext(req, adminUser)

// 			rec := httptest.NewRecorder()
// 			testApp.listTrainingEnrollmentsHandler(rec, req)

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

// func TestUpdateTrainingEnrollmentHandler(t *testing.T) {
// 	t.Log("=== Testing Update Training Enrollment Handler ===")

// 	// Create test enrollment for updating
// 	testEnrollment := createTestEnrollment(t)
// 	defer testApp.models.TrainingEnrollment.Delete(testEnrollment.ID) // Cleanup

// 	completedStatus, _ := testApp.models.ProgressStatus.GetByName("Completed")
// 	adminUser := getSeededUser(t, "admin1@police-training.bz")
// 	adminToken := createTokenForSeededUser(t, adminUser.ID)

// 	tests := []struct {
// 		name           string
// 		enrollmentID   string
// 		input          map[string]any
// 		expectedStatus int
// 		checkResponse  func(*testing.T, *http.Response)
// 	}{
// 		{
// 			name:         "Valid enrollment update",
// 			enrollmentID: strconv.FormatInt(testEnrollment.ID, 10),
// 			input: map[string]any{
// 				"progress_status_id": completedStatus.ID,
// 				"completion_date":    time.Now().Format("2006-01-02"),
// 				"certificate_issued": true,
// 				"certificate_number": "CERT-TEST-123",
// 			},
// 			expectedStatus: http.StatusOK,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating successful enrollment update")
// 				var response map[string]any
// 				err := json.NewDecoder(res.Body).Decode(&response)
// 				if err != nil {
// 					t.Fatalf("Failed to decode response: %v", err)
// 				}

// 				if response["training_enrollment"] == nil {
// 					t.Error("Expected training_enrollment object in response")
// 					return
// 				}

// 				enrollment := response["training_enrollment"].(map[string]any)
// 				if enrollment["certificate_issued"] != true {
// 					t.Errorf("Expected certificate_issued true, got %v", enrollment["certificate_issued"])
// 				}

// 				t.Logf("Step: Updated enrollment with certificate %v", enrollment["certificate_number"])
// 			},
// 		},
// 		{
// 			name:         "Non-existent enrollment",
// 			enrollmentID: "999999",
// 			input: map[string]any{
// 				"certificate_issued": true,
// 			},
// 			expectedStatus: http.StatusNotFound,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating non-existent enrollment update")
// 			},
// 		},
// 		{
// 			name:         "Invalid completion date format",
// 			enrollmentID: strconv.FormatInt(testEnrollment.ID, 10),
// 			input: map[string]any{
// 				"completion_date": "invalid-date",
// 			},
// 			expectedStatus: http.StatusBadRequest,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating invalid date format handling")
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Logf("Starting test: %s", tt.name)

// 			body, _ := json.Marshal(tt.input)
// 			path := fmt.Sprintf("/v1/training-enrollments/%s", tt.enrollmentID)
// 			req := httptest.NewRequest(http.MethodPatch, path, bytes.NewReader(body))
// 			req.Header.Set("Content-Type", "application/json")
// 			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))

// 			req = setURLParam(req, "id", tt.enrollmentID)
// 			req = setUserContext(req, adminUser)

// 			rec := httptest.NewRecorder()
// 			testApp.updateTrainingEnrollmentHandler(rec, req)

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

// func TestDeleteTrainingEnrollmentHandler(t *testing.T) {
// 	t.Log("=== Testing Delete Training Enrollment Handler ===")

// 	// Create test enrollment for deletion
// 	testEnrollment := createTestEnrollment(t)
// 	// No defer cleanup since the test will delete it

// 	adminUser := getSeededUser(t, "admin1@police-training.bz")
// 	adminToken := createTokenForSeededUser(t, adminUser.ID)

// 	tests := []struct {
// 		name           string
// 		enrollmentID   string
// 		expectedStatus int
// 		checkResponse  func(*testing.T, *http.Response)
// 	}{
// 		{
// 			name:           "Valid enrollment deletion",
// 			enrollmentID:   strconv.FormatInt(testEnrollment.ID, 10),
// 			expectedStatus: http.StatusOK,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating successful enrollment deletion")
// 				var response map[string]any
// 				err := json.NewDecoder(res.Body).Decode(&response)
// 				if err != nil {
// 					t.Fatalf("Failed to decode response: %v", err)
// 				}

// 				if response["message"] == nil {
// 					t.Error("Expected message in response")
// 				}

// 				t.Logf("Step: Enrollment deleted with message: %v", response["message"])
// 			},
// 		},
// 		{
// 			name:           "Non-existent enrollment",
// 			enrollmentID:   "999999",
// 			expectedStatus: http.StatusNotFound,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating non-existent enrollment deletion")
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Logf("Starting test: %s", tt.name)

// 			path := fmt.Sprintf("/v1/training-enrollments/%s", tt.enrollmentID)
// 			req := httptest.NewRequest(http.MethodDelete, path, nil)
// 			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))

// 			req = setURLParam(req, "id", tt.enrollmentID)
// 			req = setUserContext(req, adminUser)

// 			rec := httptest.NewRecorder()
// 			testApp.deleteTrainingEnrollmentHandler(rec, req)

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
// } // Continue from where it was cut off...

// func TestGetOfficerEnrollmentsHandler(t *testing.T) {
// 	t.Log("=== Testing Get Officer Enrollments Handler ===")

// 	// Get officer with enrollments
// 	officers, _, err := testApp.models.Officer.GetAll("", nil, nil, nil, nil, data.Filters{
// 		Page: 1, PageSize: 10, Sort: "id", SortSafelist: []string{"id"},
// 	})
// 	if err != nil || len(officers) == 0 {
// 		t.Skip("No officers found for testing")
// 	}

// 	officerID := officers[0].ID
// 	adminUser := getSeededUser(t, "admin1@police-training.bz")
// 	adminToken := createTokenForSeededUser(t, adminUser.ID)

// 	tests := []struct {
// 		name           string
// 		officerID      string
// 		expectedStatus int
// 		checkResponse  func(*testing.T, *http.Response)
// 	}{
// 		{
// 			name:           "Valid officer enrollments",
// 			officerID:      strconv.FormatInt(officerID, 10),
// 			expectedStatus: http.StatusOK,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating officer enrollments response")
// 				var response map[string]any
// 				err := json.NewDecoder(res.Body).Decode(&response)
// 				if err != nil {
// 					t.Fatalf("Failed to decode response: %v", err)
// 				}

// 				if response["training_enrollments"] == nil {
// 					t.Error("Expected training_enrollments array in response")
// 					return
// 				}

// 				enrollmentsInterface := response["training_enrollments"]
// 				if enrollmentsInterface == nil {
// 					t.Logf("Step: Retrieved 0 enrollments (empty database)")
// 					return
// 				}
// 				enrollments := enrollmentsInterface.([]any)
// 				t.Logf("Step: Retrieved %d enrollments for officer %d", len(enrollments), officerID)

// 				// Verify all enrollments belong to this officer
// 				for i, enrollmentInterface := range enrollments {
// 					enrollment := enrollmentInterface.(map[string]any)
// 					if offID, ok := enrollment["officer_id"].(float64); !ok || offID != float64(officerID) {
// 						t.Errorf("Enrollment at index %d does not belong to officer %d", i, officerID)
// 					}
// 				}
// 			},
// 		},
// 		{
// 			name:           "Non-existent officer",
// 			officerID:      "999999",
// 			expectedStatus: http.StatusOK, // Should return empty list
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating non-existent officer enrollments")
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Logf("Starting test: %s", tt.name)

// 			path := fmt.Sprintf("/v1/officers/%s/enrollments", tt.officerID)
// 			req := httptest.NewRequest(http.MethodGet, path, nil)
// 			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))

// 			req = setURLParam(req, "id", tt.officerID)
// 			req = setUserContext(req, adminUser)

// 			rec := httptest.NewRecorder()
// 			testApp.getOfficerEnrollmentsHandler(rec, req)

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

// func TestGetSessionEnrollmentsHandler(t *testing.T) {
// 	t.Log("=== Testing Get Session Enrollments Handler ===")

// 	// Get session with enrollments
// 	sessions, _, err := testApp.models.TrainingSession.GetAll(nil, nil, nil, nil, nil, nil, data.Filters{
// 		Page: 1, PageSize: 10, Sort: "id", SortSafelist: []string{"id"},
// 	})
// 	if err != nil || len(sessions) == 0 {
// 		t.Skip("No sessions found for testing")
// 	}

// 	sessionID := sessions[0].ID
// 	adminUser := getSeededUser(t, "admin1@police-training.bz")
// 	adminToken := createTokenForSeededUser(t, adminUser.ID)

// 	tests := []struct {
// 		name           string
// 		sessionID      string
// 		expectedStatus int
// 		checkResponse  func(*testing.T, *http.Response)
// 	}{
// 		{
// 			name:           "Valid session enrollments",
// 			sessionID:      strconv.FormatInt(sessionID, 10),
// 			expectedStatus: http.StatusOK,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating session enrollments response")
// 				var response map[string]any
// 				err := json.NewDecoder(res.Body).Decode(&response)
// 				if err != nil {
// 					t.Fatalf("Failed to decode response: %v", err)
// 				}

// 				if response["training_enrollments"] == nil {
// 					t.Error("Expected training_enrollments array in response")
// 					return
// 				}

// 				enrollmentsInterface := response["training_enrollments"]
// 				if enrollmentsInterface == nil {
// 					t.Logf("Step: Retrieved 0 enrollments (empty database)")
// 					return
// 				}
// 				enrollments := enrollmentsInterface.([]any)
// 				t.Logf("Step: Retrieved %d enrollments for session %d", len(enrollments), sessionID)

// 				// Verify all enrollments belong to this session
// 				for i, enrollmentInterface := range enrollments {
// 					enrollment := enrollmentInterface.(map[string]any)
// 					if sessID, ok := enrollment["session_id"].(float64); !ok || sessID != float64(sessionID) {
// 						t.Errorf("Enrollment at index %d does not belong to session %d", i, sessionID)
// 					}
// 				}
// 			},
// 		},
// 		{
// 			name:           "Non-existent session",
// 			sessionID:      "999999",
// 			expectedStatus: http.StatusOK, // Should return empty list
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating non-existent session enrollments")
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Logf("Starting test: %s", tt.name)

// 			path := fmt.Sprintf("/v1/training-sessions/%s/enrollments", tt.sessionID)
// 			req := httptest.NewRequest(http.MethodGet, path, nil)
// 			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))

// 			req = setURLParam(req, "id", tt.sessionID)
// 			req = setUserContext(req, adminUser)

// 			rec := httptest.NewRecorder()
// 			testApp.getSessionEnrollmentsHandler(rec, req)

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

// func TestIssueCertificateHandler(t *testing.T) {
// 	t.Log("=== Testing Issue Certificate Handler ===")

// 	// Create test enrollment for certificate issuance
// 	testEnrollment := createTestEnrollment(t)
// 	defer testApp.models.TrainingEnrollment.Delete(testEnrollment.ID) // Cleanup

// 	adminUser := getSeededUser(t, "admin1@police-training.bz")
// 	adminToken := createTokenForSeededUser(t, adminUser.ID)

// 	tests := []struct {
// 		name           string
// 		enrollmentID   string
// 		input          map[string]string
// 		expectedStatus int
// 		checkResponse  func(*testing.T, *http.Response)
// 	}{
// 		{
// 			name:         "Valid certificate issuance",
// 			enrollmentID: strconv.FormatInt(testEnrollment.ID, 10),
// 			input: map[string]string{
// 				"certificate_number": fmt.Sprintf("CERT-API-TEST-%d", time.Now().UnixNano()),
// 				"completion_date":    time.Now().Format("2006-01-02"),
// 			},
// 			expectedStatus: http.StatusOK,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating successful certificate issuance")
// 				var response map[string]any
// 				err := json.NewDecoder(res.Body).Decode(&response)
// 				if err != nil {
// 					t.Fatalf("Failed to decode response: %v", err)
// 				}

// 				if response["message"] == nil {
// 					t.Error("Expected message in response")
// 				}

// 				t.Logf("Step: Certificate issued with message: %v", response["message"])
// 			},
// 		},
// 		{
// 			name:         "Non-existent enrollment",
// 			enrollmentID: "999999",
// 			input: map[string]string{
// 				"certificate_number": "CERT-NONEXISTENT-123",
// 				"completion_date":    time.Now().Format("2006-01-02"),
// 			},
// 			expectedStatus: http.StatusNotFound,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating non-existent enrollment certificate")
// 			},
// 		},
// 		{
// 			name:         "Missing certificate number",
// 			enrollmentID: strconv.FormatInt(testEnrollment.ID, 10),
// 			input: map[string]string{
// 				"completion_date": time.Now().Format("2006-01-02"),
// 			},
// 			expectedStatus: http.StatusUnprocessableEntity,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating missing certificate number")
// 			},
// 		},
// 		{
// 			name:         "Invalid completion date format",
// 			enrollmentID: strconv.FormatInt(testEnrollment.ID, 10),
// 			input: map[string]string{
// 				"certificate_number": "CERT-INVALID-DATE-123",
// 				"completion_date":    "invalid-date",
// 			},
// 			expectedStatus: http.StatusUnprocessableEntity,
// 			checkResponse: func(t *testing.T, res *http.Response) {
// 				t.Log("Step: Validating invalid date format")
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Logf("Starting test: %s", tt.name)

// 			body, _ := json.Marshal(tt.input)
// 			path := fmt.Sprintf("/v1/training-enrollments/%s/certificate", tt.enrollmentID)
// 			req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
// 			req.Header.Set("Content-Type", "application/json")
// 			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))

// 			req = setURLParam(req, "id", tt.enrollmentID)
// 			req = setUserContext(req, adminUser)

// 			rec := httptest.NewRecorder()
// 			testApp.issueCertificateHandler(rec, req)

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

// func TestEnrollmentWorkflow(t *testing.T) {
// 	t.Log("=== Testing Complete Enrollment Workflow ===")

// 	// Get different officer and session to avoid conflicts
// 	officers, _, err := testApp.models.Officer.GetAll("", nil, nil, nil, nil, data.Filters{
// 		Page: 1, PageSize: 10, Sort: "id", SortSafelist: []string{"id"},
// 	})
// 	if err != nil || len(officers) < 3 {
// 		t.Skip("Need at least 3 officers for enrollment workflow testing")
// 	}

// 	sessions, _, err := testApp.models.TrainingSession.GetAll(nil, nil, nil, nil, nil, nil, data.Filters{
// 		Page: 1, PageSize: 10, Sort: "id", SortSafelist: []string{"id"},
// 	})
// 	if err != nil || len(sessions) < 3 {
// 		t.Skip("Need at least 3 sessions for enrollment workflow testing")
// 	}

// 	officerID := officers[2].ID // Use third officer
// 	sessionID := sessions[2].ID // Use third session

// 	enrollmentStatus, _ := testApp.models.EnrollmentStatus.GetByName("Enrolled")
// 	progressStatus, _ := testApp.models.ProgressStatus.GetByName("In Progress")
// 	completedStatus, _ := testApp.models.ProgressStatus.GetByName("Completed")
// 	attendanceStatus, _ := testApp.models.AttendanceStatus.GetByName("Present")

// 	adminUser := getSeededUser(t, "admin1@police-training.bz")
// 	adminToken := createTokenForSeededUser(t, adminUser.ID)

// 	// 1. Create enrollment
// 	enrollmentInput := map[string]any{
// 		"officer_id":           officerID,
// 		"session_id":           sessionID,
// 		"enrollment_status_id": enrollmentStatus.ID,
// 		"attendance_status_id": attendanceStatus.ID,
// 		"progress_status_id":   progressStatus.ID,
// 		"certificate_issued":   false,
// 	}

// 	body, _ := json.Marshal(enrollmentInput)
// 	req := httptest.NewRequest(http.MethodPost, "/v1/training-enrollments", bytes.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))
// 	req = setUserContext(req, adminUser)

// 	rec := httptest.NewRecorder()
// 	testApp.createTrainingEnrollmentHandler(rec, req)

// 	if rec.Result().StatusCode != http.StatusCreated {
// 		t.Fatalf("Enrollment creation failed with status %d", rec.Result().StatusCode)
// 	}

// 	var createResponse map[string]any
// 	_ = json.NewDecoder(rec.Result().Body).Decode(&createResponse)
// 	enrollmentResponse := createResponse["training_enrollment"].(map[string]any)
// 	enrollmentID := int64(enrollmentResponse["id"].(float64))

// 	defer testApp.models.TrainingEnrollment.Delete(enrollmentID) // Cleanup

// 	t.Logf("Step: Created enrollment ID %d", enrollmentID)

// 	// 2. Get the enrollment
// 	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/training-enrollments/%d", enrollmentID), nil)
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))
// 	req = setURLParam(req, "id", strconv.FormatInt(enrollmentID, 10))
// 	req = setUserContext(req, adminUser)

// 	rec = httptest.NewRecorder()
// 	testApp.showTrainingEnrollmentHandler(rec, req)

// 	if rec.Result().StatusCode != http.StatusOK {
// 		t.Fatalf("Enrollment retrieval failed with status %d", rec.Result().StatusCode)
// 	}

// 	t.Logf("Step: Successfully retrieved enrollment")

// 	// 3. Update the enrollment
// 	updateInput := map[string]any{
// 		"progress_status_id": completedStatus.ID,
// 		"completion_date":    time.Now().Format("2006-01-02"),
// 		"certificate_issued": true,
// 	}

// 	body, _ = json.Marshal(updateInput)
// 	req = httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/training-enrollments/%d", enrollmentID), bytes.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))
// 	req = setURLParam(req, "id", strconv.FormatInt(enrollmentID, 10))
// 	req = setUserContext(req, adminUser)

// 	rec = httptest.NewRecorder()
// 	testApp.updateTrainingEnrollmentHandler(rec, req)

// 	if rec.Result().StatusCode != http.StatusOK {
// 		t.Fatalf("Enrollment update failed with status %d", rec.Result().StatusCode)
// 	}

// 	t.Logf("Step: Successfully updated enrollment")

// 	// 4. Issue certificate
// 	certInput := map[string]string{
// 		"certificate_number": fmt.Sprintf("CERT-WORKFLOW-%d", time.Now().UnixNano()),
// 		"completion_date":    time.Now().Format("2006-01-02"),
// 	}

// 	body, _ = json.Marshal(certInput)
// 	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/v1/training-enrollments/%d/certificate", enrollmentID), bytes.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))
// 	req = setURLParam(req, "id", strconv.FormatInt(enrollmentID, 10))
// 	req = setUserContext(req, adminUser)

// 	rec = httptest.NewRecorder()
// 	testApp.issueCertificateHandler(rec, req)

// 	if rec.Result().StatusCode != http.StatusOK {
// 		t.Fatalf("Certificate issuance failed with status %d", rec.Result().StatusCode)
// 	}

// 	t.Logf("Step: Successfully issued certificate")

// 	// 5. Get officer enrollments
// 	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/officers/%d/enrollments", officerID), nil)
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))
// 	req = setURLParam(req, "id", strconv.FormatInt(officerID, 10))
// 	req = setUserContext(req, adminUser)

// 	rec = httptest.NewRecorder()
// 	testApp.getOfficerEnrollmentsHandler(rec, req)

// 	if rec.Result().StatusCode != http.StatusOK {
// 		t.Fatalf("Officer enrollments retrieval failed with status %d", rec.Result().StatusCode)
// 	}

// 	t.Logf("Step: Successfully retrieved officer enrollments")

// 	// 6. Get session enrollments
// 	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/training-sessions/%d/enrollments", sessionID), nil)
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))
// 	req = setURLParam(req, "id", strconv.FormatInt(sessionID, 10))
// 	req = setUserContext(req, adminUser)

// 	rec = httptest.NewRecorder()
// 	testApp.getSessionEnrollmentsHandler(rec, req)

// 	if rec.Result().StatusCode != http.StatusOK {
// 		t.Fatalf("Session enrollments retrieval failed with status %d", rec.Result().StatusCode)
// 	}

// 	t.Logf("Step: Successfully retrieved session enrollments")

// 	// 7. Delete the enrollment
// 	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/training-enrollments/%d", enrollmentID), nil)
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))
// 	req = setURLParam(req, "id", strconv.FormatInt(enrollmentID, 10))
// 	req = setUserContext(req, adminUser)

// 	rec = httptest.NewRecorder()
// 	testApp.deleteTrainingEnrollmentHandler(rec, req)

// 	if rec.Result().StatusCode != http.StatusOK {
// 		t.Fatalf("Enrollment deletion failed with status %d", rec.Result().StatusCode)
// 	}

// 	t.Logf("Step: Successfully deleted enrollment")

// 	t.Log("Step: Complete enrollment workflow test passed successfully!")
// }
