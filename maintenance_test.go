package gim_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aslrousta/gim"
	"github.com/gin-gonic/gin"
)

func TestMaintenance(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()

	available := true
	cb := func() (bool, error) { return !available, nil }

	r.Use(gim.Maintenance(cb))

	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	t.Run("Available", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("unexpected status: %d", w.Code)
		}
	})

	t.Run("Unavailable", func(t *testing.T) {
		available = false

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("unexpected status: %d", w.Code)
		}
	})
}
