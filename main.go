package main

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type User struct {
	UserId 		int	 		`json:"userId"`
	Username 	string  	`json:"username"`
	Email 		string  	`json:"email"`
	CreatedAt 	time.Time  	`json:"createdAt"`
}

type Response struct {
	Message 	string      `json:"message"`
	Status  	string      `json:"status"`
	Code    	int         `json:"code"`
	Data    	interface{} `json:"data,omitempty"`
}

func GetAllUser (w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodGet {
		log.Printf("Method not allowed: %s", r.Method)
		WriteJson(w, http.StatusMethodNotAllowed, "error",  nil, "Method Not Allowed") // Respond with 405
		return
	}

	db := ConnectToDB()
	defer db.Close()

	ctx := context.Background()

	sqlQuery := "SELECT user_id, username, email, created_at from users"
	rows, err := db.QueryContext(ctx, sqlQuery)
	if err != nil {
		WriteJson(w, http.StatusInternalServerError, "error",  nil, "Internal Server Error")
		return
	}

	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.UserId, &user.Username, &user.Email, &user.CreatedAt)
		if err != nil {
			WriteJson(w, http.StatusInternalServerError, "error",  nil, "Error reading data")
			return
		}
		users = append(users, user)
	}
	defer rows.Close()

	// check if there any error in iteration
    if err := rows.Err(); err != nil {
		WriteJson(w, http.StatusInternalServerError, "error",  nil, "Error iterating rows")
        return
    }

	WriteJson(w, http.StatusOK, "success",  users, "Successfully get all data from table users")

} 

func GetUserById ( w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodGet {
		log.Printf("Method not allowed: %s", r.Method)
		WriteJson(w, http.StatusMethodNotAllowed, "error",  nil, "Method Not Allowed") // Respond with 405
		return
	}

	db := ConnectToDB()
	defer db.Close()

	ctx := context.Background()
	id := r.URL.Query().Get("id")
	if id == "" {
		WriteJson(w, http.StatusBadRequest, "error", nil, "missing user id")
		return
	}

	var user User
	sqlQuery := "SELECT user_id, username, email, created_at from users where user_id = $1"
	err := db.QueryRowContext(ctx, sqlQuery, id).Scan(&user.UserId, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		WriteJson(w, http.StatusNotFound, "error", nil, "User not found")
		return
	}

	WriteJson(w, http.StatusOK, "success", user, "Successfully get detail data")

}

func CreateNewUser(w http.ResponseWriter, r *http.Request) {
    
	if r.Method != http.MethodPost {
		log.Printf("Method not allowed: %s", r.Method)
		WriteJson(w, http.StatusMethodNotAllowed, "error",  nil, "Method Not Allowed") // Respond with 405
		return
	}

	db := ConnectToDB()
    defer db.Close()

    var newUser User
    if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		WriteJson(w, http.StatusBadRequest, "error",  nil, "Invalid request payload")
        return
    }

    ctx := context.Background()

    // check if username or email already exists
    var usernameExists bool
    var emailExists bool

    checkUsernameQuery := "SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)"
    errUsername := db.QueryRowContext(ctx, checkUsernameQuery, newUser.Username).Scan(&usernameExists)
    if errUsername != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    checkEmailQuery := "SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)"
    errEmail := db.QueryRowContext(ctx, checkEmailQuery, newUser.Email).Scan(&emailExists)
    if errEmail != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    if usernameExists {
        WriteJson(w, http.StatusConflict, "error", nil, "username already exists")
        return
    }

    if emailExists {
        WriteJson(w, http.StatusConflict, "error", nil, "email already exists")
        return
    }

	// check if username or email in request is empty or not
	if newUser.Username == "" {
		WriteJson(w, http.StatusConflict, "error", nil, "username cannot be empty")
		return
	}
	if  newUser.Email == "" {
		WriteJson(w, http.StatusConflict, "error", nil, "Email cannot be empty")
		return
	}

	// ======= end validation =================

    insertQuery := "INSERT INTO users (username, email, created_at) VALUES ($1, $2, $3) RETURNING user_id, created_at"
    errInsert := db.QueryRowContext(ctx, insertQuery, newUser.Username, newUser.Email, time.Now()).Scan(&newUser.UserId, &newUser.CreatedAt)
    if errInsert != nil {
		WriteJson(w, http.StatusInternalServerError, "error", nil, "Internal Server Error")
        return
    }

    WriteJson(w, http.StatusCreated, "success", newUser, "New user created successfully")
}

