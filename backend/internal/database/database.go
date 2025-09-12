package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDatabase データベースを初期化し、テーブルを作成する
func InitDatabase() error {
	var err error
	
	// データベースファイルのパスを取得（環境変数から、またはデフォルト値）
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./todo_app.db"
	}

	// SQLiteデータベースに接続
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// 接続をテスト
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Database connected successfully")

	// テーブルを作成
	if err = createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	return nil
}

// createTables 必要なテーブルを作成する
func createTables() error {
	// ユーザーテーブル
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		name TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);`

	// タスクテーブル
	createTasksTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		priority TEXT NOT NULL DEFAULT 'medium',
		status TEXT NOT NULL DEFAULT 'pending',
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
	);`

	// リフレッシュトークンテーブル
	createRefreshTokensTable := `
	CREATE TABLE IF NOT EXISTS refresh_tokens (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		token TEXT UNIQUE NOT NULL,
		expires_at DATETIME NOT NULL,
		created_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
	);`

	// テーブルを作成
	if _, err := DB.Exec(createUsersTable); err != nil {
		return fmt.Errorf("failed to create users table: %v", err)
	}

	if _, err := DB.Exec(createTasksTable); err != nil {
		return fmt.Errorf("failed to create tasks table: %v", err)
	}

	if _, err := DB.Exec(createRefreshTokensTable); err != nil {
		return fmt.Errorf("failed to create refresh_tokens table: %v", err)
	}

	log.Println("Database tables created successfully")
	return nil
}

// CloseDatabase データベース接続を閉じる
func CloseDatabase() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
