package domain

import "time"

type TaskStatus string

const (
	TaskPending    TaskStatus = "pending"
	TaskProcessing TaskStatus = "processing"
	TaskCompleted  TaskStatus = "completed"
)

type Task struct {
	ID        int       `json:"id"`
	Links     []Link    `json:"links"`
	CreatedAt time.Time `json:"created_at"`
}