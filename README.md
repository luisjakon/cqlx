cqlx
=====

Package cqlx provides simple inobstrusive bindings for handling raw ```gocql/gocql``` queries and ```scylladb/gocqlx``` query-builders via an uniform API.

Project Website: https://github.com/luisjakon/cqlx<br>

Installation
------------

    go get github.com/luisjakon/cqlx


Features
--------

* Automatically processes both raw CQL statements and ```scylladb/gocqlx``` query builders
* Provides full pass-thru access to the underlying structs and methods built into the venerable ```gocql/gocql``` and ```scylladb/gocqlx``` client packages
* Handles single-shot autoclosing pseudo-transactions via the db.Tx(...) interface


Example of correct Tx usage:
```go
err := db.Tx(func(tx cqlx.Tx) error {
    return tx.Exec("INSERT INTO table ...;")
})
```
Example of incorrect Tx usage:
```go
err := db.Tx(func(tx cqlx.Tx) error {
    defer tx.Close()
    return tx.Exec("INSERT INTO table ...;")
})
```
Since the transaction is managed there is no need to issue a deferred closing.


Extensibility
--------

* Extended usage of cqlx package structs:

Example of setting the cqlx.DB cluster connection properties
```go
db := cqlx.Open("example", "192.168.1.225")
db.Consistency = gocql.LocalOne
db.Timeout = 3 * time.Second
db.ReconnectInterval = 6 * time.Second
```

Example of setting the cqlx.Session properties
```go
sess := db.Session()
sess.SetPageSize(-1)
```

Example of setting the cqlx.Queryx properties
```go
qry = sess.Queryx("SELECT * FROM kv WHERE key='1'").Consistency(1)
err := qry.Get(&res)
```
Since all cqlx structs derive from the base gocql and gocqlx packages, all of their implemented features are immediately available to the caller.


Example
-------

```go
//// Showcase most common use-cases and patterns when using structs and methods from this package
package main

import (
	"log"

	"github.com/luisjakon/cqlx"
	"github.com/scylladb/gocqlx/qb"
)

type kv struct {
	Key   int
	Value string
}

var (
	// pre-canned statements
	createks  = `CREATE KEYSPACE IF NOT EXISTS example WITH REPLICATION = {'class' : 'SimpleStrategy', 'replication_factor' : 1 };`
	createtbl = `CREATE TABLE IF NOT EXISTS example.kv (key int, value text, PRIMARY KEY (key));`
	dropks    = `DROP KEYSPACE IF EXISTS example;`
)

var (
	res kv
	db  *cqlx.DB
)

func main() {
	// Update the underlying *gocql.ClusterConfig keyspace
	db.Keyspace = "example"

	// Create session
	sess, err := db.Session()
	if err != nil {
		panic(err)
	}
	defer sess.Close()

	// Insert record using raw cql statements
	sess.Queryx(`INSERT INTO kv (key, value) VALUES (:key, :value)`).Put(&kv{0, "v0"})

	// Insert record using gocqlx query builder package
	sess.Queryx(qb.Insert("kv").Columns("key", "value")).Put(&kv{1, "v1"})

	// Retrieve record using raw cql statements
	sess.Queryx(`SELECT * FROM kv WHERE key=:key`, 0).Get(&res)

	// Retrieve record using gocqlx query builder package
	sess.Queryx(qb.Select("kv").Where(qb.Eq("key")), 1).Get(&res)

	// Iterate through records
	it := sess.Queryx(qb.Select("kv")).Iter()
	defer it.Close()

	for it.Next(&res) {
		log.Printf("res: %+v", res)
	}

	// Delete record using raw cql statements
	sess.Queryx(`DELETE FROM kv WHERE key=:key`, 0).Exec()

	// Delete record using gocqlx query builder package
	sess.Queryx(qb.Delete("kv").Where(qb.Eq("key")), 1).Exec()

	// Remove db
	execute(db, dropks)
}

func init() {
	db = cqlx.Open("", "192.168.1.161")
	execute(db, createks)
	execute(db, createtbl)
}

func execute(db *cqlx.DB, stmt string) {
	// Use auto-closing pseudo-Tx
	db.Update(func(tx cqlx.Tx) error {
		return tx.Query(stmt).Exec()
	})
}
```


CRUD Example (New)
-------

```go
//// Showcase CRUD-like utilities from this package
package main

import (
	"fmt"

	"github.com/luisjakon/cqlx"
)

type kv struct {
	Key   int
	Value string
}

var res kv

func main() {

	// Create session
	sess, _ := cqlx.Open("example", "192.168.1.161").Session()
	defer sess.Close()

	// Create crud struct
	kvdb := cqlx.Crud{
		`SELECT * FROM kv WHERE key=:key`,
		`INSERT INTO kv (key, value) VALUES (:key, :value)`,
		`UPDATE kv SET value=:value WHERE key=:key IF EXISTS`,
		`DELETE FROM kv WHERE key=:key`,
	}

	// Use available crud methods
	kvdb.Insert(sess, &kv{2, "val2"})
	kvdb.Update(sess, &kv{2, "val3"})
	kvdb.Get(sess, &kv{Key: 2}, &res)
	kvdb.Delete(sess, &kv{Key: 2})

	// Use explicit raw queries 
	kvdb.Query(sess, `SELECT value FROM kv WHERE key=:key`, &kv{Key: 2}).Get(&res)

	// Or gocqlx.qb queries instead
	kvdb.Query(sess, qb.Select("kv").Where(qb.Eq("key")), &kv{Key: 2}).Get(&res)

	...
}
```


Package Dependencies
---------

* [gocql](https://github.com/gocql/gocql) the primary Cassandra client library for the Go language.
* [gocqlx](https://github.com/scylladb/gocqlx) is a gocql extension that automates data binding, adds named queries support, provides flexible query builders and plays well with gocql.


License
-------

> Copyright (c) 2019 Luis Jakón. All rights reserved.