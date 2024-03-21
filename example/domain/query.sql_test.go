package domain

import (
	"context"
	"testing"

	"github.com/mypricehealth/pggen/internal/pgtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuerier_DomainOne(t *testing.T) {
	conn, cleanup := pgtest.NewPostgresSchema(t, []string{"schema.sql"})
	defer cleanup()

	q := NewQuerier(conn)
	ctx := context.Background()

	t.Run("DomainOne", func(t *testing.T) {
		postCode, err := q.DomainOne(ctx)
		require.NoError(t, err)
		assert.Equal(t, "90210", postCode)
	})
}
