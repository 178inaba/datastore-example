package repository

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/178inaba/datastore-example/entity"
)

// Task is the model used to store tasks in the datastore.
type Task struct {
	Description string    `datastore:"description"`
	Done        bool      `datastore:"done"`
	CreatedAt   time.Time `datastore:"created"`
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
	task := &Task{
		Description: description,
		CreatedAt:   createdAt,
	}

	key := datastore.IncompleteKey("Task", nil)

	return r.client.Put(ctx, key, task)
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
	query := datastore.NewQuery("Task").Order("created")
	keys, err := r.client.GetAll(ctx, query, &ts)
	if err != nil {
		return nil, err
	}

	// Repack Task into entity.Task.
	var tasks []*entity.Task
	for i, key := range keys {
		tasks = append(tasks, &entity.Task{
			ID:          key.ID,
			Description: tasks[i].Description,
			Done:        tasks[i].Done,
			CreatedAt:   tasks[i].CreatedAt,
		})
	}

	return tasks, nil
}
