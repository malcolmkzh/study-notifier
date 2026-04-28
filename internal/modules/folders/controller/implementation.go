package controller

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/malcolmkzh/study-notifier/internal/modules/folders/dto"
	"github.com/malcolmkzh/study-notifier/internal/modules/folders/service"
	"github.com/malcolmkzh/study-notifier/internal/utilities/errorutil"
	"github.com/malcolmkzh/study-notifier/internal/utilities/httpserver"
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
		return nil, errors.New("folders service is required")
	}

	m := &Implementation{
		httpServerUtility: httpServerUtility,
		service:           service,
	}

	routes := []httpserver.RegisterEndpointRequest{
		{Method: http.MethodPost, Path: "/folders", Fn: m.Create, RequireAuth: true},
		{Method: http.MethodGet, Path: "/folders", Fn: m.List, RequireAuth: true},
		{Method: http.MethodPut, Path: "/folders/:id", Fn: m.Update, RequireAuth: true},
		{Method: http.MethodDelete, Path: "/folders/:id", Fn: m.Delete, RequireAuth: true},
	}

	for _, route := range routes {
		if err := httpServerUtility.RegisterEndpoint(context.Background(), route); err != nil {
			return nil, err
		}
	}

	return m, nil
}

func (m *Implementation) Create(c *gin.Context) {
	userID, ok := httpserver.GetCurrentUserID(c)
	if !ok {
		return
	}

	var request dto.CreateFolderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(errorutil.NewWithMessage(errorutil.CodeBadRequest, "invalid request body"))
		c.Abort()
		return
	}

	folder, err := m.service.Create(c.Request.Context(), userID, request)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, folder)
}

func (m *Implementation) List(c *gin.Context) {
	userID, ok := httpserver.GetCurrentUserID(c)
	if !ok {
		return
	}

	folders, err := m.service.List(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, folders)
}

func (m *Implementation) Update(c *gin.Context) {
	userID, ok := httpserver.GetCurrentUserID(c)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errorutil.NewWithMessage(errorutil.CodeBadRequest, "invalid folder id"))
		c.Abort()
		return
	}

	var request dto.UpdateFolderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(errorutil.NewWithMessage(errorutil.CodeBadRequest, "invalid request body"))
		c.Abort()
		return
	}

	folder, err := m.service.Update(c.Request.Context(), uint(id), userID, request)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, folder)
}

func (m *Implementation) Delete(c *gin.Context) {
	userID, ok := httpserver.GetCurrentUserID(c)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errorutil.NewWithMessage(errorutil.CodeBadRequest, "invalid folder id"))
		c.Abort()
		return
	}

	if err := m.service.Delete(c.Request.Context(), uint(id), userID); err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)
}
