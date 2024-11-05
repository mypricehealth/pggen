package parser

import (
	gotok "go/token"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/mypricehealth/pggen/internal/ast"
)

func ignoreCommentPos() cmp.Option {
	return cmpopts.IgnoreFields(ast.LineComment{}, "Start")
}

func ignoreQueryPos() cmp.Option {
	return cmpopts.IgnoreFields(ast.SourceQuery{}, "Start", "EndPos")
}

func TestParseFile_Queries(t *testing.T) {
	t.Run("multiple queries", func(t *testing.T) {
		tests := []struct {
			src  string
			want []ast.Query
		}{
			{
				"-- name: BasicMulti :exec\nSELECT 1;\nSELECT 2;\n-- name: Single :one\nSELECT 3;",
				[]ast.Query{
					&ast.SourceQuery{
						Name:              "BasicMulti",
						Doc:               &ast.CommentGroup{List: []*ast.LineComment{{Text: "-- name: BasicMulti :exec"}}},
						SourceSQL:         []string{"SELECT 1;", "\nSELECT 2;"},
						ContiguousArgsSQL: []ast.ContiguousArgsSQL{{SQL: "SELECT 1;"}, {SQL: "\nSELECT 2;"}},
						PreparedSQL:       "SELECT 1;\nSELECT 2;",
						Params:            nil,
						ResultKind:        ast.ResultKindExec,
					},
					&ast.SourceQuery{
						Name:              "Single",
						Doc:               &ast.CommentGroup{List: []*ast.LineComment{{Text: "-- name: Single :one"}}},
						Start:             42,
						SourceSQL:         []string{"SELECT 3;"},
						ContiguousArgsSQL: []ast.ContiguousArgsSQL{{SQL: "SELECT 3;"}},
						PreparedSQL:       "SELECT 3;",
						ResultKind:        ":one",
						EndPos:            71,
					},
				},
			},
			{
				"-- name: MultiWithArguments :exec\nSELECT pggen.arg('arg0');\nSELECT pggen.arg('arg1'), pggen.arg('arg0');\n-- name: Single :one\nSELECT pggen.arg('arg0');",
				[]ast.Query{
					&ast.SourceQuery{
						Name:              "MultiWithArguments",
						Doc:               &ast.CommentGroup{List: []*ast.LineComment{{Text: "-- name: MultiWithArguments :exec"}}},
						SourceSQL:         []string{"SELECT pggen.arg('arg0');", "\nSELECT pggen.arg('arg1'), pggen.arg('arg0');"},
						ContiguousArgsSQL: []ast.ContiguousArgsSQL{{SQL: "SELECT $1;", Args: []int{0}, UniqueArgs: 1}, {SQL: "\nSELECT $1, $2;", Args: []int{1, 0}, UniqueArgs: 2}},
						PreparedSQL:       "SELECT $1;\nSELECT $2, $1;",
						Params:            []ast.Param{{Name: "arg0"}, {Name: "arg1"}},
						ResultKind:        ast.ResultKindExec,
					},
					&ast.SourceQuery{
						Name:              "Single",
						Doc:               &ast.CommentGroup{List: []*ast.LineComment{{Text: "-- name: Single :one"}}},
						Start:             42,
						SourceSQL:         []string{"SELECT pggen.arg('arg0');"},
						ContiguousArgsSQL: []ast.ContiguousArgsSQL{{SQL: "SELECT $1;", Args: []int{0}, UniqueArgs: 1}},
						PreparedSQL:       "SELECT $1;",
						Params:            []ast.Param{{Name: "arg0"}},
						ResultKind:        ":one",
						EndPos:            71,
					},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.src, func(t *testing.T) {
				f, err := ParseFile(gotok.NewFileSet(), "", tt.src, Trace)
				if err != nil {
					t.Fatal(err)
				}

				if diff := cmp.Diff(tt.want, f.Queries, ignoreCommentPos(), ignoreQueryPos()); diff != "" {
					t.Errorf("ParseFile() query mismatch (-want +got):\n%s", diff)
				}
			})
		}
	})

	tests := []struct {
		src  string
		want ast.Query
	}{
		{
			"-- name: Qux :many\nSELECT 1;",
			&ast.SourceQuery{
				Name:              "Qux",
				Doc:               &ast.CommentGroup{List: []*ast.LineComment{{Text: "-- name: Qux :many"}}},
				SourceSQL:         []string{"SELECT 1;"},
				ContiguousArgsSQL: []ast.ContiguousArgsSQL{{SQL: "SELECT 1;"}},
				PreparedSQL:       "SELECT 1;",
				ResultKind:        ast.ResultKindMany,
			},
		},
		{
			"-- name: Foo :one\nSELECT 1;",
			&ast.SourceQuery{
				Name:              "Foo",
				Doc:               &ast.CommentGroup{List: []*ast.LineComment{{Text: "-- name: Foo :one"}}},
				SourceSQL:         []string{"SELECT 1;"},
				ContiguousArgsSQL: []ast.ContiguousArgsSQL{{SQL: "SELECT 1;"}},
				PreparedSQL:       "SELECT 1;",
				ResultKind:        ast.ResultKindOne,
			},
		},
		{
			"-- name: Qux   :exec\nSELECT pggen.arg('Bar');",
			&ast.SourceQuery{
				Name:              "Qux",
				Doc:               &ast.CommentGroup{List: []*ast.LineComment{{Text: "-- name: Qux   :exec"}}},
				SourceSQL:         []string{"SELECT pggen.arg('Bar');"},
				ContiguousArgsSQL: []ast.ContiguousArgsSQL{{SQL: "SELECT $1;", Args: []int{0}, UniqueArgs: 1}},
				PreparedSQL:       "SELECT $1;",
				Params:            []ast.Param{{Name: "Bar"}},
				ResultKind:        ast.ResultKindExec,
			},
		},
		{
			"-- name: Qux   :exec\nSELECT pggen.arg ('Bar');",
			&ast.SourceQuery{
				Name:              "Qux",
				Doc:               &ast.CommentGroup{List: []*ast.LineComment{{Text: "-- name: Qux   :exec"}}},
				SourceSQL:         []string{"SELECT pggen.arg ('Bar');"},
				ContiguousArgsSQL: []ast.ContiguousArgsSQL{{SQL: "SELECT $1;", Args: []int{0}, UniqueArgs: 1}},
				PreparedSQL:       "SELECT $1;",
				Params:            []ast.Param{{Name: "Bar"}},
				ResultKind:        ast.ResultKindExec,
			},
		},
		{
			"-- name: Qux :one\nSELECT pggen.arg('A$_$$B123');",
			&ast.SourceQuery{
				Name:              "Qux",
				Doc:               &ast.CommentGroup{List: []*ast.LineComment{{Text: "-- name: Qux :one"}}},
				SourceSQL:         []string{"SELECT pggen.arg('A$_$$B123');"},
				ContiguousArgsSQL: []ast.ContiguousArgsSQL{{SQL: "SELECT $1;", Args: []int{0}, UniqueArgs: 1}},
				PreparedSQL:       "SELECT $1;",
				Params:            []ast.Param{{Name: "A$_$$B123"}},
				ResultKind:        ast.ResultKindOne,
			},
		},
		{
			"-- name: Qux :many\nSELECT pggen.arg('Bar'), pggen.arg('Qux'), pggen.arg('Bar');",
			&ast.SourceQuery{
				Name:              "Qux",
				Doc:               &ast.CommentGroup{List: []*ast.LineComment{{Text: "-- name: Qux :many"}}},
				SourceSQL:         []string{"SELECT pggen.arg('Bar'), pggen.arg('Qux'), pggen.arg('Bar');"},
				ContiguousArgsSQL: []ast.ContiguousArgsSQL{{SQL: "SELECT $1, $2, $1;", Args: []int{0, 1}, UniqueArgs: 2}},
				PreparedSQL:       "SELECT $1, $2, $1;",
				Params:            []ast.Param{{Name: "Bar"}, {Name: "Qux"}},
				ResultKind:        ast.ResultKindMany,
			},
		},
		{
			"-- name: Qux :many\nSELECT /*pggen.arg('Bar'),*/ pggen.arg('Qux'), pggen.arg('Bar');",
			&ast.SourceQuery{
				Name:              "Qux",
				Doc:               &ast.CommentGroup{List: []*ast.LineComment{{Text: "-- name: Qux :many"}}},
				SourceSQL:         []string{"SELECT /*pggen.arg('Bar'),*/ pggen.arg('Qux'), pggen.arg('Bar');"},
				ContiguousArgsSQL: []ast.ContiguousArgsSQL{{SQL: "SELECT /*pggen.arg('Bar'),*/ $1, $2;", Args: []int{0, 1}, UniqueArgs: 2}},
				PreparedSQL:       "SELECT /*pggen.arg('Bar'),*/ $1, $2;",
				Params:            []ast.Param{{Name: "Qux"}, {Name: "Bar"}},
				ResultKind:        ast.ResultKindMany,
			},
		},
		{
			"-- name: Qux :many proto-type=foo.Bar\nSELECT 1;",
			&ast.SourceQuery{
				Name:              "Qux",
				Doc:               &ast.CommentGroup{List: []*ast.LineComment{{Text: "-- name: Qux :many proto-type=foo.Bar"}}},
				SourceSQL:         []string{"SELECT 1;"},
				ContiguousArgsSQL: []ast.ContiguousArgsSQL{{SQL: "SELECT 1;"}},
				PreparedSQL:       "SELECT 1;",
				Params:            nil,
				ResultKind:        ast.ResultKindMany,
				Pragmas:           ast.Pragmas{ProtobufType: "foo.Bar"},
			},
		},
		{
			"-- name: Qux :many proto-type=Bar\nSELECT 1;",
			&ast.SourceQuery{
				Name:              "Qux",
				Doc:               &ast.CommentGroup{List: []*ast.LineComment{{Text: "-- name: Qux :many proto-type=Bar"}}},
				SourceSQL:         []string{"SELECT 1;"},
				ContiguousArgsSQL: []ast.ContiguousArgsSQL{{SQL: "SELECT 1;"}},
				PreparedSQL:       "SELECT 1;",
				Params:            nil,
				ResultKind:        ast.ResultKindMany,
				Pragmas:           ast.Pragmas{ProtobufType: "Bar"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.src, func(t *testing.T) {
			f, err := ParseFile(gotok.NewFileSet(), "", tt.src, Trace)
			if err != nil {
				t.Fatal(err)
			}

			if len(f.Queries) != 1 {
				t.Errorf("ParseFile() expected one query, got %d", len(f.Queries))
				return
			}

			got := f.Queries[0].(*ast.SourceQuery)

			if diff := cmp.Diff(tt.want, got, ignoreCommentPos(), ignoreQueryPos()); diff != "" {
				t.Errorf("ParseFile() query mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParseFile_Queries_Fuzz(t *testing.T) {
	tests := []struct {
		src string
	}{
		{"-- name: Qux :many\nSELECT '`\\n' as \" joe!@#$%&*()-+=\";"},
	}
	for _, tt := range tests {
		t.Run(tt.src, func(t *testing.T) {
			_, err := ParseFile(gotok.NewFileSet(), "", tt.src, Trace)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
