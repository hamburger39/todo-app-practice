package models

import (
	"time"
)

type Task struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	Title       string    `json:"title" db:"title"`
	Description *string   `json:"description,omitempty" db:"description"`
	Deadline    *time.Time `json:"deadline,omitempty" db:"deadline"`
	Priority    string    `json:"priority" db:"priority"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateTaskRequest struct {
	Title       string     `json:"title" validate:"required"`
	Description *string    `json:"description,omitempty"`
	Deadline    *time.Time `json:"deadline,omitempty"`
	Priority    string     `json:"priority" validate:"required,oneof=high medium low"`
}

type UpdateTaskRequest struct {
	Title       *string     `json:"title,omitempty"`
	Description *string     `json:"description,omitempty"`
	Deadline    *time.Time  `json:"deadline,omitempty"`
	Priority    *string     `json:"priority,omitempty" validate:"omitempty,oneof=high medium low"`
	Status      *string     `json:"status,omitempty" validate:"omitempty,oneof=pending completed"`
}

type TaskFilters struct {
	Status    *string `query:"status" validate:"omitempty,oneof=pending completed all"`
	Priority  *string `query:"priority" validate:"omitempty,oneof=high medium low all"`
	SortBy    *string `query:"sort_by" validate:"omitempty,oneof=deadline priority created_at"`
	SortOrder *string `query:"sort_order" validate:"omitempty,oneof=asc desc"`
}






