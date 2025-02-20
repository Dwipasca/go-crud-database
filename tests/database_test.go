package main

import (
	"encoding/json"
	"go-crud-database/config"
	"go-crud-database/handler"
	"go-crud-database/repository"
	"go-crud-database/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestConnectToDB(t *testing.T) {
   // Set environment variables for testing
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "pokemon")
	os.Setenv("DB_SSLMODE", "disable")

	db := config.ConnectToDB()
	if db == nil {
		t.Fatal("Expected a valid database connection, got nil")
	}

	// Check if the connection is valid
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}

	defer db.Close()
}

func TestGetAllUser(t *testing.T) {
    // Set environment variables for testing
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "pokemon")
	os.Setenv("DB_SSLMODE", "disable")

	// Create a test server and register our handler
	req, err := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	recorder := httptest.NewRecorder()

	// Initialize your database connection and repository
	db := config.ConnectToDB()
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepo)

	// Call the function we want to test
	userHandler.GetAllUser(recorder, req)

	// Check the status code
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	var response utils.Response
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the response status is "success"
	if response.Status != "success" {
		t.Errorf("handler returned wrong status: got %v want %v", response.Status, "success")
	}

	// check the Data field to ensure it contains users
	if response.Data == nil {
		t.Error("handler returned no user data")
	}
}

