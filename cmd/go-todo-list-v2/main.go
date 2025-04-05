package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kartikey1188/go-todo-list-v2/internal/config"
	"github.com/kartikey1188/go-todo-list-v2/internal/http/handlers/task"
	"github.com/kartikey1188/go-todo-list-v2/internal/storage/neondb"
)

func main() {
	// loading config

	cfg := config.MustLoad()

	// setting up database

	storage, err := neondb.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// setting up router

	router := gin.Default()

	router.POST("/api/tasks", task.New(storage))
	router.GET("/api/tasks/:id", task.GetById(storage))
	router.GET("/api/tasks", task.GetList(storage))
	// router.PUT("/api/tasks/:id", task.Update(storage))
	// router.DELETE("/api/tasks/:id", task.Delete(storage))

	//setting up server (with graceful shutdown)

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("server started", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server")
		}
	}()

	<-done

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("sever shutdown successfully")
}
