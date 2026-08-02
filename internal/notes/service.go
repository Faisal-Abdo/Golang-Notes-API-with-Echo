package notes

import (
	"context"
	"database/sql"
	"errors"
	"notes-api/internal/database"
)

type Service struct {
	queries *database.Queries
}

// Constructor function to create a new Service instance
func NewService(db *sql.DB) *Service {
	return &Service{queries: database.New(db)}
}

func (s *Service) GetAllNotes(ctx context.Context) ([]Note, error) {
	dbNotes, err := s.queries.GetAllNotes(ctx)
	if err != nil {
		return nil, err
	}

	notes := make([]Note, len(dbNotes))
	for i, dbNote := range dbNotes {
		notes[i] = toNote(dbNote)
	}

	return notes, nil
}

func (s *Service) GetNoteByID(ctx context.Context, id int) (Note, error) {
	dbNote, err := s.queries.GetNoteByID(ctx, int32(id))
	if errors.Is(err, sql.ErrNoRows) {
		return Note{}, errors.New("note not found")
	}
	if err != nil {
		return Note{}, err
	}

	return toNote(dbNote), nil
}

func (s *Service) UpdateNoteByID(ctx context.Context, id int, updatedNote Note) (Note, error) {
	dbNote, err := s.queries.UpdateNoteByID(ctx, database.UpdateNoteByIDParams{
		Title:   updatedNote.Title,
		Content: updatedNote.Content,
		ID:      int32(id),
	})
	if errors.Is(err, sql.ErrNoRows) {
		return Note{}, errors.New("note not found")
	}
	if err != nil {
		return Note{}, err
	}

	return toNote(dbNote), nil
}

func (s *Service) DeleteNoteByID(ctx context.Context, id int) error {
	rowsAffected, err := s.queries.DeleteNoteByID(ctx, int32(id))
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("note not found")
	}

	return nil
}

func (s *Service) CreateNote(ctx context.Context, note Note) (Note, error) {
	dbNote, err := s.queries.CreateNote(ctx, database.CreateNoteParams{
		Title:   note.Title,
		Content: note.Content,
	})
	if err != nil {
		return Note{}, err
	}

	return toNote(dbNote), nil
}

func toNote(dbNote database.Note) Note {
	return Note{
		ID:        int(dbNote.ID),
		Title:     dbNote.Title,
		Content:   dbNote.Content,
		CreatedAt: dbNote.CreatedAt.Time,
	}
}
