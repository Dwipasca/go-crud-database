package main

import (
	"go-crud-database/config"
	"os"
	"testing"
)

func TestConnectToDB(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "")
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
