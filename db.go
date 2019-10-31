package cqlx

import (
	"log"

	"github.com/gocql/gocql"
)

type DB struct {
	*gocql.ClusterConfig
}

func (db *DB) Open(dbkeyspace string, dbhosts ...string) (err error) {
	db.ClusterConfig = gocql.NewCluster(dbhosts...)
	db.ClusterConfig.Keyspace = dbkeyspace
	return
}

func (db *DB) View(fn func(Tx) error) error {
	s, err := db.Session()
	if err != nil {
		return err
	}
	return viewtx(s, fn)
}

func (db *DB) Update(fn func(Tx) error) error {
	s, err := db.Session()
	if err != nil {
		return err
	}
	return updatetx(s, fn)
}

func (db *DB) Session() (*Sessionx, error) {
	if db.ClusterConfig == nil {
		return nil, ErrInvalidCluster
	}
	sess, err := db.CreateSession()
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	return &Sessionx{sess}, nil
}

func (db *DB) Close() error {
	return nil
}
