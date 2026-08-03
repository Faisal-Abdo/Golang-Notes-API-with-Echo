package main

import (
	"context"
	"log"
	"net/http"
	"notes-api/internal/db"
	"notes-api/internal/notes"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	// Create a new instance of the notes service and handler
	service := notes.NewService(db)
	handler := notes.NewHandler(service)
	// Register the handler functions for the routes
	http.HandleFunc("/notes", handler.NotesHandler)
	http.HandleFunc("/notes/", handler.NoteHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
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

	err = db.Close()
	if err != nil {
		log.Println(err)
	}

	log.Println("Server Stopped")
}
