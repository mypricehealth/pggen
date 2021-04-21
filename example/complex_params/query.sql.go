// Code generated by pggen. DO NOT EDIT.

package complex_params

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

// Querier is a typesafe Go interface backed by SQL queries.
//
// Methods ending with Batch enqueue a query to run later in a pgx.Batch. After
// calling SendBatch on pgx.Conn, pgxpool.Pool, or pgx.Tx, use the Scan methods
// to parse the results.
type Querier interface {
	ParamArrayInt(ctx context.Context, ints []int) ([]int, error)
	// ParamArrayIntBatch enqueues a ParamArrayInt query into batch to be executed
	// later by the batch.
	ParamArrayIntBatch(batch *pgx.Batch, ints []int)
	// ParamArrayIntScan scans the result of an executed ParamArrayIntBatch query.
	ParamArrayIntScan(results pgx.BatchResults) ([]int, error)

	ParamNested1(ctx context.Context, dimensions Dimensions) (Dimensions, error)
	// ParamNested1Batch enqueues a ParamNested1 query into batch to be executed
	// later by the batch.
	ParamNested1Batch(batch *pgx.Batch, dimensions Dimensions)
	// ParamNested1Scan scans the result of an executed ParamNested1Batch query.
	ParamNested1Scan(results pgx.BatchResults) (Dimensions, error)

	ParamNested2(ctx context.Context, image ProductImageType) (ProductImageType, error)
	// ParamNested2Batch enqueues a ParamNested2 query into batch to be executed
	// later by the batch.
	ParamNested2Batch(batch *pgx.Batch, image ProductImageType)
	// ParamNested2Scan scans the result of an executed ParamNested2Batch query.
	ParamNested2Scan(results pgx.BatchResults) (ProductImageType, error)

	ParamNested2Array(ctx context.Context, images []ProductImageType) ([]ProductImageType, error)
	// ParamNested2ArrayBatch enqueues a ParamNested2Array query into batch to be executed
	// later by the batch.
	ParamNested2ArrayBatch(batch *pgx.Batch, images []ProductImageType)
	// ParamNested2ArrayScan scans the result of an executed ParamNested2ArrayBatch query.
	ParamNested2ArrayScan(results pgx.BatchResults) ([]ProductImageType, error)

	ParamNested3(ctx context.Context, imageSet ProductImageSetType) (ProductImageSetType, error)
	// ParamNested3Batch enqueues a ParamNested3 query into batch to be executed
	// later by the batch.
	ParamNested3Batch(batch *pgx.Batch, imageSet ProductImageSetType)
	// ParamNested3Scan scans the result of an executed ParamNested3Batch query.
	ParamNested3Scan(results pgx.BatchResults) (ProductImageSetType, error)
}

type DBQuerier struct {
	conn genericConn
}

var _ Querier = &DBQuerier{}

// genericConn is a connection to a Postgres database. This is usually backed by
// *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
type genericConn interface {
	// Query executes sql with args. If there is an error the returned Rows will
	// be returned in an error state. So it is allowed to ignore the error
	// returned from Query and handle it in Rows.
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)

	// QueryRow is a convenience wrapper over Query. Any error that occurs while
	// querying is deferred until calling Scan on the returned Row. That Row will
	// error with pgx.ErrNoRows if no rows are returned.
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row

	// Exec executes sql. sql can be either a prepared statement name or an SQL
	// string. arguments should be referenced positionally from the sql string
	// as $1, $2, etc.
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

// NewQuerier creates a DBQuerier that implements Querier. conn is typically
// *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
func NewQuerier(conn genericConn) *DBQuerier {
	return &DBQuerier{
		conn: conn,
	}
}

// WithTx creates a new DBQuerier that uses the transaction to run all queries.
func (q *DBQuerier) WithTx(tx pgx.Tx) (*DBQuerier, error) {
	return &DBQuerier{conn: tx}, nil
}

// preparer is any Postgres connection transport that provides a way to prepare
// a statement, most commonly *pgx.Conn.
type preparer interface {
	Prepare(ctx context.Context, name, sql string) (sd *pgconn.StatementDescription, err error)
}

