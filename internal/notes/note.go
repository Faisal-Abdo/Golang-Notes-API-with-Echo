package notes

type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var Notes = []Note{
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
