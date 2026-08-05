package notes

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	e.GET("/notes", h.getNotes)
	e.POST("/notes", h.createNote)
	e.GET("/notes/:id", h.getNoteByID)
	e.PUT("/notes/:id", h.updateNote)
	e.DELETE("/notes/:id", h.deleteNote)
}

func (h *Handler) getNotes(c echo.Context) error {
	notes, err := h.service.GetAllNotes(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, notes)
}

func (h *Handler) createNote(c echo.Context) error {
	var note Note

	if err := c.Bind(&note); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request body")
	}

	note, err := h.service.CreateNote(c.Request().Context(), note)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, note)
}

func (h *Handler) getNoteByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid note ID")
	}

	note, err := h.service.GetNoteByID(c.Request().Context(), id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, note)
}

func (h *Handler) updateNote(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid note ID")
	}

	updatedNote := Note{}
	if err := c.Bind(&updatedNote); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request body")
	}

	updatedNote, err = h.service.UpdateNoteByID(c.Request().Context(), id, updatedNote)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, updatedNote)
}

func (h *Handler) deleteNote(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid note ID")
	}

	if err := h.service.DeleteNoteByID(c.Request().Context(), id); err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
