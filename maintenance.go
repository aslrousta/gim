package gim

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MaintenanceFunc is a callback which returns true if server is in
// maintenance mode.
type MaintenanceFunc func() (bool, error)

// Maintenance is the middleware which checks server serviceability and
// responds with a 503 status code if the server is in maintenance mode.
func Maintenance(cb MaintenanceFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		maintenance, err := cb()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if maintenance {
			c.AbortWithStatus(http.StatusServiceUnavailable)
			return
		}

		c.Next()
	}
}
