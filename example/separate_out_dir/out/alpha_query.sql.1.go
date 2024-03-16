// Code generated by pggen. DO NOT EDIT.

package out

import (
	"sync"
	"context"
	"fmt"
)

const alphaSQL = `SELECT 'alpha' as output;`

// Alpha implements Querier.Alpha.
func (q *DBQuerier) Alpha(ctx context.Context) (string, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "Alpha")
	rows, err := q.conn.Query(ctx, alphaSQL)
	if err != nil {
		return "", fmt.Errorf("query Alpha: %w", err)
	}
	fds := rows.FieldDescriptions()
	plan0 := planScan(pgtype.TextCodec{}, fds[0], (*string)(nil))

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (string, error) {
		vals := row.RawValues()
		var item string
		if err := plan0.Scan(vals[0], &item); err != nil {
			return item, fmt.Errorf("scan Alpha.output: %w", err)
		}
		return item, nil
	})
}
