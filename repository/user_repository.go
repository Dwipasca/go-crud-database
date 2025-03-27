package repository

import (
	"context"
	"go-crud-database/models"
)

type UserRepository interface {
	GetAllUser(ctx context.Context) ([]models.User, error)
	GetUserById(ctx context.Context, id string) (models.DetailUser, error)
	Register(ctx context.Context, user *models.RegisterRequest) error
	Authentication(ctx context.Context, user *models.LoginRequest) (bool, error)
	UpdateUser(ctx context.Context, user *models.UpdateUserRequest) error
	DeleteUser(ctx context.Context, id string) error
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CheckUserExists(ctx context.Context, id string) (bool, error)
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
}
