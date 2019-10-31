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
	s := db.Session()
	defer s.Close()
	return viewtx(s, fn)
}

func (db *DB) Update(fn func(Tx) error) error {
	s := db.Session()
	defer s.Close()
	return updatetx(s, fn)
}

func (db *DB) Session() *Sessionx {
	if db.ClusterConfig == nil {
		return _NilSession
	}
	sess, err := db.CreateSession()
	if err != nil {
		log.Printf(err.Error())
		return _NilSession
	}
	return &Sessionx{sess}
}

func (db *DB) Close() error {
	return nil
}
