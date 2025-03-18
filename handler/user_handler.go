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
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

var jwtKey = []byte("my_secret_key")

func (h *UserHandler) Authentication(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJson(w, http.StatusMethodNotAllowed, "error", nil, "Method Not Allowed")
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, "error", nil, "Invalid request payload")
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
		utils.WriteJson(w, http.StatusInternalServerError, "error", nil, "Internal Server Error")
		return
	}

	// compare the password in the request with the password in the database
	// if the passwords do not match, return an error message
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		utils.WriteJson(w, http.StatusUnauthorized, "error", nil, "Invalid username or password")
		return
	}

	// create a new JWT token
	// the token will expire after 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// create a new claim with the user ID as the subject
	// the subject is the unique identifier for the token
	// the subject is stored as a string, so we need to convert the user ID to a string	
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   strconv.Itoa(storedUser.UserId),
	}

	// create a new token with the claims and the signing method
	// the signing method is HMAC with SHA-256
	// the signing method is used to sign the token with the secret key
	// the secret key is used to verify the token when it is received
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// sign the token with the secret key and get the token string
	// the token string is the token that will be sent to the client
	tokenString, err := token.SignedString(jwtKey)
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

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJson(w, http.StatusMethodNotAllowed, "error", nil, "Method Not Allowed")
		return
	}

	var newUser models.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		log.Println("error decoding request body: ", err)
		utils.WriteJson(w, http.StatusBadRequest, "error", nil, "Invalid request payload")
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

	// Check if username, email, and password in request is empty
	if newUser.Username == "" {
		utils.WriteJson(w, http.StatusConflict, "error", nil, "Username cannot be empty")
		return
	}
	
	if newUser.Email == "" {
		utils.WriteJson(w, http.StatusConflict, "error", nil, "Email cannot be empty")
		return
	}

	if newUser.Password == "" {
		utils.WriteJson(w, http.StatusConflict, "error", nil, "Password cannot be empty")
		return
	}

	// encrypt the password
	// before we hash the password, we need to convert it to a byte slice
	// because the bcrypt.GenerateFromPassword function only accepts a byte slice
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("error hashing password: ", err)
		utils.WriteJson(w, http.StatusInternalServerError, "error", nil, "Internal Server Error")
		return
	}

	// convert the byte slice to a string and store it in the newUser.Password field
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

func (h *UserHandler) UpdateDataUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteJson(w, http.StatusMethodNotAllowed, "error", nil, "Method Not Allowed")
		return
	}

	var updatedUser models.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, "error", nil, "Invalid request payload")
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
	if detailUser.Username == updatedUser.Username && detailUser.Email == updatedUser.Email {
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





