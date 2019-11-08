package gim_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aslrousta/gim"
	"github.com/gin-gonic/gin"
)

func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Permissive", func(t *testing.T) {
		opt := gim.CORSOptions{
			Permissive: true,
		}

		r := gin.New()
		r.Use(gim.CORS(&opt))

		r.GET("/", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("unexpected status: %d", w.Code)
		}
	})

	t.Run("MissingOrigin", func(t *testing.T) {
		opt := gim.CORSOptions{
			Origins: []string{"https://example.com"},
		}

		r := gin.New()
		r.Use(gim.CORS(&opt))

		r.GET("/", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("unexpected status: %d", w.Code)
		}
	})

	t.Run("InvalidOrigin", func(t *testing.T) {
		opt := gim.CORSOptions{
			Origins: []string{"https://example.com"},
		}

		r := gin.New()
		r.Use(gim.CORS(&opt))

		r.GET("/", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Origin", "https://other.com")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("unexpected status: %d", w.Code)
		}

		if w.Body.String() != "Invalid Origin" {
			t.Errorf("unexpected body: %q", w.Body.String())
		}
	})

	t.Run("Options", func(t *testing.T) {
		opt := gim.CORSOptions{
			Origins:          []string{"https://example.com"},
			Methods:          []string{http.MethodGet, http.MethodPost},
			Headers:          []string{"Content-Type", "Authorization"},
			AllowCredentials: true,
			MaxAge:           3600,
		}

		r := gin.New()
		r.Use(gim.CORS(&opt))

		r.GET("/", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodOptions, "/", nil)
		req.Header.Set("Origin", "https://example.com")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code == http.StatusForbidden {
			t.Errorf("unexpected status: %d", w.Code)
		}

		allowedOrigin := w.Header().Get("Access-Control-Allow-Origin")
		if allowedOrigin != "https://example.com" {
			t.Errorf("unexpected allowed origin: %s", allowedOrigin)
		}

		allowedMethods := w.Header().Get("Access-Control-Allow-Methods")
		if allowedMethods != "GET, POST" {
			t.Errorf("unexpected allowed methods: %s", allowedMethods)
		}

		allowedHeaders := w.Header().Get("Access-Control-Allow-Headers")
		if allowedHeaders != "Content-Type, Authorization" {
			t.Errorf("unexpected allowed headers: %s", allowedHeaders)
		}

		allowCredentials := w.Header().Get("Access-Control-Allow-Credentials")
		if allowCredentials != "true" {
			t.Errorf("unexpectedly does not allow credentials")
		}

		maxAge := w.Header().Get("Access-Control-Max-Age")
		if maxAge != "3600" {
			t.Errorf("unexpected max age: %s", maxAge)
		}
	})

	t.Run("InvalidMethod", func(t *testing.T) {
		opt := gim.CORSOptions{
			Origins: []string{"https://example.com"},
			Methods: []string{http.MethodGet, http.MethodPost},
		}

		r := gin.New()
		r.Use(gim.CORS(&opt))

		r.GET("/", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set("Origin", "https://example.com")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("unexpected status: %d", w.Code)
		}

		if w.Body.String() != "Invalid HTTP Method" {
			t.Errorf("unexpected body: %q", w.Body.String())
		}
	})

	t.Run("Valid", func(t *testing.T) {
		opt := gim.CORSOptions{
			Origins: []string{"https://example.com"},
			Methods: []string{http.MethodGet, http.MethodPost},
		}

		r := gin.New()
		r.Use(gim.CORS(&opt))

		r.GET("/", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Origin", "https://example.com")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("unexpected status: %d", w.Code)
		}

		allowedOrigin := w.Header().Get("Access-Control-Allow-Origin")
		if allowedOrigin != "https://example.com" {
			t.Errorf("unexpected allowed origin: %s", allowedOrigin)
		}
	})
}
