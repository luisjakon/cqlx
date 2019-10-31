package cqlx

import (
	"github.com/scylladb/gocqlx"
)

type Iterx struct {
	*gocqlx.Iterx
}

func (i Iterx) Next(dest interface{}) bool {
	if i.Iterx == nil {
		return false
	}
	return i.Iterx.StructScan(dest)
}

func (i Iterx) Close() error {
	if i.Iterx == nil {
		return ErrNilIter
	}
	return i.Iterx.Close()
}
