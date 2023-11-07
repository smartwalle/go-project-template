package pkg

import (
	"context"
	"github.com/smartwalle/dbr"
	"log/slog"
	"os"
)

func NewRedis(conf RedisConfig) dbr.UniversalClient {
	rClient, err := dbr.New(conf.Addr, conf.Password, conf.DB, conf.PoolSize, conf.MinIdleConns)
	if err != nil {
		slog.Error("连接 Redis 数据库发生错误", slog.Any("error", err))
		os.Exit(-1)
		return nil
	}
	if rClient == nil {
		slog.Error("连接 Redis 数据库失败", slog.Any("error", err))
		os.Exit(-1)
	}
	if _, err = rClient.Ping(context.Background()).Result(); err != nil {
		slog.Error("连接 Redis 数据库发生错误", slog.Any("error", err))
		os.Exit(-1)
	}
	return rClient
}
