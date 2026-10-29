// Package metadata 提供通用的上下文数据管理。
//
// 需要在请求处理过程中传递的通用数据（请求 id 等）在这里读写：
// 每份数据用私有的键类型，写入用 WithXxx，读取用 XxxFrom。
//
// 业务相关的数据不放这里，由各自的领域包提供，见 internal/domain。
package metadata
