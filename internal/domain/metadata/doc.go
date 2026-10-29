// Package metadata 统一管理业务上下文数据的读写。
//
// 需要在请求处理过程中传递的业务数据（当前登录用户、租户等）放这里：
// 每份数据一个文件，类型和读写函数放在一起。
//
// 通用上下文数据（请求 id）不在这里，在 internal/foundation/metadata。
package metadata
