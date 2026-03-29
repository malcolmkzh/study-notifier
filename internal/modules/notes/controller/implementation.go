package controller

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/malcolmkzh/study-notifier/internal/modules/notes/dto"
	"github.com/malcolmkzh/study-notifier/internal/modules/notes/service"
	"github.com/malcolmkzh/study-notifier/internal/utilities/httpserver"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type Implementation struct {
	httpServerUtility httpserver.Utility
	service           service.Utility
}

func NewController(httpServerUtility httpserver.Utility, service service.Utility) (*Implementation, error) {
	if httpServerUtility == nil {
		return nil, errors.New("http server utility is required")
	}
	if service == nil {
		return nil, errors.New("notes service is required")
	}

	m := &Implementation{
		httpServerUtility: httpServerUtility,
		service:           service,
	}

	routes := []httpserver.RegisterEndpointRequest{
		{Method: http.MethodPost, Path: "/notes", Fn: m.Create, RequireAuth: true},
		{Method: http.MethodGet, Path: "/notes", Fn: m.List, RequireAuth: true},
		{Method: http.MethodGet, Path: "/notes/:id", Fn: m.GetByID, RequireAuth: true},
		{Method: http.MethodPut, Path: "/notes/:id", Fn: m.Update, RequireAuth: true},
		{Method: http.MethodDelete, Path: "/notes/:id", Fn: m.Delete, RequireAuth: true},
	}

	for _, route := range routes {
		if err := httpServerUtility.RegisterEndpoint(context.Background(), route); err != nil {
			return nil, err
		}
	}

	return m, nil
}

// Create Notes
func (m *Implementation) Create(c *gin.Context) {
	userID, ok := httpserver.GetCurrentUserID(c)
	if !ok {
		return
	}

	var request dto.CreateNoteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	note, err := m.service.Create(c.Request.Context(), userID, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, note)
}

// List Notes
func (m *Implementation) List(c *gin.Context) {
	userID, ok := httpserver.GetCurrentUserID(c)
	if !ok {
		return
	}

	notes, err := m.service.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, notes)
}

// Get Note by ID
func (m *Implementation) GetByID(c *gin.Context) {
	userID, ok := httpserver.GetCurrentUserID(c)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note id"})
		return
	}

	note, err := m.service.GetByID(c.Request.Context(), uint(id), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if note == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
		return
	}

	c.JSON(http.StatusOK, note)
}

// Update Note
func (m *Implementation) Update(c *gin.Context) {
	userID, ok := httpserver.GetCurrentUserID(c)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note id"})
		return
	}

	var request dto.UpdateNoteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	note, err := m.service.Update(c.Request.Context(), uint(id), userID, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if note == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
		return
	}

	c.JSON(http.StatusOK, note)
}

// Delete Note
func (m *Implementation) Delete(c *gin.Context) {
	userID, ok := httpserver.GetCurrentUserID(c)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note id"})
		return
	}

	err = m.service.Delete(c.Request.Context(), uint(id), userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}
