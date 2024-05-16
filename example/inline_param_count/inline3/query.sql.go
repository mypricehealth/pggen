// Code generated by pggen. DO NOT EDIT.

package inline3

import (
	"context"
	"fmt"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type QueryName struct{}

// Querier is a typesafe Go interface backed by SQL queries.
type Querier interface {
	// CountAuthors returns the number of authors (zero params).
	CountAuthors(ctx context.Context) (*int, error)

	// FindAuthorById finds one (or zero) authors by ID (one param).
	FindAuthorByID(ctx context.Context, authorID int32) (FindAuthorByIDRow, error)

	// InsertAuthor inserts an author by name and returns the ID (two params).
	InsertAuthor(ctx context.Context, firstName string, lastName string) (int32, error)

	// DeleteAuthorsByFullName deletes authors by the full name (three params).
	DeleteAuthorsByFullName(ctx context.Context, firstName string, lastName string, suffix string) (pgconn.CommandTag, error)
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

// RegisterTypes should be run in config.AfterConnect to load custom types
func RegisterTypes(ctx context.Context, conn *pgx.Conn) error {
	pgxdecimal.Register(conn.TypeMap())
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

const countAuthorsSQL = `SELECT count(*) FROM author;`

// CountAuthors implements Querier.CountAuthors.
func (q *DBQuerier) CountAuthors(ctx context.Context) (*int, error) {
	ctx = context.WithValue(ctx, QueryName{}, "CountAuthors")
	rows, err := q.conn.Query(ctx, countAuthorsSQL)
	if err != nil {
		return nil, fmt.Errorf("query CountAuthors: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[*int])
	return res, q.errWrap(err)
}

const findAuthorByIDSQL = `SELECT * FROM author WHERE author_id = $1;`

type FindAuthorByIDRow struct {
	AuthorID  int32   `json:"author_id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Suffix    *string `json:"suffix"`
}

// FindAuthorByID implements Querier.FindAuthorByID.
func (q *DBQuerier) FindAuthorByID(ctx context.Context, authorID int32) (FindAuthorByIDRow, error) {
	ctx = context.WithValue(ctx, QueryName{}, "FindAuthorByID")
	rows, err := q.conn.Query(ctx, findAuthorByIDSQL, authorID)
	if err != nil {
		return FindAuthorByIDRow{}, fmt.Errorf("query FindAuthorByID: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[FindAuthorByIDRow])
	return res, q.errWrap(err)
}

const insertAuthorSQL = `INSERT INTO author (first_name, last_name)
VALUES ($1, $2)
RETURNING author_id;`

// InsertAuthor implements Querier.InsertAuthor.
func (q *DBQuerier) InsertAuthor(ctx context.Context, firstName string, lastName string) (int32, error) {
	ctx = context.WithValue(ctx, QueryName{}, "InsertAuthor")
	rows, err := q.conn.Query(ctx, insertAuthorSQL, firstName, lastName)
	if err != nil {
		return 0, fmt.Errorf("query InsertAuthor: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[int32])
	return res, q.errWrap(err)
}

const deleteAuthorsByFullNameSQL = `DELETE
FROM author
WHERE first_name = $1
  AND last_name = $2
  AND CASE WHEN $3 = '' THEN suffix IS NULL ELSE suffix = $3 END;`

// DeleteAuthorsByFullName implements Querier.DeleteAuthorsByFullName.
func (q *DBQuerier) DeleteAuthorsByFullName(ctx context.Context, firstName string, lastName string, suffix string) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, QueryName{}, "DeleteAuthorsByFullName")
	cmdTag, err := q.conn.Exec(ctx, deleteAuthorsByFullNameSQL, firstName, lastName, suffix)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query DeleteAuthorsByFullName: %w", q.errWrap(err))
	}
	return cmdTag, q.errWrap(err)
}
