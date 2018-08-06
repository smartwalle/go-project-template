package mysql

import (
	"github.com/smartwalle/dbs"
	"go-project-template/user"
	"go-project-template/user/service"
)

const (
	k_DB_USER = "user_user"
)

type userRepository struct {
	db dbs.DB
}

func NewUserRepository(db dbs.DB) service.UserRepository {
	return &userRepository{db: db}
}

func (this *userRepository) User(id int) (result *user.User, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("u.id", "u.username", "u.last_name", "u.first_name")
	sb.From(k_DB_USER, "AS u")
	sb.Where("u.id = ?", id)
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, err
}
