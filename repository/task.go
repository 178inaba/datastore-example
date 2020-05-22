package repository

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/178inaba/datastore-example/entity"
)

// Task is the model used to store tasks in the datastore.
type Task struct {
	ID          *datastore.Key `datastore:"__key__"`
	Description string
	Done        bool
	Due         time.Time
	CreatedAt   time.Time

	Metadata Metadata
}

type Metadata struct {
	URL *url.URL
}

type TaskRepository struct {
	client *datastore.Client
}

func NewTaskRepository(client *datastore.Client) *TaskRepository {
	return &TaskRepository{client: client}
}

// AddTask adds a task with the given description to the datastore,
// returning the key of the newly created entity.
func (r *TaskRepository) AddTask(ctx context.Context, description string, createdAt time.Time) (*datastore.Key, error) {
	u, err := url.Parse("https://github.com/golang/go/issues/1")
	if err != nil {
		return nil, fmt.Errorf("parse rawurl: %w", err)
	}

	task := &Task{
		Description: description,
		CreatedAt:   createdAt,

		Metadata: Metadata{
			URL: u,
		},
	}

	key := datastore.IncompleteKey("Task", nil)
	keys, err := r.client.AllocateIDs(ctx, []*datastore.Key{key})
	if err != nil {
		return nil, fmt.Errorf("allocate IDs: %w", err)
	}

	return r.client.Put(ctx, keys[0], task)
}

// MarkDone marks the task done with the given ID.
func (r *TaskRepository) MarkDone(ctx context.Context, taskID int64) error {
	// Create a key using the given integer ID.
	key := datastore.IDKey("Task", taskID, nil)

	// In a transaction load each task, set done to true and store.
	_, err := r.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var task Task
		if err := tx.Get(key, &task); err != nil {
			return err
		}

		task.Done = true
		_, err := tx.Put(key, &task)
		return err
	})

	return err
}

// DeleteTask deletes the task with the given ID.
func (r *TaskRepository) DeleteTask(ctx context.Context, taskID int64) error {
	return r.client.Delete(ctx, datastore.IDKey("Task", taskID, nil))
}

// ListTasks returns all the tasks in ascending order of creation time.
func (r *TaskRepository) ListTasks(ctx context.Context) ([]*entity.Task, error) {
	// Create a query to fetch all Task entities, ordered by "created".
	var ts []*Task
	query := datastore.NewQuery("Task").Order("CreatedAt")
	keys, err := r.client.GetAll(ctx, query, &ts)
	if err != nil {
		if _, ok := err.(*datastore.ErrFieldMismatch); !ok {
			return nil, fmt.Errorf("get all: %w", err)
		}
	}

	// Repack Task into entity.Task.
	tasks := make([]*entity.Task, len(keys))
	for i, key := range keys {
		tasks[i] = &entity.Task{
			ID:          key,
			Description: ts[i].Description,
			Done:        ts[i].Done,
			CreatedAt:   ts[i].CreatedAt,
		}
	}

	return tasks, nil
}
