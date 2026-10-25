// Package repository 实现数据访问。
//
// 每个业务模块一个文件，底层是 dbs 的通用 Repository。
package repository

import (
	"github.com/smartwalle/dbs"

	"project-template/internal/domain/userdomain"
)

// UserRepository 是用户仓储的实现。
//
// 错误原样返回，不做包装、不做转换。
type UserRepository struct {
	dbs.Repository[userdomain.User]
}

// NewUserRepository 创建仓储。db 由应用组装层统一持有和关闭。
func NewUserRepository(db dbs.Database) *UserRepository {
	return &UserRepository{dbs.NewRepository[userdomain.User](db)}
}
