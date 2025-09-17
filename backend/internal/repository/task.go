package repository

import (
	"database/sql"
	"todo-app-backend/internal/models"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		db: GetDB(),
	}
}

func (r *TaskRepository) CreateTask(task *models.Task) error {
	query := `INSERT INTO tasks (id, user_id, title, description, deadline, priority, status, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, task.ID, task.UserID, task.Title, task.Description, task.Deadline, 
		task.Priority, task.Status, task.CreatedAt, task.UpdatedAt)
	return err
}

func (r *TaskRepository) GetTasksByUserID(userID string) ([]models.Task, error) {
	query := `SELECT id, user_id, title, description, deadline, priority, status, created_at, updated_at 
			  FROM tasks WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Deadline,
			&task.Priority, &task.Status, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *TaskRepository) GetTaskByID(taskID string) (*models.Task, error) {
	query := `SELECT id, user_id, title, description, deadline, priority, status, created_at, updated_at 
			  FROM tasks WHERE id = ?`
	row := r.db.QueryRow(query, taskID)

	task := &models.Task{}
	err := row.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Deadline,
		&task.Priority, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (r *TaskRepository) UpdateTask(task *models.Task) error {
	query := `UPDATE tasks SET title = ?, description = ?, deadline = ?, priority = ?, status = ?, updated_at = ? 
			  WHERE id = ?`
	_, err := r.db.Exec(query, task.Title, task.Description, task.Deadline, task.Priority, 
		task.Status, task.UpdatedAt, task.ID)
	return err
}

func (r *TaskRepository) DeleteTask(taskID string) error {
	query := `DELETE FROM tasks WHERE id = ?`
	_, err := r.db.Exec(query, taskID)
	return err
}
