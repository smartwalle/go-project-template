// Package metadata 统一管理业务上下文数据的读写。
//
// 需要在请求处理过程中传递的业务数据（当前登录用户、租户等）放这里：
// 每份数据一个文件，只放读写的键和函数；数据类型留在各自的领域包里。
//
// 通用上下文数据（请求 id）不在这里，在 internal/foundation/metadata。
package metadata
