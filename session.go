package cqlx

import (
	"github.com/gocql/gocql"
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
	if s.Session != nil {
		s.Session.Close()
	}
	return nil
}
