package pkg

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandlerWrapper[T any](handler func(c *gin.Context) (T, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		if handler != nil {
			var data, err = handler(c)
			var response = NewResponse(int(CodeSuccess), data)
			var statusCode = http.StatusOK

			if err != nil {
				statusCode = http.StatusBadRequest

				var nErr Error
				if errors.As(err, &nErr) {
					response.Code = nErr.Code()
					response.Message = nErr.Message()
				} else {
					response.Code = int(CodeOperationFailed)
					response.Message = err.Error()
				}
			}

			c.JSON(statusCode, response)
		}
	}
}
