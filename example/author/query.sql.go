// Code generated by pggen. DO NOT EDIT.

package author

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
	// FindAuthorById finds one (or zero) authors by ID.
	FindAuthorByID(ctx context.Context, authorID int32) (FindAuthorByIDRow, error)

	// FindAuthors finds authors by first name.
	FindAuthors(ctx context.Context, firstName string) ([]FindAuthorsRow, error)

	// FindAuthorNames finds one (or zero) authors by ID.
	FindAuthorNames(ctx context.Context, authorID int32) ([]FindAuthorNamesRow, error)

	// FindFirstNames finds one (or zero) authors by ID.
	FindFirstNames(ctx context.Context, authorID int32) ([]*string, error)

	// DeleteAuthors deletes authors with a first name of "joe".
	DeleteAuthors(ctx context.Context) (pgconn.CommandTag, error)

	// DeleteAuthorsByFirstName deletes authors by first name.
	DeleteAuthorsByFirstName(ctx context.Context, firstName string) (pgconn.CommandTag, error)

	// DeleteAuthorsByFullName deletes authors by the full name.
	DeleteAuthorsByFullName(ctx context.Context, params DeleteAuthorsByFullNameParams) (pgconn.CommandTag, error)

	// InsertAuthor inserts an author by name and returns the ID.
	InsertAuthor(ctx context.Context, firstName string, lastName string) (int32, error)

	// InsertAuthorSuffix inserts an author by name and suffix and returns the
	// entire row.
	InsertAuthorSuffix(ctx context.Context, params InsertAuthorSuffixParams) (InsertAuthorSuffixRow, error)

	StringAggFirstName(ctx context.Context, authorID int32) (*string, error)

	ArrayAggFirstName(ctx context.Context, authorID int32) ([]string, error)
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

const findAuthorsSQL = `SELECT * FROM author WHERE first_name = $1;`

type FindAuthorsRow struct {
	AuthorID  int32   `json:"author_id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Suffix    *string `json:"suffix"`
}

// FindAuthors implements Querier.FindAuthors.
func (q *DBQuerier) FindAuthors(ctx context.Context, firstName string) ([]FindAuthorsRow, error) {
	ctx = context.WithValue(ctx, QueryName{}, "FindAuthors")
	rows, err := q.conn.Query(ctx, findAuthorsSQL, firstName)
	if err != nil {
		return nil, fmt.Errorf("query FindAuthors: %w", q.errWrap(err))
	}
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[FindAuthorsRow])
	return res, q.errWrap(err)
}

const findAuthorNamesSQL = `SELECT first_name, last_name FROM author ORDER BY author_id = $1;`

type FindAuthorNamesRow struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}

// FindAuthorNames implements Querier.FindAuthorNames.
func (q *DBQuerier) FindAuthorNames(ctx context.Context, authorID int32) ([]FindAuthorNamesRow, error) {
	ctx = context.WithValue(ctx, QueryName{}, "FindAuthorNames")
	rows, err := q.conn.Query(ctx, findAuthorNamesSQL, authorID)
	if err != nil {
		return nil, fmt.Errorf("query FindAuthorNames: %w", q.errWrap(err))
	}
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[FindAuthorNamesRow])
	return res, q.errWrap(err)
}

const findFirstNamesSQL = `SELECT first_name FROM author ORDER BY author_id = $1;`

// FindFirstNames implements Querier.FindFirstNames.
func (q *DBQuerier) FindFirstNames(ctx context.Context, authorID int32) ([]*string, error) {
	ctx = context.WithValue(ctx, QueryName{}, "FindFirstNames")
	rows, err := q.conn.Query(ctx, findFirstNamesSQL, authorID)
	if err != nil {
		return nil, fmt.Errorf("query FindFirstNames: %w", q.errWrap(err))
	}
	res, err := pgx.CollectRows(rows, pgx.RowTo[*string])
	return res, q.errWrap(err)
}

const deleteAuthorsSQL = `DELETE FROM author WHERE first_name = 'joe';`

// DeleteAuthors implements Querier.DeleteAuthors.
func (q *DBQuerier) DeleteAuthors(ctx context.Context) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, QueryName{}, "DeleteAuthors")
	cmdTag, err := q.conn.Exec(ctx, deleteAuthorsSQL)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query DeleteAuthors: %w", q.errWrap(err))
	}
	return cmdTag, q.errWrap(err)
}

const deleteAuthorsByFirstNameSQL = `DELETE FROM author WHERE first_name = $1;`

// DeleteAuthorsByFirstName implements Querier.DeleteAuthorsByFirstName.
func (q *DBQuerier) DeleteAuthorsByFirstName(ctx context.Context, firstName string) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, QueryName{}, "DeleteAuthorsByFirstName")
	cmdTag, err := q.conn.Exec(ctx, deleteAuthorsByFirstNameSQL, firstName)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query DeleteAuthorsByFirstName: %w", q.errWrap(err))
	}
	return cmdTag, q.errWrap(err)
}

const deleteAuthorsByFullNameSQL = `DELETE
FROM author
WHERE first_name = $1
  AND last_name = $2
  AND suffix = $3;`

type DeleteAuthorsByFullNameParams struct {
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Suffix    string `json:"Suffix"`
}

// DeleteAuthorsByFullName implements Querier.DeleteAuthorsByFullName.
func (q *DBQuerier) DeleteAuthorsByFullName(ctx context.Context, params DeleteAuthorsByFullNameParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, QueryName{}, "DeleteAuthorsByFullName")
	cmdTag, err := q.conn.Exec(ctx, deleteAuthorsByFullNameSQL, params.FirstName, params.LastName, params.Suffix)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query DeleteAuthorsByFullName: %w", q.errWrap(err))
	}
	return cmdTag, q.errWrap(err)
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

const insertAuthorSuffixSQL = `INSERT INTO author (first_name, last_name, suffix)
VALUES ($1, $2, $3)
RETURNING author_id, first_name, last_name, suffix;`

type InsertAuthorSuffixParams struct {
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Suffix    string `json:"Suffix"`
}

type InsertAuthorSuffixRow struct {
	AuthorID  int32   `json:"author_id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Suffix    *string `json:"suffix"`
}

