package data

import (
	"fmt"
	"testing"
	"time"

	"github.com/Pedro-J-Kukul/police_training/internal/validator"
)

func TestValidateOfficer(t *testing.T) {
	tests := []struct {
		name    string
		officer *Officer
		wantErr bool
	}{
		{
			name: "Valid officer",
			officer: &Officer{
				UserID:           1,
				RegulationNumber: "OFF123456",
				RankID:           1,
				PostingID:        1,
				FormationID:      1,
				RegionID:         1,
			},
			wantErr: false,
		},
		{
			name: "Missing user ID",
			officer: &Officer{
				RegulationNumber: "OFF123456",
				RankID:           1,
				PostingID:        1,
				FormationID:      1,
				RegionID:         1,
			},
			wantErr: true,
		},
		{
			name: "Empty regulation number",
			officer: &Officer{
				UserID:           1,
				RegulationNumber: "",
				RankID:           1,
				PostingID:        1,
				FormationID:      1,
				RegionID:         1,
			},
			wantErr: true,
		},
		{
			name: "Regulation number too long",
			officer: &Officer{
				UserID:           1,
				RegulationNumber: "This is a very long regulation number that exceeds fifty characters",
				RankID:           1,
				PostingID:        1,
				FormationID:      1,
				RegionID:         1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			ValidateOfficer(v, tt.officer)

			hasErrors := !v.IsEmpty()
			if hasErrors != tt.wantErr {
				t.Errorf("ValidateOfficer() error = %v, wantErr %v", hasErrors, tt.wantErr)
				if hasErrors {
					t.Logf("Validation errors: %v", v.Errors)
				}
			}
		})
	}
}

func TestInsertGetUpdateDeleteOfficer(t *testing.T) {
	db := setupTestDB(t)

	// Clean up before and after test
	cleanupTestData(t, db)
	defer cleanupTestData(t, db)

	// Create test dependencies first
	userModel := UserModel{DB: db}
	regionModel := RegionModel{DB: db}
	formationModel := FormationModel{DB: db}
	postingModel := PostingModel{DB: db}
	rankModel := RankModel{DB: db}
	officerModel := OfficerModel{DB: db}

	// Create a test user first
	testUser := &User{
		FirstName:   "Test",
		LastName:    "Officer",
		Email:       fmt.Sprintf("testofficer%d@example.com", time.Now().UnixNano()),
		Gender:      "m",
		IsActivated: true,
		IsOfficer:   true,
	}
	_ = testUser.Password.Set("TestPass123!")
	err := userModel.Insert(testUser)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create test region
	testRegion := &Region{Region: "Test Region"}
	err = regionModel.Insert(testRegion)
	if err != nil {
		t.Fatalf("Failed to create test region: %v", err)
	}

	// Create test formation
	testFormation := &Formation{Formation: "Test Formation", RegionID: testRegion.ID}
	err = formationModel.Insert(testFormation)
	if err != nil {
		t.Fatalf("Failed to create test formation: %v", err)
	}

	// Create test posting
	testPosting := &Posting{Posting: "Test Posting", Code: "TP"}
	err = postingModel.Insert(testPosting)
	if err != nil {
		t.Fatalf("Failed to create test posting: %v", err)
	}

	// Create test rank
	testRank := &Rank{Rank: "Test Rank", Code: "TR", AnnualTrainingHoursRequired: 40}
	err = rankModel.Insert(testRank)
	if err != nil {
		t.Fatalf("Failed to create test rank: %v", err)
	}

	// Now test officer operations
	officer := &Officer{
		UserID:           testUser.ID,
		RegulationNumber: fmt.Sprintf("OFF%d", time.Now().UnixNano()),
		RankID:           testRank.ID,
		PostingID:        testPosting.ID,
		FormationID:      testFormation.ID,
		RegionID:         testRegion.ID,
	}

	// Test Insert
	err = officerModel.Insert(officer)
	if err != nil {
		t.Fatalf("Insert() failed: %v", err)
	}

	if officer.ID == 0 {
		t.Error("Expected officer ID to be set after insert")
	}

	// Test Get
	fetched, err := officerModel.Get(officer.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if fetched.RegulationNumber != officer.RegulationNumber {
		t.Errorf("Expected regulation number %v, got %v", officer.RegulationNumber, fetched.RegulationNumber)
	}

	// Test GetByUserID
	fetchedByUser, err := officerModel.GetByUserID(testUser.ID)
	if err != nil {
		t.Fatalf("GetByUserID() failed: %v", err)
	}

	if fetchedByUser.ID != officer.ID {
		t.Errorf("Expected officer ID %v, got %v", officer.ID, fetchedByUser.ID)
	}

	// Test GetByRegulationNumber
	fetchedByReg, err := officerModel.GetByRegulationNumber(officer.RegulationNumber)
	if err != nil {
		t.Fatalf("GetByRegulationNumber() failed: %v", err)
	}

	if fetchedByReg.ID != officer.ID {
		t.Errorf("Expected officer ID %v, got %v", officer.ID, fetchedByReg.ID)
	}

	// Test Update
	officer.RegulationNumber = fmt.Sprintf("UPD%d", time.Now().UnixNano())
	err = officerModel.Update(officer)
	if err != nil {
		t.Fatalf("Update() failed: %v", err)
	}

	// Verify update
	updated, err := officerModel.Get(officer.ID)
	if err != nil {
		t.Fatalf("Get() after update failed: %v", err)
	}

	if updated.RegulationNumber != officer.RegulationNumber {
		t.Errorf("Expected updated regulation number %v, got %v", officer.RegulationNumber, updated.RegulationNumber)
	}

	// Test Delete
	err = officerModel.Delete(officer.ID)
	if err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}

	// Verify deletion
	_, err = officerModel.Get(officer.ID)
	if err == nil {
		t.Error("Expected officer to be deleted")
	}
	if err != ErrRecordNotFound {
		t.Errorf("Expected ErrRecordNotFound, got %v", err)
	}
}

