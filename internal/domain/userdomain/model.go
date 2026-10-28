package userdomain

import "time"

// User 是用户模块的对外模型。
//
// Model 是 Service 对外提供的结构，承载用例的输出：
// Entity 对应数据库结构，不直接对外输出，由 Service 负责转换成 Model。
// Model 怎么传输、怎么序列化，由使用它的协议层决定。
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
