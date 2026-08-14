package pggen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mypricehealth/pggen/internal/pgtest"
	"github.com/mypricehealth/pggen/internal/texts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate_Golang_Error(t *testing.T) {
	tests := []struct {
		name       string
		schema     string
		queries    string
		wantErrMsg string
	}{
		{
			name:   "duplicate query name",
			schema: "",
			queries: texts.Dedent(`
			-- name: Foo :many
			SELECT 1;
			-- name: Foo :many
			SELECT 1;
			`),
			wantErrMsg: `duplicate query name Foo`,
		},
		{
			name:   "type error",
			schema: "",
			queries: texts.Dedent(`
			-- name: Foo :one
			SELECT encode(123, 'foo'::text);
			`),
			wantErrMsg: `function encode(integer, text) does not exist`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, cleanupFunc := pgtest.NewPostgresSchemaString(t, tt.schema)
			defer cleanupFunc()
			tmpDir := t.TempDir()
			queryFile := filepath.Join(tmpDir, "query.sql")
			err := os.WriteFile(queryFile, []byte(tt.queries), 0644)
			if err != nil {
				t.Fatal(err)
			}

			err = Generate(
				GenerateOptions{
					ConnString: conn.Config().ConnString(),
					QueryFiles: []string{queryFile},
					OutputDir:  tmpDir,
					GoPackage:  "error_test",
					Language:   LangGo,
				})

			if err == nil {
				t.Fatal("expected error from generate")
			}
			assert.Contains(t, err.Error(), tt.wantErrMsg, "error message should contain substring")
		})
	}
}

func TestGenerate_Golang_Nullability(t *testing.T) {
	schema := texts.Dedent(`
		CREATE TABLE author (
			author_id serial PRIMARY KEY,
			suffix    text NULL
		);
	`)
	queries := texts.Dedent(`
		-- name: FindAuthor :one
		SELECT suffix, suffix AS "opt_suffix?" FROM author WHERE author_id = pggen.arg('ID');
	`)

	conn, cleanupFunc := pgtest.NewPostgresSchemaString(t, schema)
	defer cleanupFunc()
	tmpDir := t.TempDir()
	queryFile := filepath.Join(tmpDir, "query.sql")
	err := os.WriteFile(queryFile, []byte(queries), 0644)
	require.NoError(t, err)

	err = Generate(GenerateOptions{
		ConnString: conn.Config().ConnString(),
		QueryFiles: []string{queryFile},
		OutputDir:  tmpDir,
		GoPackage:  "nullability_test",
		Language:   LangGo,
	})
	require.NoError(t, err)

	matches, err := filepath.Glob(filepath.Join(tmpDir, "*.sql.go"))
	require.NoError(t, err)
	require.Len(t, matches, 1)
	src, err := os.ReadFile(matches[0])
	require.NoError(t, err)
	got := string(src)

	// The unmarked column stays non-nullable and needs no db tag; the marked
	// column is a pointer and carries a db tag matching the reported name.
	assert.Regexp(t, `\bSuffix\s+string\b`, got, "unmarked column should be non-nullable")
	assert.Regexp(t, `\bOptSuffix\s+\*string\s+`+"`"+`json:"opt_suffix" db:"opt_suffix\?"`+"`", got, "marked column should be a nullable pointer with a db tag")
}
