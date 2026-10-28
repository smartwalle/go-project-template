package userdomain

import "time"

// User 是用户的对外模型，用于接口返回。
//
// Entity 是数据库结构的映射，不适合直接对外输出：
// service 层取到 Entity 后转换成 User，再交给调用方（HTTP 层）返回。
type User struct {
	// Id 是用户的业务标识。
	Id int64 `json:"id"`

	// Name 是用户名称。
	Name string `json:"name"`

	// Email 是用户邮箱，在业务上唯一。
	Email string `json:"email"`

	// CreatedAt 是用户创建时间。
	CreatedAt time.Time `json:"created_at"`
}