// InsertAuthorSuffix implements Querier.InsertAuthorSuffix.
func (q *DBQuerier) InsertAuthorSuffix(ctx context.Context, params InsertAuthorSuffixParams) (InsertAuthorSuffixRow, error) {
	ctx = context.WithValue(ctx, QueryName{}, "InsertAuthorSuffix")
	rows, err := q.conn.Query(ctx, insertAuthorSuffixSQL, params.FirstName, params.LastName, params.Suffix)
	if err != nil {
		return InsertAuthorSuffixRow{}, fmt.Errorf("query InsertAuthorSuffix: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[InsertAuthorSuffixRow])
	return res, q.errWrap(err)
}

const stringAggFirstNameSQL = `SELECT string_agg(first_name, ',') AS names FROM author WHERE author_id = $1;`

// StringAggFirstName implements Querier.StringAggFirstName.
func (q *DBQuerier) StringAggFirstName(ctx context.Context, authorID int32) (*string, error) {
	ctx = context.WithValue(ctx, QueryName{}, "StringAggFirstName")
	rows, err := q.conn.Query(ctx, stringAggFirstNameSQL, authorID)
	if err != nil {
		return nil, fmt.Errorf("query StringAggFirstName: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[*string])
	return res, q.errWrap(err)
}

const arrayAggFirstNameSQL = `SELECT array_agg(first_name) AS names FROM author WHERE author_id = $1;`

// ArrayAggFirstName implements Querier.ArrayAggFirstName.
func (q *DBQuerier) ArrayAggFirstName(ctx context.Context, authorID int32) ([]string, error) {
	ctx = context.WithValue(ctx, QueryName{}, "ArrayAggFirstName")
	rows, err := q.conn.Query(ctx, arrayAggFirstNameSQL, authorID)
	if err != nil {
		return nil, fmt.Errorf("query ArrayAggFirstName: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[[]string])
	return res, q.errWrap(err)
}
