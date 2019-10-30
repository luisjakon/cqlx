package cqlx

type Tx Session

///// PSEUDO-FUNCTIONAL FAKE TRANSACTIONS
func newTx(s Session) Tx {
	return Tx(s)
}

func viewtx(s Session, fn func(tx Tx) error) error {
	tx := newTx(s)
	return fn(tx)
}

func updatetx(s Session, fn func(tx Tx) error) error {
	tx := newTx(s)
	return fn(tx)
}
