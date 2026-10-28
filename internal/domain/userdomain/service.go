package userdomain

import (
	"context"
	"fmt"
	"strings"

	"project-template/internal/foundation/errors"
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

// Service 是用户模块的业务契约，每个方法是一个业务用例。
//
// 用例返回对外模型 User，不返回实体 UserEntity：
// Entity 对应数据库结构，不直接暴露给调用方，由实现层转换成模型再返回。
//
// 入参约定：请求参数统一用 XXXRequest 命名、统一用指针传递，必须实现 Validate 方法；
// 实现方在用例入口处调用 Validate 做校验，校验不通过直接返回错误。
//
// 实现由 internal/service 提供，在应用组装层注入。
type Service interface {
	// List 返回全部用户。
	List(ctx context.Context) ([]*User, error)

	// Get 按 Id 返回用户；用户不存在时返回 ErrNotFound。
	Get(ctx context.Context, id int64) (*User, error)

	// Create 创建用户；邮箱已存在时返回 ErrAlreadyExists。
	Create(ctx context.Context, req *CreateRequest) (*User, error)
}
