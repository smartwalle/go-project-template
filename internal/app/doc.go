// Package app 负责依赖组装与生命周期协调。
//
// 应用层只做两件事：按依赖方向创建实现并注入给使用方，
// 把服务交给 bootstrap 管理启动与关闭。
package app
