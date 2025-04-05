package storage

import "github.com/kartikey1188/go-todo-list-v2/internal/types"

type Storage interface {
	CreateTask(title string, description string, deadline types.Date) (int64, error)
	GetTask(id int64) (types.Task, error)
	GetTasks() ([]types.Task, error)
	// UpdateTask(id int64, title string, description string, deadline types.Date) (types.Task, error)
	// DeleteTask(id int64) (int64, error)
}