func UpdateDataUser(w http.ResponseWriter, r *http.Request) {
    
	if r.Method != http.MethodPut {
		log.Printf("Method not allowed: %s", r.Method)
		WriteJson(w, http.StatusMethodNotAllowed, "error",  nil, "Method Not Allowed") // Respond with 405
		return
	}

	db := ConnectToDB()
    defer db.Close()

    var updatedUser User
    if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		WriteJson(w, http.StatusBadRequest, "error",  nil, "Invalid request payload")
        return
    }

    // Ensure user ID is provided
    if updatedUser.UserId == 0 {
        WriteJson(w, http.StatusBadRequest, "error",  nil, "Missing user ID")
        return
    }

    ctx := context.Background()

    var currentUser User
    fetchQuery := "SELECT user_id, username, email FROM users WHERE user_id = $1"
    err := db.QueryRowContext(ctx, fetchQuery, updatedUser.UserId).Scan(&currentUser.UserId, &currentUser.Username, &currentUser.Email)
    
    if err != nil {
        if err == sql.ErrNoRows {
            WriteJson(w, http.StatusNotFound, "error", nil, "user not found")
            return
        }
        WriteJson(w, http.StatusInternalServerError, "error", nil, "Internal Server Error")
        return
    }

    // Check if there have any changes
    if currentUser.Username == updatedUser.Username && currentUser.Email == updatedUser.Email {
        WriteJson(w, http.StatusOK, "info", nil, "No changes detected for the user")
        return
    }

    updateQuery := "UPDATE users SET username = $1, email = $2 WHERE user_id = $3"
    res, err := db.ExecContext(ctx, updateQuery, updatedUser.Username, updatedUser.Email, updatedUser.UserId)
    if err != nil {
		WriteJson(w, http.StatusInternalServerError, "error", nil, "Internal Server Error")
        return
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
		WriteJson(w, http.StatusInternalServerError, "error", nil, "Internal Server Error")
        return
    }

    if rowsAffected == 0 {
        WriteJson(w, http.StatusNotFound, "error", nil, "user not found")
        return
    }

    WriteJson(w, http.StatusOK, "success", updatedUser, "User updated successfully")
}

func DeleteDataUser(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodDelete {
		log.Printf("Method not allowed: %s", r.Method)
		WriteJson(w, http.StatusMethodNotAllowed, "error",  nil, "Method Not Allowed") // Respond with 405
		return
	}

	db := ConnectToDB()
	defer db.Close()
	
	ctx := context.Background()
	id := r.URL.Query().Get("id")
	if id == "" {
		WriteJson(w, http.StatusBadRequest, "error", nil, "missing user id")
		return
	}

	deleteQuery := "DELETE FROM users WHERE user_id = $1"
    res, err := db.ExecContext(ctx, deleteQuery, id)
    if err != nil {
		WriteJson(w, http.StatusInternalServerError, "error", nil, "Internal Server Error")
        return
    }

	rowsAffected, err := res.RowsAffected()
    if err != nil {
		WriteJson(w, http.StatusInternalServerError, "error", nil, "Internal Server Error")
        return
    }

    if rowsAffected == 0 {
        WriteJson(w, http.StatusNotFound, "error", nil, "user is not exists")
        return
    }

    WriteJson(w, http.StatusOK, "success", nil, "Id User "+id+" deleted successfully")
}



func ConnectToDB() *sql.DB {
	connStr := "user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_NAME") + " sslmode=" + os.Getenv("DB_SSLMODE")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	// Create a user table if it doesn't exist
	createTableQuery := `CREATE TABLE IF NOT EXISTS users (
		user_id serial primary key,
        username varchar(50) unique not null,
        email varchar(100) unique not null,
        created_at timestamp default current_timestamp
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	// database pooling
	db.SetMaxIdleConns(10) // jumlah minimal koneksi yg dibuat
	db.SetMaxOpenConns(100) // jumlah maksimal koneksi yg dibuat
	db.SetConnMaxIdleTime(5 * time.Minute) // jika dalam waktu tertentu tdk digunakan maka akan dihapus
	db.SetConnMaxLifetime(60 * time.Minute) // membuat koneksi baru setelah waktu yg telah ditentukan

	return db
}

func WriteJson(w http.ResponseWriter, code int,  status string, data interface{},  message string ) {
	var req = Response{
		Code:    code,
		Status:  status,
		Data:    data,
		Message: message,
	}

	res, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(res)
	if err != nil {
		panic(err)
	}
}

func loadEnv(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "#") {
            parts := strings.SplitN(line, "=", 2)
            if len(parts) == 2 {
                os.Setenv(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
            }
        }
    }

    return scanner.Err()
}

func main() {
	err := loadEnv(".env")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		switch r.Method {
			case http.MethodGet:
				if id == "" {
					GetAllUser(w,r)
				} else {
					GetUserById(w,r)
				}
			case http.MethodPost: 
				CreateNewUser(w,r)
			case http.MethodPut:
				UpdateDataUser(w,r)
			case http.MethodDelete: 
				DeleteDataUser(w,r)
		}

	})

	PORT := "8080"
	http.ListenAndServe(":"+PORT, nil)
}