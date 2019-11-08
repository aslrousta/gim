package gim_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aslrousta/gim"
	"github.com/aslrousta/ujwt"
	"github.com/gin-gonic/gin"
)

func TestHasRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()

	r.GET("/permitted", gim.Auth("somerandomkey"), gim.HasRole(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.GET("/restricted", gim.Auth("somerandomkey"), gim.HasRole("admin"), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	noRoleToken, _ := ujwt.Issue("somerandomkey", "user", "example.com", nil)
	adminToken, _ := ujwt.Issue("somerandomkey", "admin", "example.com", []string{"admin"})

	t.Run("PermissiveUnauthenticated", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/permitted", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("unexpected status = %d", w.Code)
		}
	})

	t.Run("PermissiveAuthenticated", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/permitted", nil)
		w := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", noRoleToken))

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("unexpected status = %d", w.Code)
		}
	})

	t.Run("RestrictedUnauthenticated", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/restricted", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("unexpected status = %d", w.Code)
		}
	})

	t.Run("RestrictedUnauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/restricted", nil)
		w := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", noRoleToken))

		r.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("unexpected status = %d", w.Code)
		}
	})

	t.Run("RestrictedAuthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/restricted", nil)
		w := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("unexpected status = %d", w.Code)
		}
	})
}
