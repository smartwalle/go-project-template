package app

import (
	"fmt"
	"log/slog"

	"errors"

	"project-template/internal/config"
	"project-template/internal/foundation/database"
	"project-template/internal/foundation/logging"

	"github.com/smartwalle/dbs"
)

// Environment 是应用运行期需要的资源，由它负责创建与释放。
type Environment struct {
	// Config 是应用的运行时配置。
	Config *config.Config

	// Database 是数据库连接池。
	Database *dbs.DB
}

// Prepare 校验配置并初始化应用依赖（日志、数据库）。
//
// 任何一步失败都会返回错误，已经创建的资源会被释放。
func Prepare(cfg *config.Config) (*Environment, error) {
	if cfg == nil {
		return nil, errors.New("config is required")
	}

	var env = &Environment{}
	if err := env.setup(cfg); err != nil {
		return nil, err
	}
	return env, nil
}

func (env *Environment) setup(cfg *config.Config) error {
	env.Config = cfg

	if err := env.initLogger(); err != nil {
		return err
	}
	if err := env.initDatabase(); err != nil {
		return err
	}
	return nil
}

func (env *Environment) initLogger() error {
	if err := logging.Setup(env.Config.Logging); err != nil {
		return fmt.Errorf("setup logging: %w", err)
	}
	slog.Info("日志初始化完成")
	return nil
}

func (env *Environment) initDatabase() error {
	db, err := database.OpenPostgres(env.Config.Database)
	if err != nil {
		return fmt.Errorf("initialize database: %w", err)
	}
	db.UseLogger(&logging.DBLogger{})

	env.Database = db
	slog.Info("数据库初始化完成")
	return nil
}

// Clear 释放应用占用的资源。
func (env *Environment) Clear() {
	if env.Database != nil {
		env.Database.Close()
	}
}
