cqlx
=====

Package cqlx implements a simpler, generic wrapper interface for interacting with CQL queries and native golang structs

Project Website: https://github.com/luisjakon/cqlx<br>

Installation
------------

    go get github.com/luisjakon/cqlx


Features
--------

* Integrates and makes use of the venerable gocql/gocql and scylladb/gocqlx client packages
* Simple intuitive API for executing queries using raw CQL and/or query builders from the gocqlx/qb package
* Single-shot, autoclosing transactions via the db.Tx(...) interface


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

Example
-------

```go
/* Before you execute the program, Launch `cqlsh` and execute:
create keyspace example with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
create table example.kv(key text, value text, PRIMARY KEY(id));
*/
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

var res kv

var (
	// pre-canned queries
	getall = qb.Select("kv").Columns("*")                       // `SELECT * FROM kv;`
	get    = qb.Select("kv").Where(qb.Eq("key"))                // `SELECT * FROM kv WHERE key=?;`
	put    = qb.Update("kv").Set("value").Where(qb.Eq("key"))   // `UPDATE kv SET value=? WHERE key=?;`
	delete = qb.Delete("kv").Where(qb.Eq("key"))                // `DELETE *  FROM kv WHERE key=?;`
)

func main() {
	// create connection
	db := cqlx.Open("example", "192.168.1.225")
	sess := db.Session()
	defer sess.Close()

	// Insert record
	err := sess.Queryx(put).Put(&kv{Key: "1", Value: "val1"})
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

}
```

Ecosystem
---------

* [gocql](https://github.com/gocql/gocql) the primary Cassandra client library for the Go language.
* [gocqlx](https://github.com/scylladb/gocqlx) is a gocql extension that automates data binding, adds named queries support, provides flexible query builders and plays well with gocql.


License
-------

> Copyright (c) 2019 The Luis Jakón. All rights reserved.