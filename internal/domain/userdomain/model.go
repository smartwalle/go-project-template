package userdomain

import "time"

// User 是用户模块的模型。
//
// Model 是用来整理、组合数据的结构体：
//
// 标识符规则：模型里的标识符字段类型与实体保持一致，统一使用 uint64。
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
