// Code generated by pggen. DO NOT EDIT.

package custom_types

import (
	"context"
	"fmt"
	"sync"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mypricehealth/pggen/example/custom_types/mytype"
)

type QueryName struct{}

// Querier is a typesafe Go interface backed by SQL queries.
type Querier interface {
	CustomTypes(ctx context.Context) (CustomTypesRow, error)

	CustomMyInt(ctx context.Context) (int, error)

	IntArray(ctx context.Context) ([][]int32, error)
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

const customTypesSQL = `SELECT 'some_text', 1::bigint;`

type CustomTypesRow struct {
	Column mytype.String `json:"?column?"`
	Int8   CustomInt     `json:"int8"`
}

// CustomTypes implements Querier.CustomTypes.
func (q *DBQuerier) CustomTypes(ctx context.Context) (CustomTypesRow, error) {
	ctx = context.WithValue(ctx, QueryName{}, "CustomTypes")

	err := registerTypes(ctx, q.conn)
	if err != nil {
		return CustomTypesRow{}, q.errWrap(err)
	}
	rows, err := q.conn.Query(ctx, customTypesSQL)
	if err != nil {
		return CustomTypesRow{}, fmt.Errorf("query CustomTypes: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[CustomTypesRow])
	return res, q.errWrap(err)
}

const customMyIntSQL = `SELECT '5'::my_int as int5;`

// CustomMyInt implements Querier.CustomMyInt.
func (q *DBQuerier) CustomMyInt(ctx context.Context) (int, error) {
	ctx = context.WithValue(ctx, QueryName{}, "CustomMyInt")

	err := registerTypes(ctx, q.conn)
	if err != nil {
		return 0, q.errWrap(err)
	}
	rows, err := q.conn.Query(ctx, customMyIntSQL)
	if err != nil {
		return 0, fmt.Errorf("query CustomMyInt: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[int])
	return res, q.errWrap(err)
}

const intArraySQL = `SELECT ARRAY ['5', '6', '7']::int[] as ints;`

// IntArray implements Querier.IntArray.
func (q *DBQuerier) IntArray(ctx context.Context) ([][]int32, error) {
	ctx = context.WithValue(ctx, QueryName{}, "IntArray")

	err := registerTypes(ctx, q.conn)
	if err != nil {
		return nil, q.errWrap(err)
	}
	rows, err := q.conn.Query(ctx, intArraySQL)
	if err != nil {
		return nil, fmt.Errorf("query IntArray: %w", q.errWrap(err))
	}
	res, err := pgx.CollectRows(rows, pgx.RowTo[[]int32])
	return res, q.errWrap(err)
}
