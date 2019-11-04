package cqlx

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
)

func newDBWithConfig(c *gocql.ClusterConfig) *DB {
	return &DB{c}
}

func newSession(s *gocql.Session) *Sessionx {
	return &Sessionx{s}
}

func newQueryx(q *gocqlx.Queryx, typ QueryxType) *Queryx {
	return &Queryx{q, typ}
}

func newIter(it *gocqlx.Iterx) *Iterx {
	return &Iterx{it}
}
