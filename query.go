package cqlx

import (
	"github.com/scylladb/gocqlx"
)

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
