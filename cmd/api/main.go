package main

import (
	"encoding/json"
	"net/http"
	"notes-api/internal/notes"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/notes", notesHandler)
	http.HandleFunc("/notes/", noteHandler)

	http.ListenAndServe(":8080", nil)
}

func notesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getNotes(w, r)
	case http.MethodPost:
		createNote(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// noteHandler handles requests for a specific note by ID.
func noteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getNoteByID(w, r)
	case http.MethodPut:
		updateNote(w, r)
	case http.MethodDelete:
		deleteNote(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(notes.Notes)
}

func createNote(w http.ResponseWriter, r *http.Request) {
	var note notes.Note

	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	note.ID = len(notes.Notes) + 1
	notes.Notes = append(notes.Notes, note)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(note)
}

func getNoteByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/notes/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	for _, note := range notes.Notes {
		if note.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(note)
			return
		}
	}

	http.Error(w, "Note not found", http.StatusNotFound)
}

func updateNote(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/notes/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	var updatedNote notes.Note

	err = json.NewDecoder(r.Body).Decode(&updatedNote)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for i, note := range notes.Notes {
		if note.ID == id {
			updatedNote.ID = id
			notes.Notes[i] = updatedNote

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedNote)
			return
		}
	}

	http.Error(w, "Note not found", http.StatusNotFound)
}

func deleteNote(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/notes/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	for i, note := range notes.Notes {
		if note.ID == id {
			// Rebuild the slice without the deleted note
			notes.Notes = append(notes.Notes[:i], notes.Notes[i+1:]...)

			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Note not found", http.StatusNotFound)
}
