package cqlx_test

import (
	"strconv"
	"testing"

	"github.com/scylladb/gocqlx/qb"
	"github.com/stretchr/testify/assert"
)

func TestRawIterScan(t *testing.T) {

	sess, err := db.Session()
	defer sess.Close()

	assert.Equal(t, nil, err)

	it := sess.Queryx(`SELECT * FROM kv;`).Iter() // `SELECT * FROM kv;`
	defer it.Close()

	res := kv{}
	ok := it.Next(&res)

	assert.Equal(t, true, ok)
	assert.Equal(t, "4", res.Key)
	assert.Equal(t, "val4", res.Value)
}

func TestIterScan(t *testing.T) {

	sess, err := db.Session()
	defer sess.Close()

	assert.Equal(t, nil, err)

	it := sess.Queryx(qb.Select("kv").Limit(1)).Iter() // `SELECT * FROM kv LIMIT 1;`
	defer it.Close()

	res := kv{}
	ok := it.Next(&res)

	assert.Equal(t, true, ok)
	assert.Equal(t, "4", res.Key)
	assert.Equal(t, "val4", res.Value)
}

func TestIterScanAll(t *testing.T) {

	sess, err := db.Session()
	defer sess.Close()

	assert.Equal(t, nil, err)

	it := sess.Queryx(qb.Select("kv")).Iter() // `SELECT * FROM kv;`
	defer it.Close()

	res := kv{}
	for i := 4; it.Next(&res); i-- {
		assert.Equal(t, strconv.Itoa(i), res.Key)
	}
}
