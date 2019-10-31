package cqlx_test

import (
	"testing"

	"github.com/scylladb/gocqlx/qb"
	"github.com/stretchr/testify/assert"
)

func TestSession(t *testing.T) {
	sess := db.Session()
	defer sess.Close()

	res := &kv{}
	err := sess.Query(qb.Select("kv")).Get(res)

	assert.Equal(t, nil, err)
}

func TestSessionRawQuery(t *testing.T) {
	sess := db.Session()
	defer sess.Close()

	err := sess.Query("INSERT INTO kv (key,value) VALUES ('5','val5');").Exec()
	assert.Equal(t, nil, err)

	res := kv{}
	err = sess.Query(qb.Select("kv").Where(qb.Eq("key")), "5").Get(&res)
	assert.Equal(t, nil, err)
	assert.Equal(t, res.Key, "5")
	assert.Equal(t, res.Value, "val5")

	err = sess.Query("DELETE FROM kv WHERE key = '5';").Exec()
	assert.Equal(t, nil, err)

	res = kv{}
	err = sess.Query(qb.Select("kv").Where(qb.Eq("key")), "5").Get(&res)
	assert.NotEqual(t, nil, err)
}
