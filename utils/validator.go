package utils

import (
	"go-crud-database/models"
	"regexp"
	"strings"
)

// email regex pattern to validate email
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateRegisterRequest(req models.RegisterRequest) (string, bool) {
	if strings.TrimSpace(req.Username) == "" {
		return "Username cannot be empty", false
	}

	if strings.TrimSpace(req.Email) == "" {
		return "Email cannot be empty", false
	}

	if strings.TrimSpace(req.Password) == "" {
		return "Password cannot be empty", false
	}

	if len(req.Password) < 5 {
		return "Password must be at least 5 characters", false
	}

	if !emailRegex.MatchString(req.Email) {
		return "Invalid email format", false
	}

	return "", true
}

func ValidateLoginRequest(req models.LoginRequest) (string, bool) {
	if strings.TrimSpace(req.Username) == "" {
		return "Username cannot be empty", false
	}

	if strings.TrimSpace(req.Password) == "" {
		return "Password cannot be empty", false
	}

	if len(req.Password) < 5 {
		return "Password must be at least 5 characters", false
	}

	return "", true
}

func ValidateUpdateUserRequest(req models.UpdateUserRequest) (string, bool) {

	if strings.TrimSpace(req.Username) == "" {
		return "Username cannot be empty", false
	}

	if strings.TrimSpace(req.Email) == "" {
		return "Email cannot be empty", false
	}

	if strings.TrimSpace(req.Password) == "" {
		return "Password cannot be empty", false
	}

	if len(req.Password) < 5 {
		return "Password must be at least 5 characters", false
	}

	if !emailRegex.MatchString(req.Email) {
		return "Invalid email format", false
	}

	return "", true
}
