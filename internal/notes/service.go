package notes

import (
	"errors"
)

func GetNoteByID(id int) (Note, error) {
	for _, note := range Notes {
		if note.ID == id {
			return note, nil
		}
	}
	return Note{}, errors.New("note not found")
}

func UpdateNoteByID(id int, updatedNote Note) (Note, error) {
	for i, note := range Notes {
		if note.ID == id {
			updatedNote.ID = id
			Notes[i] = updatedNote
			return updatedNote, nil
		}
	}
	return Note{}, errors.New("note not found")
}

func DeleteNoteByID(id int) error {
	for i, note := range Notes {
		if note.ID == id {
			Notes = append(Notes[:i], Notes[i+1:]...)
			return nil
		}
	}
	return errors.New("note not found")
}

func CreateNote(note Note) Note {
	note.ID = len(Notes) + 1
	Notes = append(Notes, note)
	return note
}

func GetAllNotes() []Note {
	return Notes
}
