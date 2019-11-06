package cqlx_test

import (
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

type bkv struct {
	Key   int
	Value string
}

func init_bench_db() {
	sess, err := newRawSession("", dbhost)
	defer sess.Close()
	if err != nil {
		panic(err)
	}
	var db_stms = []string{
		`DROP KEYSPACE IF EXISTS cqlx_benchtest_db;`,
		`CREATE KEYSPACE IF NOT EXISTS cqlx_benchtest_db WITH REPLICATION = {'class' : 'SimpleStrategy', 'replication_factor' : 1 };`,
		`CREATE TABLE IF NOT EXISTS cqlx_benchtest_db.bkv (
			key int,   
			value text,
			PRIMARY KEY (key)
		);`,
	}
	for _, s := range db_stms {
		if err := sess.Query(s).Exec(); err != nil {
			panic(err)
		}
	}
}

func drop_bench_db() {
	sess, err := newRawSession("", dbhost)
	defer sess.Close()
	if err != nil {
		log.Print(err.Error())
		return
	}
	if err := sess.Query(`DROP KEYSPACE IF EXISTS cqlx_benchtest_db;`).Exec(); err != nil {
		log.Print(err.Error())
		return
	}
}

func load_fixtures() []*bkv {
	var v []*bkv = make([]*bkv, 100)
	for i, _ := range v {
		v[i] = &bkv{i, strconv.Itoa(i)}
	}
	return v
}

func load_data(b *testing.B, session *gocql.Session, data []*bkv) {

	stmt, names := qb.Insert("cqlx_benchtest_db.bkv").Columns("key", "value").ToCql()
	q := gocqlx.Query(session.Query(stmt), names)

	for _, p := range data {
		if err := q.BindStruct(p).Exec(); err != nil {
			b.Fatal(err)
		}
	}
}

func createRawSession(tb testing.TB) *gocql.Session {
	c := gocql.NewCluster(strings.Split(dbhost, ",")...)
	c.Keyspace = dbkeyspace
	s, err := c.CreateSession()
	if err != nil {
		tb.Fatalf("createRawSession: %s", err.Error())
	}
	return s
}

////
//// Inserts
////

// BenchmarkGocqlInsert
func BenchmarkGocqlInsert(b *testing.B) {
	data := load_fixtures()
	sess := createRawSession(b)
	defer sess.Close()

	stmt, _ := qb.Insert("cqlx_benchtest_db.bkv").Columns("key", "value").ToCql()
	q := sess.Query(stmt)
	defer q.Release()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := data[i%len(data)]
		if err := q.Bind(p.Key, p.Value).Exec(); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGocqlxInsert
func BenchmarkGocqlxInsert(b *testing.B) {
	data := load_fixtures()
	sess := createRawSession(b)
	defer sess.Close()

	stmt, names := qb.Insert("cqlx_benchtest_db.bkv").Columns("key", "value").ToCql()
	q := gocqlx.Query(sess.Query(stmt), names)
	defer q.Release()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := data[i%len(data)]
		if err := q.BindStruct(p).Exec(); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCqlxInsert
func BenchmarkCqlxInsert(b *testing.B) {
	data := load_fixtures()
	sess, err := newRawSession(dbkeyspace, dbhost)
	if err != nil {
		b.Fatal(err)
	}
	defer sess.Close()

	stmt := qb.Insert("cqlx_benchtest_db.bkv").Columns("key", "value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q := sess.Queryx(stmt)
		p := data[i%len(data)]
		if err := q.Put(p); err != nil {
			b.Fatal(err)
		}
	}
}

////
//// Get
////

// BenchmarkGocqlGet
func BenchmarkGocqlGet(b *testing.B) {
	data := load_fixtures()
	sess := createRawSession(b)
	defer sess.Close()

	load_data(b, sess, data)

	stmt, _ := qb.Select("cqlx_benchtest_db.bkv").Columns("key", "value").Where(qb.Eq("key")).Limit(1).ToCql()
	var p bkv

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q := sess.Query(stmt)
		q.Bind(data[i%len(data)].Key)

		if err := q.Scan(&p.Key, &p.Value); err != nil {
			b.Fatal(err)
		}

		q.Release()
	}
}

// BenchmarkGocqlxGet
func BenchmarkGocqlxGet(b *testing.B) {
	data := load_fixtures()
	sess, err := newRawSession(dbkeyspace, dbhost)
	if err != nil {
		b.Fatal(err)
	}
	defer sess.Close()

	load_data(b, sess.Session, data)

	stmt, _ := qb.Select("cqlx_benchtest_db.bkv").Columns("key", "value").Where(qb.Eq("key")).Limit(1).ToCql()
	var p bkv

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q := sess.Query(stmt).Bind(data[i%len(data)].Key)

		if err := gocqlx.Query(q, nil).GetRelease(&p); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCqlxGet
func BenchmarkCqlxGet(b *testing.B) {
	data := load_fixtures()
	sess, err := newRawSession(dbkeyspace, dbhost)
	if err != nil {
		b.Fatal(err)
	}
	defer sess.Close()

	load_data(b, sess.Session, data)

	stmt := qb.Select("cqlx_benchtest_db.bkv").Columns("key", "value").Where(qb.Eq("key")).Limit(1)
	var p bkv

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		if err := sess.Queryx(stmt, &p).Exec(); err != nil {
			b.Fatal(err)
		}
	}
}

////
//// Select
////

// BenchmarkGocqlSelect.
func BenchmarkGocqlSelect(b *testing.B) {
	data := load_fixtures()
	sess := createRawSession(b)
	defer sess.Close()

	load_data(b, sess, data)

	stmt, _ := qb.Select("cqlx_benchtest_db.bkv").Columns("key", "value").Limit(100).ToCql()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := make([]*bkv, 100)
		q := sess.Query(stmt)
		i := q.Iter()

		p := new(bkv)
		for i.Scan(&p.Key, &p.Value) {
			v = append(v, p)
			p = new(bkv)
		}
		if err := i.Close(); err != nil {
			b.Fatal(err)
		}

		q.Release()
	}
}

// BenchmarkGocqlSelect.
func BenchmarkGocqlxSelect(b *testing.B) {
	data := load_fixtures()
	sess := createRawSession(b)
	defer sess.Close()

	load_data(b, sess, data)

	stmt, _ := qb.Select("cqlx_benchtest_db.bkv").Columns("key", "value").Limit(100).ToCql()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q := sess.Query(stmt)
		var v []*bkv

		if err := gocqlx.Query(q, nil).SelectRelease(&v); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCqlxSelect
func BenchmarkCqlxSelect(b *testing.B) {
	data := load_fixtures()
	sess, err := newRawSession(dbkeyspace, dbhost)
	if err != nil {
		b.Fatal(err)
	}
	defer sess.Close()

	load_data(b, sess.Session, data)

	stmt := qb.Select("cqlx_benchtest_db.bkv").Columns("key", "value").Limit(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v []*bkv

		if err := sess.Queryx(stmt).Get(&v); err != nil {
			b.Fatal(err)
		}
	}
}
