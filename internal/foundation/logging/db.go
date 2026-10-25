package logging

import (
	"context"
	"log/slog"
	"runtime"
	"time"
)

// DBLogger 把 dbs 的 SQL 日志接到 slog 上。
type DBLogger struct {
}

// Trace 按 SQL 是否出错，用 info 或 error 级别记录一条日志。
//
// depth 是 dbs 记录的调用深度，用来定位发起 SQL 的位置。
func (DBLogger) Trace(ctx context.Context, depth int, begin time.Time, sql string, args []any, rowsAffected int64, err error) {
	pc, _, _, _ := runtime.Caller(depth)

	var level = slog.LevelInfo
	var msg = "执行 SQL"
	var attrs = []any{
		slog.String("sql", sql),
		slog.Any("args", args),
		slog.Duration("elapsedTime", time.Since(begin)),
	}
	if err != nil {
		level, msg = slog.LevelError, "执行 SQL 异常"
		attrs = append(attrs, slog.Any("error", err))
	} else {
		attrs = append(attrs, slog.Int64("rowsAffected", rowsAffected))
	}

	var record = slog.NewRecord(time.Now(), level, msg, pc)
	record.Add(attrs...)
	slog.Default().Handler().Handle(ctx, record)
}
