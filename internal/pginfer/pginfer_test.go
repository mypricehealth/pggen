package pginfer

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/mypricehealth/pggen/internal/ast"
	"github.com/mypricehealth/pggen/internal/difftest"
	"github.com/mypricehealth/pggen/internal/pg"
	"github.com/mypricehealth/pggen/internal/pgtest"
	"github.com/mypricehealth/pggen/internal/texts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInferrer_InferTypes(t *testing.T) {
	conn, cleanupFunc := pgtest.NewPostgresSchemaString(t, texts.Dedent(`
		CREATE TABLE author (
			author_id  serial PRIMARY KEY,
			first_name text NOT NULL,
			last_name  text NOT NULL,
			suffix text NULL
		);

		CREATE TYPE device_type AS enum (
			'phone',
			'laptop'
		);

		CREATE DOMAIN us_postal_code AS text;
	`))
	defer cleanupFunc()
	q := pg.NewQuerier(conn)
	deviceTypeOID, err := q.FindOIDByName(context.Background(), "device_type")
	require.NoError(t, err)
	deviceTypeArrOID, err := q.FindOIDByName(context.Background(), "_device_type")
	require.NoError(t, err)

	tests := []struct {
		name     string
		query    *ast.SourceQuery
		denseSQL string
		argCount int
		want     TypedQuery
	}{
		{
			name: "literal query",
			query: &ast.SourceQuery{
				Name:       "LiteralQuery",
				ResultKind: ast.ResultKindOne,
			},
			denseSQL: "SELECT 1 as one, 'foo' as two",
			want: TypedQuery{
				Name:       "LiteralQuery",
				ResultKind: ast.ResultKindOne,
				Outputs: []OutputColumn{
					{PgName: "one", PgType: pg.Int4, Nullable: false},
					{PgName: "two", PgType: pg.Text, Nullable: false},
				},
			},
		},
		{
			name: "union one col",
			query: &ast.SourceQuery{
				Name:       "UnionOneCol",
				ResultKind: ast.ResultKindMany,
			},
			denseSQL: "SELECT 1 AS num UNION SELECT 2 AS num",
			want: TypedQuery{
				Name:       "UnionOneCol",
				ResultKind: ast.ResultKindMany,
				Outputs: []OutputColumn{
					{PgName: "num", PgType: pg.Int4, Nullable: true},
				},
			},
		},
		{
			name: "one col domain type",
			query: &ast.SourceQuery{
				Name:       "Domain",
				ResultKind: ast.ResultKindOne,
			},
			denseSQL: "SELECT '94109'::us_postal_code",
			want: TypedQuery{
				Name:       "Domain",
				ResultKind: ast.ResultKindOne,
				Outputs: []OutputColumn{{
					PgName:   "us_postal_code",
					PgType:   pg.Text,
					Nullable: false,
				}},
			},
		},
		{
			name: "one col domain type",
			query: &ast.SourceQuery{
				Name:       "UnionEnumArrays",
				ResultKind: ast.ResultKindMany,
			},
			denseSQL: texts.Dedent(`
                SELECT enum_range('phone'::device_type, 'phone'::device_type) AS device_types
                UNION ALL
                SELECT enum_range(NULL::device_type) AS device_types;
            `),
			want: TypedQuery{
				Name:       "UnionEnumArrays",
				ResultKind: ast.ResultKindMany,
				Outputs: []OutputColumn{
					{
						PgName: "device_types",
						PgType: pg.ArrayType{
							ID:   deviceTypeArrOID,
							Name: "_device_type",
							Elem: pg.EnumType{
								ID:     deviceTypeOID,
								Name:   "device_type",
								Labels: []string{"phone", "laptop"},
								Orders: []float32{1, 2},
							},
						},
						Nullable: true,
					},
				},
			},
		},
		{
			name: "find by first name",
			query: &ast.SourceQuery{
				Name:       "FindByFirstName",
				Params:     []ast.Param{{Name: "FirstName"}},
				ResultKind: ast.ResultKindMany,
				Doc:        newCommentGroup("--   Hello  ", "-- name: Foo"),
			},
			denseSQL: "SELECT first_name FROM author WHERE first_name = $1;",
			argCount: 1,
			want: TypedQuery{
				Name:       "FindByFirstName",
				ResultKind: ast.ResultKindMany,
				Doc:        []string{"Hello"},
				Inputs: []InputParam{
					{PgName: "FirstName", PgType: pg.Text},
				},
				Outputs: []OutputColumn{
					{PgName: "first_name", PgType: pg.Text, Nullable: false},
				},
			},
		},
		{
			name: "find by first name join",
			query: &ast.SourceQuery{
				Name:       "FindByFirstNameJoin",
				Params:     []ast.Param{{Name: "FirstName"}},
				ResultKind: ast.ResultKindMany,
				Doc:        newCommentGroup("--   Hello  ", "-- name: Foo"),
			},
			denseSQL: "SELECT a1.first_name FROM author a1 JOIN author a2 USING (author_id) WHERE a1.first_name = $1;",
			argCount: 1,
			want: TypedQuery{
				Name:       "FindByFirstNameJoin",
				ResultKind: ast.ResultKindMany,
				Doc:        []string{"Hello"},
				Inputs: []InputParam{
					{PgName: "FirstName", PgType: pg.Text},
				},
				Outputs: []OutputColumn{
					{PgName: "first_name", PgType: pg.Text, Nullable: true},
				},
			},
		},
		{
			name: "delete by author ID",
			query: &ast.SourceQuery{
				Name:       "DeleteAuthorByID",
				Params:     []ast.Param{{Name: "AuthorID"}},
				ResultKind: ast.ResultKindExec,
				Doc:        newCommentGroup("-- One", "--- - two", "-- name: Foo"),
			},
			denseSQL: "DELETE FROM author WHERE author_id = $1;",
			argCount: 1,
			want: TypedQuery{
				Name:       "DeleteAuthorByID",
				ResultKind: ast.ResultKindExec,
				Doc:        []string{"One", "- two"},
				Inputs: []InputParam{
					{PgName: "AuthorID", PgType: pg.Int4},
				},
				Outputs: nil,
			},
		},
		{
			name: "delete by author id returning",
			query: &ast.SourceQuery{
				Name:       "DeleteAuthorByIDReturning",
				Params:     []ast.Param{{Name: "AuthorID"}},
				ResultKind: ast.ResultKindMany,
			},
			denseSQL: "DELETE FROM author WHERE author_id = $1 RETURNING author_id, first_name, suffix;",
			argCount: 1,
			want: TypedQuery{
				Name:       "DeleteAuthorByIDReturning",
				ResultKind: ast.ResultKindMany,
				Inputs: []InputParam{
					{PgName: "AuthorID", PgType: pg.Int4},
				},
				Outputs: []OutputColumn{
					{PgName: "author_id", PgType: pg.Int4, Nullable: false},
					{PgName: "first_name", PgType: pg.Text, Nullable: false},
					{PgName: "suffix", PgType: pg.Text, Nullable: true},
				},
			},
		},
		{
			name: "update by author id returning",
			query: &ast.SourceQuery{
				Name:       "UpdateByAuthorIDReturning",
				Params:     []ast.Param{{Name: "AuthorID"}},
				ResultKind: ast.ResultKindMany,
			},
			denseSQL: "UPDATE author set first_name = 'foo' WHERE author_id = $1 RETURNING author_id, first_name, suffix;",
			argCount: 1,
			want: TypedQuery{
				Name:       "UpdateByAuthorIDReturning",
				ResultKind: ast.ResultKindMany,
				Inputs: []InputParam{
					{PgName: "AuthorID", PgType: pg.Int4},
				},
				Outputs: []OutputColumn{
					{PgName: "author_id", PgType: pg.Int4, Nullable: false},
					{PgName: "first_name", PgType: pg.Text, Nullable: false},
					{PgName: "suffix", PgType: pg.Text, Nullable: true},
				},
			},
		},
		{
			name: "void one",
			query: &ast.SourceQuery{
				Name:       "VoidOne",
				Params:     []ast.Param{},
				ResultKind: ast.ResultKindExec,
			},
			denseSQL: "SELECT ''::void;",
			want: TypedQuery{
				Name:       "VoidOne",
				ResultKind: ast.ResultKindExec,
				Inputs:     nil,
				Outputs: []OutputColumn{
					{PgName: "void", PgType: pg.Void, Nullable: false},
				},
			},
		},
		{
			name: "void two",
			query: &ast.SourceQuery{
				Name:       "VoidTwo",
				Params:     []ast.Param{},
				ResultKind: ast.ResultKindOne,
			},
			denseSQL: "SELECT 'foo' as foo, ''::void;",
			want: TypedQuery{
				Name:       "VoidTwo",
				ResultKind: ast.ResultKindOne,
				Inputs:     nil,
				Outputs: []OutputColumn{
					{PgName: "foo", PgType: pg.Text, Nullable: false},
					{PgName: "void", PgType: pg.Void, Nullable: false},
				},
			},
		},
		{
			name: "pragma proto type",
			query: &ast.SourceQuery{
				Name:       "PragmaProtoType",
				ResultKind: ast.ResultKindOne,
				Pragmas:    ast.Pragmas{ProtobufType: "foo.Bar"},
			},
			denseSQL: "SELECT 1 as one, 'foo' as two",
			want: TypedQuery{
				Name:       "PragmaProtoType",
				ResultKind: ast.ResultKindOne,
				Outputs: []OutputColumn{
					{PgName: "one", PgType: pg.Int4, Nullable: false},
					{PgName: "two", PgType: pg.Text, Nullable: false},
				},
				ProtobufType: "foo.Bar",
			},
		},
		{
			name: "aggregate non-null column has null output",
			query: &ast.SourceQuery{
				Name:       "ArrayAggFirstName",
				Params:     []ast.Param{},
				ResultKind: ast.ResultKindOne,
				Doc:        newCommentGroup("--   Hello  ", "-- name: Foo"),
			},
			denseSQL: "SELECT array_agg(first_name) AS names FROM author;",
			want: TypedQuery{
				Name:       "ArrayAggFirstName",
				ResultKind: ast.ResultKindOne,
				Doc:        []string{"Hello"},
				Inputs:     []InputParam{},
				Outputs: []OutputColumn{
					{PgName: "names", PgType: pg.TextArray, Nullable: true},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inferrer := NewInferrer(conn)

			args := make([]int, tt.argCount)
			for i := 0; i < tt.argCount; i++ {
				args[i] = i
			}

			denseSQL := []ast.DenseSQL{{SQL: tt.denseSQL, Args: args, UniqueArgs: tt.argCount}}

			tt.query.DenseSQL = denseSQL
			tt.query.PreparedSQL = tt.denseSQL

			tt.want.DenseSQL = denseSQL
            tt.want.PreparedSQL = tt.denseSQL

			got, err := inferrer.InferTypes(tt.query)
			if err != nil {
				t.Fatal(err)
			}
			opts := cmp.Options{
				cmpopts.IgnoreFields(pg.EnumType{}, "ChildOIDs"),
			}
			difftest.AssertSame(t, tt.want, got, opts)
		})
	}
}

