// Code generated by pggen. DO NOT EDIT.

package out

import (
	"context"
	"fmt"
	"sync"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type QueryName struct{}

// Querier is a typesafe Go interface backed by SQL queries.
type Querier interface {
	AlphaNested(ctx context.Context) (string, error)

	AlphaCompositeArray(ctx context.Context) ([]Alpha, error)

	Alpha(ctx context.Context) (string, error)

	Bravo(ctx context.Context) (string, error)
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

// Alpha represents the Postgres composite type "alpha".
type Alpha struct {
	Key *string `json:"key"`
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

var _ = addTypeToRegister("public.alpha")

var _ = addTypeToRegister("public._alpha")

const alphaNestedSQL = `SELECT 'alpha_nested' as output;`

// AlphaNested implements Querier.AlphaNested.
func (q *DBQuerier) AlphaNested(ctx context.Context) (string, error) {
	ctx = context.WithValue(ctx, QueryName{}, "AlphaNested")

	err := registerTypes(ctx, q.conn)
	if err != nil {
		return "", q.errWrap(err)
	}
	rows, err := q.conn.Query(ctx, alphaNestedSQL)
	if err != nil {
		return "", fmt.Errorf("query AlphaNested: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[string])
	return res, q.errWrap(err)
}

const alphaCompositeArraySQL = `SELECT ARRAY[ROW('key')]::alpha[];`

// AlphaCompositeArray implements Querier.AlphaCompositeArray.
func (q *DBQuerier) AlphaCompositeArray(ctx context.Context) ([]Alpha, error) {
	ctx = context.WithValue(ctx, QueryName{}, "AlphaCompositeArray")

	err := registerTypes(ctx, q.conn)
	if err != nil {
		return nil, q.errWrap(err)
	}
	rows, err := q.conn.Query(ctx, alphaCompositeArraySQL)
	if err != nil {
		return nil, fmt.Errorf("query AlphaCompositeArray: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[[]Alpha])
	return res, q.errWrap(err)
}
