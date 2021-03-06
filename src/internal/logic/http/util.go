package http

import (
	"bdim/src/internal/logic"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joomcode/errorx"
	"log"
	"net/http"
	"time"
)

var (
	ErrNS             = errorx.NewNamespace("error.api")
	ErrOther          = ErrNS.NewType("other")
	ErrInvalidRequest = ErrNS.NewType("invalid_request")
	ErrInternalServer = ErrNS.NewType("internal_server_error")
	ErrNotFound       = ErrNS.NewType("resource_not_found")
)

type APIError struct {
	Error    bool   `json:"error"`
	Message  string `json:"message"`
	Code     string `json:"code"`
	FullText string `json:"full_text"`
}

// MWHandleErrors creates a middleware that turns (last) error in the context into an APIError json response.
// In handlers, `c.Error(err)` can be used to attach the error to the context.
// When error is attached in the context:
// - The handler can optionally assign the HTTP status code.
// - The handler must not self-generate a response body.
func MWHandleErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err == nil {
			return
		}

		statusCode := c.Writer.Status()
		if statusCode == http.StatusOK {
			statusCode = http.StatusInternalServerError
		}

		innerErr := errorx.Cast(err.Err)
		if innerErr == nil {
			innerErr = ErrOther.WrapWithNoMessage(err.Err)
		}

		c.AbortWithStatusJSON(statusCode, APIError{
			Error:    true,
			Message:  innerErr.Error(),
			Code:     errorx.GetTypeName(innerErr),
			FullText: fmt.Sprintf("%+v", innerErr),
		})
	}
}

func RateMiddleware(limiter *logic.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		timeDur, _ := time.ParseDuration(limiter.Dur)
		if !limiter.Allow(c.ClientIP(), limiter.Count, timeDur) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "too many requests from the same IP address",
			})
			log.Println("too many requests")
			return
		}
		c.Next()
	}
}
