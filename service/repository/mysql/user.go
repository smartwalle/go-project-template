package mysql

import (
	"github.com/smartwalle/dbs"
	"go-project-template/model"
	"go-project-template/service"
)

const (
	kTblUser = "users"
)

type userRepository struct {
	db *DB
}

func NewUserRepository(db dbs.Database) service.UserRepository {
	return &userRepository{db: NewDB(db)}
}

func (this *userRepository) BeginTx() (*dbs.Tx, service.UserRepository) {
	var db, tx, err = this.db.Begin()
	if err != nil {
		panic(err)
		return nil, nil
	}
	var repo = *this
	repo.db = db
	return tx, &repo
}

func (this *userRepository) WithTx(tx *dbs.Tx) service.UserRepository {
	var repo = *this
	repo.db = repo.db.Clone(tx)
	return &repo
}

func (this *userRepository) GetUserWithId(id int64) (result *model.User, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("u.id", "u.username", "u.last_name", "u.first_name")
	sb.From(kTblUser, "AS u")
	sb.Where("u.id = ?", id)
	if err = sb.Scan(this.db, &result); err != nil && err != dbs.ErrNoRows {
		return nil, err
	}
	return result, err
}

func (this *userRepository) GetUserWithUsername(username string) (result *model.User, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("u.id", "u.username", "u.last_name", "u.first_name")
	sb.From(kTblUser, "AS u")
	sb.Where("u.username = ?", username)
	if err = sb.Scan(this.db, &result); err != nil && err != dbs.ErrNoRows {
		return nil, err
	}
	return result, err
}

func (this *userRepository) AddUser(opts service.AddUserOptions) (result int64, err error) {
	var ib = dbs.NewInsertBuilder()
	ib.Table(kTblUser)
	ib.Columns("username", "last_name", "first_name")
	ib.Values(opts.Username, opts.LastName, opts.FirstName)
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
