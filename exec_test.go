package cqlx_test

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/stretchr/testify/assert"
)

func TestCompiledQueries(t *testing.T) {

	sess, err := db.Session()
	defer sess.Close()

	assert.Equal(t, nil, err)

	var val kv = kv{"100", "val100"}
	var res kv

	err = sess.Queryx(`INSERT INTO kv (key, value) VALUES (:key, :value)`, &val).Exec()
	assert.Equal(t, nil, err)

	err = sess.Queryx(`SELECT * FROM kv WHERE key=:key`, &val).Get(&res)
	assert.Equal(t, nil, err)

	err = sess.Queryx(`DELETE FROM kv WHERE key=:key`, &val).Exec()
	assert.Equal(t, nil, err)

	err = sess.Queryx(`SELECT * FROM kv WHERE key=:key`, &val).Get(&res)
	assert.Equal(t, gocql.ErrNotFound, err)

}
