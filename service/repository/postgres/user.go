package postgres

import (
	"context"
	"errors"
	"github.com/smartwalle/dbs"
	"go-project-template/model"
	"go-project-template/pkg"
	"go-project-template/service"
)

const (
	kTblUser = "users"
)

type userRepository struct {
	db *DB
}

func NewUserRepository(db *dbs.DB) service.UserRepository {
	var sb = dbs.NewSelectBuilder()
	sb.UsePlaceholder(dbs.DollarPlaceholder)
	sb.Selects("u.id", "u.username", "u.last_name", "u.first_name")
	sb.From(kTblUser, "AS u")
	sb.Where("u.id = ?")
	pkg.PrepareStatement(context.Background(), db, "get_user_with_id", sb)

	sb = dbs.NewSelectBuilder()
	sb.UsePlaceholder(dbs.DollarPlaceholder)
	sb.Selects("u.id", "u.username", "u.last_name", "u.first_name")
	sb.From(kTblUser, "AS u")
	sb.Where("u.username = ?")
	pkg.PrepareStatement(context.Background(), db, "get_user_with_username", sb)

	return &userRepository{db: NewDB(db)}
}

func (repo *userRepository) BeginTx(ctx context.Context) (*dbs.Tx, service.UserRepository, error) {
	var db, tx, err = repo.db.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	var nRepo = *repo
	nRepo.db = db
	return tx, &nRepo, nil
}

func (repo *userRepository) WithTx(tx *dbs.Tx) service.UserRepository {
	var nRepo = *repo
	nRepo.db = nRepo.db.Clone(tx)
	return &nRepo
}

func (repo *userRepository) GetUserWithId(ctx context.Context, id int64) (result *model.User, err error) {
	result, err = dbs.Query[*model.User](ctx, repo.db, "get_user_with_id", id)
	if err != nil && !errors.Is(err, dbs.ErrNoRows) {
		return nil, err
	}
	return result, err
}

func (repo *userRepository) GetUserWithUsername(ctx context.Context, username string) (result *model.User, err error) {
	result, err = dbs.Query[*model.User](ctx, repo.db, "get_user_with_username", username)
	if err != nil && !errors.Is(err, dbs.ErrNoRows) {
		return nil, err
	}
	return result, nil
}

func (repo *userRepository) AddUser(ctx context.Context, opts service.AddUserOptions) (result int64, err error) {
	var ib = dbs.NewInsertBuilder()
	ib.UsePlaceholder(dbs.DollarPlaceholder)
	ib.Table(kTblUser)
	ib.Columns("username", "last_name", "first_name")
	ib.Values(opts.Username, opts.LastName, opts.FirstName)
	ib.Returning("id", "username")

	var aux model.User
	if err = ib.ScanContext(ctx, repo.db, &aux); err != nil {
		return 0, err
	}

	return aux.Id, err
}
