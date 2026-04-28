package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/malcolmkzh/study-notifier/internal/utilities/config"

	"github.com/gin-gonic/gin"
)

type Implementation struct {
	router *gin.Engine
	server *http.Server
	config config.Utility
}

func NewHttpServerUtility(configUtility config.Utility) *Implementation {
	router := gin.New()
	router.Use(corsMiddleware(configUtility), gin.Logger(), gin.Recovery(), addErrorResponse())

	return &Implementation{
		router: router,
		config: configUtility,
	}
}

func corsMiddleware(configUtility config.Utility) gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigin := configUtility.Config().CORSAllowedOrigin
		if strings.TrimSpace(allowedOrigin) == "" {
			allowedOrigin = "http://localhost:5173"
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func (m *Implementation) RegisterEndpoint(ctx context.Context, request RegisterEndpointRequest) error {
	_ = ctx

	//Middleware
	handlers := gin.HandlersChain{}

	//JWT validation Middleware
	if request.RequireAuth {
		handlers = append(handlers, m.parseJWT())
		handlers = append(handlers, m.requireIssuer())
	}

	handlers = append(handlers, request.Fn)
	m.router.Handle(request.Method, request.Path, handlers...)

	return nil
}

func (m *Implementation) Serve(ctx context.Context) error {
	_ = ctx

	if m.server != nil {
		return errors.New("http server already started")
	}

	port := m.config.Config().Port
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: m.router,
	}

	go func() {
		slog.Info("HTTP server started", "port", port)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server error", "error", err.Error())
		}
	}()

	m.server = server
	return nil
}

func (m *Implementation) Shutdown(ctx context.Context) error {
	if m.server == nil {
		return errors.New("server is nil")
	}
	return m.server.Shutdown(ctx)
}
