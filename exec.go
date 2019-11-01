package cqlx

import (
	"log"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

type QueryxType int

const (
	Select = QueryxType(iota + 1)
	Insert
	Update
	Delete
	Batch
	Raw
)

//// Executex
func executex(q *Queryx, item interface{}) error {
	if q == nil {
		return ErrInvalidQuery
	}
	if item == nil {
		return q.ExecRelease()
	}
	switch q.typ {
	case Select:
		if isSlice(item) {
			return q.SelectRelease(item)
		}
		return q.GetRelease(item)
	case Insert:
		return q.BindStruct(item).ExecRelease()
	case Update:
		return q.BindStruct(item).ExecRelease()
	case Delete:
		return q.ExecRelease()
	case Batch:
		return q.ExecRelease()
	case Raw:
		return q.GetRelease(item)
	default:
		log.Printf("CQLX: Unexpected query type -- %T", q.typ)
		return ErrInvalidQueryType
	}
}

//// Queryx
func queryx(sess *gocql.Session, qry interface{}, args ...interface{}) *Queryx {
	var stmt string
	var names []string
	switch q := qry.(type) {
	case *qb.SelectBuilder:
		stmt, names = q.ToCql()
	case *qb.InsertBuilder:
		stmt, names = q.ToCql()
	case *qb.UpdateBuilder:
		stmt, names = q.ToCql()
	case *qb.DeleteBuilder:
		stmt, names = q.ToCql()
	case *qb.BatchBuilder:
		stmt, names = q.ToCql()
	case string:
		stmt, names = q, nil
	default:
		return &Queryx{nil, 0}
	}
	if isMap(args...) {
		return &Queryx{gocqlx.Query(sess.Query(stmt, args...), names).BindMap(args[0].(qb.M)), queryxType(qry)}
	}
	return &Queryx{gocqlx.Query(sess.Query(stmt, args...), names), queryxType(qry)}
}

//// Iterx
func iterx(qry *Queryx) *Iterx {
	if qry.Queryx == nil {
		return &Iterx{}
	}
	return &Iterx{qry.Queryx.Iter()}
}

//// QueryxType
func queryxType(qry interface{}) QueryxType {
	switch qry.(type) {
	case *qb.SelectBuilder:
		return Select
	case *qb.InsertBuilder:
		return Insert
	case *qb.UpdateBuilder:
		return Update
	case *qb.DeleteBuilder:
		return Delete
	case *qb.BatchBuilder:
		return Batch
	case string:
		return Raw
	default:
		return 0
	}
}
