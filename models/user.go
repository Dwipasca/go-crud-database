package models

import (
	"time"
)

type User struct {
	UserId    int       `json:"userId"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	IsAdmin   bool      `json:"isAdmin"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type DetailUser struct {
	UserId    int       `json:"userId"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"isAdmin"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
