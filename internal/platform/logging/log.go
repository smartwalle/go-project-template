// Package logging 负责全局 slog logger 的初始化。
//
// 该包基于标准库 log/slog：初始化后 slog 的默认 logger 就接管了全局输出，
// 业务代码直接用 slog 记录日志，不需要再包一层日志函数。
package logging

import (
	"log/slog"
	"os"
)

// Setup 按配置初始化全局 slog logger。
//
// slog.SetDefault 会同时接管标准 log 包的输出，
// 因此第三方库通过 log.Print 打出的日志也会按同样的格式输出。
func Setup(cfg Config) error {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{AddSource: true})))
	return nil
}
