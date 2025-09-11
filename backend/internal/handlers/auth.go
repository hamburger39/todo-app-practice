package handlers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"todo-app-backend/internal/models"
)

type AuthHandler struct {
	// 後でデータベース接続を追加
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// 一時的なユーザーストレージ（後でデータベースに置き換え）
var users = make(map[string]models.User)

func (h *AuthHandler) Register(c echo.Context) error {
	var req models.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// ユーザーが既に存在するかチェック
	if _, exists := users[req.Email]; exists {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "User already exists",
		})
	}

	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to hash password",
		})
	}

	// ユーザーを作成
	user := models.User{
		ID:        generateID(),
		Email:     req.Email,
		Password:  string(hashedPassword),
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	users[req.Email] = user

	// JWTトークンを生成
	accessToken, refreshToken, err := generateTokens(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate tokens",
		})
	}

	response := models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    response,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// ユーザーを検索
	user, exists := users[req.Email]
	if !exists {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid credentials",
		})
	}

	// パスワードを検証
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid credentials",
		})
	}

	// JWTトークンを生成
	accessToken, refreshToken, err := generateTokens(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate tokens",
		})
	}

	response := models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    response,
	})
}

// ヘルパー関数
func generateID() string {
	return time.Now().Format("20060102150405")
}

func generateTokens(userID string) (string, string, error) {
	// アクセストークン（1時間有効）
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
		"type":    "access",
	})

	// リフレッシュトークン（7日間有効）
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().AddDate(0, 0, 7).Unix(),
		"type":    "refresh",
	})

	// 署名（本番環境では環境変数から取得）
	secret := []byte("your-secret-key")
	
	accessTokenString, err := accessToken.SignedString(secret)
	if err != nil {
		return "", "", err
	}

	refreshTokenString, err := refreshToken.SignedString(secret)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}





