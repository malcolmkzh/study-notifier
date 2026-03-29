package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultTimeout         = 30 * time.Second
	defaultMaxResponseSize = 5 << 20 // 5 MB
)

type Implementation struct {
	client          *http.Client
	timeout         time.Duration
	maxResponseSize int64
}

type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	if e.Body == "" {
		return http.StatusText(e.StatusCode)
	}

	return http.StatusText(e.StatusCode) + ": " + e.Body
}

func NewHTTPClientUtility() *Implementation {
	return &Implementation{
		client:          &http.Client{Timeout: defaultTimeout},
		timeout:         defaultTimeout,
		maxResponseSize: defaultMaxResponseSize,
	}
}

func (m *Implementation) Do(ctx context.Context, request Request) (*Response, error) {
	if strings.TrimSpace(request.Method) == "" {
		return nil, errors.New("http method is required")
	}
	if strings.TrimSpace(request.URL) == "" {
		return nil, errors.New("url is required")
	}

	bodyReader, contentType, err := buildRequestBody(request.Body)
	if err != nil {
		return nil, err
	}

	timeout := request.Timeout
	if timeout <= 0 {
		timeout = m.timeout
	}

	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, request.Method, request.URL, bodyReader)
	if err != nil {
		return nil, err
	}

	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	if contentType != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", contentType)
	}

	response, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(io.LimitReader(response.Body, m.maxResponseSize))
	if err != nil {
		return nil, err
	}

	result := &Response{
		StatusCode: response.StatusCode,
		Headers:    response.Header.Clone(),
		Body:       body,
	}

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return result, &HTTPError{
			StatusCode: response.StatusCode,
			Body:       strings.TrimSpace(string(body)),
		}
	}

	return result, nil
}

func (m *Implementation) DoJSON(ctx context.Context, request Request, target any) (*Response, error) {
	response, err := m.Do(ctx, request)
	if err != nil {
		return response, err
	}

	if target == nil || len(response.Body) == 0 {
		return response, nil
	}

	if err := json.Unmarshal(response.Body, target); err != nil {
		return response, err
	}

	return response, nil
}

func buildRequestBody(body any) (io.Reader, string, error) {
	if body == nil {
		return nil, "", nil
	}

	switch value := body.(type) {
	case []byte:
		return bytes.NewReader(value), "application/octet-stream", nil
	case string:
		return strings.NewReader(value), "text/plain; charset=utf-8", nil
	default:
		data, err := json.Marshal(body)
		if err != nil {
			return nil, "", err
		}
		return bytes.NewReader(data), "application/json", nil
	}
}
