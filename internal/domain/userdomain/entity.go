// Package userdomain 定义用户领域的实体、业务契约与领域错误。
//
// 每个业务模块一个包，契约都定义在这里。
package userdomain

import "time"

// UserEntity 是 users 表对应的实体。
//
// 命名规则：每张数据库表对应一个 Entity 结构体，名称是表名的驼峰形式加 Entity 后缀，
// 例如 users 表对应 UserEntity。
//
// 标识符规则：实体的标识符字段类型统一使用 uint64。
//
// sql tag 对应 users 表的列名，json tag 对应 JSON 序列化时的字段名；
// TableName 与 PrimaryKey 让它满足 dbs.Entity。
type UserEntity struct {
	// Id 是用户的业务标识，由数据库生成（default：留空时用数据库的默认值）。
	Id uint64 `sql:"id;default" json:"id"`

	// Name 是用户名称。
	Name string `sql:"name" json:"name"`

	// Email 是用户邮箱，在业务上唯一。
	Email string `sql:"email" json:"email"`

	// CreatedAt 是用户创建时间。
	CreatedAt time.Time `sql:"created_at" json:"created_at"`
}

// TableName 是 users 表名，供 dbs 的通用 Repository 使用。
func (UserEntity) TableName() string {
	return "users"
}

// PrimaryKey 是主键列名，供 dbs 的通用 Repository 使用。
func (UserEntity) PrimaryKey() string {
	return "id"
}
