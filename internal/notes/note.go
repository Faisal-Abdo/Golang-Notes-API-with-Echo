package notes

import (
	"errors"
	"strings"
	"time"
)

type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func (n Note) Validate() error {
	if strings.TrimSpace(n.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(n.Content) == "" {
		return errors.New("content is required")
	}
	return nil
}
