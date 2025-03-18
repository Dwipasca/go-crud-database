package main

import (
	"go-crud-database/config"
	"go-crud-database/handler"
	"go-crud-database/middleware"
	"go-crud-database/repository"
	"net/http"
)


func main() {
	err := config.LoadEnv(".env")
	if err != nil {
		panic(err)
	}

	db := config.ConnectToDB()
	defer db.Close()

	// Initialize the User Repository
	userRepo := repository.NewUserRepository(db)

	// Create an instance of UserHandler with the repository
	userHandler := handler.NewUserHandler(userRepo)

	// add middleware to endpoint users
	http.HandleFunc("/api/v1/users", middleware.ValidateToken(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		switch r.Method {
		case http.MethodGet:
			if id == "" {
				userHandler.GetAllUser(w, r)
			} else {
				userHandler.GetUserByID(w, r)
			}
		case http.MethodPut:
			userHandler.UpdateDataUser(w, r)
		case http.MethodDelete:
			userHandler.DeleteDataUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/api/v1/login", userHandler.Authentication)

	http.HandleFunc("/api/v1/register", userHandler.Register)

	PORT := "8080"
	http.ListenAndServe(":"+PORT, nil)
}