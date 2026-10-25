package database

import (
	"fmt"

	"github.com/smartwalle/dbs"
	"github.com/smartwalle/dbs/dialect/postgres"

	// 注册 PostgreSQL 驱动。驱动只在这一层出现。
	_ "github.com/lib/pq"
)

const (
	// DefaultMaxOpenConns 是默认的最大连接数。
	DefaultMaxOpenConns = 20
	// DefaultMaxIdleConns 是默认的最大空闲连接数。
	DefaultMaxIdleConns = 5
)

// OpenPostgres 建立 PostgreSQL 连接池，并为 dbs 设置 PostgreSQL 方言。
//
// dbs.Open 在返回前会完成一次带超时的连通性检查，检查不通过时返回错误，调用方不会拿到半成品资源。
func OpenPostgres(cfg Config) (*dbs.DB, error) {
	var maxOpenConns = cfg.MaxOpenConns
	if maxOpenConns <= 0 {
		maxOpenConns = DefaultMaxOpenConns
	}
	var maxIdleConns = cfg.MaxIdleConns
	if maxIdleConns <= 0 {
		maxIdleConns = DefaultMaxIdleConns
	}

	// 错误信息里不包含 DSN，避免把口令写进日志。
	db, err := dbs.Open(Driver, cfg.DSN, maxOpenConns, maxIdleConns)
	if err != nil {
		return nil, fmt.Errorf("open database with driver %q: %w", Driver, err)
	}
	db.UseDialect(postgres.Dialect())
	return db, nil
}
