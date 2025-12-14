package domain

import (
	"github.com/smartwalle/dbs"
)

type User struct {
	ID   int64  `sql:"id,auto_increment"`
	Name string `sql:"name"`
}

func (User) TableName() string {
	return "user"
}

func (User) PrimaryKey() string {
	return "id"
}

type UserRepository interface {
	dbs.Repositoy[User]
}
