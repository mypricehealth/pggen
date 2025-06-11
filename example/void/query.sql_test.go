package void

import (
	"context"
	"testing"

	"github.com/mypricehealth/pggen/internal/pgtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuerier(t *testing.T) {
	conn, cleanup := pgtest.NewPostgresSchema(t, []string{"schema.sql"})
	defer cleanup()

	ctx := context.Background()
	q, err := NewQuerier(ctx, conn)
	require.NoError(t, err)

	if _, err := q.VoidOnly(ctx); err != nil {
		t.Fatal(err)
	}

	if _, err := q.VoidOnlyTwoParams(ctx, 33); err != nil {
		t.Fatal(err)
	}

	{
		row, err := q.VoidTwo(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "foo", row)
	}

	{
		row, err := q.VoidThree(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, VoidThreeRow{Foo: "foo", Bar: "bar"}, row)
	}

	{
		foos, err := q.VoidThree2(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, []string{"foo"}, foos)
	}
}
