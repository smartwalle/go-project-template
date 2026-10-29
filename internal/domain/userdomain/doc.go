// Package userdomain 定义用户领域的实体、业务契约与领域错误。
//
// 文件划分：
//   - entity.go：表结构映射，一表一个 XXXEntity；
//   - model.go：整理和组合数据用的模型，不绑定数据库表；
//   - service.go：Service 契约（用例接口），以及请求、响应结构；
//   - repository.go：Repository 契约（数据访问接口）；
//   - errors.go：领域错误。
package userdomain
