package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	kHeaderRequestId = "Request-Id"
)

func MidRequestId() func(c *gin.Context) {
	return func(c *gin.Context) {
		var requestId = c.Request.Header.Get(kHeaderRequestId)
		if requestId == "" {
			requestId = uuid.New().String()
		}
		c.Writer.Header().Set(kHeaderRequestId, requestId)
		c.Next()
	}
}