func TestInferrer_InferTypes_Multi(t *testing.T) {
	query := &ast.SourceQuery{
		Name:        "MultiQuery",
		DenseSQL:    []ast.DenseSQL{{SQL: "SELECT $1::text;", Args: []int{0}}, {SQL: "SELECT $1::text[];", Args: []int{1}}},
		PreparedSQL: "SELECT $1::text; SELECT $2::text[];",
		Params:      []ast.Param{{Name: "TextParam"}, {Name: "TextArrayParam"}},
		ResultKind:  ast.ResultKindOne,
	}

	want := TypedQuery{
		Name:        "MultiQuery",
		ResultKind:  ast.ResultKindOne,
		DenseSQL:    []ast.DenseSQL{{SQL: "SELECT $1::text;", Args: []int{0}}, {SQL: "SELECT $1::text[];", Args: []int{1}}},
		PreparedSQL: "SELECT $1::text; SELECT $2::text[];",
		Inputs: []InputParam{
			{PgName: "TextParam", PgType: pg.Text},
			{PgName: "TextArrayParam", PgType: pg.TextArray},
		},
		Outputs: []OutputColumn{
			{PgName: "text", PgType: pg.TextArray, Nullable: false},
		},
	}

	conn, cleanupFunc := pgtest.NewPostgresSchemaString(t, "")
	defer cleanupFunc()

	inferrer := NewInferrer(conn)

	got, err := inferrer.InferTypes(query)
	if err != nil {
		t.Fatal(err)
	}
	opts := cmp.Options{
		cmpopts.IgnoreFields(pg.EnumType{}, "ChildOIDs"),
	}
	difftest.AssertSame(t, want, got, opts)
}

