package cqlx

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
)

func NewDBWithConfig(c *gocql.ClusterConfig) *DB {
	return &DB{c}
}

func NewDB(dbkeyspace string, dbhosts ...string) *DB {
	db := &DB{}
	db.Open(dbkeyspace, dbhosts...)
	return db
}

func NewSession(s *gocql.Session) *Session {
	return &Session{s}
}

func NewQueryx(q *gocqlx.Queryx, typ QueryxType) *Query {
	return &Query{q, typ}
}

func NewIter(it *gocqlx.Iterx) *Iter {
	return &Iter{it}
}
