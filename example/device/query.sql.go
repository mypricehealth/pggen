// Code generated by pggen. DO NOT EDIT.

package device

import (
	"sync"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgtype.FlatArray[github.com/jackc/pgx/v5/pgtype"
)

// Querier is a typesafe Go interface backed by SQL queries.
type Querier interface {
	FindDevicesByUser(ctx context.Context, id int) ([]FindDevicesByUserRow, error)

	CompositeUser(ctx context.Context) ([]CompositeUserRow, error)

	CompositeUserOne(ctx context.Context) (User, error)

	CompositeUserOneTwoCols(ctx context.Context) (CompositeUserOneTwoColsRow, error)

	CompositeUserMany(ctx context.Context) ([]User, error)

	InsertUser(ctx context.Context, userID int, name string) (pgconn.CommandTag, error)

	InsertDevice(ctx context.Context, mac pgtype.Macaddr, owner int) (pgconn.CommandTag, error)
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

// User represents the Postgres composite type "user".
type User struct {
	ID   *int    `json:"id"`
	Name *string `json:"name"`
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

const findDevicesByUserSQL = `SELECT
  id,
  name,
  (SELECT array_agg(mac) FROM device WHERE owner = id) AS mac_addrs
FROM "user"
WHERE id = $1;`

type FindDevicesByUserRow struct {
	ID       int                  `json:"id"`
	Name     string               `json:"name"`
	MacAddrs pgtype.MacaddrCodec] `json:"mac_addrs"`
}

// FindDevicesByUser implements Querier.FindDevicesByUser.
func (q *DBQuerier) FindDevicesByUser(ctx context.Context, id int) ([]FindDevicesByUserRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindDevicesByUser")
	rows, err := q.conn.Query(ctx, findDevicesByUserSQL, id)
	if err != nil {
		return nil, fmt.Errorf("query FindDevicesByUser: %w", err)
	}
	fds := rows.FieldDescriptions()
	plan0 := planScan(pgtype.TextCodec{}, fds[0], (*int)(nil))
	plan1 := planScan(pgtype.TextCodec{}, fds[1], (*string)(nil))
	plan2 := planScan(pgtype.TextCodec{}, fds[2], (*pgtype.MacaddrCodec])(nil))

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (FindDevicesByUserRow, error) {
		vals := row.RawValues()
		var item FindDevicesByUserRow
		if err := plan0.Scan(vals[0], &item); err != nil {
			return item, fmt.Errorf("scan FindDevicesByUser.id: %w", err)
		}
		if err := plan1.Scan(vals[1], &item); err != nil {
			return item, fmt.Errorf("scan FindDevicesByUser.name: %w", err)
		}
		if err := plan2.Scan(vals[2], &item); err != nil {
			return item, fmt.Errorf("scan FindDevicesByUser.mac_addrs: %w", err)
		}
		return item, nil
	})
}

const compositeUserSQL = `SELECT
  d.mac,
  d.type,
  ROW (u.id, u.name)::"user" AS "user"
FROM device d
  LEFT JOIN "user" u ON u.id = d.owner;`

type CompositeUserRow struct {
	Mac  pgtype.Macaddr `json:"mac"`
	Type DeviceType     `json:"type"`
	User User           `json:"user"`
}

// CompositeUser implements Querier.CompositeUser.
func (q *DBQuerier) CompositeUser(ctx context.Context) ([]CompositeUserRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "CompositeUser")
	rows, err := q.conn.Query(ctx, compositeUserSQL)
	if err != nil {
		return nil, fmt.Errorf("query CompositeUser: %w", err)
	}
	fds := rows.FieldDescriptions()
	plan0 := planScan(pgtype.TextCodec{}, fds[0], (*pgtype.Macaddr)(nil))
	plan1 := planScan(pgtype.TextCodec{}, fds[1], (*DeviceType)(nil))
	plan2 := planScan(pgtype.TextCodec{}, fds[2], (*User)(nil))

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (CompositeUserRow, error) {
		vals := row.RawValues()
		var item CompositeUserRow
		if err := plan0.Scan(vals[0], &item); err != nil {
			return item, fmt.Errorf("scan CompositeUser.mac: %w", err)
		}
		if err := plan1.Scan(vals[1], &item); err != nil {
			return item, fmt.Errorf("scan CompositeUser.type: %w", err)
		}
		if err := plan2.Scan(vals[2], &item); err != nil {
			return item, fmt.Errorf("scan CompositeUser.user: %w", err)
		}
		return item, nil
	})
}

