package service

// AddUserOption service 或者 repository 方法参数过多时，可以考虑使用结构体组织
type AddUserOption struct {
	Username  string // 用户名
	LastName  string // 姓
	FirstName string // 名
}
