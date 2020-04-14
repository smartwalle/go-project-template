package mysql

import (
	"github.com/smartwalle/dbs"
	"go-project-template/user"
	"go-project-template/user/service"
)

const (
	kTblUser = "users"
)

type userRepository struct {
	db dbs.DB
}

func NewUserRepository(db dbs.DB) service.UserRepository {
	return &userRepository{db: db}
}

func (this *userRepository) BeginTx() (dbs.TX, service.UserRepository) {
	var tx = dbs.MustTx(this.db)
	var repo = *this
	repo.db = tx
	return tx, &repo
}

func (this *userRepository) WithTx(tx dbs.TX) service.UserRepository {
	var repo = *this
	repo.db = tx
	return &repo
}

func (this *userRepository) GetUserWithId(id int64) (result *user.User, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("u.id", "u.username", "u.last_name", "u.first_name")
	sb.From(kTblUser, "AS u")
	sb.Where("u.id = ?", id)
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, err
}

func (this *userRepository) GetUserWithUsername(username string) (result *user.User, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("u.id", "u.username", "u.last_name", "u.first_name")
	sb.From(kTblUser, "AS u")
	sb.Where("u.username = ?", username)
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, err
}

func (this *userRepository) AddUser(user *user.AddUserParam) (result int64, err error) {
	var ib = dbs.NewInsertBuilder()
	ib.Table(kTblUser)
	ib.Columns("username", "last_name", "first_name")
	ib.Values(user.Username, user.LastName, user.FirstName)
	sResult, err := ib.Exec(this.db)
	if err != nil {
		return 0, err
	}

	nId, err := sResult.LastInsertId()
	if err != nil {
		return 0, err
	}

	return nId, err
}
