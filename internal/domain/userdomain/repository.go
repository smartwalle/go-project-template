package userdomain

import "github.com/smartwalle/dbs"

// Repository 描述用户模块的数据访问能力。
//
// 接口由使用方（Service）定义，实现见 internal/repository。
type Repository interface {
	dbs.Repository[User]
}
