package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/notes", getNotes)

	http.ListenAndServe(":8080", nil)
}

func getNotes(w http.ResponseWriter, r *http.Request) {
	notes := []string{"Learn Go", "Build Notes API"}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(notes)
}
