package logging

import (
	"context"
	"log/slog"
	"os"
	"project-template/internal/foundation/metadata"

	"github.com/smartwalle/logkit"
)

// Setup 按配置初始化全局 slog logger。
//
// slog.SetDefault 会同时接管标准 log 包的输出，
// 因此第三方库通过 log.Print 打出的日志也会按同样的格式输出。
func Setup(cfg Config) error {
	var handler = logkit.NewEnrichHandler(
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{AddSource: true}),
		func(ctx context.Context) []slog.Attr {
			var rid = metadata.RequestId(ctx)
			if rid == "" {
				return nil
			}
			return []slog.Attr{slog.String("requestId", rid)}
		},
	)

	slog.SetDefault(slog.New(handler))
	return nil
}
