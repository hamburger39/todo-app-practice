package repository

import (
	"database/sql"
	"fmt"
	"time"

	"todo-app-backend/internal/database"
	"todo-app-backend/internal/models"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		db: database.DB,
	}
}

// CreateTask 新しいタスクを作成
func (r *TaskRepository) CreateTask(task *models.Task) error {
	query := `
		INSERT INTO tasks (id, user_id, title, description, priority, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, 
		task.ID, 
		task.UserID, 
		task.Title, 
		task.Description, 
		task.Priority, 
		task.Status, 
		task.CreatedAt, 
		task.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create task: %v", err)
	}

	return nil
}

// GetTasksByUserID ユーザーのタスク一覧を取得
func (r *TaskRepository) GetTasksByUserID(userID string) ([]models.Task, error) {
	query := `
		SELECT id, user_id, title, description, priority, status, created_at, updated_at
		FROM tasks
		WHERE user_id = ?
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %v", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Title,
			&task.Description,
			&task.Priority,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %v", err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %v", err)
	}

	return tasks, nil
}

// GetTaskByID タスクIDでタスクを取得
func (r *TaskRepository) GetTaskByID(id string) (*models.Task, error) {
	query := `
		SELECT id, user_id, title, description, priority, status, created_at, updated_at
		FROM tasks
		WHERE id = ?`

	task := &models.Task{}
	err := r.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Priority,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to get task: %v", err)
	}

	return task, nil
}

// UpdateTask タスクを更新
func (r *TaskRepository) UpdateTask(task *models.Task) error {
	query := `
		UPDATE tasks
		SET title = ?, description = ?, priority = ?, status = ?, updated_at = ?
		WHERE id = ?`

	_, err := r.db.Exec(query, 
		task.Title, 
		task.Description, 
		task.Priority, 
		task.Status, 
		time.Now(), 
		task.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update task: %v", err)
	}

	return nil
}

// DeleteTask タスクを削除
func (r *TaskRepository) DeleteTask(id string) error {
	query := `DELETE FROM tasks WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}

	return nil
}

// GetTasksByUserIDAndStatus ユーザーのタスクをステータスでフィルタリング
func (r *TaskRepository) GetTasksByUserIDAndStatus(userID, status string) ([]models.Task, error) {
	query := `
		SELECT id, user_id, title, description, priority, status, created_at, updated_at
		FROM tasks
		WHERE user_id = ? AND status = ?
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by status: %v", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Title,
			&task.Description,
			&task.Priority,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %v", err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %v", err)
	}

	return tasks, nil
}

// GetTasksByUserIDAndPriority ユーザーのタスクを優先度でフィルタリング
func (r *TaskRepository) GetTasksByUserIDAndPriority(userID, priority string) ([]models.Task, error) {
	query := `
		SELECT id, user_id, title, description, priority, status, created_at, updated_at
		FROM tasks
		WHERE user_id = ? AND priority = ?
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID, priority)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by priority: %v", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Title,
			&task.Description,
			&task.Priority,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %v", err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %v", err)
	}

	return tasks, nil
}
