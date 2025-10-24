package data

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/Pedro-J-Kukul/police_training/internal/validator"
)

// Track test records for cleanup
var testWorkshopIDs []int64

func addTestWorkshopID(id int64) {
	testWorkshopIDs = append(testWorkshopIDs, id)
}

func TestValidateWorkshop(t *testing.T) {
	tests := []struct {
		name     string
		workshop *Workshop
		wantErr  bool
	}{
		{
			name: "Valid workshop",
			workshop: &Workshop{
				WorkshopName: "Test Workshop",
				CategoryID:   1,
				TypeID:       1,
				CreditHours:  40,
				IsActive:     true,
			},
			wantErr: false,
		},
		{
			name: "Missing workshop name",
			workshop: &Workshop{
				CategoryID:  1,
				TypeID:      1,
				CreditHours: 40,
				IsActive:    true,
			},
			wantErr: true,
		},
		{
			name: "Workshop name too long",
			workshop: &Workshop{
				WorkshopName: "This is a very long workshop name that exceeds the maximum allowed length of 200 characters and should trigger a validation error because it's way too long for the database field",
				CategoryID:   1,
				TypeID:       1,
				CreditHours:  40,
				IsActive:     true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			ValidateWorkshop(v, tt.workshop)

			hasErrors := !v.IsEmpty()
			if hasErrors != tt.wantErr {
				t.Errorf("ValidateWorkshop() error = %v, wantErr %v", hasErrors, tt.wantErr)
				if hasErrors {
					t.Logf("Validation errors: %v", v.Errors)
				}
			}
		})
	}
}

func cleanupTestWorkshops(t *testing.T, db *sql.DB) {
	t.Helper()

	for _, id := range testWorkshopIDs {
		_, err := db.Exec("DELETE FROM workshops WHERE id = $1", id)
		if err != nil {
			t.Logf("Warning: Failed to cleanup workshop ID %d: %v", id, err)
		}
	}
	testWorkshopIDs = nil // Clear the slice
}

func TestInsertGetUpdateDeleteWorkshop(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestWorkshops(t, db)

	workshopModel := WorkshopModel{DB: db}
	categoryModel := TrainingCategoryModel{DB: db}
	typeModel := TrainingTypeModel{DB: db}

	// Use existing seed data
	testCategory, err := categoryModel.Get(1) // Use first available category
	if err != nil {
		t.Skipf("Skipping test - seed data not available: %v", err)
	}

	testType, err := typeModel.Get(1) // Use first available type
	if err != nil {
		t.Skipf("Skipping test - seed data not available: %v", err)
	}

	// Create test workshop
	description := "Test workshop description"
	workshop := &Workshop{
		WorkshopName: fmt.Sprintf("TEST_Workshop_%d", time.Now().UnixNano()),
		CategoryID:   testCategory.ID,
		TypeID:       testType.ID,
		CreditHours:  40,
		Description:  &description,
		IsActive:     true,
	}

	// Test Insert
	err = workshopModel.Insert(workshop)
	if err != nil {
		t.Fatalf("Insert() failed: %v", err)
	}
	addTestWorkshopID(workshop.ID) // Track for cleanup

	// Test Get
	fetched, err := workshopModel.Get(workshop.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if fetched.WorkshopName != workshop.WorkshopName {
		t.Errorf("Expected workshop name %v, got %v", workshop.WorkshopName, fetched.WorkshopName)
	}

	// Test Update
	workshop.WorkshopName = fmt.Sprintf("TEST_Updated_Workshop_%d", time.Now().UnixNano())
	workshop.CreditHours = 60

	err = workshopModel.Update(workshop)
	if err != nil {
		t.Fatalf("Update() failed: %v", err)
	}

	// Test Delete (this will remove it from tracking automatically)
	err = workshopModel.Delete(workshop.ID)
	if err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}

	// Remove from tracking since we deleted it
	for i, id := range testWorkshopIDs {
		if id == workshop.ID {
			testWorkshopIDs = append(testWorkshopIDs[:i], testWorkshopIDs[i+1:]...)
			break
		}
	}
}

func TestWorkshopGetAll(t *testing.T) {
	db := setupTestDB(t)

	defer cleanupTestWorkshops(t, db)

	workshopModel := WorkshopModel{DB: db}
	categoryModel := TrainingCategoryModel{DB: db}
	typeModel := TrainingTypeModel{DB: db}

	// Use existing seed data
	testCategory, err := categoryModel.GetByName("Firearms Training")
	if err != nil {
		t.Skipf("Skipping test - seed data not available: %v", err)
	}

	testType, err := typeModel.GetByType("Mandatory")
	if err != nil {
		t.Skipf("Skipping test - seed data not available: %v", err)
	}
	// Create multiple test workshops
	numWorkshops := 3
	for i := 0; i < numWorkshops; i++ {
		workshop := &Workshop{
			WorkshopName: fmt.Sprintf("Workshop%d_%d", i, time.Now().UnixNano()),
			CategoryID:   testCategory.ID,
			TypeID:       testType.ID,
			CreditHours:  40 + i*10,
			IsActive:     true,
		}
		err = workshopModel.Insert(workshop)
		if err != nil {
			t.Fatalf("Insert() failed: %v", err)
		}
	}

	// Test GetAll
	filters := Filters{
		Page:         1,
		PageSize:     10,
		Sort:         "id",
		SortSafelist: []string{"id", "-id", "workshop_name", "-workshop_name"},
	}

	workshops, metadata, err := workshopModel.GetAll("", nil, nil, nil, filters)
	if err != nil {
		t.Fatalf("GetAll() failed: %v", err)
	}

	// CHeck if we get 10 workshops in the database and the ones we just added
	if len(workshops) < numWorkshops {
		t.Errorf("Expected at least %d workshops, got %d", numWorkshops, len(workshops))
	}
	// Check if we got the expected number of workshops
	if len(workshops) != numWorkshops { // +7 for existing seed data
		t.Errorf("Expected %d workshops, got %d", numWorkshops, len(workshops))
	}

	if metadata.TotalRecords != numWorkshops {
		t.Errorf("Expected total records %d, got %d", numWorkshops, metadata.TotalRecords)
	}
}
