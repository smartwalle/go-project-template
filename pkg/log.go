package pkg

import (
	"fmt"
	"github.com/smartwalle/nlog"
	"github.com/smartwalle/nlog/rfile"
	"log/slog"
	"os"
)

func InitDefaultLog(conf ServerConfig) func() {
	var mHandler = nlog.NewMultiHandler()

	if conf.LogStdout {
		mHandler.Add(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	var delaySync = func() {}

	if conf.LogFile {
		var file, err = rfile.New("./logs/temp.log", rfile.WithBuffer(1*1024*1024), rfile.WithMaxAge(7*24*60*60), rfile.WithMaxSize(10*1024*1024))
		if err != nil {
			fmt.Println("初始化日志发生错误:", err)
			os.Exit(-1)
		}
		mHandler.Add(slog.NewTextHandler(file, &slog.HandlerOptions{Level: slog.LevelDebug}))
		delaySync = func() {
			file.Close()
		}
	}

	var logger = slog.New(mHandler)
	slog.SetDefault(logger)
	return delaySync
}
