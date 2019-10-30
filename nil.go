package cqlx

var (
	_NilSession = &nilSession{}
	_NilQuery   = &nilQuery{}
	_NilIter    = &nilIter{}
)

//// NilSession
type nilSession struct{}

func (s *nilSession) Query(query interface{}, args ...interface{}) Queryx {
	return _NilQuery
}

func (s *nilSession) Exec(stmt interface{}) error {
	return ErrInvalidSession
}

func (s *nilSession) Close() error {
	return nil
}

//// NilQuery
type nilQuery struct{}

func (q *nilQuery) Exec() error {
	return ErrInvalidQuery
}

func (q *nilQuery) Put(newitem interface{}) error {
	return ErrInvalidQuery
}

func (q *nilQuery) Get(res interface{}) error {
	return ErrInvalidQuery
}

func (q *nilQuery) Iter() Iterx {
	return _NilIter
}

//// NilIter
type nilIter struct{}

func (i nilIter) Next(dest interface{}) bool {
	return false
}

func (i nilIter) Close() error {
	return ErrNilIter
}
