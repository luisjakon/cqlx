package cqlx_test

import (
	"testing"

	"github.com/scylladb/gocqlx/qb"
	"github.com/stretchr/testify/assert"
)

func TestSession(t *testing.T) {
	var res kv

	sess, err := db.Session()
	assert.Equal(t, nil, err)
	defer sess.Close()

	err = sess.Queryx(qb.Select("kv")).Get(&res)
	assert.Equal(t, nil, err)
}
