// Package main 是应用进程的入口：加载配置、准备运行环境、启动服务。
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"project-template/internal/app"
	"project-template/internal/config"
)

func main() {
	if err := run(); err != nil {
		slog.Error("application failed", "error", err)
		os.Exit(1)
	}
}

// run 加载配置、准备运行环境并启动服务。
//
// 环境在这里创建并释放
func run() error {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	env, err := app.Prepare(cfg)
	if err != nil {
		return fmt.Errorf("prepare environment: %w", err)
	}
	defer env.Clear()

	if err = app.Run(context.Background(), env); err != nil {
		return err
	}
	return nil
}
