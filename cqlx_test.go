package cqlx_test

import (
	"flag"
	"log"
	"os"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/luisjakon/cqlx"
)

var dbhost string
var dbkeyspace string
var db *cqlx.DB

type kv struct {
	Key   string
	Value string
}

func init() {
	flag.StringVar(&dbhost, "dbhost", "192.168.10.135", "cassandra host(s)")
	flag.StringVar(&dbkeyspace, "ks", "cqlx_test_db", "cassandra keyspace")
}

func init_flags() {
	if !flag.Parsed() {
		flag.Parse()
	}
}

func init_db() {
	sess, err := newRawSession("", dbhost)
	defer sess.Close()
	if err != nil {
		panic(err)
	}
	var db_stms = []string{
		`DROP KEYSPACE IF EXISTS cqlx_test_db;`,
		`CREATE KEYSPACE IF NOT EXISTS cqlx_test_db WITH REPLICATION = {'class' : 'SimpleStrategy', 'replication_factor' : 1 };`,
		`CREATE TABLE IF NOT EXISTS cqlx_test_db.kv (
			key text,   
			value text,
			PRIMARY KEY (key)
		);`,
		`INSERT INTO cqlx_test_db.kv (key,value) VALUES ('1','val1');`,
		`INSERT INTO cqlx_test_db.kv (key,value) VALUES ('2','val2');`,
		`INSERT INTO cqlx_test_db.kv (key,value) VALUES ('3','val3');`,
		`INSERT INTO cqlx_test_db.kv (key,value) VALUES ('4','val4');`,
	}
	for _, s := range db_stms {
		if err := sess.Query(s).Exec(); err != nil {
			panic(s + ", err: " + err.Error())
		}
	}
	db = cqlx.Open(dbkeyspace, dbhost)
}

func drop_db() {
	sess, err := newRawSession("", dbhost)
	defer sess.Close()
	if err != nil {
		log.Print(err.Error())
		return
	}
	if err := sess.Query(`DROP KEYSPACE IF EXISTS cqlx_test_db;`).Exec(); err != nil {
		log.Print(err.Error())
		return
	}
}

func newRawSession(keyspace, host string) (*cqlx.Sessionx, error) {
	db := cqlx.Open(keyspace, host)
	db.Keyspace = keyspace
	db.Consistency = gocql.LocalOne
	db.ReconnectInterval = 6 * time.Second
	db.Timeout = 3 * time.Second
	return db.Session()
}

func TestMain(m *testing.M) {
	init_flags()
	init_db()
	code := m.Run()
	drop_db()
	os.Exit(code)
}
