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

func (s *Service) GetNoteByID(id int) (Note, error) {
	query := `
		SELECT id, title, content, created_at
		FROM notes
		WHERE id = $1;
	`

	var note Note
	err := s.db.QueryRow(query, id).Scan(
		&note.ID,
		&note.Title,
		&note.Content,
		&note.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Note{}, errors.New("note not found")
	}
	if err != nil {
		return Note{}, err
	}

	return note, nil
}

func (s *Service) UpdateNoteByID(id int, updatedNote Note) (Note, error) {
	query := `
		UPDATE notes
		SET title = $1, content = $2
		WHERE id = $3
		RETURNING id, title, content, created_at;
	`

	var note Note
	err := s.db.QueryRow(query, updatedNote.Title, updatedNote.Content, id).Scan(
		&note.ID,
		&note.Title,
		&note.Content,
		&note.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Note{}, errors.New("note not found")
	}
	if err != nil {
		return Note{}, err
	}

	return note, nil
}

func (s *Service) DeleteNoteByID(id int) error {
	query := `DELETE FROM notes WHERE id = $1;`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("note not found")
	}

	return nil
}

func (s *Service) CreateNote(note Note) (Note, error) {
	query := `
		INSERT INTO notes (title, content)
		VALUES ($1, $2)
		RETURNING id, title, content, created_at;
	`

	var created Note
	err := s.db.QueryRow(query, note.Title, note.Content).Scan(
		&created.ID,
		&created.Title,
		&created.Content,
		&created.CreatedAt,
	)
	if err != nil {
		return Note{}, err
	}

	return created, nil
}
