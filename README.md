cqlx
=====
Package cqlx provides unified bindings for handling [gocql/gocql](https://github.com/gocql/gocql) raw cql statements and [scylladb/gocqlx](https://github.com/scylladb/gocqlx/tree/master/qb) query-builders

Project Website: https://github.com/luisjakon/cqlx<br>

Installation
------------

    go get github.com/luisjakon/cqlx


Features
--------

* Delivers an uniform abstraction layer along with an ultra-simple-to-use, inobstrusive development API 
* Enables simple, easy-to-use CRUD-like interactions via ```cqlx.Crud``` structs
* Works with either ```gocql/gocql``` raw cql statements and/or ```scylladb/gocqlx``` query builders indiscriminately
* Provides full (pass-thru) access to all underlying ```gocql/gocql``` and ```scylladb/gocqlx``` structs and methods
* Handles single-shot autoclosing pseudo-transactions via the ```db.Tx(...)``` interface


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


 CRUD Interaction
-------

Example of using ```cqlx.Crud``` structs and functions
```go
// Create crud struct
kvdb := cqlx.Crud{
	`SELECT * FROM kv WHERE key=:key`,
	`INSERT INTO kv (key, value) VALUES (:key, :value)`,
	`UPDATE kv SET value=:value WHERE key=:key IF EXISTS`,
	`DELETE FROM kv WHERE key=:key`,
}

// Open a cluster connection
sess, _ := cqlx.Open("example", "192.168.1.161").Session()

// Use available crud methods
kvdb.Insert(sess, &kv{2, "val2"})
kvdb.Update(sess, &kv{2, "val3"})
kvdb.Select(sess, &kv{Key: 2}).Get(&res)
kvdb.Delete(sess, &kv{Key: 2})

// Use explicit raw cql queries
kvdb.Query(sess, `SELECT value FROM kv WHERE key=:key`, &kv{Key: 2}).Get(&res)

// Use gocqlx.qb queries
kvdb.Query(sess, qb.Select("kv").Columns("value").Where(qb.Eq("key")), &kv{Key: 2}).Get(&res)
```
While not a necessity, making use of ```cqlx.Crud``` in your code can significantly simplify and standardize access patterns to otherwise complex cql tables.


Extensibility
--------

* Extended usage of package structs:

Example of setting the ```cqlx.DB``` connection properties
```go
db := cqlx.Open("example", "192.168.1.225")
db.Consistency = gocql.LocalOne
db.Timeout = 3 * time.Second
db.ReconnectInterval = 6 * time.Second
```

Example of setting the ```cqlx.Session``` properties
```go
sess := db.Session()
sess.SetPageSize(-1)
```

Example of setting the ```cqlx.Queryx``` properties
```go
qry := sess.Queryx("SELECT * FROM kv WHERE key='1'").Consistency(1)
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


Package Dependencies
---------

* [gocql](https://github.com/gocql/gocql) the primary Cassandra client library for the Go language.
* [gocqlx](https://github.com/scylladb/gocqlx) is a gocql extension that automates data binding, adds named queries support, provides flexible query builders and plays well with gocql.


License
-------

> Copyright (c) 2019 Luis Jakón. All rights reserved.