// PrepareAllQueries executes a PREPARE statement for all pggen generated SQL
// queries in querier files. Typical usage is as the AfterConnect callback
// for pgxpool.Config
//
// pgx will use the prepared statement if available. Calling PrepareAllQueries
// is an optional optimization to avoid a network round-trip the first time pgx
// runs a query if pgx statement caching is enabled.
func PrepareAllQueries(ctx context.Context, p preparer) error {
	if _, err := p.Prepare(ctx, paramArrayIntSQL, paramArrayIntSQL); err != nil {
		return fmt.Errorf("prepare query 'ParamArrayInt': %w", err)
	}
	if _, err := p.Prepare(ctx, paramNested1SQL, paramNested1SQL); err != nil {
		return fmt.Errorf("prepare query 'ParamNested1': %w", err)
	}
	if _, err := p.Prepare(ctx, paramNested2SQL, paramNested2SQL); err != nil {
		return fmt.Errorf("prepare query 'ParamNested2': %w", err)
	}
	if _, err := p.Prepare(ctx, paramNested2ArraySQL, paramNested2ArraySQL); err != nil {
		return fmt.Errorf("prepare query 'ParamNested2Array': %w", err)
	}
	if _, err := p.Prepare(ctx, paramNested3SQL, paramNested3SQL); err != nil {
		return fmt.Errorf("prepare query 'ParamNested3': %w", err)
	}
	return nil
}

// assignProductImageTypeArray returns all elements for the Postgres '_product_image_type' array type as a
// slice of interface{} for use with the pgtype.Value Set method.
func assignProductImageTypeArray(ps []ProductImageType) []interface{} {
	elems := make([]interface{}, len(ps))
	for i, p := range ps {
		elems[i] = assignProductImageTypeComposite(p)
	}
	return elems
}

// newProductImageTypeArrayDecoder creates a new decoder for the Postgres '_product_image_type' array type.
func newProductImageTypeArrayDecoder() pgtype.ValueTranscoder {
	return pgtype.NewArrayType("_product_image_type", ignoredOID, newProductImageTypeDecoder)
}

// newProductImageTypeArrayEncoder creates a new encoder for the Postgres '_product_image_type' array type query params.
func newProductImageTypeArrayEncoder(ps []ProductImageType) textEncoder {
	dec := newProductImageTypeArrayDecoder()
	if err := dec.Set(assignProductImageTypeArray(ps)); err != nil {
		panic("encode []ProductImageType: " + err.Error()) // should always succeed
	}
	return textEncoder{ValueTranscoder: dec}
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

// assignDimensionsComposite returns all composite fields for the Postgres
// 'dimensions' composite type as a slice of interface{} for use with the
// pgtype.Value Set method.
func assignDimensionsComposite(p Dimensions) []interface{} {
	return []interface{}{
		p.Width,
		p.Height,
	}
}

// assignProductImageSetTypeComposite returns all composite fields for the Postgres
// 'product_image_set_type' composite type as a slice of interface{} for use with the
// pgtype.Value Set method.
func assignProductImageSetTypeComposite(p ProductImageSetType) []interface{} {
	return []interface{}{
		p.Name,
		assignProductImageTypeComposite(p.OrigImage),
		assignProductImageTypeArray(p.Images),
	}
}

// assignProductImageTypeComposite returns all composite fields for the Postgres
// 'product_image_type' composite type as a slice of interface{} for use with the
// pgtype.Value Set method.
func assignProductImageTypeComposite(p ProductImageType) []interface{} {
	return []interface{}{
		p.Source,
		assignDimensionsComposite(p.Dimensions),
	}
}

// newDimensionsDecoder creates a new decoder for the Postgres 'dimensions' composite type.
func newDimensionsDecoder() pgtype.ValueTranscoder {
	return newCompositeType(
		"dimensions",
		[]string{"width", "height"},
		&pgtype.Int4{},
		&pgtype.Int4{},
	)
}

// newProductImageSetTypeDecoder creates a new decoder for the Postgres 'product_image_set_type' composite type.
func newProductImageSetTypeDecoder() pgtype.ValueTranscoder {
	return newCompositeType(
		"product_image_set_type",
		[]string{"name", "orig_image", "images"},
		&pgtype.Text{},
		newProductImageTypeDecoder(),
		newProductImageTypeArrayDecoder(),
	)
}

// newProductImageTypeDecoder creates a new decoder for the Postgres 'product_image_type' composite type.
func newProductImageTypeDecoder() pgtype.ValueTranscoder {
	return newCompositeType(
		"product_image_type",
		[]string{"source", "dimensions"},
		&pgtype.Text{},
		newDimensionsDecoder(),
	)
}

// newDimensionsEncoder creates a new encoder for the Postgres 'dimensions' composite type query params.
func newDimensionsEncoder(p Dimensions) textEncoder {
	dec := newDimensionsDecoder()
	if err := dec.Set(assignDimensionsComposite(p)); err != nil {
		panic("encode Dimensions: " + err.Error()) // should always succeed
	}
	return textEncoder{ValueTranscoder: dec}
}

// newProductImageSetTypeEncoder creates a new encoder for the Postgres 'product_image_set_type' composite type query params.
func newProductImageSetTypeEncoder(p ProductImageSetType) textEncoder {
	dec := newProductImageSetTypeDecoder()
	if err := dec.Set(assignProductImageSetTypeComposite(p)); err != nil {
		panic("encode ProductImageSetType: " + err.Error()) // should always succeed
	}
	return textEncoder{ValueTranscoder: dec}
}

// newProductImageTypeEncoder creates a new encoder for the Postgres 'product_image_type' composite type query params.
func newProductImageTypeEncoder(p ProductImageType) textEncoder {
	dec := newProductImageTypeDecoder()
	if err := dec.Set(assignProductImageTypeComposite(p)); err != nil {
		panic("encode ProductImageType: " + err.Error()) // should always succeed
	}
	return textEncoder{ValueTranscoder: dec}
}

// ignoredOID means we don't know or care about the OID for a type. This is okay
// because pgx only uses the OID to encode values and lookup a decoder. We only
// use ignoredOID for decoding and we always specify a concrete decoder for scan
// methods.
const ignoredOID = 0

// textEncoder wraps a pgtype.ValueTranscoder and sets the preferred encoding
// format to text instead binary (the default). pggen must use the text format
// because the Postgres binary format requires the type OID but pggen doesn't
// necessarily know the OIDs of the types, hence ignoredOID.
type textEncoder struct {
	pgtype.ValueTranscoder
}

// PreferredParamFormat implements pgtype.ParamFormatPreferrer.
func (t textEncoder) PreferredParamFormat() int16 { return pgtype.TextFormatCode }

func newCompositeType(name string, fieldNames []string, vals ...pgtype.ValueTranscoder) *pgtype.CompositeType {
	fields := make([]pgtype.CompositeTypeField, len(fieldNames))
	for i, name := range fieldNames {
		fields[i] = pgtype.CompositeTypeField{Name: name, OID: ignoredOID}
	}
	// Okay to ignore error because it's only thrown when the number of field
	// names does not equal the number of ValueTranscoders.
	rowType, _ := pgtype.NewCompositeTypeValues(name, fields, vals)
	return rowType
}

const paramArrayIntSQL = `SELECT $1::bigint[];`

// ParamArrayInt implements Querier.ParamArrayInt.
func (q *DBQuerier) ParamArrayInt(ctx context.Context, ints []int) ([]int, error) {
	row := q.conn.QueryRow(ctx, paramArrayIntSQL, ints)
	item := []int{}
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query ParamArrayInt: %w", err)
	}
	return item, nil
}

