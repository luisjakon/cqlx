package cqlx

type Tx = *Session

///// PSEUDO-FUNCTIONAL FAKE TRANSACTIONS
func viewtx(s *Session, fn func(tx Tx) error) error {
	return fn(Tx(s))
}

func updatetx(s *Session, fn func(tx Tx) error) error {
	return fn(Tx(s))
}
