package internal

import (
	"github.com/smartwalle/dbs"
)

type Enviroment struct {
	Config Config
	DB     *dbs.DB
}

func NewEnviroment() *Enviroment {
	var env = &Enviroment{}
	return env
}

func (env *Enviroment) Prepare() (err error) {
	// 加载配置文件

	// 连接数据库
	//env.DB, err = dbs.Open(
	//	env.Config.SQL.Driver,
	//	env.Config.SQL.URL,
	//	env.Config.SQL.MaxOpen,
	//	env.Config.SQL.MaxIdle,
	//)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (env *Enviroment) Clear() {
	if env.DB != nil {
		env.DB.Close()
	}
}
