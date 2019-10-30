package cqlx_test

import (
	"testing"

	"github.com/luisjakon/cqlx"
	"github.com/scylladb/gocqlx/qb"
	"github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {
	sess, err := newRawSession(dbkeyspace, dbhost)
	defer sess.Close()
	if err != nil {
		t.Errorf("%T: %s", err, err.Error())
	}

	keyspaceMetadata, err := sess.KeyspaceMetadata(dbkeyspace)
	if err != nil {
		t.Errorf("%T: %s", err, err.Error())
	}

	for _, table := range keyspaceMetadata.Tables {
		// Check keyspace exists
		if table.Keyspace != "cqlx_test_db" {
			t.Error("Invalid keyspace")
		}
		// Check kv table exists
		if table.Name != "kv" {
			t.Error("Invalid table")
		}
	}
}

func TestDBViewTx(t *testing.T) {

	get := qb.Select("kv").Where(qb.Eq("key")).Limit(1)

	db.View(func(tx cqlx.Tx) error {

		res := kv{}
		err := tx.Query(get, "4").Get(&res)

		assert.Equal(t, nil, err)
		assert.Equal(t, "4", res.Key)
		assert.Equal(t, "val4", res.Value)

		return err
	})

}

func TestDBUpdateTx(t *testing.T) {

	get := qb.Select("kv").Where(qb.Eq("key"))
	update := qb.Update("kv").Where(qb.Eq("key")).Set("value")
	delete := qb.Delete("kv").Where(qb.Eq("key"))

	// Upsert record
	db.Update(func(tx cqlx.Tx) error {

		err := tx.Query(update).Put(&kv{Key: "5", Value: "val5"})
		assert.Nil(t, err)

		return err
	})

	// Retrieve and compare
	db.View(func(tx cqlx.Tx) error {

		res := kv{}
		err := tx.Query(get, "5").Get(&res)

		assert.Equal(t, nil, err)
		assert.Equal(t, "5", res.Key)
		assert.Equal(t, "val5", res.Value)

		return err
	})

	// Remove record
	db.Update(func(tx cqlx.Tx) error {

		err := tx.Query(delete, "5").Exec()
		assert.Equal(t, nil, err)

		return err
	})

}
