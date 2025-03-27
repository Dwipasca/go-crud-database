package main

import (
	"go-crud-database/models"
	"go-crud-database/utils"
	"testing"
)

func TestValidateRegisterRequest(t *testing.T) {

	testCases := []struct {
		name     string
		input    models.RegisterRequest
		wantMsg  string
		wantBool bool
	}{
		{
			name:     "Valid input",
			input:    models.RegisterRequest{Username: "user123", Email: "user@example.com", Password: "password"},
			wantMsg:  "",
			wantBool: true,
		},
		{
			name:     "Empty username",
			input:    models.RegisterRequest{Username: "", Email: "user@example.com", Password: "password"},
			wantMsg:  "Username cannot be empty",
			wantBool: false,
		},
		{
			name:     "Empty email",
			input:    models.RegisterRequest{Username: "user123", Email: "", Password: "password"},
			wantMsg:  "Email cannot be empty",
			wantBool: false,
		},
		{
			name:     "Empty password",
			input:    models.RegisterRequest{Username: "user123", Email: "user@example.com", Password: ""},
			wantMsg:  "Password cannot be empty",
			wantBool: false,
		},
		{
			name:     "Short password",
			input:    models.RegisterRequest{Username: "user123", Email: "user@example.com", Password: "123"},
			wantMsg:  "Password must be at least 5 characters",
			wantBool: false,
		},
		{
			name:     "Invalid email format (missing @)",
			input:    models.RegisterRequest{Username: "user123", Email: "userexample.com", Password: "password"},
			wantMsg:  "Invalid email format",
			wantBool: false,
		},
		{
			name:     "Invalid email format (missing domain)",
			input:    models.RegisterRequest{Username: "user123", Email: "user@.com", Password: "password"},
			wantMsg:  "Invalid email format",
			wantBool: false,
		},
		{
			name:     "Invalid email format (missing TLD)",
			input:    models.RegisterRequest{Username: "user123", Email: "user@example", Password: "password"},
			wantMsg:  "Invalid email format",
			wantBool: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			gotMsg, gotBool := utils.ValidateRegisterRequest(test.input)
			if gotMsg != test.wantMsg || gotBool != test.wantBool {
				t.Errorf("ValidateRegisterRequest(%v) = (%v, %v), want (%v, %v)", test.input, gotMsg, gotBool, test.wantMsg, test.wantBool)
			}
		})
	}
}

func TestValidateLoginRequest(t *testing.T) {

	testCases := []struct {
		name     string
		input    models.LoginRequest
		wantMsg  string
		wantBool bool
	}{
		{
			name:     "Valid input",
			input:    models.LoginRequest{Username: "user123", Password: "password"},
			wantMsg:  "",
			wantBool: true,
		},
		{
			name:     "Empty username",
			input:    models.LoginRequest{Username: "", Password: "password"},
			wantMsg:  "Username cannot be empty",
			wantBool: false,
		},
		{
			name:     "Empty password",
			input:    models.LoginRequest{Username: "user123", Password: ""},
			wantMsg:  "Password cannot be empty",
			wantBool: false,
		},
		{
			name:     "Short password",
			input:    models.LoginRequest{Username: "user123", Password: "123"},
			wantMsg:  "Password must be at least 5 characters",
			wantBool: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			gotMsg, gotBool := utils.ValidateLoginRequest(test.input)
			if gotMsg != test.wantMsg || gotBool != test.wantBool {
				t.Errorf("ValidateLoginRequest(%v) = (%v, %v), want (%v, %v)", test.input, gotMsg, gotBool, test.wantMsg, test.wantBool)
			}
		})
	}

}

func TestValidateUpdateUserRequest(t *testing.T) {

	tests := []struct {
		name     string
		input    models.UpdateUserRequest
		wantMsg  string
		wantBool bool
	}{
		{
			name:     "Valid input",
			input:    models.UpdateUserRequest{Username: "user123", Email: "user@example.com", Password: "password"},
			wantMsg:  "",
			wantBool: true,
		},
		{
			name:     "Empty username",
			input:    models.UpdateUserRequest{Username: "", Email: "user@example.com", Password: "password"},
			wantMsg:  "Username cannot be empty",
			wantBool: false,
		},
		{
			name:     "Empty email",
			input:    models.UpdateUserRequest{Username: "user123", Email: "", Password: "password"},
			wantMsg:  "Email cannot be empty",
			wantBool: false,
		},
		{
			name:     "Empty password",
			input:    models.UpdateUserRequest{Username: "user123", Email: "user@example.com", Password: ""},
			wantMsg:  "Password cannot be empty",
			wantBool: false,
		},
		{
			name:     "Short password",
			input:    models.UpdateUserRequest{Username: "user123", Email: "user@example.com", Password: "123"},
			wantMsg:  "Password must be at least 5 characters",
			wantBool: false,
		},
		{
			name:     "Invalid email format (missing @)",
			input:    models.UpdateUserRequest{Username: "user123", Email: "userexample.com", Password: "password"},
			wantMsg:  "Invalid email format",
			wantBool: false,
		},
		{
			name:     "Invalid email format (missing domain)",
			input:    models.UpdateUserRequest{Username: "user123", Email: "user@.com", Password: "password"},
			wantMsg:  "Invalid email format",
			wantBool: false,
		},
		{
			name:     "Invalid email format (missing TLD)",
			input:    models.UpdateUserRequest{Username: "user123", Email: "user@example", Password: "password"},
			wantMsg:  "Invalid email format",
			wantBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMsg, gotBool := utils.ValidateUpdateUserRequest(tt.input)
			if gotMsg != tt.wantMsg || gotBool != tt.wantBool {
				t.Errorf("ValidateUpdateUserRequest(%v) = (%v, %v), want (%v, %v)", tt.input, gotMsg, gotBool, tt.wantMsg, tt.wantBool)
			}
		})
	}
}
