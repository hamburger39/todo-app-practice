package services

import (
	"errors"
	"time"

	"todo-app-backend/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey string
	expiryHours int
}

type Claims struct {
	UserID string `json:"user_id"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

func NewJWTService(cfg *config.Config) *JWTService {
	return &JWTService{
		secretKey: cfg.JWTSecret,
		expiryHours: cfg.JWTExpiryHours,
	}
}

func (s *JWTService) GenerateTokens(userID string) (string, string, error) {
	// アクセストークンを生成
	accessToken, err := s.generateToken(userID, "access", time.Duration(s.expiryHours)*time.Hour)
	if err != nil {
		return "", "", err
	}

	// リフレッシュトークンを生成（7日間有効）
	refreshToken, err := s.generateToken(userID, "refresh", 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *JWTService) generateToken(userID, tokenType string, duration time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *JWTService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *JWTService) RefreshAccessToken(refreshToken string) (string, string, error) {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	// リフレッシュトークンかチェック
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return "", "", errors.New("invalid token type")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", "", errors.New("invalid user ID")
	}

	// 新しいトークンペアを生成
	return s.GenerateTokens(userID)
}

func (s *JWTService) RevokeRefreshToken(refreshToken string) error {
	// 実際の実装では、リフレッシュトークンのブラックリストを管理する
	// ここでは簡単な実装として、トークンの有効性をチェックするだけ
	_, err := s.ValidateToken(refreshToken)
	return err
}
