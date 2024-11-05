// Code generated by pggen. DO NOT EDIT.

package composite

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
	SearchScreenshots(ctx context.Context, params SearchScreenshotsParams) ([]SearchScreenshotsRow, error)

	SearchScreenshotsOneCol(ctx context.Context, params SearchScreenshotsOneColParams) ([][]Blocks, error)

	InsertScreenshotBlocks(ctx context.Context, screenshotID int, body string) (InsertScreenshotBlocksRow, error)

	ArraysInput(ctx context.Context, arrays Arrays) (Arrays, error)

	UserEmails(ctx context.Context) (UserEmail, error)
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

// Arrays represents the Postgres composite type "arrays".
type Arrays struct {
	Texts  []string   `json:"texts"`
	Int8s  []*int     `json:"int8s"`
	Bools  []bool     `json:"bools"`
	Floats []*float64 `json:"floats"`
}

// Blocks represents the Postgres composite type "blocks".
type Blocks struct {
	ID           int    `json:"id"`
	ScreenshotID int    `json:"screenshot_id"`
	Body         string `json:"body"`
}

// UserEmail represents the Postgres composite type "user_email".
type UserEmail struct {
	ID    string      `json:"id"`
	Email pgtype.Text `json:"email"`
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

var _ = addTypeToRegister("public.arrays")

var _ = addTypeToRegister("public.blocks")

var _ = addTypeToRegister("public.user_email")

var _ = addTypeToRegister("public._blocks")

const searchScreenshotsSQL = `SELECT
  ss.id,
  array_agg(bl) AS blocks
FROM screenshots ss
  JOIN blocks bl ON bl.screenshot_id = ss.id
WHERE bl.body LIKE $1 || '%'
GROUP BY ss.id
ORDER BY ss.id
LIMIT $2 OFFSET $3;`

type SearchScreenshotsParams struct {
	Body   string `json:"Body"`
	Limit  int    `json:"Limit"`
	Offset int    `json:"Offset"`
}

type SearchScreenshotsRow struct {
	ID     int      `json:"id"`
	Blocks []Blocks `json:"blocks"`
}

// SearchScreenshots implements Querier.SearchScreenshots.
func (q *DBQuerier) SearchScreenshots(ctx context.Context, params SearchScreenshotsParams) ([]SearchScreenshotsRow, error) {
	err := registerTypes(ctx, q.conn)
	if err != nil {
		return nil, fmt.Errorf("registering types failed: %w", q.errWrap(err))
	}

	ctx = context.WithValue(ctx, QueryName{}, "SearchScreenshots")
	rows, err := q.conn.Query(ctx, searchScreenshotsSQL, params.Body, params.Limit, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("query SearchScreenshots: %w", q.errWrap(err))
	}
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[SearchScreenshotsRow])
	return res, q.errWrap(err)
}

const searchScreenshotsOneColSQL = `SELECT
  array_agg(bl) AS blocks
FROM screenshots ss
  JOIN blocks bl ON bl.screenshot_id = ss.id
WHERE bl.body LIKE $1 || '%'
GROUP BY ss.id
ORDER BY ss.id
LIMIT $2 OFFSET $3;`

type SearchScreenshotsOneColParams struct {
	Body   string `json:"Body"`
	Limit  int    `json:"Limit"`
	Offset int    `json:"Offset"`
}

// SearchScreenshotsOneCol implements Querier.SearchScreenshotsOneCol.
func (q *DBQuerier) SearchScreenshotsOneCol(ctx context.Context, params SearchScreenshotsOneColParams) ([][]Blocks, error) {
	err := registerTypes(ctx, q.conn)
	if err != nil {
		return nil, fmt.Errorf("registering types failed: %w", q.errWrap(err))
	}

	ctx = context.WithValue(ctx, QueryName{}, "SearchScreenshotsOneCol")
	rows, err := q.conn.Query(ctx, searchScreenshotsOneColSQL, params.Body, params.Limit, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("query SearchScreenshotsOneCol: %w", q.errWrap(err))
	}
	res, err := pgx.CollectRows(rows, pgx.RowTo[[]Blocks])
	return res, q.errWrap(err)
}

const insertScreenshotBlocksSQL = `WITH screens AS (
  INSERT INTO screenshots (id) VALUES ($1)
    ON CONFLICT DO NOTHING
)
INSERT
INTO blocks (screenshot_id, body)
VALUES ($1, $2)
RETURNING id, screenshot_id, body;`

type InsertScreenshotBlocksRow struct {
	ID           int    `json:"id"`
	ScreenshotID int    `json:"screenshot_id"`
	Body         string `json:"body"`
}

// InsertScreenshotBlocks implements Querier.InsertScreenshotBlocks.
func (q *DBQuerier) InsertScreenshotBlocks(ctx context.Context, screenshotID int, body string) (InsertScreenshotBlocksRow, error) {
	err := registerTypes(ctx, q.conn)
	if err != nil {
		return InsertScreenshotBlocksRow{}, fmt.Errorf("registering types failed: %w", q.errWrap(err))
	}

	ctx = context.WithValue(ctx, QueryName{}, "InsertScreenshotBlocks")
	rows, err := q.conn.Query(ctx, insertScreenshotBlocksSQL, screenshotID, body)
	if err != nil {
		return InsertScreenshotBlocksRow{}, fmt.Errorf("query InsertScreenshotBlocks: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[InsertScreenshotBlocksRow])
	return res, q.errWrap(err)
}

const arraysInputSQL = `SELECT $1::arrays;`

// ArraysInput implements Querier.ArraysInput.
func (q *DBQuerier) ArraysInput(ctx context.Context, arrays Arrays) (Arrays, error) {
	err := registerTypes(ctx, q.conn)
	if err != nil {
		return Arrays{}, fmt.Errorf("registering types failed: %w", q.errWrap(err))
	}

	ctx = context.WithValue(ctx, QueryName{}, "ArraysInput")
	rows, err := q.conn.Query(ctx, arraysInputSQL, arrays)
	if err != nil {
		return Arrays{}, fmt.Errorf("query ArraysInput: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[Arrays])
	return res, q.errWrap(err)
}

const userEmailsSQL = `SELECT ('foo', 'bar@example.com')::user_email;`

// UserEmails implements Querier.UserEmails.
func (q *DBQuerier) UserEmails(ctx context.Context) (UserEmail, error) {
	err := registerTypes(ctx, q.conn)
	if err != nil {
		return UserEmail{}, fmt.Errorf("registering types failed: %w", q.errWrap(err))
	}

	ctx = context.WithValue(ctx, QueryName{}, "UserEmails")
	rows, err := q.conn.Query(ctx, userEmailsSQL)
	if err != nil {
		return UserEmail{}, fmt.Errorf("query UserEmails: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[UserEmail])
	return res, q.errWrap(err)
}
