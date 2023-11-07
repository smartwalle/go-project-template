package pkg

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/smartwalle/dbs"
	"log/slog"
	"os"
)

func NewSQL(conf SQLConfig) *dbs.DB {
	db, err := dbs.Open(conf.Driver, conf.URL, conf.MaxOpen, conf.MaxIdle)
	if err != nil {
		slog.Error("连接 SQL 数据库发生错误", slog.Any("error", err))
		os.Exit(-1)
	}
	if err = db.Ping(); err != nil {
		slog.Error("连接 SQL 数据库发生错误", slog.Any("error", err))
		os.Exit(-1)
	}
	return dbs.New(db)
}

func PrepareStatement(ctx context.Context, db *dbs.DB, key string, clause dbs.SQLClause) error {
	var query, _, err = clause.SQL()
	if err != nil {
		return err
	}
	return db.PrepareStatement(ctx, key, query)
}
