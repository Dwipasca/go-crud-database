package config

import (
	"bufio"
	"database/sql"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func ConnectToDB() *sql.DB {

	connStr := "user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_NAME") + " sslmode=" + os.Getenv("DB_SSLMODE") + " host=" + os.Getenv("DB_HOST") + " port=" + os.Getenv("DB_PORT")

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
		password varchar(255) not null,
		is_admin boolean default false,
        created_at timestamp default current_timestamp,
		updated_at timestamp default current_timestamp
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	// database pooling
	db.SetMaxIdleConns(10)                  // jumlah minimal koneksi yg dibuat
	db.SetMaxOpenConns(100)                 // jumlah maksimal koneksi yg dibuat
	db.SetConnMaxIdleTime(5 * time.Minute)  // jika dalam waktu tertentu tdk digunakan maka akan dihapus
	db.SetConnMaxLifetime(60 * time.Minute) // membuat koneksi baru setelah waktu yg telah ditentukan

	return db
}

func LoadEnv(filename string) error {
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
