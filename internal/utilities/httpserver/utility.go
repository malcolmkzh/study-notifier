package httpserver

import (
	"context"

	"github.com/gin-gonic/gin"
)

type RegisterEndpointRequest struct {
	Method      string
	Path        string
	Fn          gin.HandlerFunc
	RequireAuth bool
}

type Utility interface {
	RegisterEndpoint(ctx context.Context, request RegisterEndpointRequest) error
	Serve(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
