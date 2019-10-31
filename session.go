package cqlx

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
)

type Sessionx struct {
	*gocql.Session
}

func (s *Sessionx) Queryx(qry interface{}, args ...interface{}) *Queryx {
	return queryx(s.Session, qry, args...)
}

func (s *Sessionx) Exec(stmt interface{}, args ...interface{}) error {
	return s.Queryx(stmt, args...).Exec()
}

func (s *Sessionx) Close() error {
	s.Session.Close()
	return nil
}

type Queryx struct {
	*gocqlx.Queryx
	typ QueryxType
}

func (q *Queryx) Get(res interface{}) error {
	return executex(q, res)
}

func (q *Queryx) Put(newitem interface{}) error {
	return executex(q, newitem)
}

func (q *Queryx) Exec() error {
	return executex(q, nil)
}

func (q *Queryx) Iter() *Iterx {
	return iterx(q)
}
