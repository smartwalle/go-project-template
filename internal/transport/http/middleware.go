package http

import (
	"log/slog"
	"runtime/debug"
	"time"
	"uuid"

	"project-template/internal/foundation/errors"

	"github.com/smartwalle/rx"
)

// MidRecovery 捕获 panic，避免连接被直接断开。
func MidRecovery() rx.HandlerFunc {
	return func(c *rx.Context) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("panic recovered", "panic", r, "stack", string(debug.Stack()))
				c.AbortWithError(errors.New(errors.CodeInternalError, "internal server error"))
			}
		}()
		c.Next()
	}
}

// MidLogger 是基于 slog 的访问日志中间件。
func MidLogger() rx.HandlerFunc {
	return func(c *rx.Context) {
		var start = time.Now()
		c.Next()
		slog.Info("http request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration", time.Since(start).String(),
		)
	}
}

// HeaderRequestId 是请求 ID 的请求头名称。
const HeaderRequestId = "Request-Id"

// MidRequestId 保证每个请求都有请求 ID。
//
// 请求头带了 Request-Id 就沿用，没带就生成一个 UUID v4，
// 最终生效的请求 ID 都会写进响应头。
func MidRequestId() rx.HandlerFunc {
	return func(c *rx.Context) {
		var rid = c.Request.Header.Get(HeaderRequestId)
		if rid == "" {
			rid = uuid.New().String()
		}
		c.Writer.Header().Set(HeaderRequestId, rid)
		c.Next()
	}
}
