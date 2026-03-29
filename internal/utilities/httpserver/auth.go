package httpserver

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	gojwt "github.com/golang-jwt/jwt/v5"
)

const (
	ContextKeyJWTSubject = "jwt_subject"
	ContextKeyJWTIssuer  = "jwt_issuer"
	ContextKeyJWTEmail   = "jwt_email"
	ContextKeyJWTRole    = "jwt_role"
)

type jwtClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	gojwt.RegisteredClaims
}

func (m *Implementation) parseJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		tokenString, err := extractBearerToken(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		claims := &jwtClaims{}
		token, err := gojwt.ParseWithClaims(tokenString, claims, func(token *gojwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*gojwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(m.config.Config().JWTSecretKey), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set(ContextKeyJWTSubject, claims.Subject)
		c.Set(ContextKeyJWTIssuer, claims.Issuer)
		c.Set(ContextKeyJWTEmail, claims.Email)
		c.Set(ContextKeyJWTRole, claims.Role)
		c.Next()
	}
}

func (m *Implementation) requireIssuer() gin.HandlerFunc {
	return func(c *gin.Context) {
		issuer := c.GetString(ContextKeyJWTIssuer)
		if issuer == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token issuer"})
			return
		}
		if issuer != m.config.Config().JWTIssuer {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token issuer"})
			return
		}

		c.Next()
	}
}

func extractBearerToken(authHeader string) (string, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
		return "", errors.New("invalid bearer token")
	}

	return strings.TrimSpace(parts[1]), nil
}
