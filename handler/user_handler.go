package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"go-crud-database/models"
	"go-crud-database/repository"
	"go-crud-database/utils"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func (h *UserHandler) Authentication(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJson(w, http.StatusMethodNotAllowed, "error", nil, "Method Not Allowed")
		return
	}

	var user models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, "error", nil, "Invalid request payload")
		return
	}

	if msg, isValid := utils.ValidateLoginRequest(user); !isValid {
		utils.WriteJson(w, http.StatusConflict, "error", nil, msg)
		return
	}

	ctx := context.Background()
	storedUser, err := h.repo.GetUserByUsername(ctx, user.Username)
	if err != nil {
		// if the user is not found, return an error message
		if err == sql.ErrNoRows {
			utils.WriteJson(w, http.StatusUnauthorized, "error", nil, "Invalid username or password")
			return
		}
		// if the error is not sql.ErrNoRows, return an internal server error message
		log.Println("error getting user by username: ", err)
		utils.WriteJson(w, http.StatusInternalServerError, "error", nil, "Internal Server Error")
		return
	}

	if !utils.CheckPassword(storedUser.Password, user.Password) {
		utils.WriteJson(w, http.StatusUnauthorized, "error", nil, "Invalid username or password")
		return
	}

	// token will expire after 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)

	// create a new token with the claims and the signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":  storedUser.UserId,
		"isAdmin": storedUser.IsAdmin,
		"exp":     expirationTime.Unix(),
	})

	// sign the token with the secret key and get the token string
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, "error", nil, "Internal Server Error")
		return
	}

	utils.WriteJson(w, http.StatusOK, "success", tokenString, "Authentication successful")
}

func (h *UserHandler) GetAllUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJson(w, http.StatusMethodNotAllowed, "error", nil, "Method Not Allowed")
		return
	}

	// check role user
	isAdmin, ok := r.Context().Value("isAdmin").(bool)
	if !ok || !isAdmin {
		utils.WriteJson(w, http.StatusUnauthorized, "error", nil, "Unauthorized only admin can access")
		return
	}

	ctx := context.Background()
	users, err := h.repo.GetAllUser(ctx)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, "error", nil, "Internal Server Error")
		return
	}

	if len(users) == 0 {
		utils.WriteJson(w, http.StatusOK, "info", "No users found", "No users found")
		return
	}

	utils.WriteJson(w, http.StatusOK, "success", users, "Successfully retrieved all users")
}

