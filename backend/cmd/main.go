package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"todo-app-backend/internal/config"
	"todo-app-backend/internal/handlers"
	authmiddleware "todo-app-backend/internal/middleware"
)

func main() {
	// 設定を読み込み
	cfg := config.Load()

	// Echoインスタンスを作成
	e := echo.New()

	// ミドルウェアを設定
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	
	// CORS設定
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// ハンドラーを初期化
	authHandler := handlers.NewAuthHandler(cfg)
	taskHandler := handlers.NewTaskHandler()

	// ルートを設定
	setupRoutes(e, authHandler, taskHandler, cfg)

	// サーバーを起動
	log.Printf("Server starting on port %s", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRoutes(e *echo.Echo, authHandler *handlers.AuthHandler, taskHandler *handlers.TaskHandler, cfg *config.Config) {
	// 認証不要のルート
	e.POST("/api/auth/register", authHandler.Register)
	e.POST("/api/auth/login", authHandler.Login)
	e.POST("/api/auth/refresh", authHandler.RefreshToken)
	e.POST("/api/auth/logout", authHandler.Logout)

	// 認証が必要なルート
	api := e.Group("/api")
	api.Use(authmiddleware.JWTAuth(cfg))

	// タスク関連のルート
	api.GET("/tasks", taskHandler.GetTasks)
	api.POST("/tasks", taskHandler.CreateTask)
	api.PUT("/tasks/:id", taskHandler.UpdateTask)
	api.DELETE("/tasks/:id", taskHandler.DeleteTask)

	// ヘルスチェック
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})
}
