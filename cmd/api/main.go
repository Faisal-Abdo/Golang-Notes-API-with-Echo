package main

import (
	"context"
	"log"
	"net/http"
	"notes-api/internal/auth"
	"notes-api/internal/config"
	"notes-api/internal/db"
	"notes-api/internal/notes"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Only a local-dev convenience; in Docker, env vars come from the environment
	// directly, so a missing .env here is not fatal.
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using existing environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	authMiddleware, err := auth.NewAuthMiddleware(cfg.KeycloakURL, cfg.KeycloakRealm, cfg.KeycloakClientID)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	service := notes.NewService(dbConn)
	handler := notes.NewHandler(service, authMiddleware)
	handler.RegisterRoutes(e)

	server := &http.Server{
		Addr:    ":8080",
		Handler: e, // Use the Echo instance as the HTTP handler
	}

	log.Println("Server running on :8080")
	// Start the server in a separate goroutine to allow for graceful shutdown
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	// Create a channel to listen for OS signals for graceful shutdown
	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGTERM,
	)
	<-quit
	log.Println("Server shutting down...")
	// Create a context with a timeout for the shutdown process
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()
	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	err = dbConn.Close()
	if err != nil {
		log.Println(err)
	}

	log.Println("Server Stopped")
}
