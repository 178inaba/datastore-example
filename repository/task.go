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
	Description string         `datastore:",omitempty"`
	Text        string         `datastore:",omitempty"`
	Done        bool
	Due         time.Time `datastore:",omitempty"`
	CreatedAt   time.Time

	Metadata metadata
}

type metadata struct {
	URL *url.URL
}

// TaskRepository is task repository.
type TaskRepository struct {
	client *datastore.Client
}

// NewTaskRepository returns task repository.
func NewTaskRepository(client *datastore.Client) *TaskRepository {
	return &TaskRepository{client: client}
}

// AddTask adds a task with the given description to the datastore,
// returning the key of the newly created entity.
func (r *TaskRepository) AddTask(ctx context.Context, description, text string, createdAt time.Time) (*datastore.Key, error) {
	u, err := url.Parse("https://github.com/golang/go/issues/1")
	if err != nil {
		return nil, fmt.Errorf("parse rawurl: %w", err)
	}

	task := Task{
		Description: description,
		Text:        text,
		CreatedAt:   createdAt,

		Metadata: metadata{
			URL: u,
		},
	}

	key := datastore.IncompleteKey("Task", nil)
	keys, err := r.client.AllocateIDs(ctx, []*datastore.Key{key})
	if err != nil {
		return nil, fmt.Errorf("allocate IDs: %w", err)
	}

	//tasks := []*Task{&task}
	//ks, err := r.client.PutMulti(ctx, keys, tasks)
	//if err != nil {
	//	return nil, fmt.Errorf("put multi: %w", err)
	//}

	//return ks[0], nil
	return r.client.Put(ctx, keys[0], &task)
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
			Text:        ts[i].Text,
			Done:        ts[i].Done,
			Due:         ts[i].Due,
			CreatedAt:   ts[i].CreatedAt,
		}
	}

	return tasks, nil
}

// FilterKey returns tasks by key.
func (r *TaskRepository) FilterKey(ctx context.Context, key *datastore.Key) ([]*entity.Task, error) {
	// Create a query to fetch Task entities by key.
	var ts []*Task
	keys, err := r.client.GetAll(ctx, datastore.NewQuery("Task").Filter("__key__ =", key), &ts)
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
			Due:         ts[i].Due,
			CreatedAt:   ts[i].CreatedAt,
		}
	}

	return tasks, nil
}

// FilterDescription returns tasks by description.
func (r *TaskRepository) FilterDescription(ctx context.Context, description string) ([]*entity.Task, error) {
	// Create a query to fetch Task entities by key.
	var ts []*Task
	keys, err := r.client.GetAll(ctx, datastore.NewQuery("Task").Filter("Description =", description), &ts)
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
			Due:         ts[i].Due,
			CreatedAt:   ts[i].CreatedAt,
		}
	}

	return tasks, nil
}

// GetTaskIDsFilterDescription returns task IDs by description.
func (r *TaskRepository) GetTaskIDsFilterDescription(ctx context.Context, description string) ([]entity.TaskID, error) {
	keys, err := r.client.GetAll(ctx, datastore.NewQuery("Task").Filter("Description =", description), &[]*Task{})
	if err != nil {
		if _, ok := err.(*datastore.ErrFieldMismatch); !ok {
			return nil, fmt.Errorf("get all: %w", err)
		}
	}

	ids := make([]entity.TaskID, len(keys))
	for i, key := range keys {
		ids[i] = key
	}

	return ids, nil
}

// CountAll returns all task count.
func (r *TaskRepository) CountAll(ctx context.Context) (int, error) {
	return r.client.Count(ctx, datastore.NewQuery("Task"))
}

// CountDescNotNull returns description not null task count.
func (r *TaskRepository) CountDescNotNull(ctx context.Context) (int, error) {
	return r.client.Count(ctx, datastore.NewQuery("Task").Filter("Description >", ""))
}
