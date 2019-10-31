package cqlx

import (
	"errors"
	"reflect"

	"github.com/scylladb/gocqlx/qb"
)

var (
	ErrInvalidCluster   = errors.New("Invalid Cluster.")
	ErrInvalidSession   = errors.New("Invalid Session.")
	ErrInvalidQueryType = errors.New("Invalid Query Type.")
	ErrInvalidQuery     = errors.New("Invalid Query.")
	ErrNilIter          = errors.New("Invalid Iterator.")
)

type Closeable interface {
	Close() error
}

type DB interface {
	Open(dbkeyspace string, dbhosts ...string) error
	View(func(Tx) error) error
	Update(func(Tx) error) error
	Session() Sessionx
}

type Sessionx interface {
	Query(query interface{}, args ...interface{}) Queryx
	Exec(stmnt interface{}) error
	Closeable
}

type Queryx interface {
	Get(res interface{}) error
	Put(item interface{}) error
	Exec() error
	Iter() Iterx
}

type Iterx interface {
	Next(interface{}) bool
	Closeable
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
