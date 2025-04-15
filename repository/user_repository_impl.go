package repository

import (
	"context"
	"database/sql"
	"go-crud-database/models"
	"go-crud-database/utils"
)

type userRepositoryImpl struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryImpl{DB: db}
}

func (r *userRepositoryImpl) GetAllUser(ctx context.Context, limit, offset int) ([]models.User, error) {

	sqlQuery := "SELECT user_id, username, email, password, is_admin, created_at, updated_at FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2"

	rows, err := r.DB.QueryContext(ctx, sqlQuery, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.UserId, &user.Username, &user.Email, &user.Password, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
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

func (r *userRepositoryImpl) GetUserById(ctx context.Context, tx *sql.Tx, id string) (models.DetailUser, error) {
	sqlQuery := "SELECT user_id, username, email, is_admin, created_at, updated_at from users where user_id = $1"
	
	var user models.DetailUser
	var row *sql.Row

	if tx != nil {
		row = tx.QueryRowContext(ctx, sqlQuery, id)
	} else {
		row = r.DB.QueryRowContext(ctx, sqlQuery, id)
	}

	err := row.Scan(&user.UserId, &user.Username, &user.Email, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	
	return user, err
}

func (r *userRepositoryImpl) Authentication(ctx context.Context, user *models.LoginRequest) (bool, error) {
	sqlQuery := "SELECT password FROM users WHERE username = $1"

	var hashedPassword string
	err := r.DB.QueryRowContext(ctx, sqlQuery, user.Username).Scan(&hashedPassword)
	if err != nil {
		return false, err
	}

	exists := utils.CheckPassword(hashedPassword, user.Password)

	return exists, nil
}

func (r *userRepositoryImpl) Register(ctx context.Context, tx *sql.Tx, user *models.RegisterRequest) error {
	sqlQuery := "INSERT INTO users(username, email, password, is_admin) VALUES ($1, $2, $3, $4)"

	_, err := tx.ExecContext(ctx, sqlQuery, user.Username, user.Email, user.Password, user.IsAdmin)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepositoryImpl) UpdateUser(ctx context.Context, tx *sql.Tx, user *models.UpdateUserRequest) error {
	sqlQuery := "UPDATE users SET username = $1, email = $2, is_admin = $3 where user_id = $4"

	_, err := tx.ExecContext(ctx, sqlQuery, user.Username, user.Email, user.IsAdmin, user.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepositoryImpl) DeleteUser(ctx context.Context, tx *sql.Tx, id string) error {
	sqlQuery := "DELETE FROM users where user_id = $1"
	
	_, err := tx.ExecContext(ctx, sqlQuery, id)
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

func (r *userRepositoryImpl) GetUserByUsername(ctx context.Context, tx *sql.Tx, username string) (models.User, error) {
	query := "SELECT user_id, username, email, password, is_admin, created_at, updated_at FROM users WHERE username = $1"
	
	var user models.User
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRowContext(ctx, query, username)
	} else {
		row = r.DB.QueryRowContext(ctx, query, username)
	}

	err := row.Scan(&user.UserId, &user.Username, &user.Email, &user.Password, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}

func (r *userRepositoryImpl) CountUser(ctx context.Context) (int, error) {
	sqlQuery := "SELECT COUNT(*) FROM users"

	var count int
	err := r.DB.QueryRowContext(ctx, sqlQuery).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
