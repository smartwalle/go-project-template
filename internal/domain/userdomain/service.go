package userdomain

import (
	"context"
	"fmt"
	"strings"

	"project-template/internal/foundation/errors"
	"project-template/internal/foundation/pagination"
)

// CreateRequest 是创建用户用例的请求参数。
//
// 请求参数统一命名为 XXXRequest，XXX 是用例名，
// 统一使用指针传递，并实现 Validate 方法。
type CreateRequest struct {
	// Name 是用户名称。
	Name string `json:"name"`

	// Email 是用户邮箱。
	Email string `json:"email"`
}

// Validate 校验请求参数。
//
// 请求为 nil 时返回 MissingArgument。
// 字段会先去除首尾空白，再进行校验；
// 校验通过后，清洗结果会写回请求字段。
//
// 参数缺失返回 MissingArgument，格式不合法返回 InvalidArgument。
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

// GetRequest 是获取用户用例的请求参数。
type GetRequest struct {
	// Id 是用户标识。
	Id int64 `json:"id"`
}

// Validate 校验请求参数。
func (r *GetRequest) Validate() error {
	if r == nil {
		return errors.MissingArgument("request is required")
	}
	if r.Id <= 0 {
		return errors.InvalidArgument("id must be greater than 0")
	}

	return nil
}

// ListRequest 是列出用户用例的请求参数。
//
// 分页参数内嵌 pagination.Page，客户端看到的字段仍然是 page、page_size；
// Name 和 Email 是模糊匹配条件，为空时表示不过滤。
type ListRequest struct {
	// Page 是分页参数。
	pagination.Page

	// Name 是用户名称的匹配条件，为空表示不过滤。
	Name string `json:"name"`

	// Email 是用户邮箱的匹配条件，为空表示不过滤。
	Email string `json:"email"`
}

// Validate 校验请求参数。
//
// 分页参数由分页包归一化；
// 过滤条件去除首尾空白，为空时表示不过滤。
func (r *ListRequest) Validate() error {
	if r == nil {
		return errors.MissingArgument("request is required")
	}

	r.Page.Normalize()
	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.TrimSpace(r.Email)

	return nil
}

// UpdateRequest 是更新用户用例的请求参数。
type UpdateRequest struct {
	// Id 是用户标识。
	Id int64 `json:"id"`

	// Name 是用户名称。
	Name string `json:"name"`

	// Email 是用户邮箱。
	Email string `json:"email"`
}

// Validate 校验请求参数。
func (r *UpdateRequest) Validate() error {
	if r == nil {
		return errors.MissingArgument("request is required")
	}

	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.TrimSpace(r.Email)

	if r.Id <= 0 {
		return errors.InvalidArgument("id must be greater than 0")
	}
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

// ListResponse 是列出用户用例的返回数据。
//
// 用例只需要返回模型时直接返回模型；
// 需要附带总数、分页等额外信息时，使用 Response 包装。
type ListResponse struct {
	// Total 是符合查询条件的用户总数。
	Total int `json:"total"`

	// UserList 是用户列表。
	UserList []*User `json:"user_list"`
}

// Service 是用户模块的业务契约，每个方法对应一个业务用例。
//
// 返回约定：
//   - 返回模型 User，不直接返回数据库实体 UserEntity。
//   - 只需要返回模型时直接返回模型。
//   - 需要附带总数、分页等额外信息时，使用 XXXResponse 包装。
//
// 入参约定：
//   - 请求参数统一命名为 XXXRequest，使用指针传递。
//   - 每个 Request 必须实现 Validate 方法。
//   - 实现方在用例入口处调用 Validate，校验失败直接返回错误。
//
// 结构约定：
//   - 请求参数和响应结构统一定义在本文件。
//   - 请求和响应结构可以引用 model.go 中定义的模型。
//
// 实现由 internal/service 提供，并在应用组装层注入。
type Service interface {
	// List 返回符合条件的用户列表。
	// Total 表示符合查询条件的总数，不是当前页的记录数。
	List(ctx context.Context, req *ListRequest) (*ListResponse, error)

	// Get 根据用户标识获取用户。
	// 用户不存在时返回 ErrNotFound。
	Get(ctx context.Context, req *GetRequest) (*User, error)

	// Create 创建用户。
	// 邮箱已存在时返回 AlreadyExists 错误。
	Create(ctx context.Context, req *CreateRequest) (*User, error)

	// Update 更新用户信息。
	// 用户不存在时返回 ErrNotFound，邮箱已被占用时返回 AlreadyExists 错误。
	Update(ctx context.Context, req *UpdateRequest) (*User, error)
}
