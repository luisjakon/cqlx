package cqlx

import (
	"context"

	"github.com/gocql/gocql"
)

func (q *Queryx) Consistency(c gocql.Consistency) *Queryx {
	q.Queryx.Consistency(c)
	return q
}

func (q *Queryx) CustomPayload(customPayload map[string][]byte) *Queryx {
	q.Queryx.CustomPayload(customPayload)
	return q
}

func (q *Queryx) Trace(trace gocql.Tracer) *Queryx {
	q.Queryx.Trace(trace)
	return q
}

func (q *Queryx) Observer(observer gocql.QueryObserver) *Queryx {
	q.Queryx.Observer(observer)
	return q
}

func (q *Queryx) PageSize(n int) *Queryx {
	q.Queryx.PageSize(n)
	return q
}

func (q *Queryx) DefaultTimestamp(enable bool) *Queryx {
	q.Queryx.DefaultTimestamp(enable)
	return q
}

func (q *Queryx) WithTimestamp(timestamp int64) *Queryx {
	q.Queryx.WithTimestamp(timestamp)
	return q
}

func (q *Queryx) RoutingKey(routingKey []byte) *Queryx {
	q.Queryx.RoutingKey(routingKey)
	return q
}

func (q *Queryx) WithContext(ctx context.Context) *Queryx {
	q.Queryx = q.Queryx.WithContext(ctx)
	return q
}

func (q *Queryx) Prefetch(p float64) *Queryx {
	q.Queryx.Prefetch(p)
	return q
}

func (q *Queryx) RetryPolicy(r gocql.RetryPolicy) *Queryx {
	q.Queryx.RetryPolicy(r)
	return q
}

func (q *Queryx) SetSpeculativeExecutionPolicy(sp gocql.SpeculativeExecutionPolicy) *Queryx {
	q.Queryx.SetSpeculativeExecutionPolicy(sp)
	return q
}

func (q *Queryx) Idempotent(value bool) *Queryx {
	q.Queryx.Idempotent(value)
	return q
}

func (q *Queryx) Bind(v ...interface{}) *Queryx {
	q.Queryx.Bind(v...)
	return q
}

func (q *Queryx) SerialConsistency(cons gocql.SerialConsistency) *Queryx {
	q.Queryx.SerialConsistency(cons)
	return q
}

func (q *Queryx) PageState(state []byte) *Queryx {
	q.Queryx.PageState(state)
	return q
}

func (q *Queryx) NoSkipMetadata() *Queryx {
	q.Queryx.NoSkipMetadata()
	return q
}
