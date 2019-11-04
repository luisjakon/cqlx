package cqlx

import (
	"github.com/gocql/gocql"
)

type DB struct {
	*gocql.ClusterConfig
}

func (db *DB) Session() (*Sessionx, error) {
	if db.ClusterConfig == nil {
		return nil, ErrInvalidCluster
	}
	sess, err := db.CreateSession()
	return &Sessionx{sess}, err
}

func (db *DB) View(fn func(Tx) error) error {
	return viewtx(db, fn)
}

func (db *DB) Update(fn func(Tx) error) error {
	return updatetx(db, fn)
}

func (db *DB) Close() error {
	return nil
}
