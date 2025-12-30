package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/leehai1107/chophimco-server/pkg/apiwrapper"
)

// AuthMiddleware verifies JWT token and sets user info in context
func AuthMiddleware(jwtService IJWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get token from cookie first
		tokenString, err := c.Cookie("jwt_token")
		if err != nil || tokenString == "" {
			// Fallback to Authorization header for backward compatibility
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenString = parts[1]
				}
			}
		}

		if tokenString == "" {
			apiwrapper.SendUnauthorized(c, "authentication required")
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			apiwrapper.SendUnauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		// Set user info in context for handlers to use
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// OptionalAuthMiddleware sets user info if token is provided, but doesn't require it
func OptionalAuthMiddleware(jwtService IJWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get token from cookie first
		tokenString, err := c.Cookie("jwt_token")
		if err != nil || tokenString == "" {
			// Fallback to Authorization header
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenString = parts[1]
				}
			}
		}

		if tokenString != "" {
			claims, err := jwtService.ValidateToken(tokenString)
			if err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("user_email", claims.Email)
				c.Set("user_role", claims.Role)
			}
		}

		c.Next()
	}
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			apiwrapper.SendUnauthorized(c, "authentication required")
			c.Abort()
			return
		}

		role := userRole.(string)
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		apiwrapper.SendBadRequest(c, "insufficient permissions")
		c.Abort()
	}
}
