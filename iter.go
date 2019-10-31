package cqlx

import (
	"github.com/scylladb/gocqlx"
)

type Iterx struct {
	*gocqlx.Iterx
}

func (i Iterx) Next(dest interface{}) bool {
	return i.Iterx.StructScan(dest)
}

func (i Iterx) Close() error {
	return i.Iterx.Close()
}
