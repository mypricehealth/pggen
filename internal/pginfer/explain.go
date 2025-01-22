package pginfer

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/mypricehealth/pggen/internal/ast"
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
	Type     PlanType
	Relation string   // target relation if any
	Outputs  []string // the output expressions if any
}

type ExplainQueryResultRow struct {
	Plan            map[string]interface{} `json:"Plan,omitempty"`
	QueryIdentifier *uint64                `json:"QueryIdentifier,omitempty"`
}

// explainQuery executes explain plan to get the node plan type and the format
// of the output columns.
//
// Expects the prior SQL to have already been run in the current connection.
func (inf *Inferrer) explainQuery(query *ast.SourceQuery) (Plan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	if len(query.ContiguousArgsSQL) == 0 {
		return Plan{}, fmt.Errorf("no sql to explain")
	}

	tx, err := inf.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return Plan{}, fmt.Errorf("begin transaction failure: %w", err)
	}
	defer tx.Rollback(context.WithoutCancel(ctx))

	lastSQL := query.ContiguousArgsSQL[len(query.ContiguousArgsSQL)-1]

	explainQuery := `EXPLAIN (VERBOSE, FORMAT JSON) ` + lastSQL.SQL
	row := inf.conn.QueryRow(ctx, explainQuery, createParamArgs(len(lastSQL.Args))...)
	explain := make([]ExplainQueryResultRow, 0, 1)
	if err := row.Scan(&explain); err != nil {
		return Plan{}, fmt.Errorf("explain prepared query: %w", err)
	}
	if len(explain) == 0 {
		return Plan{}, fmt.Errorf("no explain output")
	}
	plan := explain[0].Plan
	if len(plan) == 0 {
		return Plan{}, fmt.Errorf("explain output had no 'Plan' node")
	}

	// Node type
	node, ok := plan["Node Type"]
	if !ok {
		return Plan{}, fmt.Errorf("explain output had no 'Plan[Node Type]' node")
	}
	strNode, ok := node.(string)
	if !ok {
		return Plan{}, fmt.Errorf("explain output 'Plan[Node Type]' is not string; got type %T for value %v", node, node)
	}

	// Relation
	relation := plan["Relation Name"]
	relationStr, _ := relation.(string)

	// Outputs
	rawOuts := plan["Output"]
	outs, _ := rawOuts.([]interface{})
	strOuts := make([]string, len(outs))
	for i, out := range outs {
		out, ok := out.(string)
		if !ok {
			return Plan{}, fmt.Errorf("explain output 'Plan.Output[%d]' was not a string; got type %T for value %v", i, out, out)
		}
		strOuts[i] = out
	}
	return Plan{
		Type:     PlanType(strNode),
		Relation: relationStr,
		Outputs:  strOuts,
	}, nil
}
