package main

import (
	"go-crud-database/utils"
	"testing"
	"time"
)

func TestHashPassword_NormalPassword(t *testing.T) {
	password := "password"
	hashedPassword, err := utils.EncryptPassword(password)
	if err != nil {
		t.Errorf("Error hashing password: %v", err)
	}

	if hashedPassword == password {
		t.Errorf("Hashed password should not be the same as the original password")
	}
}

func TestHashPassword_ShouldBeUnique(t *testing.T) {
	password := "password"
	hashedPassword1, _ := utils.EncryptPassword(password)
	hashedPassword2, _ := utils.EncryptPassword(password)

	if hashedPassword1 == hashedPassword2 {
		t.Errorf("Hashed password should be unique")
	}
}

func TestHashPassword_Performance(t *testing.T) {
	password := "PasswordPerformanceTest123"

	start := time.Now()
	_, err := utils.EncryptPassword(password)
	duration := time.Since(start)

	// fmt.Printf("Hashing password took %v\n", duration)

	if err != nil {
		t.Errorf("Error hashing password: %v", err)
	}

	// 500 milliseconds = 0.5 seconds
	if duration > 500*time.Millisecond {
		t.Fatalf("Hashing password took too long: %v", duration)
	}
}

func TestCheckPassword_WrongPassword(t *testing.T) {
	password := "correctPassword"
	hashedPassword, _ := utils.EncryptPassword(password)

	if utils.CheckPassword(hashedPassword, "wrongPassword") {
		t.Errorf("CheckPassword should return false for wrong password")
	}
}
