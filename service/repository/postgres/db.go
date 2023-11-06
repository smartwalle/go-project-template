package postgres

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

func (db *DB) Begin(ctx context.Context) (*DB, *dbs.Tx, error) {
	return db.BeginTx(ctx, nil)
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*DB, *dbs.Tx, error) {
	var tx, err = db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, nil, err
	}
	var ndb = *db
	ndb.Session = tx
	return &ndb, tx, nil
}

func (db *DB) Clone(session dbs.Session) *DB {
	var ndb = *db
	ndb.Session = session
	return &ndb
}
