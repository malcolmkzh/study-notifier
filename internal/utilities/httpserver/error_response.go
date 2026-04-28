package httpserver

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/malcolmkzh/study-notifier/internal/utilities/errorutil"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func addErrorResponse() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		lastError := c.Errors.Last()
		if lastError == nil || c.Writer.Written() {
			return
		}

		var appError *errorutil.AppError
		if errors.As(lastError.Err, &appError) {
			c.JSON(appError.HTTPStatus, ErrorResponse{
				Code:    int(appError.Code),
				Message: appError.Message,
			})
			return
		}

		slog.ErrorContext(c.Request.Context(),
			"unhandled request error",
			"error", lastError.Err,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
		)

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    int(errorutil.CodeInternalServerError),
			Message: "internal server error",
		})
	}
}