// ParamArrayIntBatch implements Querier.ParamArrayIntBatch.
func (q *DBQuerier) ParamArrayIntBatch(batch *pgx.Batch, ints []int) {
	batch.Queue(paramArrayIntSQL, ints)
}

// ParamArrayIntScan implements Querier.ParamArrayIntScan.
func (q *DBQuerier) ParamArrayIntScan(results pgx.BatchResults) ([]int, error) {
	row := results.QueryRow()
	item := []int{}
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan ParamArrayIntBatch row: %w", err)
	}
	return item, nil
}

const paramNested1SQL = `SELECT $1::dimensions;`

// ParamNested1 implements Querier.ParamNested1.
func (q *DBQuerier) ParamNested1(ctx context.Context, dimensions Dimensions) (Dimensions, error) {
	row := q.conn.QueryRow(ctx, paramNested1SQL, newDimensionsEncoder(dimensions))
	var item Dimensions
	dimensionsRow := newDimensionsDecoder()
	if err := row.Scan(dimensionsRow); err != nil {
		return item, fmt.Errorf("query ParamNested1: %w", err)
	}
	if err := dimensionsRow.AssignTo(&item); err != nil {
		return item, fmt.Errorf("assign ParamNested1 row: %w", err)
	}
	return item, nil
}

// ParamNested1Batch implements Querier.ParamNested1Batch.
func (q *DBQuerier) ParamNested1Batch(batch *pgx.Batch, dimensions Dimensions) {
	batch.Queue(paramNested1SQL, newDimensionsEncoder(dimensions))
}

// ParamNested1Scan implements Querier.ParamNested1Scan.
func (q *DBQuerier) ParamNested1Scan(results pgx.BatchResults) (Dimensions, error) {
	row := results.QueryRow()
	var item Dimensions
	dimensionsRow := newDimensionsDecoder()
	if err := row.Scan(dimensionsRow); err != nil {
		return item, fmt.Errorf("scan ParamNested1Batch row: %w", err)
	}
	if err := dimensionsRow.AssignTo(&item); err != nil {
		return item, fmt.Errorf("assign ParamNested1 row: %w", err)
	}
	return item, nil
}

const paramNested2SQL = `SELECT $1::product_image_type;`

