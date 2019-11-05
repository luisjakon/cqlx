package cqlx

import (
	"github.com/scylladb/gocqlx/qb"
)

type Crud struct {
	SelectQuery interface{}
	InsertQuery interface{}
	UpdateQuery interface{}
	DeleteQuery interface{}
}

func (c Crud) Query(sess *Sessionx, stmt interface{}, args ...interface{}) *Queryx {
	return sess.Queryx(stmt, args...)
}

func (c Crud) Select(sess *Sessionx, args ...interface{}) *Queryx {
	switch s := c.SelectQuery.(type) {
	case *qb.SelectBuilder:
		return sess.Queryx(s, args...)
	case string:
		return sess.Queryx(s, args...)
	}
	return _NilQuery
}

func (c Crud) Insert(sess *Sessionx, newitem interface{}) error {
	switch s := c.InsertQuery.(type) {
	case *qb.InsertBuilder:
		return sess.Queryx(s).Put(newitem)
	case string:
		return sess.Queryx(s).Put(newitem)
	}
	return ErrInvalidQuery
}

func (c Crud) Update(sess *Sessionx, item interface{}) error {
	switch s := c.UpdateQuery.(type) {
	case *qb.UpdateBuilder:
		return sess.Queryx(s.Existing()).Put(item)
	case string:
		return sess.Queryx(s).Put(item)
	}
	return ErrInvalidQuery
}

func (c Crud) Delete(sess *Sessionx, item interface{}) error {
	switch s := c.DeleteQuery.(type) {
	case *qb.DeleteBuilder:
		return sess.Queryx(s, item).Exec()
	case string:
		return sess.Queryx(s, item).Exec()
	}
	return ErrInvalidQuery
}
