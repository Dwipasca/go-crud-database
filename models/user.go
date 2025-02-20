package models

import (
	"time"
)

type User struct {
	UserId    int       `json:"userId"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}