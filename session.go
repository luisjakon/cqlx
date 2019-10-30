package cqlx

import (
	"github.com/gocql/gocql"
)

type session struct {
	*gocql.Session
}

func (s *session) Query(qry interface{}, args ...interface{}) Queryx {
	return &query{s.Session, qry, args} // clone a new session
}

func (s *session) Exec(stmt interface{}) error {
	return s.Query(stmt, nil).Exec()
}

func (s *session) Close() error {
	s.Session.Close()
	return nil
}

type query struct {
	sess  *gocql.Session
	query interface{}
	args  []interface{}
}

func (q *query) Exec() error {
	return execute(q.sess, nil, q.query, q.args...)
}

func (q *query) Put(newitem interface{}) error {
	return smart_put(q.sess, newitem, q.query, q.args...)
}

func (q *query) Get(res interface{}) error {
	return smart_get(q.sess, res, q.query, q.args...)
}

func (q *query) Iter() Iterx {
	return iterx(queryx(q.sess, q.query, q.args...))
}
