// Code generated by pggen. DO NOT EDIT.

package enums

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type QueryName struct{}

// Querier is a typesafe Go interface backed by SQL queries.
type Querier interface {
	FindAllDevices(ctx context.Context) ([]FindAllDevicesRow, error)

	InsertDevice(ctx context.Context, mac pgtype.Macaddr, typePg DeviceType) (pgconn.CommandTag, error)

	// Select an array of all device_type enum values.
	FindOneDeviceArray(ctx context.Context) ([]DeviceType, error)

	// Select many rows of device_type enum values.
	FindManyDeviceArray(ctx context.Context) ([][]DeviceType, error)

	// Select many rows of device_type enum values with multiple output columns.
	FindManyDeviceArrayWithNum(ctx context.Context) ([]FindManyDeviceArrayWithNumRow, error)

	// Regression test for https://github.com/jschaf/pggen/issues/23.
	EnumInsideComposite(ctx context.Context) (Device, error)
}

var _ Querier = &DBQuerier{}

type DBQuerier struct {
	conn genericConn
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

// Device represents the Postgres composite type "device".
type Device struct {
	Mac  pgtype.Macaddr `json:"mac"`
	Type DeviceType     `json:"type"`
}

// newDeviceTypeEnum creates a new pgtype.ValueTranscoder for the
// Postgres enum type 'device_type'.
func newDeviceTypeEnum() pgtype.ValueTranscoder {
	return pgtype.NewEnumType(
		"device_type",
		[]string{
			string(DeviceTypeUndefined),
			string(DeviceTypePhone),
			string(DeviceTypeLaptop),
			string(DeviceTypeIpad),
			string(DeviceTypeDesktop),
			string(DeviceTypeIot),
		},
	)
}

// DeviceType represents the Postgres enum "device_type".
type DeviceType string

const (
	DeviceTypeUndefined DeviceType = "undefined"
	DeviceTypePhone     DeviceType = "phone"
	DeviceTypeLaptop    DeviceType = "laptop"
	DeviceTypeIpad      DeviceType = "ipad"
	DeviceTypeDesktop   DeviceType = "desktop"
	DeviceTypeIot       DeviceType = "iot"
)

func (d DeviceType) String() string { return string(d) }

const findAllDevicesSQL = `SELECT mac, type
FROM device;`

type FindAllDevicesRow struct {
	Mac  pgtype.Macaddr `json:"mac"`
	Type DeviceType     `json:"type"`
}

// FindAllDevices implements Querier.FindAllDevices.
func (q *DBQuerier) FindAllDevices(ctx context.Context) ([]FindAllDevicesRow, error) {
	ctx = context.WithValue(ctx, QueryName{}, "FindAllDevices")
	rows, err := q.conn.Query(ctx, findAllDevicesSQL)
	if err != nil {
		return nil, fmt.Errorf("query FindAllDevices: %w", err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[FindAllDevicesRow])
}

const insertDeviceSQL = `INSERT INTO device (mac, type)
VALUES ($1, $2);`

// InsertDevice implements Querier.InsertDevice.
func (q *DBQuerier) InsertDevice(ctx context.Context, mac pgtype.Macaddr, typePg DeviceType) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, QueryName{}, "InsertDevice")
	cmdTag, err := q.conn.Exec(ctx, insertDeviceSQL, mac, typePg)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query InsertDevice: %w", err)
	}
	return cmdTag, err
}

const findOneDeviceArraySQL = `SELECT enum_range(NULL::device_type) AS device_types;`

// FindOneDeviceArray implements Querier.FindOneDeviceArray.
func (q *DBQuerier) FindOneDeviceArray(ctx context.Context) ([]DeviceType, error) {
	ctx = context.WithValue(ctx, QueryName{}, "FindOneDeviceArray")
	rows, err := q.conn.Query(ctx, findOneDeviceArraySQL)
	if err != nil {
		return nil, fmt.Errorf("query FindOneDeviceArray: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowTo[[]DeviceType])
}

const findManyDeviceArraySQL = `SELECT enum_range('ipad'::device_type, 'iot'::device_type) AS device_types
UNION ALL
SELECT enum_range(NULL::device_type) AS device_types;`

// FindManyDeviceArray implements Querier.FindManyDeviceArray.
func (q *DBQuerier) FindManyDeviceArray(ctx context.Context) ([][]DeviceType, error) {
	ctx = context.WithValue(ctx, QueryName{}, "FindManyDeviceArray")
	rows, err := q.conn.Query(ctx, findManyDeviceArraySQL)
	if err != nil {
		return nil, fmt.Errorf("query FindManyDeviceArray: %w", err)
	}

	return pgx.CollectRows(rows, pgx.RowTo[[]DeviceType])
}

const findManyDeviceArrayWithNumSQL = `SELECT 1 AS num, enum_range('ipad'::device_type, 'iot'::device_type) AS device_types
UNION ALL
SELECT 2 as num, enum_range(NULL::device_type) AS device_types;`

type FindManyDeviceArrayWithNumRow struct {
	Num         *int32       `json:"num"`
	DeviceTypes []DeviceType `json:"device_types"`
}

// FindManyDeviceArrayWithNum implements Querier.FindManyDeviceArrayWithNum.
func (q *DBQuerier) FindManyDeviceArrayWithNum(ctx context.Context) ([]FindManyDeviceArrayWithNumRow, error) {
	ctx = context.WithValue(ctx, QueryName{}, "FindManyDeviceArrayWithNum")
	rows, err := q.conn.Query(ctx, findManyDeviceArrayWithNumSQL)
	if err != nil {
		return nil, fmt.Errorf("query FindManyDeviceArrayWithNum: %w", err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[FindManyDeviceArrayWithNumRow])
}

const enumInsideCompositeSQL = `SELECT ROW('08:00:2b:01:02:03'::macaddr, 'phone'::device_type) ::device;`

// EnumInsideComposite implements Querier.EnumInsideComposite.
func (q *DBQuerier) EnumInsideComposite(ctx context.Context) (Device, error) {
	ctx = context.WithValue(ctx, QueryName{}, "EnumInsideComposite")
	rows, err := q.conn.Query(ctx, enumInsideCompositeSQL)
	if err != nil {
		return Device{}, fmt.Errorf("query EnumInsideComposite: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowTo[Device])
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
