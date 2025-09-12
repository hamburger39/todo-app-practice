package repository

import (
	"database/sql"
	"fmt"
	"time"

	"todo-app-backend/internal/database"
)

type RefreshToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type RefreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository() *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db: database.DB,
	}
}

// CreateRefreshToken 新しいリフレッシュトークンを作成
func (r *RefreshTokenRepository) CreateRefreshToken(token *RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, token.ID, token.UserID, token.Token, token.ExpiresAt, token.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create refresh token: %v", err)
	}

	return nil
}

// GetRefreshToken トークン文字列でリフレッシュトークンを取得
func (r *RefreshTokenRepository) GetRefreshToken(token string) (*RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM refresh_tokens
		WHERE token = ?`

	refreshToken := &RefreshToken{}
	err := r.db.QueryRow(query, token).Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.Token,
		&refreshToken.ExpiresAt,
		&refreshToken.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("refresh token not found")
		}
		return nil, fmt.Errorf("failed to get refresh token: %v", err)
	}

	return refreshToken, nil
}

// DeleteRefreshToken リフレッシュトークンを削除
func (r *RefreshTokenRepository) DeleteRefreshToken(token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = ?`

	_, err := r.db.Exec(query, token)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %v", err)
	}

	return nil
}

// DeleteUserRefreshTokens ユーザーのすべてのリフレッシュトークンを削除
func (r *RefreshTokenRepository) DeleteUserRefreshTokens(userID string) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = ?`

	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user refresh tokens: %v", err)
	}

	return nil
}

// CleanupExpiredTokens 期限切れのリフレッシュトークンを削除
func (r *RefreshTokenRepository) CleanupExpiredTokens() error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < ?`

	_, err := r.db.Exec(query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to cleanup expired tokens: %v", err)
	}

	return nil
}
