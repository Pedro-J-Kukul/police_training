package data

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

var testDB *sql.DB

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "postgres://police:police@localhost/police_training_testing?sslmode=disable")
	if err != nil {
		t.Fatalf("failed to connect test db: %v", err)
	}
	return db
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
	defer db.Close()
	model := UserModel{DB: db}

	user := &User{
		FirstName:     "Alice",
		LastName:      "Walker",
		Email:         "alice@example.com",
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

	// Update
	user.FirstName = "Alicia"
	err = model.Update(user)
	if err != nil {
		t.Fatalf("Update() failed: %v", err)
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
	defer db.Close()
	model := UserModel{DB: db}

	user := &User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "testpw@example.com",
		Gender:    "m",
	}
	_ = user.Password.Set("OldPass123!")
	_ = model.Insert(user)

	time.Sleep(100 * time.Millisecond)

	err := model.UpdatePassword(user.ID, "NewPass456@")
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
