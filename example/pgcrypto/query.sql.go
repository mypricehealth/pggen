// Code generated by pggen. DO NOT EDIT.

package pgcrypto

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
	CreateUser(ctx context.Context, email string, password string) (pgconn.CommandTag, error)

	FindUser(ctx context.Context, email string) (FindUserRow, error)
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

const createUserSQL = `INSERT INTO "user" (email, pass)
VALUES ($1, crypt($2, gen_salt('bf')));`

// CreateUser implements Querier.CreateUser.
func (q *DBQuerier) CreateUser(ctx context.Context, email string, password string) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, QueryName{}, "CreateUser")

	err := registerTypes(ctx, q.conn)
	if err != nil {
		return pgconn.CommandTag{}, q.errWrap(err)
	}
	cmdTag, err := q.conn.Exec(ctx, createUserSQL, email, password)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query CreateUser: %w", q.errWrap(err))
	}
	return cmdTag, q.errWrap(err)
}

const findUserSQL = `SELECT email, pass from "user"
where email = $1;`

type FindUserRow struct {
	Email string `json:"email"`
	Pass  string `json:"pass"`
}

// FindUser implements Querier.FindUser.
func (q *DBQuerier) FindUser(ctx context.Context, email string) (FindUserRow, error) {
	ctx = context.WithValue(ctx, QueryName{}, "FindUser")

	err := registerTypes(ctx, q.conn)
	if err != nil {
		return FindUserRow{}, q.errWrap(err)
	}
	rows, err := q.conn.Query(ctx, findUserSQL, email)
	if err != nil {
		return FindUserRow{}, fmt.Errorf("query FindUser: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[FindUserRow])
	return res, q.errWrap(err)
}
