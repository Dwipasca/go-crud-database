package repository

import (
	"context"
	"database/sql"
	"go-crud-database/models"
)

type userRepositoryImpl struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryImpl{DB: db}
}

func (r *userRepositoryImpl) GetAllUser(ctx context.Context) ([]models.User, error) {

	sqlQuery := "SELECT user_id, username, email, created_at from users"
	rows, err := r.DB.QueryContext(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.UserId, &user.Username, &user.Email, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepositoryImpl) GetUserById(ctx context.Context, id string) (models.User, error) {
	sqlQuery := "SELECT user_id, username, email, created_at from users where user_id = $1"
	var user models.User
	err := r.DB.QueryRowContext(ctx, sqlQuery, id).Scan(&user.UserId, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *userRepositoryImpl) CreateUser(ctx context.Context, user *models.User)  error {
	sqlQuery := "INSERT INTO users(username, email) VALUES ($1, $2)"
	_, err := r.DB.ExecContext(ctx, sqlQuery, user.Username, user.Email)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepositoryImpl) UpdateUser(ctx context.Context, user *models.User) error {
	sqlQuery := "UPDATE users SET username = $1, email = $2 where user_id = $3"
	_, err := r.DB.ExecContext(ctx, sqlQuery, user.Username, user.Email, user.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepositoryImpl) DeleteUser(ctx context.Context, id string) error {
	sqlQuery := "DELETE FROM users where user_id = $1"
	_, err := r.DB.ExecContext(ctx, sqlQuery, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepositoryImpl) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	var emailExists bool
	sqlQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	err := r.DB.QueryRowContext(ctx, sqlQuery, email).Scan(&emailExists)
	if err != nil {
		return false, err
	}

	return emailExists, nil

}

func (r *userRepositoryImpl) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	var usernameExists bool
	sqlQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)"
	err := r.DB.QueryRowContext(ctx, sqlQuery, username).Scan(&usernameExists)
	if err != nil {
		return false, err
	}

	return usernameExists, nil
}

func (r *userRepositoryImpl) CheckUserExists(ctx context.Context, id string) (bool, error) {
	var userIdExists bool
	sqlQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)"
	err := r.DB.QueryRowContext(ctx, sqlQuery, id).Scan(&userIdExists)
	if err != nil {
		return false, err
	}

	return userIdExists, nil
}