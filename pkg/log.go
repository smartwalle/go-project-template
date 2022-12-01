package pkg

import (
	"fmt"

	"github.com/smartwalle/log4go"
)

func InitDefaultLog(conf ServerConfig) {
	log4go.SetPrefix(fmt.Sprintf("[%s] ", conf.Name))
	if conf.LogStdout == false {
		log4go.RemoveWriter("stdout")
	}
	if conf.LogFile {
		log4go.AddWriter("file", log4go.NewFileWriter(log4go.LevelTrace, log4go.WithLogDir("./logs"), log4go.WithMaxAge(60*60*24*30)))
	}
}
