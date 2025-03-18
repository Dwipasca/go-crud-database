package models

import (
	"time"
)

type User struct {
	UserId    int       `json:"userId"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdateAt  time.Time `json:"updateAt"`
}

type DetailUser struct {
	UserId   int       `json:"userId"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdateAt  time.Time `json:"updateAt"`
}