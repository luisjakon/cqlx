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
func executex(q *Query, item interface{}) error {
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
		return q.BindStruct(item).GetRelease(item)
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
func queryx(sess *gocql.Session, qry interface{}, args ...interface{}) Queryx {
	var stmt string
	var names []string
	switch query := qry.(type) {
	case *qb.SelectBuilder:
		stmt, names = query.ToCql()
	case *qb.InsertBuilder:
		stmt, names = query.ToCql()
	case *qb.UpdateBuilder:
		stmt, names = query.ToCql()
	case *qb.DeleteBuilder:
		stmt, names = query.ToCql()
	case *qb.BatchBuilder:
		stmt, names = query.ToCql()
	case string:
		stmt, names = query, nil
	default:
		return _NilQuery
	}
	if isMap(args...) {
		return &Query{gocqlx.Query(sess.Query(stmt, args...), names).BindMap(args[0].(qb.M)), queryxType(qry)}
	}
	return &Query{gocqlx.Query(sess.Query(stmt, args...), names), queryxType(qry)}
}

//// Iterx
func iterx(qry *gocqlx.Queryx) Iterx {
	if qry == nil {
		return _NilIter
	}
	return &Iter{qry.Iter()}
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
