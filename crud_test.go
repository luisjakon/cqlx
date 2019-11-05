package cqlx_test

import (
	"testing"

	"github.com/luisjakon/cqlx"
	"github.com/scylladb/gocqlx/qb"
	"github.com/stretchr/testify/assert"
)

func TestRawCrud(t *testing.T) {

	sess, err := db.Session()
	defer sess.Close()

	assert.Equal(t, nil, err)

	var kvdb = cqlx.Crud{
		`SELECT * FROM kv WHERE key=:key`,
		`INSERT INTO kv (key, value) VALUES (:key, :value)`,
		`UPDATE kv SET value=:value WHERE key=:key IF EXISTS`,
		`DELETE FROM kv WHERE key=:key`,
	}

	val := &kv{"200", "val200"}
	var res kv

	err = kvdb.Update(sess, val)
	assert.Equal(t, nil, err)

	err = kvdb.Insert(sess, val)
	assert.Equal(t, nil, err)

	err = kvdb.Select(sess, val).Get(&res)
	assert.Equal(t, nil, err)

	assert.Equal(t, "200", res.Key)
	assert.Equal(t, "val200", res.Value)

	err = kvdb.Delete(sess, val)
	assert.Equal(t, nil, err)

}

func TestQbCrud(t *testing.T) {

	sess, err := db.Session()
	defer sess.Close()

	assert.Equal(t, nil, err)

	var kvdb = cqlx.Crud{
		SelectQuery: qb.Select("kv").Where(qb.Eq("key")),
		InsertQuery: qb.Insert("kv").Columns("key", "value"),
		UpdateQuery: qb.Update("kv").Set("value").Where(qb.Eq("key")).Existing(),
		DeleteQuery: qb.Delete("kv").Where(qb.Eq("key")),
	}

	val := &kv{"200", "val200"}
	var res kv

	err = kvdb.Update(sess, val)
	assert.Equal(t, nil, err)

	err = kvdb.Insert(sess, val)
	assert.Equal(t, nil, err)

	err = kvdb.Select(sess, val).Get(&res)
	assert.Equal(t, nil, err)

	assert.Equal(t, "200", res.Key)
	assert.Equal(t, "val200", res.Value)

	err = kvdb.Delete(sess, val)
	assert.Equal(t, nil, err)

}

func TestRawCrudDirect(t *testing.T) {

	sess, err := db.Session()
	defer sess.Close()

	assert.Equal(t, nil, err)

	var kvdb = cqlx.Crud{
		`SELECT * FROM kv WHERE key=:key`,
		`INSERT INTO kv (key, value) VALUES (:key, :value)`,
		`UPDATE kv SET value=:value WHERE key=:key IF EXISTS`,
		`DELETE FROM kv WHERE key=:key`,
	}

	var res kv

	err = kvdb.Delete(sess, "201")
	assert.Equal(t, nil, err)

	err = kvdb.Update(sess, &kv{"201", "val201"})
	assert.Equal(t, nil, err)

	err = kvdb.Insert(sess, &kv{"201", "val201"})
	assert.Equal(t, nil, err)

	err = kvdb.Select(sess, &kv{Key: "201"}).Get(&res)
	assert.Equal(t, nil, err)

	assert.Equal(t, "201", res.Key)
	assert.Equal(t, "val201", res.Value)

	err = kvdb.Delete(sess, "201")
	assert.Equal(t, nil, err)

}

func TestQbCrudDirect(t *testing.T) {

	sess, err := db.Session()
	defer sess.Close()

	assert.Equal(t, nil, err)

	var kvdb = cqlx.Crud{
		SelectQuery: qb.Select("kv").Where(qb.Eq("key")),
		InsertQuery: qb.Insert("kv").Columns("key", "value"),
		UpdateQuery: qb.Update("kv").Set("value").Where(qb.Eq("key")).Existing(),
		DeleteQuery: qb.Delete("kv").Where(qb.Eq("key")),
	}

	var res kv

	err = kvdb.Delete(sess, "201")
	assert.Equal(t, nil, err)

	err = kvdb.Update(sess, &kv{"201", "val201"})
	assert.Equal(t, nil, err)

	err = kvdb.Insert(sess, &kv{"201", "val201"})
	assert.Equal(t, nil, err)

	err = kvdb.Select(sess, &kv{Key: "201"}).Get(&res)
	assert.Equal(t, nil, err)

	assert.Equal(t, "201", res.Key)
	assert.Equal(t, "val201", res.Value)

	err = kvdb.Delete(sess, "201")
	assert.Equal(t, nil, err)

}

func TestCrudQueryRaw(t *testing.T) {

	sess, err := db.Session()
	defer sess.Close()

	assert.Equal(t, nil, err)

	var kvdb cqlx.Crud
	var res kv

	err = kvdb.Query(sess, `INSERT INTO kv (key, value) VALUES (:key, :value)`).Put(&kv{"201", "val201"})
	assert.Equal(t, nil, err)

	err = kvdb.Query(sess, `SELECT * FROM kv WHERE key=:key`, &kv{Key: "201"}).Get(&res)
	assert.Equal(t, nil, err)

	assert.Equal(t, "201", res.Key)
	assert.Equal(t, "val201", res.Value)

	err = kvdb.Query(sess, `DELETE FROM kv WHERE key=:key`, "201").Exec()
	assert.Equal(t, nil, err)

}
