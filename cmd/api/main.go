package main

import (
	"encoding/json"
	"net/http"
	"notes-api/internal/notes"
)

func main() {
	http.HandleFunc("/notes", getNotes)

	http.ListenAndServe(":8080", nil)
}

func getNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(notes.Notes)
}
