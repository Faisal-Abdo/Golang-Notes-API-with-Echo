package notes

import "time"

type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

var notes = []Note{
	{
		ID:      1,
		Title:   "Learn Go",
		Content: "Build a CRUD API",
	},
	{
		ID:      2,
		Title:   "Learn net/http",
		Content: "Understand handlers and routing",
	},
}
