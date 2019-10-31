package cqlx

import (
	"log"

	"github.com/gocql/gocql"
)

type db struct {
	*gocql.ClusterConfig
}

func OpenWithConfig(c *gocql.ClusterConfig) DB {
	return &db{c}
}

func Open(dbkeyspace string, dbhosts ...string) DB {
	db := &db{}
	db.Open(dbkeyspace, dbhosts...)
	return db
}

func (db *db) Open(dbkeyspace string, dbhosts ...string) (err error) {
	db.ClusterConfig = gocql.NewCluster(dbhosts...)
	db.ClusterConfig.Keyspace = dbkeyspace
	return
}

func (db *db) View(fn func(Tx) error) error {
	s := db.Session()
	defer s.Close()
	return viewtx(s, fn)
}

func (db *db) Update(fn func(Tx) error) error {
	s := db.Session()
	defer s.Close()
	return updatetx(s, fn)
}

func (db *db) Session() Sessionx {
	if db.ClusterConfig == nil {
		return _NilSession
	}
	sess, err := db.CreateSession()
	if err != nil {
		log.Printf(err.Error())
		return _NilSession
	}
	return &session{sess}
}

func (db *db) Close() error {
	return nil
}
