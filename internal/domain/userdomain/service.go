package userdomain

import "context"

// CreateRequest 是"创建用户"用例的输入契约。
type CreateRequest struct {
	// Name 是用户名称。
	Name string `json:"name"`

	// Email 是用户邮箱。
	Email string `json:"email"`
}

// Service 是用户模块的业务契约，每个方法是一个业务用例。
//
// 用例返回对外模型 User，不返回实体 UserEntity：
// Entity 对应数据库结构，不直接暴露给调用方，由实现层转换成模型再返回。
//
// 实现由 internal/service 提供，在应用组装层注入。
type Service interface {
	// List 返回全部用户。
	List(ctx context.Context) ([]*User, error)

	// Get 按 Id 返回用户；用户不存在时返回 ErrNotFound。
	Get(ctx context.Context, id int64) (*User, error)

	// Create 创建用户；邮箱已存在时返回 ErrAlreadyExists。
	Create(ctx context.Context, req CreateRequest) (*User, error)
}
