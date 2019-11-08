package gim

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HasRole is the middleware for allowing specific roles to access a resource.
//
// If the authenticated user does not have at-least one role in the allowed
// roles, it will be reported back by an access forbidden error. It has no
// effects if roles are empty.
//
// Note: HasRole middleware must always come after Auth middleware.
func HasRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(roles) == 0 {
			c.Next()
			return
		}

		claims, ok := RequestClaims(c)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		for _, role := range roles {
			if claims.HasRole(role) {
				c.Next()
				return
			}
		}

		c.AbortWithStatus(http.StatusForbidden)
	}
}
