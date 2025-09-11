package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"todo-app-backend/internal/handlers"
	authMiddleware "todo-app-backend/internal/middleware"
)

func main() {
	// Echoインスタンスを作成
	e := echo.New()

	// ミドルウェアを設定
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	// ハンドラーを初期化
	authHandler := handlers.NewAuthHandler()
	taskHandler := handlers.NewTaskHandler()

	// ルートを設定
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Todo App Backend API",
			"version": "1.0.0",
		})
	})

	// ヘルスチェックエンドポイント
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
		})
	})

	// 認証ルート
	auth := e.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	// タスクルート（認証が必要）
	tasks := e.Group("/tasks")
	tasks.Use(authMiddleware.JWTAuth())
	tasks.GET("", taskHandler.GetTasks)
	tasks.POST("", taskHandler.CreateTask)
	tasks.PUT("/:id", taskHandler.UpdateTask)
	tasks.DELETE("/:id", taskHandler.DeleteTask)

	// サーバーを起動
	log.Println("Starting server on :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
