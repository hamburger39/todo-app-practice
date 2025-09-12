package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"todo-app-backend/internal/config"
	"todo-app-backend/internal/models"
	"todo-app-backend/internal/repository"
	"todo-app-backend/internal/services"
)

type AuthHandler struct {
	userRepo  *repository.UserRepository
	jwtService *services.JWTService
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		userRepo:   repository.NewUserRepository(),
		jwtService: services.NewJWTService(cfg),
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req models.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// ユーザーが既に存在するかチェック
	exists, err := h.userRepo.UserExists(req.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to check user existence",
		})
	}
	if exists {
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

	// データベースにユーザーを保存
	if err := h.userRepo.CreateUser(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create user",
		})
	}

	// JWTトークンを生成
	accessToken, refreshToken, err := h.jwtService.GenerateTokens(user.ID)
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
	user, err := h.userRepo.GetUserByEmail(req.Email)
	if err != nil {
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
	accessToken, refreshToken, err := h.jwtService.GenerateTokens(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate tokens",
		})
	}

	response := models.AuthResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    response,
	})
}

// RefreshToken リフレッシュトークンを使用してアクセストークンを更新
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// リフレッシュトークンを使用して新しいトークンを生成
	newAccessToken, newRefreshToken, err := h.jwtService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid refresh token",
		})
	}

	response := map[string]interface{}{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    response,
	})
}

// Logout ログアウト（リフレッシュトークンを無効化）
func (h *AuthHandler) Logout(c echo.Context) error {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// リフレッシュトークンを無効化
	if err := h.jwtService.RevokeRefreshToken(req.RefreshToken); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to revoke token",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Logged out successfully",
	})
}

// ヘルパー関数
func generateID() string {
	return time.Now().Format("20060102150405")
}






