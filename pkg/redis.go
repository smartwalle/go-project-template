package pkg

import (
	"context"
	"github.com/smartwalle/dbr"
	"github.com/smartwalle/log4go"
	"os"
)

func NewRedis(conf RedisConfig) dbr.UniversalClient {
	rClient, err := dbr.New(conf.Addr, conf.Password, conf.DB, conf.PoolSize, conf.MinIdleConns)
	if err != nil {
		log4go.Errorln("连接 Redis 数据库发生错误:", err)
		os.Exit(-1)
		return nil
	}
	if rClient == nil {
		log4go.Errorln("连接 Redis 数据库失败")
		os.Exit(-1)
	}
	if _, err := rClient.Ping(context.TODO()).Result(); err != nil {
		log4go.Errorln("连接 Redis 数据库发生错误:", err)
		os.Exit(-1)
	}
	return rClient
}
