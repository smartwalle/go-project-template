package pkg

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/log4go"
	"os"
)

func NewSQL(conf SQLConfig) *dbs.DB {
	db, err := dbs.Open(conf.Driver, conf.URL, conf.MaxOpen, conf.MaxIdle)
	if err != nil {
		log4go.Errorln("连接MySQL数据库发生错误: ", err)
		os.Exit(-1)
	}
	if err = db.Ping(); err != nil {
		log4go.Errorln("连接 SQL 数据库发生错误:", err)
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
