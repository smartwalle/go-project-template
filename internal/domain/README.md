# internal/domain

领域契约：业务领域的接口和结构体定义在这里，不写实现。

**本文是这一层的规范表述**：下面的代码是标准写法，`internal/domain/userdomain` 是它的可运行实例。
两边不一致时以本文为准，并同步更新实例。

一个业务领域一个包（`internal/domain/<模块>/`），包内按角色分文件：

| 文件 | 放什么 |
| --- | --- |
| `doc.go` | 包注释，说明这个领域的文件划分 |
| `entity.go` | 表结构映射，一个表一个 XXXEntity |
| `model.go` | 整理、组合数据的结构体 |
| `service.go` | Service 契约（用例接口），以及请求、响应结构 |
| `repository.go` | Repository 契约（数据访问接口） |
| `errors.go` | 领域专属错误 |

依赖方向：domain 不 import transport、service、repository，也不 import 数据库驱动；
仓储契约基于 dbs 的通用 Repository。

领域包之间可以根据业务关系建立依赖，但应遵循以下约定：

- 领域包只依赖业务上确实需要的其他领域包，不为了方便而随意引用。
- 优先通过业务标识（ID）建立关联，避免不必要地依赖其他领域的内部结构。
- 确实需要引用其他领域的公开模型或类型时，可以直接依赖对应领域包。
- 基础业务领域不应依赖组合业务领域。
- 禁止循环依赖。
- 不允许因为两个领域之间存在调用关系，就将它们合并成一个领域包。

## 领域共用的内容

这个目录本身是 `package domain`，用于定义多个业务领域共用的类型和契约。

适合放在这里的内容包括：

- 多个领域都会使用的枚举类型。
- 多个领域都会使用的值对象。
- 多个领域都会使用的公共业务定义。

只在一个业务领域中使用的定义，应放在对应的领域包中。

不要仅仅因为某个类型看起来通用，就将它放到这里。
是否放在根包，应以是否真正被多个业务领域共同使用为依据。

```go
package domain

// DeleteStatus 表示数据的逻辑删除状态。
//
// 正常状态表示数据未被删除，删除状态表示数据已被逻辑删除。
type DeleteStatus int8

const (
	// DeleteStatusNormal 表示数据正常，未被删除。
	DeleteStatusNormal DeleteStatus = 0

	// DeleteStatusDeleted 表示数据已被逻辑删除。
	DeleteStatusDeleted DeleteStatus = 1
)