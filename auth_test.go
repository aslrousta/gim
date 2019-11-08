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

func TestAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/", gim.Auth("somerandomstring"), func(c *gin.Context) {
		if _, ok := gim.RequestClaims(c); ok {
			c.String(http.StatusOK, "authenticated")
		} else {
			c.String(http.StatusOK, "unauthenticated")
		}
	})

	t.Run("Unauthenticated", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("unexpected status = %d", w.Code)
		}

		body := w.Body.String()
		if body != "unauthenticated" {
			t.Errorf("unexpected body = %q", body)
		}
	})

	t.Run("WrongAuthentication", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		req.Header.Set("Authorization", "Bearer unimportant")

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("unexpected status = %d", w.Code)
		}

		body := w.Body.String()
		if body != "unauthenticated" {
			t.Errorf("unexpected body = %q", body)
		}
	})

	token, _ := ujwt.Issue("somerandomstring", "user", "example.com", nil)

	t.Run("Authenticated", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("unexpected status = %d", w.Code)
		}

		body := w.Body.String()
		if body != "authenticated" {
			t.Errorf("unexpected body = %q", body)
		}
	})
}
