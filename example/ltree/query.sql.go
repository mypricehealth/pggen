// Code generated by pggen. DO NOT EDIT.

package ltree

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
	FindTopScienceChildren(ctx context.Context) ([]pgtype.Text, error)

	FindTopScienceChildrenAgg(ctx context.Context) (pgtype.TextArray, error)

	InsertSampleData(ctx context.Context) (pgconn.CommandTag, error)

	FindLtreeInput(ctx context.Context, inLtree pgtype.Text, inLtreeArray []string) (FindLtreeInputRow, error)
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
				registerErr = err
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

const findTopScienceChildrenSQL = `SELECT path
FROM test
WHERE path <@ 'Top.Science';`

// FindTopScienceChildren implements Querier.FindTopScienceChildren.
func (q *DBQuerier) FindTopScienceChildren(ctx context.Context) ([]pgtype.Text, error) {
	err := registerTypes(ctx, q.conn)
	if err != nil {
		return nil, fmt.Errorf("registering types failed: %w", q.errWrap(err))
	}

	ctx = context.WithValue(ctx, QueryName{}, "FindTopScienceChildren")
	rows, err := q.conn.Query(ctx, findTopScienceChildrenSQL)
	if err != nil {
		return nil, fmt.Errorf("query FindTopScienceChildren: %w", q.errWrap(err))
	}
	res, err := pgx.CollectRows(rows, pgx.RowTo[pgtype.Text])
	return res, q.errWrap(err)
}

const findTopScienceChildrenAggSQL = `SELECT array_agg(path)
FROM test
WHERE path <@ 'Top.Science';`

// FindTopScienceChildrenAgg implements Querier.FindTopScienceChildrenAgg.
func (q *DBQuerier) FindTopScienceChildrenAgg(ctx context.Context) (pgtype.TextArray, error) {
	err := registerTypes(ctx, q.conn)
	if err != nil {
		return TextArray{}, fmt.Errorf("registering types failed: %w", q.errWrap(err))
	}

	ctx = context.WithValue(ctx, QueryName{}, "FindTopScienceChildrenAgg")
	rows, err := q.conn.Query(ctx, findTopScienceChildrenAggSQL)
	if err != nil {
		return TextArray{}, fmt.Errorf("query FindTopScienceChildrenAgg: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[pgtype.TextArray])
	return res, q.errWrap(err)
}

const insertSampleDataSQL = `INSERT INTO test
VALUES ('Top'),
       ('Top.Science'),
       ('Top.Science.Astronomy'),
       ('Top.Science.Astronomy.Astrophysics'),
       ('Top.Science.Astronomy.Cosmology'),
       ('Top.Hobbies'),
       ('Top.Hobbies.Amateurs_Astronomy'),
       ('Top.Collections'),
       ('Top.Collections.Pictures'),
       ('Top.Collections.Pictures.Astronomy'),
       ('Top.Collections.Pictures.Astronomy.Stars'),
       ('Top.Collections.Pictures.Astronomy.Galaxies'),
       ('Top.Collections.Pictures.Astronomy.Astronauts');`

// InsertSampleData implements Querier.InsertSampleData.
func (q *DBQuerier) InsertSampleData(ctx context.Context) (pgconn.CommandTag, error) {
	err := registerTypes(ctx, q.conn)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("registering types failed: %w", q.errWrap(err))
	}

	ctx = context.WithValue(ctx, QueryName{}, "InsertSampleData")
	cmdTag, err := q.conn.Exec(ctx, insertSampleDataSQL)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query InsertSampleData: %w", q.errWrap(err))
	}
	return cmdTag, q.errWrap(err)
}

const findLtreeInputSQL = `SELECT
  $1::ltree                   AS ltree,
  -- This won't work, but I'm not quite sure why.
  -- Postgres errors with "wrong element type (SQLSTATE 42804)"
  -- All caps because we use regex to find pggen.arg and it confuses pggen.
  -- PGGEN.arg('in_ltree_array_direct')::ltree[]    AS direct_arr,

  -- The parenthesis around the text[] cast are important. They signal to pggen
  -- that we need a text array that Postgres then converts to ltree[].
  ($2::text[])::ltree[] AS text_arr;`

type FindLtreeInputRow struct {
	Ltree   pgtype.Text      `json:"ltree"`
	TextArr pgtype.TextArray `json:"text_arr"`
}

// FindLtreeInput implements Querier.FindLtreeInput.
func (q *DBQuerier) FindLtreeInput(ctx context.Context, inLtree pgtype.Text, inLtreeArray []string) (FindLtreeInputRow, error) {
	err := registerTypes(ctx, q.conn)
	if err != nil {
		return FindLtreeInputRow{}, fmt.Errorf("registering types failed: %w", q.errWrap(err))
	}

	ctx = context.WithValue(ctx, QueryName{}, "FindLtreeInput")
	rows, err := q.conn.Query(ctx, findLtreeInputSQL, inLtree, inLtreeArray)
	if err != nil {
		return FindLtreeInputRow{}, fmt.Errorf("query FindLtreeInput: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[FindLtreeInputRow])
	return res, q.errWrap(err)
}
