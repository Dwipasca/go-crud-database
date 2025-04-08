package main

import (
	"context"
	"database/sql"
	"go-crud-database/config"
	"go-crud-database/models"
	"go-crud-database/repository"
	"log"
	"os"
	"strconv"
	"testing"

	_ "github.com/lib/pq"
)

var testDB *sql.DB
var userRepo repository.UserRepository

func TestMain(m *testing.M) {
	// Set environment variables for test DB
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "")
	os.Setenv("DB_NAME", "users_test")
	os.Setenv("DB_SSLMODE", "disable")

	// Connect to test DB
	testDB = config.ConnectToDB()
	if testDB == nil {
		log.Fatal("Failed to connect to test database")
	}

	// Assign repository
	userRepo = repository.NewUserRepository(testDB)

	// Run tests
	code := m.Run()

	// Close DB connection after all tests
	testDB.Close()

	os.Exit(code)
}

func TestConnectToDB(t *testing.T) {
	if err := testDB.Ping(); err != nil {
		t.Fatalf("Expected valid DB connection, got error: %v", err)
	}
}

func TestGetAllUser(t *testing.T) {
	// Buat context kosong, bisa diganti kalau handler punya context tambahan
	ctx := context.Background()

	users, err := userRepo.GetAllUser(ctx)
	if err != nil {
		t.Fatalf("Failed to get users: %v", err)
	}

	// Contoh verifikasi: pastikan slice tidak nil (jumlah bisa kosong tergantung data)
	if users == nil {
		t.Error("Expected users slice, got nil")
	}
}

func TestRegisterAndGetUserByUsername(t *testing.T) {
	ctx := context.Background()

	tx, err := testDB.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Dummy data
	newUser := &models.RegisterRequest{
		Username: "testuser_integration",
		Email:    "testuser_integration@example.com",
		Password: "hashedpassword",
		IsAdmin:  false,
	}

	txRepo := repository.NewUserRepository(testDB)

	err = txRepo.Register(ctx, tx, newUser)
	if err != nil {
		t.Fatalf("Failed to register user in TX: %v", err)
	}

	user, err := txRepo.GetUserByUsername(ctx, tx, newUser.Username)
	if err != nil {
		t.Fatalf("Failed to get user by username in TX: %v", err)
	}

	if user.Username != newUser.Username {
		t.Errorf("Expected username %s, got %s", newUser.Username, user.Username)
	}
}

func TestUpdateAndGetUserById(t *testing.T) {
	ctx := context.Background()

	tx, err := testDB.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	newUser := &models.RegisterRequest{
		Username: "testuser_integration",
		Email:    "testuser_integration@example.com",
		Password: "hashedpassword",
		IsAdmin:  false,
	}

	txRepo := repository.NewUserRepository(testDB)

	// create a new user
	err = txRepo.Register(ctx, tx, newUser)
	if err != nil {
		t.Fatalf("Failed to register user in TX: %v", err)
	}

	// get data user by username
	user, err := txRepo.GetUserByUsername(ctx, tx, newUser.Username)
	if err != nil {
		t.Fatalf("Failed to get user by username in TX: %v", err)
	}

	updateUser := &models.UpdateUserRequest{
		UserId:   user.UserId,
		Username: "updateduser_integration",
		Email: "updateuser_integration",
		IsAdmin:  true,
	}

	// Update user
	err = txRepo.UpdateUser(ctx, tx, updateUser)
	if err != nil {
		t.Fatalf("Failed to update user in TX: %v", err)
	}

	// Get user by ID
	detailUser, err := txRepo.GetUserById(ctx, tx, strconv.Itoa(user.UserId))
	if err != nil {
		t.Fatalf("Failed to get user by ID in TX: %v", err)
	}

	// check if the updated user data is correct
	if detailUser.Username != updateUser.Username {
		t.Errorf("Expected username %s, got %s", updateUser.Username, detailUser.Username)
	}
}
