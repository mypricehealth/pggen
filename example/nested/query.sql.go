// Code generated by pggen. DO NOT EDIT.

package nested

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
	ArrayNested2(ctx context.Context) ([]ProductImageType, error)

	Nested3(ctx context.Context) ([]ProductImageSetType, error)
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

// Dimensions represents the Postgres composite type "dimensions".
type Dimensions struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// ProductImageSetType represents the Postgres composite type "product_image_set_type".
type ProductImageSetType struct {
	Name      string             `json:"name"`
	OrigImage ProductImageType   `json:"orig_image"`
	Images    []ProductImageType `json:"images"`
}

// ProductImageType represents the Postgres composite type "product_image_type".
type ProductImageType struct {
	Source     string     `json:"source"`
	Dimensions Dimensions `json:"dimensions"`
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

var _ = addTypeToRegister("public.dimensions")

var _ = addTypeToRegister("public.product_image_set_type")

var _ = addTypeToRegister("public.product_image_type")

var _ = addTypeToRegister("public._product_image_type")

const arrayNested2SQL = `SELECT
  ARRAY [
    ROW ('img2', ROW (22, 22)::dimensions)::product_image_type,
    ROW ('img3', ROW (33, 33)::dimensions)::product_image_type
    ] AS images;`

// ArrayNested2 implements Querier.ArrayNested2.
func (q *DBQuerier) ArrayNested2(ctx context.Context) ([]ProductImageType, error) {
	ctx = context.WithValue(ctx, QueryName{}, "ArrayNested2")

	err := registerTypes(ctx, q.conn)
	if err != nil {
		return nil, q.errWrap(err)
	}
	rows, err := q.conn.Query(ctx, arrayNested2SQL)
	if err != nil {
		return nil, fmt.Errorf("query ArrayNested2: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[[]ProductImageType])
	return res, q.errWrap(err)
}

const nested3SQL = `SELECT
  ROW (
    'name', -- name
    ROW ('img1', ROW (11, 11)::dimensions)::product_image_type, -- orig_image
    ARRAY [ --images
      ROW ('img2', ROW (22, 22)::dimensions)::product_image_type,
      ROW ('img3', ROW (33, 33)::dimensions)::product_image_type
      ]
    )::product_image_set_type;`

// Nested3 implements Querier.Nested3.
func (q *DBQuerier) Nested3(ctx context.Context) ([]ProductImageSetType, error) {
	ctx = context.WithValue(ctx, QueryName{}, "Nested3")

	err := registerTypes(ctx, q.conn)
	if err != nil {
		return nil, q.errWrap(err)
	}
	rows, err := q.conn.Query(ctx, nested3SQL)
	if err != nil {
		return nil, fmt.Errorf("query Nested3: %w", q.errWrap(err))
	}
	res, err := pgx.CollectRows(rows, pgx.RowTo[ProductImageSetType])
	return res, q.errWrap(err)
}
