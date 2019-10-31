package cqlx

type Tx Sessionx

///// PSEUDO-FUNCTIONAL FAKE TRANSACTIONS
func newTx(s Sessionx) Tx {
	return Tx(s)
}

func viewtx(s Sessionx, fn func(tx Tx) error) error {
	tx := newTx(s)
	return fn(tx)
}

func updatetx(s Sessionx, fn func(tx Tx) error) error {
	tx := newTx(s)
	return fn(tx)
}
