package main

import (
	"net/http"
	"notes-api/internal/notes"
)

func main() {
	http.HandleFunc("/notes", notes.NotesHandler)
	http.HandleFunc("/notes/", notes.NoteHandler)

	http.ListenAndServe(":8080", nil)
}
