// Code generated by pggen. DO NOT EDIT.

package slices

import (
	"context"
	"fmt"
	"sync"
	"time"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type QueryName struct{}

// Querier is a typesafe Go interface backed by SQL queries.
type Querier interface {
	GetBools(ctx context.Context, data []bool) ([]bool, error)

	GetOneTimestamp(ctx context.Context, data *time.Time) (*time.Time, error)

	GetManyTimestamptzs(ctx context.Context, data []time.Time) ([]*time.Time, error)

	GetManyTimestamps(ctx context.Context, data []*time.Time) ([]*time.Time, error)
}

var _ Querier = &DBQuerier{}

type DBQuerier struct {
	conn    genericConn
	errWrap func(err error) error
}

// genericConn is a connection like *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
type genericConn interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	TypeMap() *pgtype.Map
	LoadType(ctx context.Context, typeName string) (*pgtype.Type, error)
}

// NewQuerier creates a DBQuerier that implements Querier.
func NewQuerier(conn genericConn) *DBQuerier {
	return &DBQuerier{
		conn: conn,
		errWrap: func(err error) error {
			return err
		},
	}
}

var registerOnce sync.Once
var registerErr error

func registerTypes(ctx context.Context, conn genericConn) error {
	registerOnce.Do(func() {
		typeMap := conn.TypeMap()

		pgxdecimal.Register(typeMap)
		for _, typ := range typesToRegister {
			dt, err := conn.LoadType(ctx, typ)
			if err != nil {
				registerErr = fmt.Errorf("could not register type %q: %w", typ, err)
				return
			}
			typeMap.RegisterType(dt)
		}
	})

	return registerErr
}

var typesToRegister = []string{}

func addTypeToRegister(typ string) struct{} {
	typesToRegister = append(typesToRegister, typ)
	return struct{}{}
}

const getBoolsSQL = `SELECT $1::boolean[];`

// GetBools implements Querier.GetBools.
func (q *DBQuerier) GetBools(ctx context.Context, data []bool) ([]bool, error) {
	ctx = context.WithValue(ctx, QueryName{}, "GetBools")

	err := registerTypes(ctx, q.conn)
	if err != nil {
		return nil, q.errWrap(err)
	}
	rows, err := q.conn.Query(ctx, getBoolsSQL, data)
	if err != nil {
		return nil, fmt.Errorf("query GetBools: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[[]bool])
	return res, q.errWrap(err)
}

const getOneTimestampSQL = `SELECT $1::timestamp;`

// GetOneTimestamp implements Querier.GetOneTimestamp.
func (q *DBQuerier) GetOneTimestamp(ctx context.Context, data *time.Time) (*time.Time, error) {
	ctx = context.WithValue(ctx, QueryName{}, "GetOneTimestamp")

	err := registerTypes(ctx, q.conn)
	if err != nil {
		return nil, q.errWrap(err)
	}
	rows, err := q.conn.Query(ctx, getOneTimestampSQL, data)
	if err != nil {
		return nil, fmt.Errorf("query GetOneTimestamp: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[*time.Time])
	return res, q.errWrap(err)
}

const getManyTimestamptzsSQL = `SELECT *
FROM unnest($1::timestamptz[]);`

// GetManyTimestamptzs implements Querier.GetManyTimestamptzs.
func (q *DBQuerier) GetManyTimestamptzs(ctx context.Context, data []time.Time) ([]*time.Time, error) {
	ctx = context.WithValue(ctx, QueryName{}, "GetManyTimestamptzs")

	err := registerTypes(ctx, q.conn)
	if err != nil {
		return nil, q.errWrap(err)
	}
	rows, err := q.conn.Query(ctx, getManyTimestamptzsSQL, data)
	if err != nil {
		return nil, fmt.Errorf("query GetManyTimestamptzs: %w", q.errWrap(err))
	}
	res, err := pgx.CollectRows(rows, pgx.RowTo[*time.Time])
	return res, q.errWrap(err)
}

const getManyTimestampsSQL = `SELECT *
FROM unnest($1::timestamp[]);`

// GetManyTimestamps implements Querier.GetManyTimestamps.
func (q *DBQuerier) GetManyTimestamps(ctx context.Context, data []*time.Time) ([]*time.Time, error) {
	ctx = context.WithValue(ctx, QueryName{}, "GetManyTimestamps")

	err := registerTypes(ctx, q.conn)
	if err != nil {
		return nil, q.errWrap(err)
	}
	rows, err := q.conn.Query(ctx, getManyTimestampsSQL, data)
	if err != nil {
		return nil, fmt.Errorf("query GetManyTimestamps: %w", q.errWrap(err))
	}
	res, err := pgx.CollectRows(rows, pgx.RowTo[*time.Time])
	return res, q.errWrap(err)
}
