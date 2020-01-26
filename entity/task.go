package entity

import "time"

// Task is the task entity.
type Task struct {
	ID          int64
	Description string
	Done        bool
	CreatedAt   time.Time
}
