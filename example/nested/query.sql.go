// Code generated by pggen. DO NOT EDIT.

package nested

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
	ArrayNested2(ctx context.Context) ([]ProductImageType, error)

	Nested3(ctx context.Context) ([]ProductImageSetType, error)
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

const arrayNested2SQL = `SELECT
  ARRAY [
    ROW ('img2', ROW (22, 22)::dimensions)::product_image_type,
    ROW ('img3', ROW (33, 33)::dimensions)::product_image_type
    ] AS images;`

// ArrayNested2 implements Querier.ArrayNested2.
func (q *DBQuerier) ArrayNested2(ctx context.Context) ([]ProductImageType, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "ArrayNested2")
	rows, err := q.conn.Query(ctx, arrayNested2SQL)
	if err != nil {
		return nil, fmt.Errorf("query ArrayNested2: %w", err)
	}
	fds := rows.FieldDescriptions()
	plan0 := planScan(pgtype.TextCodec{}, fds[0], (*[]ProductImageType)(nil))

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) ([]ProductImageType, error) {
		vals := row.RawValues()
		var item []ProductImageType
		if err := plan0.Scan(vals[0], &item.Images); err != nil {
			return item, fmt.Errorf("scan ArrayNested2.images: %w", err)
		}
		return item, nil
	})
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
	ctx = context.WithValue(ctx, "pggen_query_name", "Nested3")
	rows, err := q.conn.Query(ctx, nested3SQL)
	if err != nil {
		return nil, fmt.Errorf("query Nested3: %w", err)
	}
	fds := rows.FieldDescriptions()
	plan0 := planScan(pgtype.TextCodec{}, fds[0], (*ProductImageSetType)(nil))

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (ProductImageSetType, error) {
		vals := row.RawValues()
		var item ProductImageSetType
		if err := plan0.Scan(vals[0], &item); err != nil {
			return item, fmt.Errorf("scan Nested3.row: %w", err)
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