func (h *UserHandler) UpdateDataUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteJson(w, http.StatusMethodNotAllowed, "error", nil, "Method Not Allowed")
		return
	}

	// check role user
	isAdmin, ok := r.Context().Value("isAdmin").(bool)
	if !ok || !isAdmin {
		utils.WriteJson(w, http.StatusUnauthorized, "error", nil, "Unauthorized only admin can access")
		return
	}

	var updatedUser models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, "error", nil, "Invalid request payload")
		return
	}

	if msg, isValid := utils.ValidateUpdateUserRequest(updatedUser); !isValid {
		utils.WriteJson(w, http.StatusConflict, "error", nil, msg)
		return
	}

	ctx := context.Background()

	if updatedUser.UserId == 0 {
		utils.WriteJson(w, http.StatusBadRequest, "error", nil, "Missing user ID")
		return
	}

	// Check if the user exists
	detailUser, err := h.repo.GetUserById(ctx, strconv.Itoa(updatedUser.UserId))
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJson(w, http.StatusNotFound, "error", nil, "User not found")
			return
		}
		utils.WriteJson(w, http.StatusInternalServerError, "error", nil, "Internal Server Error")
		return
	}

	// Check if the data has changed
	if detailUser.Username == updatedUser.Username && detailUser.Email == updatedUser.Email && detailUser.IsAdmin == updatedUser.IsAdmin {
		utils.WriteJson(w, http.StatusOK, "info", nil, "No changes detected for the user")
		return
	}

	// Check if the username already exists
	if updatedUser.Username != detailUser.Username {
		usernameExists, err := h.repo.CheckUsernameExists(ctx, updatedUser.Username)
		if err != nil {
			utils.WriteJson(w, http.StatusInternalServerError, "error", nil, "Internal Server Error")
			return
		}
		if usernameExists {
			utils.WriteJson(w, http.StatusConflict, "error", nil, "Username already exists")
			return
		}
	}

	// Proceed with the update
	err = h.repo.UpdateUser(ctx, &updatedUser)
	if err != nil {
		log.Printf("Error updating user with ID %d: %v", updatedUser.UserId, err)
		utils.WriteJson(w, http.StatusInternalServerError, "error", nil, "update error")
		return
	}

	utils.WriteJson(w, http.StatusOK, "success", nil, "User updated successfully")

}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJson(w, http.StatusMethodNotAllowed, "error", nil, "Method Not Allowed")
		return
	}

	var newUser models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		log.Println("error decoding request body: ", err)
		utils.WriteJson(w, http.StatusBadRequest, "error", nil, "Invalid request payload")
		return
	}

	if msg, isValid := utils.ValidateRegisterRequest(newUser); !isValid {
		utils.WriteJson(w, http.StatusConflict, "error", nil, msg)
		return
	}

	ctx := context.Background()

	// Check if username already exists
	usernameExists, err := h.repo.CheckUsernameExists(ctx, newUser.Username)
	if err != nil {
		log.Println("error checking username exists: ", err)
		utils.WriteJson(w, http.StatusInternalServerError, "error", nil, "Internal Server Error")
		return
	}

	// check if email already exists
	emailExists, err := h.repo.CheckEmailExists(ctx, newUser.Email)
	if err != nil {
		log.Println("error checking email exists: ", err)
		utils.WriteJson(w, http.StatusInternalServerError, "error", nil, "Internal Server Error")
		return
	}

	if usernameExists {
		utils.WriteJson(w, http.StatusConflict, "error", nil, "Username already exists")
		return
	}

	if emailExists {
		utils.WriteJson(w, http.StatusConflict, "error", nil, "Email already exists")
		return
	}

	passwordHash, err := utils.EncryptPassword(newUser.Password)
	if err != nil {
		log.Println("error hashing password: ", err)
		utils.WriteJson(w, http.StatusInternalServerError, "error", nil, "Internal Server Error")
		return
	}

	// convert the byte slice to a string
	newUser.Password = string(passwordHash)

	err = h.repo.Register(ctx, &newUser)
	if err != nil {
		log.Println("error creating new user: ", err)
		utils.WriteJson(w, http.StatusInternalServerError, "error", nil, "Internal Server Error")
		return
	}

	utils.WriteJson(w, http.StatusCreated, "success", nil, "New user created successfully")
}

func (h *UserHandler) DeleteDataUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteJson(w, http.StatusMethodNotAllowed, "error", nil, "Method Not Allowed")
		return
	}

	// check role user
	isAdmin, ok := r.Context().Value("isAdmin").(bool)
	if !ok || !isAdmin {
		utils.WriteJson(w, http.StatusUnauthorized, "error", nil, "Unauthorized only admin can access")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		utils.WriteJson(w, http.StatusBadRequest, "error", nil, "missing user id")
		return
	}

	ctx := context.Background()

	// Check if the user exists
	userIdExists, err := h.repo.CheckUserExists(ctx, id)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, "error", nil, "Internal Server Error")
		return
	}

	if !userIdExists {
		utils.WriteJson(w, http.StatusConflict, "error", nil, "user not found")
		return
	}

	err = h.repo.DeleteUser(ctx, id)
	if err != nil {
		log.Printf("Error deleting user with ID %s: %v", id, err)
		utils.WriteJson(w, http.StatusInternalServerError, "error", nil, "Internal Server Error during deletion")
		return
	}

	utils.WriteJson(w, http.StatusOK, "success", nil, "User deleted successfully")
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJson(w, http.StatusMethodNotAllowed, "error", nil, "Method Not Allowed")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		utils.WriteJson(w, http.StatusBadRequest, "error", nil, "missing user id")
		return
	}

	ctx := context.Background()
	user, err := h.repo.GetUserById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJson(w, http.StatusNotFound, "error", nil, "User not found")
			return
		}
		utils.WriteJson(w, http.StatusInternalServerError, "error", nil, "Internal Server Error")
		return
	}

	utils.WriteJson(w, http.StatusOK, "success", user, "Successfully retrieved user details")
}
