package mysql

import (
	"context"
	"database/sql"
	"github.com/smartwalle/dbs"
)

type DB struct {
	dbs.Session
	db dbs.Database
}

func NewDB(db dbs.Database) *DB {
	return &DB{db: db, Session: db}
}

func (this *DB) Begin() (*DB, *dbs.Tx, error) {
	return this.BeginTx(context.Background(), nil)
}

func (this *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*DB, *dbs.Tx, error) {
	var tx, err = this.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, nil, err
	}
	var db = *this
	db.Session = tx
	return &db, tx, nil
}

func (this *DB) Clone(session dbs.Session) *DB {
	var db = *this
	db.Session = session
	return &db
}
