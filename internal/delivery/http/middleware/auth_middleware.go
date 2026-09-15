package middleware

import (
	"strings"

	"mkp-cinema-ticketing/internal/config"
	"mkp-cinema-ticketing/pkg/jwt"
	"mkp-cinema-ticketing/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	CtxUserIDKey = "userID"
	CtxEmailKey  = "email"
	CtxRoleKey   = "role"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization token is required (Header: Authorization: Bearer <token>)")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Unauthorized(c, "Invalid Authorization header format. Expected 'Bearer <token>'")
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(parts[1])
		claims, err := jwt.ValidateToken(tokenString, cfg.JWTSecret)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token: "+err.Error())
			c.Abort()
			return
		}

		c.Set(CtxUserIDKey, claims.UserID)
		c.Set(CtxEmailKey, claims.Email)
		c.Set(CtxRoleKey, claims.Role)

		c.Next()
	}
}
