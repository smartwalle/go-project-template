package userdomain

import (
	"context"
	"fmt"
	"strings"

	"project-template/internal/foundation/errors"
	"project-template/internal/foundation/pagination"
)

// CreateRequest 是"创建用户"用例的请求参数。
type CreateRequest struct {
	// Name 是用户名称。
	Name string `json:"name"`

	// Email 是用户邮箱。
	Email string `json:"email"`
}

// Validate 校验请求参数：去掉字段首尾空白后校验，清洗结果写回字段。
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

// ListRequest 是"列出用户"用例的请求参数。
//
// 分页参数内嵌 pagination.Page。
type ListRequest struct {
	// Page 是分页参数：页码和每页条数。
	pagination.Page

	// Name 是用户名称的匹配条件，为空表示不过滤。
	Name string `json:"name"`

	// Email 是用户邮箱的匹配条件，为空表示不过滤。
	Email string `json:"email"`
}

// Validate 校验请求参数：分页参数归一化，过滤条件去掉首尾空白。
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
// 用例只需要返回模型本身时直接返回模型（Get、Create 返回 *User）；
// 需要带上模型之外的信息（总数、分页等）时用 Response 包装。
type ListResponse struct {
	// Total 是用户总数。
	Total int `json:"total"`

	// UserList 是用户列表。
	UserList []*User `json:"user_list"`
}

// GetRequest 是"查询用户"用例的请求参数。
type GetRequest struct {
	// Id 是要查询的用户标识。
	Id uint64 `json:"id"`
}

// Validate 校验请求参数：标识缺失返回 MissingArgument。
func (r *GetRequest) Validate() error {
	if r == nil {
		return errors.MissingArgument("request is required")
	}
	if r.Id <= 0 {
		return errors.MissingArgument("id is required")
	}
	return nil
}

// UpdateRequest 是"更新用户"用例的请求参数。
type UpdateRequest struct {
	// Id 是要更新的用户标识。
	Id uint64 `json:"id"`

	// Name 是新的用户名称。
	Name string `json:"name"`

	// Email 是新的用户邮箱。
	Email string `json:"email"`
}

// Validate 校验请求参数：去掉字段首尾空白后校验，清洗结果写回字段。
//
// 参数缺失返回 MissingArgument，格式不合法返回 InvalidArgument。
func (r *UpdateRequest) Validate() error {
	if r == nil {
		return errors.MissingArgument("request is required")
	}

	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.TrimSpace(r.Email)

	if r.Id <= 0 {
		return errors.MissingArgument("id is required")
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

// Service 是用户模块的业务契约，每个方法是一个业务用例。
//
// 入参约定：请求参数统一命名为 XXXRequest（XXX 是用例名）、统一用指针传递，
// 和响应数据一起定义在本文件；每个 Request 必须实现 Validate，
// 实现方在用例入口处调用它做校验，校验不通过直接返回错误。
//
// 返回约定：用例返回模型 User，不建议直接返回实体 UserEntity：
// Entity 是表映射，直接返回等于把表结构带到调用方，通常由实现层转换成模型再返回。
// 只需要返回模型时直接返回模型，需要带上额外信息时用 XXXResponse 包装。
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
