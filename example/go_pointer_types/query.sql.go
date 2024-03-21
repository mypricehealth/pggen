// Code generated by pggen. DO NOT EDIT.

package go_pointer_types

import (
	"sync"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

// Querier is a typesafe Go interface backed by SQL queries.
type Querier interface {
	GenSeries1(ctx context.Context) (*int, error)

	GenSeries(ctx context.Context) ([]*int, error)

	GenSeriesArr1(ctx context.Context) ([]int, error)

	GenSeriesArr(ctx context.Context) ([][]int, error)

	GenSeriesStr1(ctx context.Context) (*string, error)

	GenSeriesStr(ctx context.Context) ([]*string, error)
}

var _ Querier = &DBQuerier{}

type DBQuerier struct {
	conn  genericConn
}

// genericConn is a connection like *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
type genericConn interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

// NewQuerier creates a DBQuerier that implements Querier.
func NewQuerier(conn genericConn) *DBQuerier {
	return &DBQuerier{conn: conn}
}

const genSeries1SQL = `SELECT n
FROM generate_series(0, 2) n
LIMIT 1;`

// GenSeries1 implements Querier.GenSeries1.
func (q *DBQuerier) GenSeries1(ctx context.Context) (*int, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GenSeries1")
	rows, err := q.conn.Query(ctx, genSeries1SQL)
	if err != nil {
		return nil, fmt.Errorf("query GenSeries1: %w", err)
	}
	fds := rows.FieldDescriptions()
	plan0 := planScan(pgtype.TextCodec{}, fds[0], (**int)(nil))

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (*int, error) {
		vals := row.RawValues()
		var item *int
		if err := plan0.Scan(vals[0], &item); err != nil {
			return item, fmt.Errorf("scan GenSeries1.n: %w", err)
		}
		return item, nil
	})
}

const genSeriesSQL = `SELECT n
FROM generate_series(0, 2) n;`

// GenSeries implements Querier.GenSeries.
func (q *DBQuerier) GenSeries(ctx context.Context) ([]*int, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GenSeries")
	rows, err := q.conn.Query(ctx, genSeriesSQL)
	if err != nil {
		return nil, fmt.Errorf("query GenSeries: %w", err)
	}
	fds := rows.FieldDescriptions()
	plan0 := planScan(pgtype.TextCodec{}, fds[0], (**int)(nil))

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*int, error) {
		vals := row.RawValues()
		var item *int
		if err := plan0.Scan(vals[0], &item); err != nil {
			return item, fmt.Errorf("scan GenSeries.n: %w", err)
		}
		return item, nil
	})
}

const genSeriesArr1SQL = `SELECT array_agg(n)
FROM generate_series(0, 2) n;`

// GenSeriesArr1 implements Querier.GenSeriesArr1.
func (q *DBQuerier) GenSeriesArr1(ctx context.Context) ([]int, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GenSeriesArr1")
	rows, err := q.conn.Query(ctx, genSeriesArr1SQL)
	if err != nil {
		return nil, fmt.Errorf("query GenSeriesArr1: %w", err)
	}
	fds := rows.FieldDescriptions()
	plan0 := planScan(pgtype.TextCodec{}, fds[0], (*[]int)(nil))

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) ([]int, error) {
		vals := row.RawValues()
		var item []int
		if err := plan0.Scan(vals[0], &item.ArrayAgg); err != nil {
			return item, fmt.Errorf("scan GenSeriesArr1.array_agg: %w", err)
		}
		return item, nil
	})
}

const genSeriesArrSQL = `SELECT array_agg(n)
FROM generate_series(0, 2) n;`

