package gim_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aslrousta/gim"
	"github.com/gin-gonic/gin"
)

func TestLang(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(gim.Lang("en-US", "fa-IR"))

	r.GET("/", func(c *gin.Context) {
		lang := gim.RequestLang(c)
		c.String(http.StatusOK, lang.String())
	})

	t.Run("Default", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		body := w.Body.String()
		if body != "en-US" {
			t.Errorf("unexpected lang: %s", body)
		}
	})

	t.Run("Explicit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		req.Header.Add("Accept-Language", "fa-IR")

		r.ServeHTTP(w, req)

		body := w.Body.String()
		if body != "fa-IR" {
			t.Errorf("unexpected lang: %s", body)
		}
	})
}
