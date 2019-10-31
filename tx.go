package cqlx

type Tx = *Sessionx

///// AUTO_CLOSING PSEUDO-FUNCTIONAL FAKE TRANSACTIONS
func viewtx(db *DB, fn func(tx Tx) error) error {
	s, err := db.Session()
	if err != nil {
		return err
	}
	defer s.Close()
	return fn(Tx(s))
}

func updatetx(db *DB, fn func(tx Tx) error) error {
	s, err := db.Session()
	if err != nil {
		return err
	}
	defer s.Close()
	return fn(Tx(s))
}
