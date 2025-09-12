package repository

import (
	"database/sql"
	"fmt"
	"time"

	"todo-app-backend/internal/database"
	"todo-app-backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.DB,
	}
}

// CreateUser 新しいユーザーを作成
func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (id, email, password, name, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, user.ID, user.Email, user.Password, user.Name, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}

// GetUserByEmail メールアドレスでユーザーを取得
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, password, name, created_at, updated_at
		FROM users
		WHERE email = ?`

	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

// GetUserByID IDでユーザーを取得
func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
	query := `
		SELECT id, email, password, name, created_at, updated_at
		FROM users
		WHERE id = ?`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

// UpdateUser ユーザー情報を更新
func (r *UserRepository) UpdateUser(user *models.User) error {
	query := `
		UPDATE users
		SET email = ?, name = ?, updated_at = ?
		WHERE id = ?`

	_, err := r.db.Exec(query, user.Email, user.Name, time.Now(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
}

// DeleteUser ユーザーを削除
func (r *UserRepository) DeleteUser(id string) error {
	query := `DELETE FROM users WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}

// UserExists ユーザーが存在するかチェック
func (r *UserRepository) UserExists(email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = ?`

	var count int
	err := r.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %v", err)
	}

	return count > 0, nil
}
