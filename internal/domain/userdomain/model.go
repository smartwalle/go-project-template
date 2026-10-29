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
