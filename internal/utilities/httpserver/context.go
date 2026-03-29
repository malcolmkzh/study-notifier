package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCurrentUserID(c *gin.Context) (string, bool) {
	userID := c.GetString(ContextKeyJWTSubject)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated user"})
		return "", false
	}

	return userID, true
}