func TestInferrer_InferTypes_Error(t *testing.T) {
	conn, cleanupFunc := pgtest.NewPostgresSchema(t, []string{
		"../../example/author/schema.sql",
	})
	defer cleanupFunc()

	tests := []struct {
		query *ast.SourceQuery
		want  error
	}{
		{
			&ast.SourceQuery{
				Name:        "DeleteAuthorByIDMany",
				DenseSQL:    []ast.DenseSQL{{SQL: "DELETE FROM author WHERE author_id = $1;", Args: []int{0}}},
				PreparedSQL: "DELETE FROM author WHERE author_id = $1;",
				Params:      []ast.Param{{Name: "AuthorID"}},
				ResultKind:  ast.ResultKindMany,
			},
			errors.New("query DeleteAuthorByIDMany has incompatible result kind :many; " +
				"the query doesn't return any columns; " +
				"use :exec or :setup if the query shouldn't return any columns"),
		},
		{
			&ast.SourceQuery{
				Name:        "DeleteAuthorByIDOne",
				DenseSQL:    []ast.DenseSQL{{SQL: "DELETE FROM author WHERE author_id = $1;", Args: []int{0}}},
				PreparedSQL: "DELETE FROM author WHERE author_id = $1;",
				Params:      []ast.Param{{Name: "AuthorID"}},
				ResultKind:  ast.ResultKindOne,
			},
			errors.New(
				"query DeleteAuthorByIDOne has incompatible result kind :one; " +
					"the query doesn't return any columns; " +
					"use :exec or :setup if the query shouldn't return any columns"),
		},
		{
			&ast.SourceQuery{
				Name:        "VoidOne",
				DenseSQL:    []ast.DenseSQL{{SQL: "SELECT ''::void;"}},
				PreparedSQL: "SELECT ''::void;",
				Params:      nil,
				ResultKind:  ast.ResultKindMany,
			},
			errors.New(
				"query VoidOne has incompatible result kind :many; " +
					"the query only has void columns; " +
					"use :exec or :setup if the query shouldn't return any columns"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.query.Name, func(t *testing.T) {
			inferrer := NewInferrer(conn)
			got, err := inferrer.InferTypes(tt.query)
			assert.Equal(t, TypedQuery{}, got, "InferTypes should error and return empty TypedQuery struct")
			assert.Equal(t, tt.want, err, "InferType error should match")
		})
	}
}

func newCommentGroup(lines ...string) *ast.CommentGroup {
	cs := make([]*ast.LineComment, len(lines))
	for i, line := range lines {
		cs[i] = &ast.LineComment{Text: line}
	}
	return &ast.CommentGroup{List: cs}
}
