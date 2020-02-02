package entity

import (
	"time"

	"cloud.google.com/go/datastore"
)

// TaskID is task ID.
type TaskID *datastore.Key

// Task is the task entity.
type Task struct {
	ID          TaskID
	Description string
	Done        bool
	CreatedAt   time.Time
}
