package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// GetUserIDFromContext extracts user ID from gin context
func GetUserIDFromContext(c *gin.Context) (int, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, errors.New("user ID not found in context")
	}

	id, ok := userID.(int)
	if !ok {
		return 0, errors.New("invalid user ID type")
	}

	return id, nil
}

// GetUserEmailFromContext extracts user email from gin context
func GetUserEmailFromContext(c *gin.Context) (string, error) {
	email, exists := c.Get("user_email")
	if !exists {
		return "", errors.New("user email not found in context")
	}

	emailStr, ok := email.(string)
	if !ok {
		return "", errors.New("invalid user email type")
	}

	return emailStr, nil
}

// GetUserRoleFromContext extracts user role from gin context
func GetUserRoleFromContext(c *gin.Context) (string, error) {
	role, exists := c.Get("user_role")
	if !exists {
		return "", errors.New("user role not found in context")
	}

	roleStr, ok := role.(string)
	if !ok {
		return "", errors.New("invalid user role type")
	}

	return roleStr, nil
}
