package jsoncodec

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"project-template/internal/foundation/errors"

	"github.com/smartwalle/rx"
)

// Handler 绑定请求体、调用业务函数并输出结果。
//
// 业务函数的返回值就是响应里的 data；错误交给 ErrorEncoder 统一编码成响应，
// 业务函数不用关心 HTTP 细节。
func Handler[T any, R any](fn func(c *rx.Context, req T) (R, error)) rx.HandlerFunc {
	return func(c *rx.Context) {
		var req T
		if err := decode(c, &req); err != nil {
			// 错误直接终止，需要在 rx.Engine 的 ErrorEncoder 中统一处理
			c.AbortWithError(errors.InvalidArgument("invalid argument"))
			return
		}

		var data, err = fn(c, req)
		if err != nil {
			// 错误直接终止，需要在 rx.Engine 的 ErrorEncoder 中统一处理
			c.AbortWithError(err)
			return
		}
		c.JSON(http.StatusOK, Response{Data: data})
	}
}

// decode 绑定请求体。
//
// 只有带请求体的方法才解码：POST、PUT、PATCH。
// 其它方法的请求参数为空，不解码。
func decode(c *rx.Context, req any) error {
	switch c.Request.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
	default:
		return nil
	}

	if c.Request.Body == nil || c.Request.Body == http.NoBody || c.Request.ContentLength == 0 {
		return io.EOF
	}
	defer func() {
		_ = c.Request.Body.Close()
	}()

	return json.NewDecoder(c.Request.Body).Decode(req)
}

// ErrorEncoder 把错误编码成响应。
//
// 状态码没有统一规范，由项目自己定：这里业务错误返回 400、其它错误返回 500，
// 要区分就按项目自己的规范改这里。body 里的 code 是错误标识，与状态码没有对应关系。
func ErrorEncoder() rx.ErrorEncoder {
	return func(c *rx.Context, err error) {
		slog.ErrorContext(c.Context(), "请求处理失败",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Any("error", err),
		)

		// 业务错误统一处理
		var nErr errors.Error
		if errors.As(err, &nErr) {
			c.JSON(http.StatusBadRequest, Response{Code: nErr.Code(), Message: nErr.Message()})
			return
		}

		// 其它不可识别为业务错误的类型，统一响应“服务端内部错误”
		c.JSON(http.StatusInternalServerError, Response{Code: errors.CodeInternalServerError, Message: "internal server error"})
	}
}

// NotFoundHandler 是路由未命中时的处理函数，返回 NotFound 业务错误。
func NotFoundHandler() rx.HandlerFunc {
	return func(c *rx.Context) {
		c.AbortWithError(errors.NotFound("route not found"))
	}
}
