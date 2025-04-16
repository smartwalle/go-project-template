package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	kHeaderRequestID = "Request-ID"
)

func MidRequestID() func(c *gin.Context) {
	return func(c *gin.Context) {
		var rid = c.Request.Header.Get(kHeaderRequestID)
		if rid == "" {
			rid = uuid.New().String()
		}
		c.Writer.Header().Set(kHeaderRequestID, rid)
		c.Next()
	}
}
