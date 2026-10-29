package authdomain

// AuthUser 是当前登录用户的信息。
//
// 标识符规则：标识符字段类型统一使用 uint64。
type AuthUser struct {
	// Id 是用户的业务标识。
	Id uint64 `json:"id"`

	// Name 是用户名称。
	Name string `json:"name"`
}
