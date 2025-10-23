package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Pedro-J-Kukul/police_training/internal/data"
	_ "github.com/lib/pq"
)

var testApp *appDependencies

func TestMain(m *testing.M) {
	db, err := sql.Open("postgres", "postgres://police:police@localhost/police_training_testing?sslmode=disable")
	if err != nil {
		panic(err)
	}

	testApp = &appDependencies{
		models: data.Models{User: data.UserModel{DB: db}},
	}

	code := m.Run()
	db.Close()
	os.Exit(code)
}

func TestRegisterUserHandler(t *testing.T) {
	input := map[string]any{
		"first_name":     "Test",
		"last_name":      "User",
		"email":          "testuser@example.com",
		"gender":         "m",
		"password":       "StrongPass1!",
		"is_facilitator": false,
		"is_officer":     true,
	}

	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	testApp.registerUserHandler(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201 Created; got %d", res.StatusCode)
	}

	var response map[string]any
	_ = json.NewDecoder(res.Body).Decode(&response)
	if response["user"] == nil {
		t.Error("expected user object in response")
	}
}
