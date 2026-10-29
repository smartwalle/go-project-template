# internal/domain

领域契约：业务模块的**接口和结构体**定义在这里，不写实现。

**本文是这一层的规范表述**：下面的代码是标准写法，`internal/domain/userdomain` 是它的可运行实例。
两边不一致时以本文为准，并同步更新实例。

一个业务模块一个包（`internal/domain/<模块>/`），包内按角色分文件：

| 文件 | 放什么 |
| --- | --- |
| `doc.go` | 包注释，说明这个模块的文件划分 |
| `entity.go` | 表结构映射，一表一个 XXXEntity |
| `model.go` | 整理、组合数据的结构体 |
| `service.go` | Service 契约（用例接口），以及请求、响应结构 |
| `repository.go` | Repository 契约（数据访问接口） |
| `errors.go` | 领域专属错误 |
| `context.go` | 需要在上下文里传递的业务数据（按需添加） |

依赖方向：domain 不 import transport、service、repository，也不 import 数据库驱动；
仓储契约基于 dbs 的通用 Repository。

## 领域共用的内容

这个目录本身是 `package domain`，放所有领域共用的内容；
各业务模块的领域包在子目录里，彼此独立。

多个模块都要用的定义（枚举、值对象、公共契约等）放在这里，
避免各模块各写一份；只在一个模块里用的，定义在各自的领域包中。

```go
package domain

// DeleteStatus 是逻辑删除状态。
//
// 表里用 deleted 列记录删除标记：0-正常 1-删除。
// 查询时带上 deleted = 0，就是只取没被删除的数据。
type DeleteStatus int8

const (
	// DeleteStatusNormal 是正常状态，未删除。
	DeleteStatusNormal DeleteStatus = 0

	// DeleteStatusDeleted 是已删除状态。
	DeleteStatusDeleted DeleteStatus = 1
)
```

## entity.go —— 表结构映射

```go
package userdomain

import (
	"time"

	"project-template/internal/domain"
)

// UserEntity 是 users 表对应的实体。
//
// 命名规则：每张数据库表对应一个 Entity 结构体，名称是表名的驼峰形式加 Entity 后缀，
// 例如 users 表对应 UserEntity。
//
// 标识符规则：实体的标识符字段类型统一使用 uint64。
//
// sql tag 对应 users 表的列名，json tag 是 JSON 序列化时的字段名；
// TableName 与 PrimaryKey 让它满足 dbs.Entity。
type UserEntity struct {
	// Id 是用户的业务标识，由数据库生成（default：留空时用数据库的默认值）。
	Id uint64 `sql:"id;default" json:"id"`

	// Name 是用户名称。
	Name string `sql:"name" json:"name"`

	// Email 是用户邮箱，在业务上唯一。
	Email string `sql:"email" json:"email"`

	// CreatedAt 是用户创建时间。
	CreatedAt time.Time `sql:"created_at" json:"created_at"`

	// Deleted 是逻辑删除状态：0-正常 1-删除。
	Deleted domain.DeleteStatus `sql:"deleted" json:"deleted"`
}

// TableName 是 users 表名，供 dbs 的通用 Repository 使用。
func (UserEntity) TableName() string {
	return "users"
}

// PrimaryKey 是主键列名，供 dbs 的通用 Repository 使用。
func (UserEntity) PrimaryKey() string {
	return "id"
}
```

## model.go —— 整理数据的结构体

```go
package userdomain

import "time"

// User 是用户模块的模型。
type User struct {
	// Id 是用户的业务标识。
	Id uint64 `json:"id"`

	// Name 是用户名称。
	Name string `json:"name"`

	// Email 是用户邮箱，在业务上唯一。
	Email string `json:"email"`

	// CreatedAt 是用户创建时间。
	CreatedAt time.Time `json:"created_at"`
}
```

## service.go —— 用例契约、请求与响应

