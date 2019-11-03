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
	ErrNoHosts          = errors.New("Invalid Host(s) List.")
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

func isMap(arg interface{}) bool {
	return reflect.Map == reflect.Indirect(reflect.ValueOf(arg)).Kind()
}

func isStruct(arg interface{}) bool {
	return reflect.Struct == reflect.Indirect(reflect.ValueOf(arg)).Kind()
}

func asMap(args ...interface{}) qb.M {
	if len(args) == 1 {
		switch m := args[0].(type) {
		case qb.M:
			return m
		case *qb.M:
			return *m
		case map[string]interface{}:
			return qb.M(m)
		case *map[string]interface{}:
			return qb.M(*m)
		}
	}
	return nil
}

func asStruct(args ...interface{}) interface{} {
	if len(args) == 1 && isStruct(args[0]) {
		return args[0]
	}
	return nil
}
