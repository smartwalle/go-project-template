// Package jsonrender 是 rx 之上的 JSON 渲染层：请求绑定、业务错误映射与统一响应输出。
//
// 业务 Handler 只需要写"拿到入参、调用用例、返回数据或错误"，
// 其余都交给这里的 HandleJSON / HandleQuery。
package jsonrender

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/smartwalle/rx"

	"project-template/internal/foundation/errors"
)

// HandleJSON 绑定请求体、调用业务函数并输出结果。
//
// 业务函数的返回值就是响应里的 data；错误统一在这里映射成状态码与错误响应，
// 业务函数不用关心 HTTP 细节。
func HandleJSON[T any, R any](fn func(c *rx.Context, req T) (R, error)) rx.HandlerFunc {
	return func(c *rx.Context) {
		var req T
		if err := bindJSON(c, &req); err != nil {
			c.AbortWithError(errors.New(CodeInvalidJSON, "request body must be valid JSON"))
			return
		}

		var data, err = fn(c, req)
		if err != nil {
			c.AbortWithError(err)
			return
		}
		Success(c, http.StatusOK, "", "", data)
	}
}

// HandleQuery 调用业务函数并输出结果，请求不带 body。
func HandleQuery[R any](fn func(c *rx.Context) (R, error)) rx.HandlerFunc {
	return func(c *rx.Context) {
		var data, err = fn(c)
		if err != nil {
			c.AbortWithError(err)
			return
		}
		Success(c, http.StatusOK, "", "", data)
	}
}

// bindJSON 绑定请求体。
func bindJSON(c *rx.Context, req any) error {
	if c.Request.Body == nil {
		return io.EOF
	}
	defer func() {
		_ = c.Request.Body.Close()
	}()

	return json.NewDecoder(c.Request.Body).Decode(req)
}
