package gim

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSOptions are parameters used to configure CORS middleware.
type CORSOptions struct {
	// Origins are the origins which are allowed to access resources on
	// this server. A nil or empty list accepts all origins.
	Origins []string

	// Methods is the list of methods which are acceptable.
	Methods []string

	// Methods is the list of headers which are acceptable.
	Headers []string

	// AllowCredentials determines if clients are allowed to send cookies or
	// HTTP authentication.
	AllowCredentials bool

	// MaxAge is the maximum duration in seconds which CORS headers are valid.
	MaxAge int

	// Permissive determines if requests without `Origin` should be accepted.
	// Permitted requests will not receive CORS headers.
	Permissive bool
}

func (opt *CORSOptions) allowsOrigin(origin string) bool {
	if origin == "" {
		return false
	}

	if len(opt.Origins) == 0 {
		return true
	}

	for _, o := range opt.Origins {
		if o == origin {
			return true
		}
	}

	return false
}

func (opt *CORSOptions) allowsMethod(method string) bool {
	for _, m := range opt.Methods {
		if strings.ToUpper(m) == strings.ToUpper(method) {
			return true
		}
	}

	return false
}

func (opt *CORSOptions) allowsHeader(header string) bool {
	for _, h := range opt.Headers {
		if strings.ToUpper(h) == strings.ToUpper(header) {
			return true
		}
	}

	return false
}

// CORS is a middleware to handle cross-origin resource sharing control.
func CORS(opt *CORSOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" && opt.Permissive {
			c.Next()
			return
		}

		if !opt.allowsOrigin(origin) {
			c.String(http.StatusForbidden, "Invalid Origin")
			c.Abort()
			return
		}

		if c.Request.Method == http.MethodOptions {
			writeCORSHeaders(c, origin, opt)
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		if !opt.allowsMethod(c.Request.Method) {
			c.String(http.StatusForbidden, "Invalid HTTP Method")
			c.Abort()
			return
		}

		// TODO check headers

		c.Header("Access-Control-Allow-Origin", origin)
		c.Next()
	}
}

func writeCORSHeaders(c *gin.Context, origin string, opt *CORSOptions) {
	var allowsCredentials string
	if opt.AllowCredentials {
		allowsCredentials = "true"
	} else {
		allowsCredentials = "false"
	}

	c.Header("Access-Control-Allow-Origin", origin)
	c.Header("Access-Control-Allow-Methods", strings.Join(opt.Methods, ", "))
	c.Header("Access-Control-Allow-Headers", strings.Join(opt.Headers, ", "))
	c.Header("Access-Control-Allow-Credentials", allowsCredentials)
	c.Header("Access-Control-Max-Age", strconv.Itoa(opt.MaxAge))
}
