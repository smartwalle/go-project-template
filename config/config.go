package config

import (
	"github.com/smartwalle/ini4go"
	"go-project-template/pkg"
)

type Config struct {
	Server pkg.ServerConfig `ini:"server"`
	SQL    pkg.SQLConfig    `ini:"sql"`
	Redis  pkg.RedisConfig  `ini:"redis"`
	HTTP   pkg.HTTPConfig   `ini:"http"`
}

func LoadIni(file string) (*Config, error) {
	// 读取主配置文件
	var mIni, err = loadIni(file)
	if err != nil {
		return nil, err
	}
	var config *Config
	if err = mIni.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func loadIni(file string) (*ini4go.Ini, error) {
	var ini = ini4go.New(false)
	if err := ini.LoadFiles(file); err != nil {
		return nil, err
	}
	return ini, nil
}
