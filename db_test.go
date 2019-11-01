package cqlx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {
	sess, err := newRawSession(dbkeyspace, dbhost)
	defer sess.Close()
	assert.Equal(t, nil, err)

	keyspaceMetadata, err := sess.KeyspaceMetadata(dbkeyspace)
	assert.Equal(t, nil, err)

	for _, table := range keyspaceMetadata.Tables {
		assert.Equal(t, "cqlx_test_db", table.Keyspace)
		assert.Equal(t, "kv", table.Name)
	}
}
