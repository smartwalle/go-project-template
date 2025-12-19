package internal

import (
	"github.com/smartwalle/dbs"
	"gopkg.in/yaml.v3"
	"os"
)

type Enviroment struct {
	Config Config
	DB     *dbs.DB
}

func NewEnviroment() *Enviroment {
	var env = &Enviroment{}
	return env
}

func (env *Enviroment) Prepare(filepath string) (err error) {
	// 加载配置文件
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = yaml.NewDecoder(file).Decode(&env.Config); err != nil {
		return err
	}

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
