package http

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"project-template/internal/platform/errors"
	"project-template/internal/transport/http/jsonrender"

	"github.com/smartwalle/rx"
)

// MidRecovery 捕获 panic，避免连接被直接断开。
func MidRecovery() rx.HandlerFunc {
	return func(c *rx.Context) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("panic recovered", "panic", r, "stack", string(debug.Stack()))
				jsonrender.Error(c, http.StatusInternalServerError, errors.CodeInternalError, "internal server error")
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
