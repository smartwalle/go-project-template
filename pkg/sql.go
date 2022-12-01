package pkg

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/log4go"
	"os"
)

func NewSQL(conf SQLConfig) dbs.DB {
	// 数据库日志
	var dbLogger = log4go.New()
	dbLogger.AddWriter("file", log4go.NewFileWriter(log4go.LevelTrace, log4go.WithLogDir("./logs_dbs"), log4go.WithMaxAge(60*60*24*30)))
	dbs.SetLogger(dbLogger)

	db, err := dbs.New(conf.Driver, conf.URL, conf.MaxOpen, conf.MaxIdle)
	if err != nil {
		log4go.Errorln("连接MySQL数据库发生错误: ", err)
		os.Exit(-1)
	}
	if err = db.Ping(); err != nil {
		log4go.Errorln("连接 SQL 数据库发生错误:", err)
		os.Exit(-1)
	}
	return db
}
