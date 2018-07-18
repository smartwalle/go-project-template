package mysql

import (
	"github.com/smartwalle/dbs"
	"go-projet-template/user"
	"go-projet-template/user/service"
)

const (
	k_DB_USER = "user_user"
)

type UserRepository struct {
	DB dbs.DB
}

func NewUserRepository(db dbs.DB) service.IUserRepository {
	return &UserRepository{DB: db}
}

func (this *UserRepository) User(id int) (result *user.User, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("u.id", "u.username", "u.last_name", "u.first_name")
	sb.From(k_DB_USER, "AS u")
	sb.Where("u.id = ?", id)
	if err = sb.Scan(this.DB, &result); err != nil {
		return nil, err
	}
	return result, err
}
