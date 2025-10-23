package data

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	// Try multiple environment variable names for flexibility
	dbDSN := os.Getenv("TEST_DATABASE_DSN")
	if dbDSN == "" {
		dbDSN = os.Getenv("TEST_DB_DSN")
	}
	if dbDSN == "" {
		// Fallback to match your .envrc
		dbDSN = "postgres://police:police@localhost/police_training_testing?sslmode=disable"
	}

	var err error
	testDB, err = sql.Open("postgres", dbDSN)
	if err != nil {
		panic(fmt.Sprintf("Failed to open database connection: %v", err))
	}

	// Test the connection
	if err = testDB.Ping(); err != nil {
		panic(fmt.Sprintf("Could not connect to test database: %v\nUsing DSN: %s", err, dbDSN))
	}

	code := m.Run()
	testDB.Close()
	os.Exit(code)
}

func setupTestDB(t *testing.T) *sql.DB {
	return testDB
}

// Helper function to clean up test data
func cleanupTestData(t *testing.T, db *sql.DB) {
	t.Helper()

	// Clean up in reverse order of dependencies
	_, err := db.Exec("TRUNCATE TABLE tokens CASCADE")
	if err != nil {
		t.Logf("Warning: Failed to cleanup tokens: %v", err)
	}

	_, err = db.Exec("TRUNCATE TABLE roles_users CASCADE")
	if err != nil {
		t.Logf("Warning: Failed to cleanup roles_users: %v", err)
	}

	_, err = db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	if err != nil {
		t.Logf("Warning: Failed to cleanup users: %v", err)
	}

	_, err = db.Exec("TRUNCATE TABLE officers RESTART IDENTITY CASCADE")
	if err != nil {
		t.Logf("Warning: Failed to cleanup officers: %v", err)
	}
}

func TestPasswordSetAndMatch(t *testing.T) {
	var p Password
	err := p.Set("StrongPass123!")
	if err != nil {
		t.Fatalf("Set() failed: %v", err)
	}

	match, _ := p.Matches("StrongPass123!")
	if !match {
		t.Error("expected password to match")
	}

	noMatch, _ := p.Matches("WrongPassword!")
	if noMatch {
		t.Error("expected mismatch")
	}
}

func TestInsertGetUpdateDeleteUser(t *testing.T) {
	db := setupTestDB(t)

	// Clean up before and after test
	cleanupTestData(t, db)
	defer cleanupTestData(t, db)

	model := UserModel{DB: db}

	user := &User{
		FirstName:     "Alice",
		LastName:      "Walker",
		Email:         fmt.Sprintf("alice-%d@example.com", time.Now().UnixNano()),
		Gender:        "f",
		IsActivated:   true,
		IsFacilitator: false,
		IsOfficer:     true,
	}
	_ = user.Password.Set("TestPass123!")

	// Insert
	err := model.Insert(user)
	if err != nil {
		t.Fatalf("Insert() failed: %v", err)
	}

	// Get
	fetched, err := model.Get(user.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}
	if fetched.Email != user.Email {
		t.Errorf("expected %v, got %v", user.Email, fetched.Email)
	}

	// Update - Store the current version
	originalVersion := user.Version
	user.FirstName = "Alicia"
	err = model.Update(user)
	if err != nil {
		t.Fatalf("Update() failed: %v", err)
	}

	// Verify the version was updated
	if user.Version <= originalVersion {
		t.Error("Expected version to be incremented after update")
	}

	// Soft delete
	err = model.SoftDelete(user.ID)
	if err != nil {
		t.Fatalf("SoftDelete() failed: %v", err)
	}

	// Restore
	err = model.Restore(user.ID)
	if err != nil {
		t.Fatalf("Restore() failed: %v", err)
	}

	// Hard delete
	err = model.HardDelete(user.ID)
	if err != nil {
		t.Fatalf("HardDelete() failed: %v", err)
	}
}

func TestUpdatePassword(t *testing.T) {
	db := setupTestDB(t)

	// Clean up before and after test
	cleanupTestData(t, db)
	defer cleanupTestData(t, db)

	model := UserModel{DB: db}

	user := &User{
		FirstName: "Test",
		LastName:  "User",
		Email:     fmt.Sprintf("testpw-%d@example.com", time.Now().UnixNano()),
		Gender:    "m",
	}
	_ = user.Password.Set("OldPass123!")
	err := model.Insert(user)
	if err != nil {
		t.Fatalf("Insert() failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	err = model.UpdatePassword(user.ID, "NewPass456@")
	if err != nil {
		t.Fatalf("UpdatePassword() failed: %v", err)
	}

	fetched, err := model.Get(user.ID)
	if err != nil {
		t.Fatalf("Get() failed after UpdatePassword: %v", err)
	}

	match, err := fetched.Password.Matches("NewPass456@")
	if err != nil {
		t.Fatalf("Matches() failed: %v", err)
	}
	if !match {
		t.Error("expected new password to match")
	}
}
