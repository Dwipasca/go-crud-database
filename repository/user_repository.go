package repository

import (
	"context"
	"go-crud-database/models"
)

type UserRepository interface {
	GetAllUser(ctx context.Context) ([]models.User, error)
	GetUserById(ctx context.Context, id string) (models.DetailUser, error)
	Register(ctx context.Context, user *models.User) error
	Authentication(ctx context.Context, username, password string) (bool, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id string) error
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CheckUserExists(ctx context.Context, id string) (bool, error)
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
}



