package cqlx

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

//// BASE EXECUTOR
func execute(sess *gocql.Session, item interface{}, qry interface{}, args ...interface{}) error {
	switch query := qry.(type) {
	case *qb.SelectBuilder:
		return smart_get(sess, item, query, args...)
	case *qb.InsertBuilder:
		return insert(sess, item, query)
	case *qb.UpdateBuilder:
		return update(sess, item, query, args...)
	case *qb.DeleteBuilder:
		return delete(sess, query, args...)
	case *qb.BatchBuilder:
		return batch(sess, query, args...)
	case string:
		return raw(sess, item, query, args...)
	default:
		log.Printf("CQLX: Unexpected query type -- %T", query)
		return ErrInvalidQueryType
	}
}

//// SMART EXECUTORS
func smart_get(sess *gocql.Session, res interface{}, qry interface{}, args ...interface{}) error {
	switch query := qry.(type) {
	case *qb.SelectBuilder:
		if len(args) == 1 {
			if arg, ok := args[0].(qb.M); ok {
				if isSlice(res) {
					return getMapped(sess, res, query, arg)
				}
				return getMapped(sess, res, query.Limit(1), arg)
			}
		}
		if isSlice(res) {
			return get(sess, res, query, args...)
		}
		return get(sess, res, query.Limit(1), args...)
	default:
		log.Printf("CQLX: Unexpected query type -- %T", query)
		return ErrInvalidQueryType
	}
}

func smart_put(sess *gocql.Session, newitem interface{}, qry interface{}, args ...interface{}) error {
	switch query := qry.(type) {
	case *qb.InsertBuilder:
		return insert(sess, newitem, query)
	case *qb.UpdateBuilder:
		return update(sess, newitem, query, args...)
	default:
		log.Printf("CQLX: Unexpected query type -- %T", query)
		return ErrInvalidQueryType
	}
}

//// DUMB EXECUTORS
func getMapped(sess *gocql.Session, res interface{}, query *qb.SelectBuilder, args qb.M) error {
	return queryx(sess, query, args).BindMap(args).SelectRelease(res)
}

func getStruct(sess *gocql.Session, newitem interface{}, query *qb.SelectBuilder) error {
	return queryx(sess, query).BindStruct(newitem).GetRelease(newitem)
}

func get(sess *gocql.Session, res interface{}, query *qb.SelectBuilder, args ...interface{}) error {
	return queryx(sess, query, args...).GetRelease(res)
}

func insert(sess *gocql.Session, newitem interface{}, query *qb.InsertBuilder) error {
	return queryx(sess, query).BindStruct(newitem).GetRelease(newitem)
}

func update(sess *gocql.Session, item interface{}, query *qb.UpdateBuilder, args ...interface{}) error {
	return queryx(sess, query, args...).BindStruct(item).ExecRelease()
}

func delete(sess *gocql.Session, query *qb.DeleteBuilder, args ...interface{}) error {
	return queryx(sess, query, args...).ExecRelease()
}

func batch(sess *gocql.Session, query *qb.BatchBuilder, args ...interface{}) error {
	return queryx(sess, query, args...).ExecRelease()
}

func raw(sess *gocql.Session, item interface{}, query string, args ...interface{}) error {
	if len(args) > 0 {
		sess.Query(fmt.Sprintf(query, args...)).Exec()
	}
	return sess.Query(query).Exec()
}

//// QUERYX + ITERX GENERATORS
func queryx(sess *gocql.Session, qry interface{}, args ...interface{}) *gocqlx.Queryx {
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
		return nil
	}
	return gocqlx.Query(sess.Query(stmt, args...), names)
}

func iterx(qry *gocqlx.Queryx) Iterx {
	if qry == nil {
		return _NilIter
	}
	return &iter{qry.Iter()}
}
