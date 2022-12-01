package model

// User 数据实体
type User struct {
	Id        int64  `sql:"id"`         // id
	Username  string `sql:"username"`   // 用户名
	LastName  string `sql:"last_name"`  // 姓
	FirstName string `sql:"first_name"` // 名
}
