// Package config 负责配置文件的加载与校验。
//
// 配置结构由各功能模块自己维护：传输层维护 http 段，平台层维护 database 与 logging 段，
// 这里只负责把它们组合成配置文件的结构，并完成解码与校验。
package config
