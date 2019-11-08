package gim

import (
	"net/http"
	"strings"

	"github.com/aslrousta/ujwt"
	"github.com/gin-gonic/gin"
)

// claimsKey is the key by-which claims are stored in the request context.
const claimsKey = "__claims"

// Auth is the middleware for handling authentication.
func Auth(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		parts := strings.Split(header, " ")

		if len(parts) < 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next()
			return
		}

		var claims ujwt.Claims

		switch err := ujwt.Parse(secretKey, parts[1], &claims); err {
		case ujwt.ErrInvalidSecretKey:
			c.AbortWithError(http.StatusInternalServerError, err)
		case nil:
			c.Set(claimsKey, claims)
			c.Next()
		default:
			c.Next()
		}
	}
}

// RequestClaims returns the request claims.
func RequestClaims(c *gin.Context) (ujwt.Claims, bool) {
	v, ok := c.Get(claimsKey)
	if !ok {
		return ujwt.Claims{}, false
	}

	claims, ok := v.(ujwt.Claims)
	if !ok {
		return ujwt.Claims{}, false
	}

	return claims, true
}
