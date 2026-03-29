package httpclient

import (
	"context"
	"net/http"
	"time"
)

type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    any
	Timeout time.Duration
}

type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

type Utility interface {
	Do(ctx context.Context, request Request) (*Response, error)
	DoJSON(ctx context.Context, request Request, target any) (*Response, error)
}
