package domain

import "time"

type TaskStatus string

const (
	TaskPending    TaskStatus = "pending"
	TaskProcessing TaskStatus = "processing"
	TaskCompleted  TaskStatus = "completed"
)

type Task struct {
	ID        string
	Links     []Link
	Status    TaskStatus
	LinksNum  int
	CreatedAt time.Time
	UpdatedAt time.Time
}