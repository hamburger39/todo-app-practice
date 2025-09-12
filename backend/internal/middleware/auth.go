package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"todo-app-backend/internal/config"
	"todo-app-backend/internal/services"
)

func JWTAuth(cfg *config.Config) echo.MiddlewareFunc {
	jwtService := services.NewJWTService(cfg)
	
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Authorizationヘッダーを取得
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authorization header required",
				})
			}

			// Bearerトークンを抽出
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid authorization header format",
				})
			}

			// JWTトークンを検証
			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid token",
				})
			}

			// トークンタイプを確認（アクセストークンのみ許可）
			tokenType, ok := claims["type"].(string)
			if !ok || tokenType != "access" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid token type",
				})
			}

			userID, ok := claims["user_id"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid token claims",
				})
			}

			// デバッグ用ログ
			fmt.Printf("JWT Auth - UserID: %s, Path: %s\n", userID, c.Request().URL.Path)

			// コンテキストにユーザーIDを設定
			c.Set("user_id", userID)
			return next(c)
		}
	}
}






