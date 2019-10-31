package cqlx

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
)

type session struct {
	*gocql.Session
}

func (s *session) Query(qry interface{}, args ...interface{}) Queryx {
	return queryx(s.Session, qry, args...) // clone a new session
}

func (s *session) Exec(stmt interface{}) error {
	return s.Query(stmt, nil).Exec()
}

func (s *session) Close() error {
	s.Session.Close()
	return nil
}

type query struct {
	*gocqlx.Queryx
	typ QueryxType
}

func (q *query) Exec() error {
	return executex(q, nil)
}

func (q *query) Put(newitem interface{}) error {
	return executex(q, newitem)
}

func (q *query) Get(res interface{}) error {
	return executex(q, res)
}

func (q *query) Iter() Iterx {
	return iterx(q.Queryx)
}