func TestOfficerGetAll(t *testing.T) {
	db := setupTestDB(t)

	// Clean up before and after test

	// Create test dependencies
	userModel := UserModel{DB: db}
	regionModel := RegionModel{DB: db}
	formationModel := FormationModel{DB: db}
	postingModel := PostingModel{DB: db}
	rankModel := RankModel{DB: db}
	officerModel := OfficerModel{DB: db}

	// Create test data
	testRegion := &Region{Region: "Test Resdsdsdsdasdfsdfsdfgi12on"}
	err := regionModel.Insert(testRegion)
	if err != nil {
		t.Fatalf("Failed to create test region: %v", err)
	}

	testFormation := &Formation{Formation: "Test Formsdsdssdasdfsdfsdfasdsdsdtion", RegionID: testRegion.ID}
	err = formationModel.Insert(testFormation)
	if err != nil {
		t.Fatalf("Failed to create test formation: %v", err)
	}

	testPosting := &Posting{Posting: "Test Posdsdsdstasdfsdfsdfing", Code: "sdasdasdfsdfsdfTP"}
	err = postingModel.Insert(testPosting)
	if err != nil {
		t.Fatalf("Failed to create test posting: %v", err)
	}

	testRank := &Rank{Rank: "Test Ranasdasdasdfsdfsdfssdsdk", Code: "TasdasdsasdfsdfsdfadasdR", AnnualTrainingHoursRequired: 40}
	err = rankModel.Insert(testRank)
	if err != nil {
		t.Fatalf("Failed to create test rank: %v", err)
	}

	// Create multiple test officers
	numOfficers := 3
	for i := 0; i < numOfficers; i++ {
		testUser := &User{
			FirstName:   fmt.Sprintf("Officer%d", i),
			LastName:    "Test",
			Email:       fmt.Sprintf("officer%d_%d@example.com", i, time.Now().UnixNano()),
			Gender:      "m",
			IsActivated: true,
			IsOfficer:   true,
		}
		_ = testUser.Password.Set("TestPass123!")
		err := userModel.Insert(testUser)
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
		testUserID := testUser.ID

		officer := &Officer{
			UserID:           testUserID,
			RegulationNumber: fmt.Sprintf("OFF%d_%d", i, time.Now().UnixNano()),
			RankID:           testRank.ID,
			PostingID:        testPosting.ID,
			FormationID:      testFormation.ID,
			RegionID:         testRegion.ID,
		}
		err = officerModel.Insert(officer)
		if err != nil {
			t.Fatalf("Insert() failed: %v", err)
		}
	}

	// Test GetAll
	filters := Filters{
		Page:         1,
		PageSize:     10,
		Sort:         "id",
		SortSafelist: []string{"id", "-id"},
	}

	officers, metadata, err := officerModel.GetAll("", nil, nil, nil, nil, filters)
	if err != nil {
		t.Fatalf("GetAll() failed: %v", err)
	}

	if len(officers) != numOfficers {
		t.Errorf("Expected %d officers, got %d", numOfficers, len(officers))
	}

	if metadata.TotalRecords != numOfficers {
		t.Errorf("Expected total records %d, got %d", numOfficers, metadata.TotalRecords)
	}
}

