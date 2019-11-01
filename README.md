cqlx
=====

Package cqlx implements simple bindings for working with raw CQL and other query-builders via a uniform API

Project Website: https://github.com/luisjakon/cqlx<br>

Installation
------------

    go get github.com/luisjakon/cqlx


Features
--------

* Handles both raw CQL statements and ```scylladb/gocqlx``` query builders indiscriminately
* Provides full pass-thru access to the underlying structs and methods built into the venerable ```gocql/gocql``` and ```scylladb/gocqlx``` client packages
* Processes single-shot, autoclosing transactions via the db.Tx(...) interface


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

Example of setting the client connection timeout
```go
db := cqlx.Open("example", "192.168.1.225")
db.Timeout = 3 * time.Second
```

Example of setting the session page size
```go
sess := db.Session()
sess.SetPageSize(-1)
```

Example of adjusting a queryx consistency level
```go
qry = sess.Queryx("SELECT * FROM kv WHERE key='1'")
err := qry.Consistency(1).Get(&res)
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
	Key   string
	Value string
}

var (
	createks  = `CREATE KEYSPACE IF NOT EXISTS example WITH REPLICATION = {'class' : 'SimpleStrategy', 'replication_factor' : 1 };`
	createtbl = `CREATE TABLE IF NOT EXISTS example.kv (key text, value text, PRIMARY KEY (key));`
	dropks    = `DROP KEYSPACE IF EXISTS example;`
)

var (
	// pre-canned queries
	getall = qb.Select("kv").Columns("*")                     // `SELECT * FROM kv;`
	get    = qb.Select("kv").Where(qb.Eq("key"))              // `SELECT * FROM kv WHERE key=?;`
	put    = qb.Update("kv").Set("value").Where(qb.Eq("key")) // `UPDATE kv SET value=? WHERE key=?;`
	delete = qb.Delete("kv").Where(qb.Eq("key"))              // `DELETE *  FROM kv WHERE key=?;`
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

	// Insert record
	err = sess.Queryx(put).Put(&kv{Key: "1", Value: "val1"})
	if err != nil {
		log.Fatal("insert error:", err)
	}

	// Retrieve record
	err = sess.Queryx(get, "1").Get(&res)
	if err != nil {
		log.Fatal("get error:", err)
	}

	// Iterate through records
	it := sess.Queryx(getall).Iter()
	defer it.Close()

	for it.Next(&res) {
		log.Printf("res: %+v", res)
	}

	// Delete record
	err = sess.Queryx(delete, "1").Exec()
	if err != nil {
		log.Fatal("delete err:", err)
	}

	// Remove db
	execute(db, dropks)
}

func init() {
	db = cqlx.Open("", "192.168.10.135")
	execute(db, createks)
	execute(db, createtbl)
}

func execute(db *cqlx.DB, stmt string) {
	// Tap onto the underlying gocql.Query(...) interface via an auto-closing Tx
	db.Update(func(tx cqlx.Tx) error {
		return tx.Query(stmt).Exec()
	})
}

```

Package Dependencies
---------

* [gocql](https://github.com/gocql/gocql) the primary Cassandra client library for the Go language.
* [gocqlx](https://github.com/scylladb/gocqlx) is a gocql extension that automates data binding, adds named queries support, provides flexible query builders and plays well with gocql.


License
-------

> Copyright (c) 2019 Luis Jakón. All rights reserved.