package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/_k1tasun_/Task_Project_Manager/internal/app"
)

func main() {
	cfg := app.LoadConfigFromEnv()
	if err := app.ValidateConfig(cfg); err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	app.SetCurrentConfig(cfg)
	app.PrintConfig(cfg)

	if err := app.LoadProjects(); err != nil {
		log.Fatalf("failed to load projects: %v", err)
	}
	if err := app.LoadTasks(); err != nil {
		log.Fatalf("failed to load tasks: %v", err)
	}
	if err := app.LoadUsers(); err != nil {
		log.Fatalf("failed to load users: %v", err)
	}
	if err := app.EnsureDefaultAdmin(); err != nil {
		log.Fatalf("failed to create default admin: %v", err)
	}

	mux := http.NewServeMux()
	app.RegisterRoutes(mux)

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		log.Printf("Server listening on http://%s", addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Graceful shutdown failed: %v", err)
	} else {
		log.Println("Server stopped cleanly")
	}
}
