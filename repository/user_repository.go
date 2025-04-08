package repository

import (
	"context"
	"database/sql"
	"go-crud-database/models"
)

type UserRepository interface {
	GetAllUser(ctx context.Context, limit, offset int) ([]models.User, error)
	GetUserById(ctx context.Context, tx *sql.Tx, id string) (models.DetailUser, error)
	GetUserByUsername(ctx context.Context, tx *sql.Tx, username string) (models.User, error)
	Register(ctx context.Context, tx *sql.Tx, user *models.RegisterRequest) error
	Authentication(ctx context.Context, user *models.LoginRequest) (bool, error)
	UpdateUser(ctx context.Context, tx *sql.Tx, user *models.UpdateUserRequest) error
	DeleteUser(ctx context.Context, tx *sql.Tx, id string) error
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CheckUserExists(ctx context.Context, id string) (bool, error)
	CountUser(ctx context.Context) (int, error)
}
