package controller

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/malcolmkzh/study-notifier/internal/modules/questions/service"
	"github.com/malcolmkzh/study-notifier/internal/utilities/errorutil"
	"github.com/malcolmkzh/study-notifier/internal/utilities/httpserver"

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
		return nil, errors.New("questions service is required")
	}

	m := &Implementation{
		httpServerUtility: httpServerUtility,
		service:           service,
	}

	err := httpServerUtility.RegisterEndpoint(context.Background(), httpserver.RegisterEndpointRequest{
		Method:      http.MethodPost,
		Path:        "/notes/:id/generate-questions",
		Fn:          m.GenerateQuestions,
		RequireAuth: true,
	})
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Implementation) GenerateQuestions(c *gin.Context) {
	userID, ok := httpserver.GetCurrentUserID(c)
	if !ok {
		return
	}

	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errorutil.NewWithMessage(errorutil.CodeBadRequest, "invalid note id"))
		c.Abort()
		return
	}

	response, err := m.service.GenerateQuestions(c.Request.Context(), uint(noteID), userID)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}
	if response == nil {
		_ = c.Error(errorutil.NewWithMessage(errorutil.CodeNotFound, "note not found"))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response)
}
