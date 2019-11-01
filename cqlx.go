package cqlx

import (
	"errors"
	"reflect"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

var (
	ErrInvalidCluster   = errors.New("Invalid Cluster Config.")
	ErrInvalidSession   = errors.New("Invalid Session.")
	ErrInvalidQueryType = errors.New("Invalid Query Type.")
	ErrInvalidQuery     = errors.New("Invalid Query.")
	ErrNilIter          = errors.New("Invalid Iterator.")
)

func OpenWithConfig(c *gocql.ClusterConfig) *DB {
	return newDBWithConfig(c)
}

func Open(dbkeyspace string, dbhosts ...string) *DB {
	return newDB(dbkeyspace, dbhosts...)
}

func Session(s *gocql.Session) *Sessionx {
	return newSession(s)
}

func Query(qry *gocqlx.Queryx, typ QueryxType) *Queryx {
	return newQueryx(qry, typ)
}

func Iter(it *gocqlx.Iterx) *Iterx {
	return newIter(it)
}

func isSlice(res interface{}) bool {
	switch reflect.Indirect(reflect.ValueOf(res)).Kind() {
	case reflect.Slice, reflect.Array:
		return true
	}
	return false
}

func isMap(args ...interface{}) bool {
	if len(args) != 1 {
		return false
	}
	_, ok := args[0].(qb.M)
	return ok
}
