package main

import (
	"log"
	"net/http"
	"notes-api/internal/db"
	"notes-api/internal/notes"
)

func main() {
	db, err := db.Open()
	log.Println("Server running on port 8080")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// Create a new instance of the notes service and handler
	service := notes.NewService(db)
	handler := notes.NewHandler(service)
	// Register the handler functions for the routes
	http.HandleFunc("/notes", handler.NotesHandler)
	http.HandleFunc("/notes/", handler.NoteHandler)

	http.ListenAndServe(":8080", nil)
}
