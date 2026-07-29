package notes

import (
	"errors"
)

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

func GetAllnotes() []Note {
	return notes
}
