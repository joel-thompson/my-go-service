package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joel-thompson/my-go-service/api/server"
	"github.com/joel-thompson/my-go-service/cmd/server/setup"
)

func main() {
	ctx := context.Background()

	// Initialize application
	app, err := setup.NewApp(ctx)
	if err != nil {
		log.Fatal("Failed to initialize app:", err)
	}
	defer app.Close()

	// Setup API server
	api := server.New(app.Logger, app.DB)
	router := api.SetupRoutes()

	// Create HTTP server
	srv := &http.Server{
		Addr:    app.Config.ServerAddr,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		app.Logger.Info("Starting server", "addr", app.Config.ServerAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.Logger.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	app.Logger.Info("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		app.Logger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	app.Logger.Info("Server exited")
}
