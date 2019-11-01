package cqlx_test

import (
	"testing"
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