```go
package userdomain

import (
	"context"
	"fmt"
	"strings"

	"project-template/internal/foundation/errors"
	"project-template/internal/foundation/pagination"
)

// CreateRequest 是"创建用户"用例的请求参数。
//
// 命名规则：请求参数统一命名为 XXXRequest，XXX 是用例名（创建用户 → CreateRequest），
// 并统一用指针传递。
// 每个 Request 必须实现 Validate，参数校验写在 Validate 里；
// service 在用例入口处调用 Validate，校验不通过直接返回错误。
type CreateRequest struct {
	// Name 是用户名称。
	Name string `json:"name"`

	// Email 是用户邮箱。
	Email string `json:"email"`
}

// Validate 校验请求参数。
//
// 请求为 nil 时直接返回 MissingArgument；
// 否则先去掉字段首尾空白再校验，并把清洗结果写回请求字段：
// 校验通过后请求里的字段已是规范化后的值，不必再清洗一次。
//
// 参数缺失返回 MissingArgument，参数格式不合法返回 InvalidArgument。
func (r *CreateRequest) Validate() error {
	if r == nil {
		return errors.MissingArgument("request is required")
	}

	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.TrimSpace(r.Email)

	if r.Name == "" {
		return errors.MissingArgument("name is required")
	}
	if r.Email == "" {
		return errors.MissingArgument("email is required")
	}
	if !strings.Contains(r.Email, "@") {
		return errors.InvalidArgument(fmt.Sprintf("email %q is not a valid address", r.Email))
	}
	return nil
}

// ListRequest 是"列出用户"用例的请求参数。
//
// 分页参数内嵌 pagination.Page，客户端看到的字段还是 page、page_size；
// 过滤条件是模糊匹配，为空表示不过滤。
type ListRequest struct {
	// Page 是分页参数：页码和每页条数。
	pagination.Page

	// Name 是用户名称的匹配条件，为空表示不过滤。
	Name string `json:"name"`

	// Email 是用户邮箱的匹配条件，为空表示不过滤。
	Email string `json:"email"`
}

// Validate 校验请求参数。
//
// 请求为 nil 时直接返回 MissingArgument；
// 分页参数由分页包归一化；过滤条件去掉首尾空白，为空表示不过滤。
func (r *ListRequest) Validate() error {
	if r == nil {
		return errors.MissingArgument("request is required")
	}

	r.Page.Normalize()
	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.TrimSpace(r.Email)
	return nil
}

// ListResponse 是"列出用户"用例的返回数据。
//
// 命名规则：返回数据统一命名为 XXXResponse，XXX 是用例名（列出用户 → ListResponse）。
// 用例只需要返回模型本身时直接返回模型（Get、Create 返回 *User）；
// 需要带上模型之外的信息（总数、分页等）时用 Response 包装。
type ListResponse struct {
	// Total 是用户总数。
	Total int `json:"total"`

	// UserList 是用户列表。
	UserList []*User `json:"user_list"`
}

// GetRequest、UpdateRequest 与 CreateRequest 同构，字段不同而已：
// GetRequest 只有 Id；UpdateRequest 有 Id、Name、Email。
// 校验规则同样是 nil 判断 + 字段校验 + 归一化写回，
// 标识符用 `r.Id <= 0` 判断无效（不要用 == 0）。

// Service 是用户模块的业务契约，每个方法是一个业务用例。
//
// 返回约定：用例返回模型 User，不建议直接返回实体 UserEntity：
// Entity 是表映射，直接返回等于把表结构带到调用方，通常由实现层转换成模型再返回。
// 只需要返回模型时直接返回模型，需要带上额外信息时用 XXXResponse 包装。
//
// 入参约定：请求参数统一用 XXXRequest 命名、统一用指针传递，必须实现 Validate 方法；
// 实现方在用例入口处调用 Validate 做校验，校验不通过直接返回错误。
//
// 结构约定：请求参数和响应数据统一定义在本文件，命名为 XXXRequest、XXXResponse；
// 结构里可以引用 model.go 里定义的模型。
//
// 实现由 internal/service 提供，在应用组装层注入。
type Service interface {
	// List 返回用户列表。
	List(ctx context.Context, req *ListRequest) (*ListResponse, error)

	// Get 按 Id 返回用户；用户不存在时返回 ErrNotFound。
	Get(ctx context.Context, req *GetRequest) (*User, error)

	// Create 创建用户；邮箱已存在时返回 ErrAlreadyExists。
	Create(ctx context.Context, req *CreateRequest) (*User, error)

	// Update 更新用户信息（要更新的用户标识在请求参数里）；
	// 用户不存在时返回 ErrNotFound，邮箱已被占用时返回 ErrAlreadyExists。
	Update(ctx context.Context, req *UpdateRequest) (*User, error)
}
```

## context.go —— 上下文里的业务数据

业务相关的上下文数据（当前登录用户、租户等）由对应的领域包提供读写函数。

```go
package authdomain

import "context"

// authUserKey 是承载当前登录用户信息的 context 键类型。
//
// 用途：HTTP 中间件鉴权通过后，把当前登录用户写进 Context，
// 后续业务代码用 AuthUserFrom 读取，
// 不必把用户信息一层层往下传。
//
// 使用私有空结构体作为键，与其它包的键互不冲突。
type authUserKey struct{}

// WithAuthUser 返回一个携带当前登录用户信息的 Context。
//
// 由鉴权中间件在鉴权通过后调用；原 Context 不会被修改，
// 若 Context 中已有用户信息，本次注入的值会覆盖它。
func WithAuthUser(ctx context.Context, user AuthUser) context.Context {
	return context.WithValue(ctx, authUserKey{}, user)
}

// AuthUserFrom 从 Context 中取出当前登录用户信息。
//
// 供业务代码读取登录用户；Context 中没有用户信息
// （例如接口没有经过鉴权中间件）时返回 AuthUser 的零值。
func AuthUserFrom(ctx context.Context) AuthUser {
	user, _ := ctx.Value(authUserKey{}).(AuthUser)
	return user
}
```

## repository.go —— 数据访问契约

```go
package userdomain

import "github.com/smartwalle/dbs"

// Repository 是 users 表对应的仓储接口。
//
// 约定：每张数据库表对应一个仓储接口，放在表所属的领域包内；
// 接口由使用方（Service）定义，实现见 internal/repository。
//
// 查询语义：查不到数据时返回 nil，不返回空记录；命中的记录主键大于 0。
type Repository interface {
	dbs.Repository[UserEntity]
}
```

## errors.go —— 领域错误

```go
package userdomain

import "project-template/internal/foundation/errors"

// 领域错误：通用语义（资源不存在、参数缺失、参数不合法等）直接用 foundation/errors 预置的错误码，
// 业务专属的错误码才定义在领域包里。
var (
	// ErrNotFound 表示目标用户不存在。
	ErrNotFound = errors.NotFound("user not found")
)
```
