package cqlx

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
)

type Session struct {
	*gocql.Session
}

func (s *Session) Queryx(qry interface{}, args ...interface{}) *Queryx {
	return queryx(s.Session, qry, args...)
}

func (s *Session) Exec(stmt interface{}) error {
	return s.Queryx(stmt, nil).Exec()
}

func (s *Session) Close() error {
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
