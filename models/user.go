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

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	UserId   int    `json:"userId"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
}

type PaginationMeta struct {
	CurrentPage int `json:"currentPage"`
	Limit 	 	int `json:"limit"`
	TotalItems 	int `json:"totalItems"`
	TotalPage   int `json:"totalPage"`
}

type PaginatedResponse[T any] struct {
	Message 	string      	`json:"message"`
	Status  	string      	`json:"status"`
	Code    	int         	`json:"code"`
	Data    	interface{} 	`json:"data,omitempty"`
	Pagination 	PaginationMeta 	`json:"pagination,omitempty"`
}