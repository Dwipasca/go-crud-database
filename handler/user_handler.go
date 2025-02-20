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
)

type UserHandler struct {
	repo repository.UserRepository
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

func (h *UserHandler) CreateNewUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJson(w, http.StatusMethodNotAllowed, "error", nil, "Method Not Allowed")
		return
	}

	var newUser models.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, "error", nil, "Invalid request payload")
		return
	}

	ctx := context.Background()

	// Check if username or email already exists
	usernameExists, err := h.repo.CheckUsernameExists(ctx, newUser.Username)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, "error", nil, "Internal Server Error")
		return
	}

	emailExists, err := h.repo.CheckEmailExists(ctx, newUser.Email)
	if err != nil {
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

	// Check if username or email in request is empty
	if newUser.Username == "" {
		utils.WriteJson(w, http.StatusConflict, "error", nil, "Username cannot be empty")
		return
	}
	if newUser.Email == "" {
		utils.WriteJson(w, http.StatusConflict, "error", nil, "Email cannot be empty")
		return
	}

	err = h.repo.CreateUser(ctx, &newUser)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, "error", nil, "Internal Server Error")
		return
	}

	utils.WriteJson(w, http.StatusCreated, "success", nil, "New user created successfully")
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

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
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

	utils.WriteJson(w, http.StatusOK, "success", users, "Successfully retrieved all users")
}