func TestOfficerGetWithDetails(t *testing.T) {
	db := setupTestDB(t)

	// Create test dependencies and officer (similar to above)
	userModel := UserModel{DB: db}
	regionModel := RegionModel{DB: db}
	formationModel := FormationModel{DB: db}
	postingModel := PostingModel{DB: db}
	rankModel := RankModel{DB: db}
	officerModel := OfficerModel{DB: db}

	// Create test data
	testUser := &User{
		FirstName:   "Detailed",
		LastName:    "Officer",
		Email:       fmt.Sprintf("detailed%d@example.com", time.Now().UnixNano()),
		Gender:      "f",
		IsActivated: true,
		IsOfficer:   true,
	}
	_ = testUser.Password.Set("TestPass123!")
	_ = userModel.Insert(testUser)

	testRegion := &Region{Region: "Detailed Region"}
	_ = regionModel.Insert(testRegion)

	testFormation := &Formation{Formation: "Detailed Formation", RegionID: testRegion.ID}
	_ = formationModel.Insert(testFormation)

	testPosting := &Posting{Posting: "Detailed Posting", Code: "DP"}
	_ = postingModel.Insert(testPosting)

	testRank := &Rank{Rank: "Detailed Rank", Code: "DR", AnnualTrainingHoursRequired: 50}
	_ = rankModel.Insert(testRank)

	officer := &Officer{
		UserID:           testUser.ID,
		RegulationNumber: fmt.Sprintf("DET%d", time.Now().UnixNano()),
		RankID:           testRank.ID,
		PostingID:        testPosting.ID,
		FormationID:      testFormation.ID,
		RegionID:         testRegion.ID,
	}
	_ = officerModel.Insert(officer)

	// Test GetWithDetails
	detailed, err := officerModel.GetWithDetails(officer.ID)
	if err != nil {
		t.Fatalf("GetWithDetails() failed: %v", err)
	}

	// Verify related data is populated
	if detailed.User == nil {
		t.Error("Expected user data to be populated")
	} else if detailed.User.FirstName != testUser.FirstName {
		t.Errorf("Expected user first name %v, got %v", testUser.FirstName, detailed.User.FirstName)
	}

	if detailed.Rank == nil {
		t.Error("Expected rank data to be populated")
	} else if detailed.Rank.Rank != testRank.Rank {
		t.Errorf("Expected rank %v, got %v", testRank.Rank, detailed.Rank.Rank)
	}

	if detailed.Posting == nil {
		t.Error("Expected posting data to be populated")
	} else if detailed.Posting.Posting != testPosting.Posting {
		t.Errorf("Expected posting %v, got %v", testPosting.Posting, detailed.Posting.Posting)
	}

	if detailed.Formation == nil {
		t.Error("Expected formation data to be populated")
	} else if detailed.Formation.Formation != testFormation.Formation {
		t.Errorf("Expected formation %v, got %v", testFormation.Formation, detailed.Formation.Formation)
	}

	if detailed.Region == nil {
		t.Error("Expected region data to be populated")
	} else if detailed.Region.Region != testRegion.Region {
		t.Errorf("Expected region %v, got %v", testRegion.Region, detailed.Region.Region)
	}
}

