package cqlx_test

import (
	"testing"

	"github.com/scylladb/gocqlx/qb"
	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	var res kv

	sess, err := db.Session()
	assert.Equal(t, nil, err)
	defer sess.Close()

	err = sess.Queryx("SELECT * FROM kv;").Get(&res)
	assert.Equal(t, nil, err)
}

func TestRawCQLQuery(t *testing.T) {
	var res kv

	sess, err := db.Session()
	assert.Equal(t, nil, err)
	defer sess.Close()

	err = sess.Queryx("SELECT * FROM kv WHERE key='5';").Get(&res)
	assert.Equal(t, "not found", err.Error())

	err = sess.Queryx("INSERT INTO kv (key,value) VALUES ('5','val5');").Exec()
	assert.Equal(t, nil, err)

	err = sess.Queryx("SELECT * FROM kv WHERE key='5';").Get(&res)
	assert.Equal(t, nil, err)
	assert.Equal(t, "5", res.Key)
	assert.Equal(t, "val5", res.Value)
	assert.Equal(t, kv{"5", "val5"}, res)

	err = sess.Queryx("DELETE FROM kv WHERE key='5';").Exec()
	assert.Equal(t, nil, err)

	res = kv{}
	err = sess.Queryx("SELECT * FROM kv WHERE key='5';").Get(&res)
	assert.NotEqual(t, nil, err)
}

func TestGoCqlxQueryBuilder(t *testing.T) {
	var res kv

	sess, err := db.Session()
	assert.Equal(t, nil, err)
	defer sess.Close()

	err = sess.Queryx(qb.Select("kv").Where(qb.Eq("key")), "6").Get(&res)
	assert.Equal(t, "not found", err.Error())

	err = sess.Queryx(qb.Insert("kv").Columns("key", "value")).Put(&kv{"6", "val6"})
	assert.Equal(t, nil, err)

	err = sess.Queryx(qb.Select("kv").Where(qb.Eq("key")), "6").Get(&res)
	assert.Equal(t, nil, err)
	assert.Equal(t, "6", res.Key)
	assert.Equal(t, "val6", res.Value)
	assert.Equal(t, kv{"6", "val6"}, res)

	err = sess.Queryx(qb.Delete("kv").Where(qb.Eq("key")), "6").Exec()
	assert.Equal(t, nil, err)

	res = kv{}
	err = sess.Queryx(qb.Select("kv").Where(qb.Eq("key")), "6").Get(&res)
	assert.NotEqual(t, nil, err)
}
