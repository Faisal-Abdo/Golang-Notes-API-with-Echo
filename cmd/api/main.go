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

	http.HandleFunc("/notes", notes.NotesHandler)
	http.HandleFunc("/notes/", notes.NoteHandler)

	http.ListenAndServe(":8080", nil)
}
