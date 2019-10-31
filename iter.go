package cqlx

import (
	"github.com/scylladb/gocqlx"
)

type Iter struct {
	*gocqlx.Iterx
}

func (i Iter) Next(dest interface{}) bool {
	return i.Iterx.StructScan(dest)
}

func (i Iter) Close() error {
	return i.Iterx.Close()
}
