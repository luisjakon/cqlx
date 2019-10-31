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

func TestSessionRawQuery(t *testing.T) {
	var res kv

	sess, err := db.Session()
	assert.Equal(t, nil, err)
	defer sess.Close()

	err = sess.Queryx("INSERT INTO kv (key,value) VALUES ('5','val5');").Exec()
	assert.Equal(t, nil, err)

	err = sess.Queryx(qb.Select("kv").Where(qb.Eq("key")), "5").Get(&res)
	assert.Equal(t, nil, err)
	assert.Equal(t, res.Key, "5")
	assert.Equal(t, res.Value, "val5")

	err = sess.Queryx("DELETE FROM kv WHERE key = '5';").Exec()
	assert.Equal(t, nil, err)

	res = kv{}
	err = sess.Queryx(qb.Select("kv").Where(qb.Eq("key")), "5").Get(&res)
	assert.NotEqual(t, nil, err)
}
