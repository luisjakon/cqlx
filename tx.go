package cqlx

type Tx = *Sessionx

///// PSEUDO-FUNCTIONAL FAKE TRANSACTIONS
func viewtx(s *Sessionx, fn func(tx Tx) error) error {
	defer s.Close()
	return fn(Tx(s))
}

func updatetx(s *Sessionx, fn func(tx Tx) error) error {
	defer s.Close()
	return fn(Tx(s))
}
