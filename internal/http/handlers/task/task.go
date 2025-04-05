package task

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kartikey1188/go-todo-list-v2/internal/storage"
	"github.com/kartikey1188/go-todo-list-v2/internal/types"
	"github.com/kartikey1188/go-todo-list-v2/internal/utils/response"
)

func New(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {

		var task types.Task

		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		lastId, err := storage.CreateTask(
			task.Title,
			task.Description,
			task.Deadline,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("task created successfully", slog.String("Task ID", fmt.Sprint(lastId)))

		c.JSON(http.StatusCreated, gin.H{
			"status": "OK",
			"ID":     lastId,
		})
	}
}

func GetById(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, response.GeneralError(fmt.Errorf("no task id passed")))
		}

		slog.Info("getting a task", slog.String("id", id))

		intID, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid id %s", id)))
			return
		}

		task, err := storage.GetTask(intID)

		if err != nil {
			slog.Error("failed to get task", slog.String("Task ID:", id))
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, task)
	}
}

func GetList(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Info("getting all tasks")

		tasks, err := storage.GetTasks()

		if err != nil {
			slog.Error("failed to get tasks")
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
		}

		if len(tasks) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"Error": "No tasks found"})
			return
		}

		c.JSON(http.StatusOK, tasks)
	}
}
