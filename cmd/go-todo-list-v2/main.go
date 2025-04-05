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
	"github.com/kartikey1188/go-todo-list-v2/internal/storage/neondb"
)

func main() {
	// load config

	cfg := config.MustLoad()

	// setup database

	_, err := neondb.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// setup router
	router := gin.Default()

	// router.POST("/api/tasks", task.New(storage))
	// router.GET("/api/tasks/:id", task.GetById(storage)) // Use :id for path param
	// router.GET("/api/tasks", task.GetList(storage))
	// router.PUT("/api/tasks/:id", task.Update(storage))    // Use :id for path param
	// router.DELETE("/api/tasks/:id", task.Delete(storage)) // Use :id for path param

	//setup server

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
