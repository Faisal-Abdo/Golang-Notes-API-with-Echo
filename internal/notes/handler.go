package notes

import (
	"errors"
	"log"
	"net/http"
	"notes-api/internal/auth"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service        *Service
	authMiddleware *auth.AuthMiddleware
}

func NewHandler(service *Service, authMiddleware *auth.AuthMiddleware) *Handler {
	return &Handler{service: service, authMiddleware: authMiddleware}
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	notes := e.Group("/notes", h.authMiddleware.Authenticate)

	notes.GET("", h.getNotes)
	notes.POST("", h.createNote)
	notes.GET("/:id", h.getNoteByID)
	notes.PUT("/:id", h.updateNote)
	notes.DELETE("/:id", h.deleteNote, auth.RequireRole("Admin"))
}

func respondError(c echo.Context, status int, message string) error {
	return c.JSON(status, map[string]string{"error": message})
}

// Logs the real cause server-side; the client only ever sees a generic message.
func respondInternalError(c echo.Context, err error) error {
	log.Printf("internal error: %v", err)
	return respondError(c, http.StatusInternalServerError, "internal server error")
}

func (h *Handler) getNotes(c echo.Context) error {
	notes, err := h.service.GetAllNotes(c.Request().Context())
	if err != nil {
		return respondInternalError(c, err)
	}
	return c.JSON(http.StatusOK, notes)
}

func (h *Handler) createNote(c echo.Context) error {
	var note Note

	if err := c.Bind(&note); err != nil {
		return respondError(c, http.StatusBadRequest, "invalid request body")
	}
	if err := note.Validate(); err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}

	note, err := h.service.CreateNote(c.Request().Context(), note)
	if err != nil {
		return respondInternalError(c, err)
	}

	return c.JSON(http.StatusCreated, note)
}

func (h *Handler) getNoteByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return respondError(c, http.StatusBadRequest, "invalid note id")
	}

	note, err := h.service.GetNoteByID(c.Request().Context(), id)
	if errors.Is(err, ErrNotFound) {
		return respondError(c, http.StatusNotFound, err.Error())
	}
	if err != nil {
		return respondInternalError(c, err)
	}

	return c.JSON(http.StatusOK, note)
}

func (h *Handler) updateNote(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return respondError(c, http.StatusBadRequest, "invalid note id")
	}

	updatedNote := Note{}
	if err := c.Bind(&updatedNote); err != nil {
		return respondError(c, http.StatusBadRequest, "invalid request body")
	}
	if err := updatedNote.Validate(); err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}

	updatedNote, err = h.service.UpdateNoteByID(c.Request().Context(), id, updatedNote)
	if errors.Is(err, ErrNotFound) {
		return respondError(c, http.StatusNotFound, err.Error())
	}
	if err != nil {
		return respondInternalError(c, err)
	}

	return c.JSON(http.StatusOK, updatedNote)
}

func (h *Handler) deleteNote(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return respondError(c, http.StatusBadRequest, "invalid note id")
	}

	err = h.service.DeleteNoteByID(c.Request().Context(), id)
	if errors.Is(err, ErrNotFound) {
		return respondError(c, http.StatusNotFound, err.Error())
	}
	if err != nil {
		return respondInternalError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