// ParamNested2 implements Querier.ParamNested2.
func (q *DBQuerier) ParamNested2(ctx context.Context, image ProductImageType) (ProductImageType, error) {
	row := q.conn.QueryRow(ctx, paramNested2SQL, newProductImageTypeEncoder(image))
	var item ProductImageType
	productImageTypeRow := newProductImageTypeDecoder()
	if err := row.Scan(productImageTypeRow); err != nil {
		return item, fmt.Errorf("query ParamNested2: %w", err)
	}
	if err := productImageTypeRow.AssignTo(&item); err != nil {
		return item, fmt.Errorf("assign ParamNested2 row: %w", err)
	}
	return item, nil
}

// ParamNested2Batch implements Querier.ParamNested2Batch.
func (q *DBQuerier) ParamNested2Batch(batch *pgx.Batch, image ProductImageType) {
	batch.Queue(paramNested2SQL, newProductImageTypeEncoder(image))
}

// ParamNested2Scan implements Querier.ParamNested2Scan.
func (q *DBQuerier) ParamNested2Scan(results pgx.BatchResults) (ProductImageType, error) {
	row := results.QueryRow()
	var item ProductImageType
	productImageTypeRow := newProductImageTypeDecoder()
	if err := row.Scan(productImageTypeRow); err != nil {
		return item, fmt.Errorf("scan ParamNested2Batch row: %w", err)
	}
	if err := productImageTypeRow.AssignTo(&item); err != nil {
		return item, fmt.Errorf("assign ParamNested2 row: %w", err)
	}
	return item, nil
}

const paramNested2ArraySQL = `SELECT $1::product_image_type[];`

// ParamNested2Array implements Querier.ParamNested2Array.
func (q *DBQuerier) ParamNested2Array(ctx context.Context, images []ProductImageType) ([]ProductImageType, error) {
	row := q.conn.QueryRow(ctx, paramNested2ArraySQL, newProductImageTypeArrayEncoder(images))
	item := []ProductImageType{}
	productImageTypeArray := newProductImageTypeArrayDecoder()
	if err := row.Scan(productImageTypeArray); err != nil {
		return item, fmt.Errorf("query ParamNested2Array: %w", err)
	}
	if err := productImageTypeArray.AssignTo(&item); err != nil {
		return item, fmt.Errorf("assign ParamNested2Array row: %w", err)
	}
	return item, nil
}

// ParamNested2ArrayBatch implements Querier.ParamNested2ArrayBatch.
func (q *DBQuerier) ParamNested2ArrayBatch(batch *pgx.Batch, images []ProductImageType) {
	batch.Queue(paramNested2ArraySQL, newProductImageTypeArrayEncoder(images))
}

// ParamNested2ArrayScan implements Querier.ParamNested2ArrayScan.
func (q *DBQuerier) ParamNested2ArrayScan(results pgx.BatchResults) ([]ProductImageType, error) {
	row := results.QueryRow()
	item := []ProductImageType{}
	productImageTypeArray := newProductImageTypeArrayDecoder()
	if err := row.Scan(productImageTypeArray); err != nil {
		return item, fmt.Errorf("scan ParamNested2ArrayBatch row: %w", err)
	}
	if err := productImageTypeArray.AssignTo(&item); err != nil {
		return item, fmt.Errorf("assign ParamNested2Array row: %w", err)
	}
	return item, nil
}

const paramNested3SQL = `SELECT $1::product_image_set_type;`

// ParamNested3 implements Querier.ParamNested3.
func (q *DBQuerier) ParamNested3(ctx context.Context, imageSet ProductImageSetType) (ProductImageSetType, error) {
	row := q.conn.QueryRow(ctx, paramNested3SQL, newProductImageSetTypeEncoder(imageSet))
	var item ProductImageSetType
	productImageSetTypeRow := newProductImageSetTypeDecoder()
	if err := row.Scan(productImageSetTypeRow); err != nil {
		return item, fmt.Errorf("query ParamNested3: %w", err)
	}
	if err := productImageSetTypeRow.AssignTo(&item); err != nil {
		return item, fmt.Errorf("assign ParamNested3 row: %w", err)
	}
	return item, nil
}

// ParamNested3Batch implements Querier.ParamNested3Batch.
func (q *DBQuerier) ParamNested3Batch(batch *pgx.Batch, imageSet ProductImageSetType) {
	batch.Queue(paramNested3SQL, newProductImageSetTypeEncoder(imageSet))
}

// ParamNested3Scan implements Querier.ParamNested3Scan.
func (q *DBQuerier) ParamNested3Scan(results pgx.BatchResults) (ProductImageSetType, error) {
	row := results.QueryRow()
	var item ProductImageSetType
	productImageSetTypeRow := newProductImageSetTypeDecoder()
	if err := row.Scan(productImageSetTypeRow); err != nil {
		return item, fmt.Errorf("scan ParamNested3Batch row: %w", err)
	}
	if err := productImageSetTypeRow.AssignTo(&item); err != nil {
		return item, fmt.Errorf("assign ParamNested3 row: %w", err)
	}
	return item, nil
}
