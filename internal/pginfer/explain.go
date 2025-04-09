package pginfer

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// PlanType is the top-level node plan type that Postgres plans for executing
// query. https://www.postgresql.org/docs/13/executor.html
type PlanType string

const (
	PlanResult      PlanType = "Result"      // select statement
	PlanLimit       PlanType = "Limit"       // select statement with a limit
	PlanModifyTable PlanType = "ModifyTable" // update, insert, or delete statement
)

// Plan is the plan output from an EXPLAIN query.
type Plan struct {
	Type      PlanType `json:"Node Type"`
	Operation string   `json:"Operation"`     // the operation, set for `ModifyTable`
	Relation  string   `json:"Relation Name"` // target relation if any
	Outputs   []string `json:"Output"`        // the output expressions if any
}

type ExplainQueryResultRow struct {
	Plan            Plan    `json:"Plan,omitempty"`
	QueryIdentifier *uint64 `json:"QueryIdentifier,omitempty"`
}

// explainQuery executes explain plan to get the node plan type and the format
// of the output columns.
func (inf *Inferrer) explainQuery(query string, argCount int) (Plan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	tx, err := inf.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return Plan{}, fmt.Errorf("begin transaction failure: %w", err)
	}
	defer tx.Rollback(context.WithoutCancel(ctx))

	explainQuery := `EXPLAIN (VERBOSE, FORMAT JSON) ` + query
	row := inf.conn.QueryRow(ctx, explainQuery, createParamArgs(argCount)...)
	explain := make([]ExplainQueryResultRow, 0, 1)
	if err := row.Scan(&explain); err != nil {
		return Plan{}, fmt.Errorf("explain prepared query: %w", err)
	}

	if len(explain) != 1 {
		return Plan{}, fmt.Errorf("expected 1 plan but got %d plans", len(explain))
	}

	return explain[0].Plan, nil
}
