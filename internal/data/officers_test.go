package data

import (
	"database/sql"
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

// Helper function to clean up officer test data only
func CleanupOfficerTestData(t *testing.T, db *sql.DB) {
	t.Helper()

	_, err := db.Exec("TRUNCATE TABLE officers RESTART IDENTITY CASCADE")
	if err != nil {
		t.Logf("Warning: Failed to cleanup officers: %v", err)
	}
}

// Helper function to create officer test data using existing seed data
func CreateOfficerTestData(t *testing.T, db *sql.DB) *Officer {
	t.Helper()

	// Create models
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

	// Get existing seed data by name using the new GetByName methods
	region, err := regionModel.GetByName("Northern Region")
	if err != nil {
		t.Fatalf("Failed to get region: %v", err)
	}

	// Get a formation from Northern Region (should be ID 1 or 2)
	formation, err := formationModel.GetByName("Corozal Police Formation") // Corozal Police Formation
	if err != nil {
		t.Fatalf("Failed to get formation: %v", err)
	}

	// Get a posting
	posting, err := postingModel.GetByName("Relief")
	if err != nil {
		t.Fatalf("Failed to get posting: %v", err)
	}

	// Get a rank
	rank, err := rankModel.GetByName("Constable")
	if err != nil {
		t.Fatalf("Failed to get rank: %v", err)
	}

	officer := &Officer{
		UserID:           testUser.ID,
		RegulationNumber: fmt.Sprintf("OFF%d", time.Now().UnixNano()),
		RankID:           rank.ID,
		PostingID:        posting.ID,
		FormationID:      formation.ID,
		RegionID:         region.ID,
	}

	err = officerModel.Insert(officer)
	if err != nil {
		t.Fatalf("Failed to create test officer: %v", err)
	}

	return officer
}

func TestOfficerGetAll(t *testing.T) {
	db := setupTestDB(t)

	// Clean up before and after test - officer-specific
	CleanupOfficerTestData(t, db)
	defer CleanupOfficerTestData(t, db)

	// Create test dependencies
	userModel := UserModel{DB: db}
	regionModel := RegionModel{DB: db}
	formationModel := FormationModel{DB: db}
	postingModel := PostingModel{DB: db}
	rankModel := RankModel{DB: db}
	officerModel := OfficerModel{DB: db}

	// Use existing seed data instead of creating new ones
	testRegion, err := regionModel.GetByName("Northern Region")
	if err != nil {
		t.Fatalf("Failed to get test region: %v", err)
	}

	testFormation, err := formationModel.GetByName("Corozal Police Formation")
	if err != nil {
		t.Fatalf("Failed to get test formation: %v", err)
	}

	testPosting, err := postingModel.GetByName("Relief")
	if err != nil {
		t.Fatalf("Failed to get test posting: %v", err)
	}

	testRank, err := rankModel.GetByName("Constable")
	if err != nil {
		t.Fatalf("Failed to get test rank: %v", err)
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

		officer := &Officer{
			UserID:           testUser.ID,
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
func TestOfficerConstraints(t *testing.T) {
	db := setupTestDB(t)

	// Clean up before and after test
	CleanupOfficerTestData(t, db)
	defer CleanupOfficerTestData(t, db)

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
	err := userModel.Insert(testUser2)
	if err != nil {
		t.Fatalf("Failed to create test user2: %v", err)
	}

	testRegion, err := regionModel.GetByName("Southern Region")
	if err != nil {
		t.Fatalf("Failed to get test region: %v", err)
	}

	testFormation, err := formationModel.GetByName("Ladyville Police Sub-Formation")
	if err != nil {
		t.Fatalf("Failed to get test formation: %v", err)
	}

	testPosting, err := postingModel.GetByName("Special Branch")
	if err != nil {
		t.Fatalf("Failed to get test posting: %v", err)
	}

	testRank, err := rankModel.GetByName("Constable")
	if err != nil {
		t.Fatalf("Failed to get test rank: %v", err)
	}

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
	err = officerModel.Insert(officer2)
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

func TestInsertGetUpdateDeleteOfficer(t *testing.T) {
	db := setupTestDB(t)

	// Clean up before and after test - officer-specific
	CleanupOfficerTestData(t, db)
	defer CleanupOfficerTestData(t, db)

	// Create test dependencies
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

	// Use existing seed data instead of creating new ones
	testRegion, err := regionModel.GetByName("Eastern Division")
	if err != nil {
		t.Fatalf("Failed to get test region: %v", err)
	}

	testFormation, err := formationModel.Get(8) // Police Headquarters - Eastern Division
	if err != nil {
		t.Fatalf("Failed to get test formation: %v", err)
	}

	testPosting, err := postingModel.GetByName("Station Manager")
	if err != nil {
		t.Fatalf("Failed to get test posting: %v", err)
	}

	testRank, err := rankModel.GetByName("Corporal")
	if err != nil {
		t.Fatalf("Failed to get test rank: %v", err)
	}

	// Create test officer
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

	// Test Update - change to a different rank and posting
	updatedRank, err := rankModel.GetByName("Inspector of Police")
	if err != nil {
		t.Fatalf("Failed to get updated rank: %v", err)
	}

	updatedPosting, err := postingModel.GetByName("Crimes Investigation Branch")
	if err != nil {
		t.Fatalf("Failed to get updated posting: %v", err)
	}

	officer.RegulationNumber = fmt.Sprintf("UPD%d", time.Now().UnixNano())
	officer.RankID = updatedRank.ID
	officer.PostingID = updatedPosting.ID

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

	if updated.RankID != updatedRank.ID {
		t.Errorf("Expected updated rank ID %v, got %v", updatedRank.ID, updated.RankID)
	}

	if updated.PostingID != updatedPosting.ID {
		t.Errorf("Expected updated posting ID %v, got %v", updatedPosting.ID, updated.PostingID)
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
func TestOfficerGetWithDetails(t *testing.T) {
	db := setupTestDB(t)

	// Clean up before and after test - officer-specific
	CleanupOfficerTestData(t, db)
	defer CleanupOfficerTestData(t, db)

	// Create test dependencies
	userModel := UserModel{DB: db}
	regionModel := RegionModel{DB: db}
	formationModel := FormationModel{DB: db}
	postingModel := PostingModel{DB: db}
	rankModel := RankModel{DB: db}
	officerModel := OfficerModel{DB: db}

	// Create test user
	testUser := &User{
		FirstName:   "Detailed",
		LastName:    "Officer",
		Email:       fmt.Sprintf("detailed%d@example.com", time.Now().UnixNano()),
		Gender:      "f",
		IsActivated: true,
		IsOfficer:   true,
	}
	_ = testUser.Password.Set("TestPass123!")
	err := userModel.Insert(testUser)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Use existing seed data
	testRegion, err := regionModel.GetByName("Western Region")
	if err != nil {
		t.Fatalf("Failed to get test region: %v", err)
	}

	testFormation, err := formationModel.Get(3) // Police Headquarters - Belmopan
	if err != nil {
		t.Fatalf("Failed to get test formation: %v", err)
	}

	testPosting, err := postingModel.GetByName("Staff Duties")
	if err != nil {
		t.Fatalf("Failed to get test posting: %v", err)
	}

	testRank, err := rankModel.GetByName("Sergeant")
	if err != nil {
		t.Fatalf("Failed to get test rank: %v", err)
	}

	// Create test officer
	officer := &Officer{
		UserID:           testUser.ID,
		RegulationNumber: fmt.Sprintf("DET%d", time.Now().UnixNano()),
		RankID:           testRank.ID,
		PostingID:        testPosting.ID,
		FormationID:      testFormation.ID,
		RegionID:         testRegion.ID,
	}
	err = officerModel.Insert(officer)
	if err != nil {
		t.Fatalf("Failed to insert test officer: %v", err)
	}

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