const compositeUserOneSQL = `SELECT ROW (15, 'qux')::"user" AS "user";`

// CompositeUserOne implements Querier.CompositeUserOne.
func (q *DBQuerier) CompositeUserOne(ctx context.Context) (User, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "CompositeUserOne")
	rows, err := q.conn.Query(ctx, compositeUserOneSQL)
	if err != nil {
		return User{}, fmt.Errorf("query CompositeUserOne: %w", err)
	}
	fds := rows.FieldDescriptions()
	plan0 := planScan(pgtype.TextCodec{}, fds[0], (*User)(nil))

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (User, error) {
		vals := row.RawValues()
		var item User
		if err := plan0.Scan(vals[0], &item); err != nil {
			return item, fmt.Errorf("scan CompositeUserOne.user: %w", err)
		}
		return item, nil
	})
}

const compositeUserOneTwoColsSQL = `SELECT 1 AS num, ROW (15, 'qux')::"user" AS "user";`

type CompositeUserOneTwoColsRow struct {
	Num  int32 `json:"num"`
	User User  `json:"user"`
}

// CompositeUserOneTwoCols implements Querier.CompositeUserOneTwoCols.
func (q *DBQuerier) CompositeUserOneTwoCols(ctx context.Context) (CompositeUserOneTwoColsRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "CompositeUserOneTwoCols")
	rows, err := q.conn.Query(ctx, compositeUserOneTwoColsSQL)
	if err != nil {
		return CompositeUserOneTwoColsRow{}, fmt.Errorf("query CompositeUserOneTwoCols: %w", err)
	}
	fds := rows.FieldDescriptions()
	plan0 := planScan(pgtype.TextCodec{}, fds[0], (*int32)(nil))
	plan1 := planScan(pgtype.TextCodec{}, fds[1], (*User)(nil))

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (CompositeUserOneTwoColsRow, error) {
		vals := row.RawValues()
		var item CompositeUserOneTwoColsRow
		if err := plan0.Scan(vals[0], &item); err != nil {
			return item, fmt.Errorf("scan CompositeUserOneTwoCols.num: %w", err)
		}
		if err := plan1.Scan(vals[1], &item); err != nil {
			return item, fmt.Errorf("scan CompositeUserOneTwoCols.user: %w", err)
		}
		return item, nil
	})
}

const compositeUserManySQL = `SELECT ROW (15, 'qux')::"user" AS "user";`

// CompositeUserMany implements Querier.CompositeUserMany.
func (q *DBQuerier) CompositeUserMany(ctx context.Context) ([]User, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "CompositeUserMany")
	rows, err := q.conn.Query(ctx, compositeUserManySQL)
	if err != nil {
		return nil, fmt.Errorf("query CompositeUserMany: %w", err)
	}
	fds := rows.FieldDescriptions()
	plan0 := planScan(pgtype.TextCodec{}, fds[0], (*User)(nil))

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (User, error) {
		vals := row.RawValues()
		var item User
		if err := plan0.Scan(vals[0], &item); err != nil {
			return item, fmt.Errorf("scan CompositeUserMany.user: %w", err)
		}
		return item, nil
	})
}

const insertUserSQL = `INSERT INTO "user" (id, name)
VALUES ($1, $2);`

// InsertUser implements Querier.InsertUser.
func (q *DBQuerier) InsertUser(ctx context.Context, userID int, name string) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertUser")
	cmdTag, err := q.conn.Exec(ctx, insertUserSQL, userID, name)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query InsertUser: %w", err)
	}
	return cmdTag, err
}

const insertDeviceSQL = `INSERT INTO device (mac, owner)
VALUES ($1, $2);`

// InsertDevice implements Querier.InsertDevice.
func (q *DBQuerier) InsertDevice(ctx context.Context, mac pgtype.Macaddr, owner int) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertDevice")
	cmdTag, err := q.conn.Exec(ctx, insertDeviceSQL, mac, owner)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query InsertDevice: %w", err)
	}
	return cmdTag, err
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