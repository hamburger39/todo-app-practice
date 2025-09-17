package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"todo-app-backend/internal/models"
	"todo-app-backend/internal/repository"
	"todo-app-backend/internal/utils"
)

type TaskHandler struct {
	taskRepo *repository.TaskRepository
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		taskRepo: repository.NewTaskRepository(),
	}
}

func (h *TaskHandler) GetTasks(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	// ユーザーのタスクを取得
	userTasks, err := h.taskRepo.GetTasksByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get tasks",
		})
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
		ID:          utils.GenerateID(),
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Deadline:    req.Deadline,
		Priority:    req.Priority,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// データベースにタスクを保存
	if err := h.taskRepo.CreateTask(&task); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create task",
		})
	}

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
	task, err := h.taskRepo.GetTaskByID(taskID)
	if err != nil {
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

	// データベースでタスクを更新
	if err := h.taskRepo.UpdateTask(task); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update task",
		})
	}

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
	
	// タスクが存在するかチェック
	task, err := h.taskRepo.GetTaskByID(taskID)
	if err != nil {
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

	// データベースからタスクを削除
	if err := h.taskRepo.DeleteTask(taskID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete task",
		})
	}

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

