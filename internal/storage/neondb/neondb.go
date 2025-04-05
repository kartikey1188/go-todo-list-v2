package neondb

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver for database/sql

	"github.com/kartikey1188/go-todo-list-v2/internal/config"
	"github.com/kartikey1188/go-todo-list-v2/internal/types"
)

type Postgres struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Postgres, error) {
	db, err := sql.Open("pgx", cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		 id SERIAL PRIMARY KEY,
		 title TEXT NOT NULL,
		 description TEXT NOT NULL,
		 deadline DATE 
	)`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &Postgres{
		Db: db,
	}, nil
}

func (p *Postgres) CreateTask(title string, description string, deadline types.Date) (int64, error) {
	formattedDate := deadline.Time.Format("2006-01-02")

	var lastId int64
	err := p.Db.QueryRow(
		"INSERT INTO tasks (title, description, deadline) VALUES ($1, $2, $3) RETURNING id",
		title, description, formattedDate,
	).Scan(&lastId)
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (s *Postgres) GetTask(id int64) (types.Task, error) {
	stmt, err := s.Db.Prepare("SELECT id, title, description, deadline FROM tasks WHERE id = $1 LIMIT 1")
	if err != nil {
		return types.Task{}, err
	}
	defer stmt.Close()

	var task types.Task
	var deadline time.Time

	err = stmt.QueryRow(id).Scan(&task.ID, &task.Title, &task.Description, &deadline)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Task{}, fmt.Errorf("no task found with id %s", fmt.Sprint(id))
		}
		return types.Task{}, fmt.Errorf("query error: %w", err)
	}

	task.Deadline = types.Date{Time: deadline}
	return task, nil
}

func (s *Postgres) GetTasks() ([]types.Task, error) {
	stmt, err := s.Db.Prepare("SELECT id, title, description, deadline FROM tasks")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []types.Task

	for rows.Next() {
		var task types.Task
		var deadline time.Time

		err := rows.Scan(&task.ID, &task.Title, &task.Description, &deadline)
		if err != nil {
			return nil, err
		}
		task.Deadline = types.Date{Time: deadline}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *Postgres) UpdateTask(id int64, title string, description string, deadline types.Date) (types.Task, error) {
	formattedDate := deadline.Time.Format("2006-01-02")

	result, err := s.Db.Exec(
		"UPDATE tasks SET title = $1, description = $2, deadline = $3 WHERE id = $4",
		title, description, formattedDate, id,
	)
	if err != nil {
		return types.Task{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return types.Task{}, err
	}

	if rowsAffected == 0 {
		return types.Task{}, fmt.Errorf("no task found with id %d", id)
	}

	return types.Task{
		ID:          id,
		Title:       title,
		Description: description,
		Deadline:    deadline,
	}, nil
}

func (s *Postgres) DeleteTask(id int64) (int64, error) {
	res, err := s.Db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	if rowsAffected == 0 {
		return 0, fmt.Errorf("no task found with id %d", id)
	}

	return rowsAffected, nil
}
