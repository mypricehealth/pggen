// Code generated by pggen. DO NOT EDIT.

package complex_params

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type QueryName struct{}

// Querier is a typesafe Go interface backed by SQL queries.
type Querier interface {
	ParamArrayInt(ctx context.Context, ints []int) ([]int, error)

	ParamNested1(ctx context.Context, dimensions Dimensions) (Dimensions, error)

	ParamNested2(ctx context.Context, image ProductImageType) (ProductImageType, error)

	ParamNested2Array(ctx context.Context, images []ProductImageType) ([]ProductImageType, error)

	ParamNested3(ctx context.Context, imageSet ProductImageSetType) (ProductImageSetType, error)
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

// RegisterTypes should be run in config.AfterConnect to load custom types
func RegisterTypes(ctx context.Context, conn *pgx.Conn) error {
	for _, typ := range typesToRegister {
		dt, err := conn.LoadType(ctx, typ)
		if err != nil {
			return err
		}
		conn.TypeMap().RegisterType(dt)
	}
	return nil
}

var typesToRegister = []string{}

func addTypeToRegister(typ string) struct{} {
	typesToRegister = append(typesToRegister, typ)
	return struct{}{}
}

var _ = addTypeToRegister("dimensions")

var _ = addTypeToRegister("product_image_set_type")

var _ = addTypeToRegister("product_image_type")

var _ = addTypeToRegister("_product_image_type")

const paramArrayIntSQL = `SELECT $1::bigint[];`

// ParamArrayInt implements Querier.ParamArrayInt.
func (q *DBQuerier) ParamArrayInt(ctx context.Context, ints []int) ([]int, error) {
	ctx = context.WithValue(ctx, QueryName{}, "ParamArrayInt")
	rows, err := q.conn.Query(ctx, paramArrayIntSQL, ints)
	if err != nil {
		return nil, fmt.Errorf("query ParamArrayInt: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[[]int])
	return res, q.errWrap(err)
}

const paramNested1SQL = `SELECT $1::dimensions;`

// ParamNested1 implements Querier.ParamNested1.
func (q *DBQuerier) ParamNested1(ctx context.Context, dimensions Dimensions) (Dimensions, error) {
	ctx = context.WithValue(ctx, QueryName{}, "ParamNested1")
	rows, err := q.conn.Query(ctx, paramNested1SQL, dimensions)
	if err != nil {
		return Dimensions{}, fmt.Errorf("query ParamNested1: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[Dimensions])
	return res, q.errWrap(err)
}

const paramNested2SQL = `SELECT $1::product_image_type;`

// ParamNested2 implements Querier.ParamNested2.
func (q *DBQuerier) ParamNested2(ctx context.Context, image ProductImageType) (ProductImageType, error) {
	ctx = context.WithValue(ctx, QueryName{}, "ParamNested2")
	rows, err := q.conn.Query(ctx, paramNested2SQL, image)
	if err != nil {
		return ProductImageType{}, fmt.Errorf("query ParamNested2: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[ProductImageType])
	return res, q.errWrap(err)
}

const paramNested2ArraySQL = `SELECT $1::product_image_type[];`

// ParamNested2Array implements Querier.ParamNested2Array.
func (q *DBQuerier) ParamNested2Array(ctx context.Context, images []ProductImageType) ([]ProductImageType, error) {
	ctx = context.WithValue(ctx, QueryName{}, "ParamNested2Array")
	rows, err := q.conn.Query(ctx, paramNested2ArraySQL, images)
	if err != nil {
		return nil, fmt.Errorf("query ParamNested2Array: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[[]ProductImageType])
	return res, q.errWrap(err)
}

const paramNested3SQL = `SELECT $1::product_image_set_type;`

// ParamNested3 implements Querier.ParamNested3.
func (q *DBQuerier) ParamNested3(ctx context.Context, imageSet ProductImageSetType) (ProductImageSetType, error) {
	ctx = context.WithValue(ctx, QueryName{}, "ParamNested3")
	rows, err := q.conn.Query(ctx, paramNested3SQL, imageSet)
	if err != nil {
		return ProductImageSetType{}, fmt.Errorf("query ParamNested3: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[ProductImageSetType])
	return res, q.errWrap(err)
}
