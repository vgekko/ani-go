package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	Generator func() string
	Handler   func(c *gin.Context, requestID string)
)

var headerXRequestID = "X-Request-ID"

// config defines the config for RequestID middleware
type config struct {
	// Generator defines a function to generate an ID.
	// Optional. Default: func() string {
	//   return uuid.New().String()
	// }
	generator Generator
	handler   Handler
}

// RequestID initializes the RequestID middleware.
func RequestID() gin.HandlerFunc {
	cfg := &config{
		generator: func() string {
			return uuid.New().String()
		},
	}

	return func(c *gin.Context) {
		// Get id from request
		rid := c.GetHeader(headerXRequestID)
		if rid == "" {
			rid = cfg.generator()
			c.Request.Header.Add(headerXRequestID, rid)
		}
		if cfg.handler != nil {
			cfg.handler(c, rid)
		}
		// Set the id to ensure that the request_id is in the response
		c.Header(headerXRequestID, rid)
		c.Next()
	}
}

// GetRequestID returns the request identifier
func GetRequestID(c *gin.Context) string {
	return c.GetHeader(headerXRequestID)
}
