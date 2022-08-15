package service

// AddUserOptions service 或者 repository 方法参数过多时，可以考虑使用结构体组织
type AddUserOptions struct {
	Username  string // 用户名
	LastName  string // 姓
	FirstName string // 名
}
