package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"todo-app-backend/internal/models"
)

type TaskHandler struct {
	// 後でデータベース接続を追加
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{}
}

// 一時的なタスクストレージ（後でデータベースに置き換え）
var tasks = make(map[string]models.Task)

func (h *TaskHandler) GetTasks(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	// ユーザーのタスクを取得
	var userTasks []models.Task
	for _, task := range tasks {
		if task.UserID == userID {
			userTasks = append(userTasks, task)
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    userTasks,
	})
}

func (h *TaskHandler) CreateTask(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	var req models.CreateTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// タスクを作成
	task := models.Task{
		ID:          generateID(),
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Deadline:    req.Deadline,
		Priority:    req.Priority,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tasks[task.ID] = task

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    task,
	})
}

func (h *TaskHandler) UpdateTask(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	taskID := c.Param("id")
	task, exists := tasks[taskID]
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Task not found",
		})
	}

	// タスクがユーザーのものかチェック
	if task.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Access denied",
		})
	}

	var req models.UpdateTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// タスクを更新
	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = req.Description
	}
	if req.Deadline != nil {
		task.Deadline = req.Deadline
	}
	if req.Priority != nil {
		task.Priority = *req.Priority
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	task.UpdatedAt = time.Now()

	tasks[taskID] = task

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    task,
	})
}

func (h *TaskHandler) DeleteTask(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	taskID := c.Param("id")
	task, exists := tasks[taskID]
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Task not found",
		})
	}

	// タスクがユーザーのものかチェック
	if task.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Access denied",
		})
	}

	delete(tasks, taskID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Task deleted successfully",
	})
}

// ヘルパー関数
func getUserIDFromContext(c echo.Context) string {
	// JWTからユーザーIDを取得
	userID := c.Get("user_id")
	if userID == nil {
		return ""
	}
	return userID.(string)
}
