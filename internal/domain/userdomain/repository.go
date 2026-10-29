package userdomain

import "github.com/smartwalle/dbs"

// Repository 是 users 表对应的仓储接口。
//
// 约定：每张数据库表对应一个仓储接口，放在表所属的领域包内；
// 接口由使用方（Service）定义，实现见 internal/repository。
//
// 查询语义：查不到数据时返回 nil，不返回空记录；命中的记录主键大于 0。
type Repository interface {
	dbs.Repository[UserEntity]
}