// GenSeriesArr implements Querier.GenSeriesArr.
func (q *DBQuerier) GenSeriesArr(ctx context.Context) ([][]int, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GenSeriesArr")
	rows, err := q.conn.Query(ctx, genSeriesArrSQL)
	if err != nil {
		return nil, fmt.Errorf("query GenSeriesArr: %w", err)
	}
	fds := rows.FieldDescriptions()
	plan0 := planScan(pgtype.TextCodec{}, fds[0], (*[]int)(nil))

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) ([]int, error) {
		vals := row.RawValues()
		var item []int
		if err := plan0.Scan(vals[0], &item.ArrayAgg); err != nil {
			return item, fmt.Errorf("scan GenSeriesArr.array_agg: %w", err)
		}
		return item, nil
	})
}

const genSeriesStr1SQL = `SELECT n::text
FROM generate_series(0, 2) n
LIMIT 1;`

// GenSeriesStr1 implements Querier.GenSeriesStr1.
func (q *DBQuerier) GenSeriesStr1(ctx context.Context) (*string, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GenSeriesStr1")
	rows, err := q.conn.Query(ctx, genSeriesStr1SQL)
	if err != nil {
		return nil, fmt.Errorf("query GenSeriesStr1: %w", err)
	}
	fds := rows.FieldDescriptions()
	plan0 := planScan(pgtype.TextCodec{}, fds[0], (**string)(nil))

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (*string, error) {
		vals := row.RawValues()
		var item *string
		if err := plan0.Scan(vals[0], &item); err != nil {
			return item, fmt.Errorf("scan GenSeriesStr1.n: %w", err)
		}
		return item, nil
	})
}

const genSeriesStrSQL = `SELECT n::text
FROM generate_series(0, 2) n;`

// GenSeriesStr implements Querier.GenSeriesStr.
func (q *DBQuerier) GenSeriesStr(ctx context.Context) ([]*string, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GenSeriesStr")
	rows, err := q.conn.Query(ctx, genSeriesStrSQL)
	if err != nil {
		return nil, fmt.Errorf("query GenSeriesStr: %w", err)
	}
	fds := rows.FieldDescriptions()
	plan0 := planScan(pgtype.TextCodec{}, fds[0], (**string)(nil))

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*string, error) {
		vals := row.RawValues()
		var item *string
		if err := plan0.Scan(vals[0], &item); err != nil {
			return item, fmt.Errorf("scan GenSeriesStr.n: %w", err)
		}
		return item, nil
	})
}

type scanCacheKey struct {
	oid      uint32
	format   int16
	typeName string
}

var (
	plans   = make(map[scanCacheKey]pgtype.ScanPlan, 16)
	plansMu sync.RWMutex
)

func planScan(codec pgtype.Codec, fd pgconn.FieldDescription, target any) pgtype.ScanPlan {
	key := scanCacheKey{fd.DataTypeOID, fd.Format, fmt.Sprintf("%T", target)}
	plansMu.RLock()
	plan := plans[key]
	plansMu.RUnlock()
	if plan != nil {
		return plan
	}
	plan = codec.PlanScan(nil, fd.DataTypeOID, fd.Format, target)
	plansMu.Lock()
	plans[key] = plan
	plansMu.Unlock()
	return plan
}

type ptrScanner[T any] struct {
	basePlan pgtype.ScanPlan
}

func (s ptrScanner[T]) Scan(src []byte, dst any) error {
	if src == nil {
		return nil
	}
	d := dst.(**T)
	*d = new(T)
	return s.basePlan.Scan(src, *d)
}

func planPtrScan[T any](codec pgtype.Codec, fd pgconn.FieldDescription, target *T) pgtype.ScanPlan {
	key := scanCacheKey{fd.DataTypeOID, fd.Format, fmt.Sprintf("*%T", target)}
	plansMu.RLock()
	plan := plans[key]
	plansMu.RUnlock()
	if plan != nil {
		return plan
	}
	basePlan := planScan(codec, fd, target)
	ptrPlan := ptrScanner[T]{basePlan}
	plansMu.Lock()
	plans[key] = plan
	plansMu.Unlock()
	return ptrPlan
}