package notes

import (
	"database/sql"
	"errors"
)

type Service struct {
	db *sql.DB
}

// Constructor function to create a new Service instance
func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// Go on and do the same for the other CRUD operations, using the Service struct to interact with the database.
func (s *Service) GetAllNotes() ([]Note, error) {

	query := `
		SELECT id, title, content, created_at
		FROM notes
		ORDER BY id;
	`
	// Execute the query and retrieve the rows
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	// Iterate through the rows and scan the data into Note structs
	for rows.Next() {

		var note Note
		// Scan the data from the current row into the Note struct
		err := rows.Scan(
			&note.ID,
			&note.Title,
			&note.Content,
			&note.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}
	// Check for any errors that occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

func GetNoteByID(id int) (Note, error) {
	for _, note := range notes {
		if note.ID == id {
			return note, nil
		}
	}
	return Note{}, errors.New("note not found")
}

func UpdateNoteByID(id int, updatedNote Note) (Note, error) {
	for i, note := range notes {
		if note.ID == id {
			updatedNote.ID = id
			notes[i] = updatedNote
			return updatedNote, nil
		}
	}
	return Note{}, errors.New("note not found")
}

func DeleteNoteByID(id int) error {
	for i, note := range notes {
		if note.ID == id {
			notes = append(notes[:i], notes[i+1:]...)
			return nil
		}
	}
	return errors.New("note not found")
}

func CreateNote(note Note) Note {
	note.ID = len(notes) + 1
	notes = append(notes, note)
	return note
}