func TestOfficerConstraints(t *testing.T) {
	db := setupTestDB(t)

	// Clean up before and after test
	cleanupTestData(t, db)
	defer cleanupTestData(t, db)

	userModel := UserModel{DB: db}
	regionModel := RegionModel{DB: db}
	formationModel := FormationModel{DB: db}
	postingModel := PostingModel{DB: db}
	rankModel := RankModel{DB: db}
	officerModel := OfficerModel{DB: db}

	// Create test dependencies
	testUser1 := &User{
		FirstName:   "User1",
		LastName:    "Test",
		Email:       fmt.Sprintf("user1_%d@example.com", time.Now().UnixNano()),
		Gender:      "m",
		IsActivated: true,
		IsOfficer:   true,
	}
	_ = testUser1.Password.Set("TestPass123!")
	_ = userModel.Insert(testUser1)

	testUser2 := &User{
		FirstName:   "User2",
		LastName:    "Test",
		Email:       fmt.Sprintf("user2_%d@example.com", time.Now().UnixNano()),
		Gender:      "f",
		IsActivated: true,
		IsOfficer:   true,
	}
	_ = testUser2.Password.Set("TestPass123!")
	_ = userModel.Insert(testUser2)

	testRegion := &Region{Region: "Constraint Test Region"}
	_ = regionModel.Insert(testRegion)

	testFormation := &Formation{Formation: "Constraint Test Formation", RegionID: testRegion.ID}
	_ = formationModel.Insert(testFormation)

	testPosting := &Posting{Posting: "Constraint Test Posting", Code: "CTP"}
	_ = postingModel.Insert(testPosting)

	testRank := &Rank{Rank: "Constraint Test Rank", Code: "CTR", AnnualTrainingHoursRequired: 45}
	_ = rankModel.Insert(testRank)

	regNumber := fmt.Sprintf("CONST%d", time.Now().UnixNano())

	// Test unique user_id constraint
	officer1 := &Officer{
		UserID:           testUser1.ID,
		RegulationNumber: regNumber + "1",
		RankID:           testRank.ID,
		PostingID:        testPosting.ID,
		FormationID:      testFormation.ID,
		RegionID:         testRegion.ID,
	}
	_ = officerModel.Insert(officer1)

	// Try to create another officer with the same user_id
	officer2 := &Officer{
		UserID:           testUser1.ID, // Same user ID
		RegulationNumber: regNumber + "2",
		RankID:           testRank.ID,
		PostingID:        testPosting.ID,
		FormationID:      testFormation.ID,
		RegionID:         testRegion.ID,
	}
	err := officerModel.Insert(officer2)
	if err == nil {
		t.Error("Expected error when inserting officer with duplicate user_id")
	}
	if err != ErrDuplicateValue {
		t.Errorf("Expected ErrDuplicateValue, got %v", err)
	}

	// Test unique regulation_number constraint
	officer3 := &Officer{
		UserID:           testUser2.ID,
		RegulationNumber: regNumber + "1", // Same regulation number as officer1
		RankID:           testRank.ID,
		PostingID:        testPosting.ID,
		FormationID:      testFormation.ID,
		RegionID:         testRegion.ID,
	}
	err = officerModel.Insert(officer3)
	if err == nil {
		t.Error("Expected error when inserting officer with duplicate regulation_number")
	}
	if err != ErrDuplicateValue {
		t.Errorf("Expected ErrDuplicateValue, got %v", err)
	}
}
