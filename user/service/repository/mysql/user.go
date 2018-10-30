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

func (this *userRepository) GetUserWithId(id int) (result *user.User, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("u.id", "u.username", "u.last_name", "u.first_name")
	sb.From(k_DB_USER, "AS u")
	sb.Where("u.id = ?", id)
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, err
}

func (this *userRepository) GetUserWithUsername(username string) (result *user.User, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("u.id", "u.username", "u.last_name", "u.first_name")
	sb.From(k_DB_USER, "AS u")
	sb.Where("u.username = ?", username)
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, err
}

func (this *userRepository) AddUser(user *user.AddUserParam) (result *user.User, err error) {
	var tx = dbs.MustTx(this.db)

	var ib = dbs.NewInsertBuilder()
	ib.Table(k_DB_USER)
	ib.Columns("username", "last_name", "first_name")
	ib.Values(user.Username, user.LastName, user.FirstName)
	sResult, err := ib.Exec(tx)
	if err != nil {
		return nil, err
	}

	nId, err := sResult.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var sb = dbs.NewSelectBuilder()
	sb.Selects("u.id", "u.username", "u.last_name", "u.first_name")
	sb.From(k_DB_USER, "AS u")
	sb.Where("u.id = ?", nId)
	if err = sb.Scan(tx, &result); err != nil {
		return nil, err
	}

	tx.Commit()

	return result, err
}
