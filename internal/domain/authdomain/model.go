package authdomain

type AuthUser struct {
	// Id 是用户的业务标识。
	Id int64 `json:"id"`

	// Name 是用户名称。
	Name string `json:"name"`
}
