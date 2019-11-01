package cqlx_test

import (
	"testing"

	"github.com/luisjakon/cqlx"
	"github.com/scylladb/gocqlx/qb"
	"github.com/stretchr/testify/assert"
)

func TestViewTx(t *testing.T) {

	get := qb.Select("kv").Where(qb.Eq("key")).Limit(1)

	db.View(func(tx cqlx.Tx) error {

		res := kv{}
		err := tx.Queryx(get, "4").Get(&res)

		assert.Equal(t, nil, err)
		assert.Equal(t, "4", res.Key)
		assert.Equal(t, "val4", res.Value)

		return err
	})

}

func TestUpdateTx(t *testing.T) {

	get := qb.Select("kv").Where(qb.Eq("key"))
	update := qb.Update("kv").Where(qb.Eq("key")).Set("value")
	delete := qb.Delete("kv").Where(qb.Eq("key"))

	// Upsert record
	db.Update(func(tx cqlx.Tx) error {

		err := tx.Queryx(update).Put(&kv{Key: "5", Value: "val5"})
		assert.Nil(t, err)

		return err
	})

	// Retrieve and compare
	db.View(func(tx cqlx.Tx) error {

		res := kv{}
		err := tx.Queryx(get, "5").Get(&res)

		assert.Equal(t, nil, err)
		assert.Equal(t, "5", res.Key)
		assert.Equal(t, "val5", res.Value)

		return err
	})

	// Remove record
	db.Update(func(tx cqlx.Tx) error {

		err := tx.Queryx(delete, "5").Exec()
		assert.Equal(t, nil, err)

		return err
	})

}
