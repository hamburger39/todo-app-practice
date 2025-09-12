package services

import (
	"fmt"
	"time"

	"todo-app-backend/internal/config"
	"todo-app-backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	config                *config.Config
	refreshTokenRepo      *repository.RefreshTokenRepository
}

func NewJWTService(cfg *config.Config) *JWTService {
	return &JWTService{
		config:           cfg,
		refreshTokenRepo: repository.NewRefreshTokenRepository(),
	}
}

// GenerateTokens アクセストークンとリフレッシュトークンを生成
func (s *JWTService) GenerateTokens(userID string) (string, string, error) {
	// アクセストークンを生成
	accessToken, err := s.generateAccessToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %v", err)
	}

	// リフレッシュトークンを生成
	refreshToken, err := s.generateRefreshToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %v", err)
	}

	return accessToken, refreshToken, nil
}

// generateAccessToken アクセストークンを生成
func (s *JWTService) generateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.config.JWTAccessTokenExpiry).Unix(),
		"type":    "access",
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecretKey))
}

// generateRefreshToken リフレッシュトークンを生成
func (s *JWTService) generateRefreshToken(userID string) (string, error) {
	// JWTトークンを生成
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.config.JWTRefreshTokenExpiry).Unix(),
		"type":    "refresh",
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWTSecretKey))
	if err != nil {
		return "", err
	}

	// データベースにリフレッシュトークンを保存
	refreshToken := &repository.RefreshToken{
		ID:        generateID(),
		UserID:    userID,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(s.config.JWTRefreshTokenExpiry),
		CreatedAt: time.Now(),
	}

	if err := s.refreshTokenRepo.CreateRefreshToken(refreshToken); err != nil {
		return "", fmt.Errorf("failed to save refresh token: %v", err)
	}

	return tokenString, nil
}

// ValidateToken トークンを検証
func (s *JWTService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 署名方法を確認
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// RefreshAccessToken リフレッシュトークンを使用してアクセストークンを更新
func (s *JWTService) RefreshAccessToken(refreshTokenString string) (string, string, error) {
	// リフレッシュトークンを検証
	claims, err := s.ValidateToken(refreshTokenString)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %v", err)
	}

	// トークンタイプを確認
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return "", "", fmt.Errorf("invalid token type")
	}

	// データベースでリフレッシュトークンの存在を確認
	refreshToken, err := s.refreshTokenRepo.GetRefreshToken(refreshTokenString)
	if err != nil {
		return "", "", fmt.Errorf("refresh token not found in database")
	}

	// 期限切れチェック
	if time.Now().After(refreshToken.ExpiresAt) {
		// 期限切れのトークンを削除
		s.refreshTokenRepo.DeleteRefreshToken(refreshTokenString)
		return "", "", fmt.Errorf("refresh token expired")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid user ID in token")
	}

	// 新しいアクセストークンとリフレッシュトークンを生成
	newAccessToken, err := s.generateAccessToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new access token: %v", err)
	}

	// 古いリフレッシュトークンを削除
	s.refreshTokenRepo.DeleteRefreshToken(refreshTokenString)

	// 新しいリフレッシュトークンを生成
	newRefreshToken, err := s.generateRefreshToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new refresh token: %v", err)
	}

	return newAccessToken, newRefreshToken, nil
}

// RevokeRefreshToken リフレッシュトークンを無効化
func (s *JWTService) RevokeRefreshToken(refreshTokenString string) error {
	return s.refreshTokenRepo.DeleteRefreshToken(refreshTokenString)
}

// RevokeAllUserTokens ユーザーのすべてのリフレッシュトークンを無効化
func (s *JWTService) RevokeAllUserTokens(userID string) error {
	return s.refreshTokenRepo.DeleteUserRefreshTokens(userID)
}

// generateID ユニークなIDを生成
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
