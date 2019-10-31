package cqlx

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
)

type Session struct {
	*gocql.Session
}

func (s *Session) Query(qry interface{}, args ...interface{}) Queryx {
	return queryx(s.Session, qry, args...) // clone a new session
}

func (s *Session) Exec(stmt interface{}) error {
	return s.Query(stmt, nil).Exec()
}

func (s *Session) Close() error {
	s.Session.Close()
	return nil
}

type Query struct {
	*gocqlx.Queryx
	typ QueryxType
}

func (q *Query) Exec() error {
	return executex(q, nil)
}

func (q *Query) Put(newitem interface{}) error {
	return executex(q, newitem)
}

func (q *Query) Get(res interface{}) error {
	return executex(q, res)
}

func (q *Query) Iter() Iterx {
	return iterx(q.Queryx)
}
