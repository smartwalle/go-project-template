// Package database 负责数据库连接、连接池参数与关闭。
//
// 该包是唯一允许出现数据库驱动与 DSN 的位置。
package database

import (
	"fmt"
	"strings"
	"time"
)

// Driver 是默认的 database/sql 驱动名。
const Driver = "postgres"

// Config 是数据库连接池的配置，也是配置文件里 database 段的定义。
//
// 连接参数和校验规则由平台层自己维护。
type Config struct {
	// DSN 是连接字符串，包含口令，不要提交到仓库。
	DSN string `yaml:"dsn"`

	// MaxOpenConns 是最大打开连接数。
	MaxOpenConns int `yaml:"max_open_conns"`

	// MaxIdleConns 是最大空闲连接数。
	MaxIdleConns int `yaml:"max_idle_conns"`

	ConnectTimeout time.Duration `yaml:"connect_timeout"`
}

// Validate 校验配置。错误信息带具体配置路径，便于定位问题。
func (c Config) Validate() error {
	if strings.TrimSpace(c.DSN) == "" {
		return fmt.Errorf("database.dsn is required")
	}
	if c.MaxOpenConns <= 0 {
		return fmt.Errorf("database.max_open_conns must be positive")
	}
	if c.MaxIdleConns < 0 {
		return fmt.Errorf("database.max_idle_conns must not be negative")
	}
	if c.MaxIdleConns > c.MaxOpenConns {
		return fmt.Errorf("database.max_idle_conns must not exceed database.max_open_conns")
	}
	if c.ConnectTimeout <= 0 {
		return fmt.Errorf("database.connect_timeout must be positive")
	}
	return nil
}
