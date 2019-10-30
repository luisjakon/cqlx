package cqlx

import (
	"github.com/scylladb/gocqlx"
)

type iter struct {
	*gocqlx.Iterx
}

func (i iter) Next(dest interface{}) bool {
	return i.Iterx.StructScan(dest)
}

func (i iter) Close() error {
	return i.Iterx.Close()
}
