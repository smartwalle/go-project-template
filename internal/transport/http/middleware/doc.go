// Package middleware 提供 HTTP 中间件。
//
// 一个中间件一个文件，函数返回 rx.HandlerFunc，在 transport/http 的 Server 里注册。
package middleware
