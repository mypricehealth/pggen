package pginfer

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jackc/pgx/v5"
	"github.com/mypricehealth/pggen/internal/ast"
	"github.com/mypricehealth/pggen/internal/pg"
)

const defaultTimeout = 3 * time.Second

// TypedQuery is an enriched form of ast.SourceQuery after running it on
// Postgres to get information about the ast.SourceQuery.
type TypedQuery struct {
	// Name of the query, from the comment preceding the query. Like 'FindAuthors'
	// in the source SQL: "-- name: FindAuthors :many"
	Name string
	// The result output kind, :one, :many, or :exec.
	ResultKind ast.ResultKind
	// The comment lines preceding the query, without the SQL comment syntax and
	// excluding the :name line.
	Doc []string
	// The SQL query, with pggen functions replaced with Postgres syntax. Ready
	// to run on Postgres with the PREPARE statement.
	ContiguousArgsSQL []ast.ContiguousArgsSQL
	PreparedSQL       string
	// The input parameters to the query.
	Inputs []InputParam
	// The output columns of the query.
	Outputs []OutputColumn
	// Qualified protocol buffer message type to use for each output row, like
	// "erp.api.Product". If empty, generate our own Row type.
	ProtobufType string
}

// InputParam is an input parameter for a prepared query.
type InputParam struct {
	// Name of the param, like 'FirstName' in pggen.arg('FirstName').
	PgName string
	// The postgres type of this param as reported by Postgres.
	PgType pg.Type
	// Whether the input parameter is optional
	IsOptional bool
}

func (p InputParam) IsEmpty() bool {
	return p.PgName == "" && p.PgType == nil && !p.IsOptional
}

// OutputColumn is a single column output from a select query or returning
// clause in an update, insert, or delete query.
type OutputColumn struct {
	// Name of an output column, named by Postgres, like "foo" in "SELECT 1 as foo".
	PgName string
	// The postgres type of the column as reported by Postgres.
	PgType pg.Type
	// If the type can be null; depends on the query. A column defined
	// with a NOT NULL constraint can still be null in the output with a left
	// join. Nullability is determined using rudimentary control-flow analysis.
	Nullable bool
}

type Inferrer struct {
	conn        *pgx.Conn
	typeFetcher *pg.TypeFetcher
}

// NewInferrer infers information about a query by running the query on
// Postgres and extracting information from the catalog tables.
func NewInferrer(conn *pgx.Conn) *Inferrer {
	return &Inferrer{
		conn:        conn,
		typeFetcher: pg.NewTypeFetcher(conn),
	}
}

func (inf *Inferrer) RunSetup(sql string) (pgconn.CommandTag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	return inf.conn.Exec(ctx, sql)
}

func (inf *Inferrer) InferTypes(query *ast.SourceQuery) (TypedQuery, error) {
	inputs, outputs, err := inf.prepareTypes(query)
	if err != nil {
		return TypedQuery{}, fmt.Errorf("infer output types for query: %w", err)
	}
	if query.ResultKind != ast.ResultKindExec && query.ResultKind != ast.ResultKindSetup && query.ResultKind != ast.ResultKindString {
		if len(outputs) == 0 {
			return TypedQuery{}, fmt.Errorf(
				"query %s has incompatible result kind %s; the query doesn't return any columns; "+
					"use :exec, :setup, or :string if the query shouldn't return any columns",
				query.Name, query.ResultKind)
		}
		if countVoids(outputs) == len(outputs) {
			return TypedQuery{}, fmt.Errorf(
				"query %s has incompatible result kind %s; the query only has void columns; "+
					"use :exec, :setup, or :string if the query shouldn't return any columns",
				query.Name, query.ResultKind)
		}
	}
	doc := extractDoc(query)
	return TypedQuery{
		Name:              query.Name,
		ResultKind:        query.ResultKind,
		Doc:               doc,
		ContiguousArgsSQL: query.ContiguousArgsSQL,
		PreparedSQL:       query.PreparedSQL,
		Inputs:            inputs,
		Outputs:           outputs,
		ProtobufType:      query.Pragmas.ProtobufType,
	}, nil
}

func (inf *Inferrer) prepareTypes(query *ast.SourceQuery) (_a []InputParam, _ []OutputColumn, mErr error) {
	// Execute the query to get field descriptions of the output columns.
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	tx, err := inf.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("begin transaction failure: %w", err)
	}
	defer tx.Rollback(context.WithoutCancel(ctx))

	needsOutputPlan := query.ResultKind == ast.ResultKindMany || query.ResultKind == ast.ResultKindOne || query.ResultKind == ast.ResultKindRows

	statementsData, err := inf.getStatementsData(ctx, query, needsOutputPlan)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get statement data: %w", err)
	}

	inputParams, err := inf.getInputParamTypes(statementsData.inputStatements, query)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get inputs: %w", err)
	}

	if len(inputParams) != len(query.Params) {
		return nil, nil, fmt.Errorf("expected %d parameter types for query; got %d", len(query.Params), len(inputParams))
	}

	var outputColumns []OutputColumn
	if len(statementsData.inputStatements) > 0 {
		// The output is based upon the final statement.
		finalStatement := statementsData.inputStatements[len(statementsData.inputStatements)-1]

		// Resolve type names of output column data type OIDs.
		outputOIDs := make([]uint32, len(finalStatement.Fields))
		for i, desc := range finalStatement.Fields {
			outputOIDs[i] = desc.DataTypeOID
		}
		outputTypes, err := inf.typeFetcher.FindTypesByOIDs(outputOIDs...)
		if err != nil {
			return nil, nil, fmt.Errorf("fetch oid types: %w", err)
		}

		// Output nullability.
		nullables, err := inf.inferOutputNullability(query, finalStatement.Fields, statementsData.outputPlan)
		if err != nil {
			return nil, nil, fmt.Errorf("infer output type nullability: %w", err)
		}

		// Create output columns
		for i, desc := range finalStatement.Fields {
			pgType, ok := outputTypes[uint32(desc.DataTypeOID)]
			if !ok {
				return nil, nil, fmt.Errorf("no postgrestype name found for column %s with oid %d", string(desc.Name), desc.DataTypeOID)
			}
			outputColumns = append(outputColumns, OutputColumn{
				PgName:   string(desc.Name),
				PgType:   pgType,
				Nullable: nullables[i],
			})
		}
	}

	return inputParams, outputColumns, nil
}

type statementsData struct {
	inputStatements []*pgconn.StatementDescription
	outputPlan      Plan
}

func (inf *Inferrer) getStatementsData(ctx context.Context, query *ast.SourceQuery, needsOutputPlan bool) (statementsData, error) {
	// The final statement description is used for the output, if any.
	inputStatements := make([]*pgconn.StatementDescription, 0, len(query.ContiguousArgsSQL))
	var outputPlan Plan

	for i, sql := range query.ContiguousArgsSQL {
		stmtDesc, err := inf.prepareQuery(ctx, sql.SQL, query.ResultKind)
		if err != nil {
			return statementsData{}, fmt.Errorf("failed to prepare query: %w", err)
		}

		if len(stmtDesc.ParamOIDs) != sql.UniqueArgs {
			fundamentalError := fmt.Errorf("%d unique args are written in the query but got oids for %d", sql.UniqueArgs, len(stmtDesc.ParamOIDs))
			if len(stmtDesc.ParamOIDs) < sql.UniqueArgs {
				return statementsData{}, fmt.Errorf("%s, this likely indicates that an argument can be parsed but not bound, for example `CREATE VIEW pggen.arg('view_name') AS (SELECT 1)` would result in this error", fundamentalError)
			}

			return statementsData{}, fmt.Errorf("%s, this either indicates a bug or a positional parameter being used directly in a query instead of `pggen.arg('arg_name')`", fundamentalError)
		}

		inputStatements = append(inputStatements, stmtDesc)

		shouldExecute := true

		isLastStatement := i == len(query.ContiguousArgsSQL)-1

		// If this is the last statement, there's no need to actually execute it.
		if isLastStatement {
			shouldExecute = false
		}

		if isLastStatement {
			if needsOutputPlan {
				plan, err := inf.explainQuery(sql.SQL, len(sql.Args))
				if err != nil {
					return statementsData{}, fmt.Errorf("could not explain query: %w", err)
				}

				outputPlan = plan
			}
		}

		if shouldExecute {
			err := inf.execute(ctx, stmtDesc, sql)
			if err != nil {
				return statementsData{}, fmt.Errorf("could not execute query: %w", err)
			}
		}

	}

	data := statementsData{inputStatements, outputPlan}

	return data, nil
}

// Executes the query so that side effects that the next query may depend on are applied.
func (inf *Inferrer) execute(ctx context.Context, stmtDesc *pgconn.StatementDescription, sql ast.ContiguousArgsSQL) error {
	paramValues := make([][]byte, len(sql.Args))
	paramFormats := make([]int16, len(sql.Args))
	for i := range len(sql.Args) {
		paramValues[i] = nil
	}

	rr := inf.conn.PgConn().ExecPrepared(ctx, stmtDesc.Name, paramValues, paramFormats, nil)
	_, err := rr.Close()
	if err != nil {
		return fmt.Errorf("failed to execute prepared query: %w", err)
	}

	return nil
}

func (inf *Inferrer) prepareQuery(ctx context.Context, sql string, resultKind ast.ResultKind) (*pgconn.StatementDescription, error) {
	// If paramOIDs is null, Postgres infers the type for each parameter.
	var paramOIDs []uint32
	stmtDesc, err := inf.conn.PgConn().Prepare(ctx, "", sql, paramOIDs)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			msg := "fetch field descriptions: " + pgErr.Message
			if pgErr.Where != "" {
				msg += "\n    WHERE: " + pgErr.Where
			}
			if pgErr.Detail != "" {
				msg += "\n    DETAIL: " + pgErr.Detail
			}
			if pgErr.Hint != "" {
				msg += "\n    HINT: " + pgErr.Hint
			}
			if pgErr.DataTypeName != "" {
				msg += "\n    DataType: " + pgErr.DataTypeName
			}
			if pgErr.TableName != "" {
				msg += "\n    TableName: " + pgErr.TableName
			}
			// Provide hint to use a returning clause. pggen ignores most errors but
			// only if there's output columns. If the user has an UPDATE or INSERT
			// without a RETURNING clause, pggen will surface the null constraint
			// errors because len(descriptions) == 0.
			if strings.Contains(strings.ToLower(sql), "update") ||
				strings.Contains(strings.ToLower(sql), "insert") {
				msg += "\n    HINT: if the main statement is an UPDATE or INSERT ensure that you have"
				msg += "\n          a RETURNING clause (this query is marked " + string(resultKind) + ")."
				msg += "\n          Use :exec if you don't need the query output."
			}
			return nil, fmt.Errorf(msg+"\n    %w", pgErr)
		}
		return nil, fmt.Errorf("prepare query to infer types: %w", err)
	}

	return stmtDesc, nil
}

func (inf *Inferrer) getInputParamTypes(statements []*pgconn.StatementDescription, query *ast.SourceQuery) ([]InputParam, error) {
	var oids []uint32
	oidToParamIndex := make(map[uint32][]int)

	setParams := make([]bool, len(query.Params))

	if len(statements) != len(query.ContiguousArgsSQL) {
		return nil, fmt.Errorf("expected %d statement descriptions; got %d", len(query.ContiguousArgsSQL), len(statements))
	}

	for i, statement := range statements {
		contiguousArgs := query.ContiguousArgsSQL[i]

		if len(statement.ParamOIDs) != contiguousArgs.UniqueArgs {
			return nil, fmt.Errorf("expected %d parameter oids; got %d", contiguousArgs.UniqueArgs, len(statement.ParamOIDs))
		}

		// To be analyzed, the parameters are required to be dense, i.e. `$1`, `$2`, `$3` etc.
		// However across multiple statements an individual statement might only have say `$4` and `$2`.
		for j, oid := range statement.ParamOIDs {
			if _, ok := oidToParamIndex[oid]; !ok {
				oids = append(oids, oid)
			}

			paramIndex := contiguousArgs.Args[j]
			indexes := oidToParamIndex[oid]

			// This slice should be small enough that a linear search is fine (and even faster than a set).
			// A query with 100 parameters is already quite large.
			if !slices.Contains(indexes, paramIndex) {
				oidToParamIndex[oid] = append(indexes, paramIndex)

				if paramIndex >= len(setParams) {
					return nil, fmt.Errorf("has input params count %d but got parameter index %d", len(setParams), paramIndex)
				}

				setParams[paramIndex] = true
			}
		}
	}

	unsetParams := &strings.Builder{}
	for i, set := range setParams {
		if !set {
			if len(unsetParams.String()) > 0 {
				unsetParams.WriteString(", ")
			}
			unsetParams.WriteString(strconv.Itoa(i))
		}
	}

	if unsetParams.Len() > 0 {
		return nil, fmt.Errorf("cannot find oids for parameters at index %s", unsetParams.String())
	}

	paramTypes := make([]InputParam, len(setParams))
	if len(oids) > 0 {
		types, err := inf.typeFetcher.FindTypesByOIDs(oids...)
		if err != nil {
			return nil, fmt.Errorf("fetch oid types: %w", err)
		}

		for _, oid := range oids {
			paramIndexes := oidToParamIndex[oid]
			for _, paramIndex := range paramIndexes {
				param := query.Params[paramIndex]

				inputType, ok := types[oid]
				if !ok {
					return nil, fmt.Errorf("no postgres type name found for parameter %s with oid %d", param.Name, oid)
				}

				oldParam := paramTypes[paramIndex]

				newParam := InputParam{
					PgName:     param.Name,
					PgType:     inputType,
					IsOptional: param.IsOptional,
				}
				paramTypes[paramIndex] = newParam

				if oldParam.IsEmpty() {
					continue
				}

				currentOID := newParam.PgType.OID()
				oldOID := oldParam.PgType.OID()
				if newParam.PgName != oldParam.PgName || currentOID != oldOID {
					return nil, fmt.Errorf("parameter %s has conflicting types: %s (OID %d) and %s (OID %d), note that this analysis may be overly pessimistic", param.Name, newParam.PgType, currentOID, oldParam.PgType, oldOID)
				}

				paramTypes[paramIndex].IsOptional = oldParam.IsOptional && newParam.IsOptional
			}
		}
	}

	return paramTypes, nil
}

// inferOutputNullability infers which of the output columns produced by the
// query and described by descs can be null.
func (inf *Inferrer) inferOutputNullability(query *ast.SourceQuery, descs []pgconn.FieldDescription, outputPlan Plan) ([]bool, error) {
	if len(descs) == 0 {
		return nil, nil
	}

	columnKeys := make([]pg.ColumnKey, len(descs))
	for i, desc := range descs {
		if desc.TableOID > 0 {
			columnKeys[i] = pg.ColumnKey{
				TableOID: uint32(desc.TableOID),
				Number:   desc.TableAttributeNumber,
			}
		}
	}
	cols, err := pg.FetchColumns(inf.conn, columnKeys)
	if err != nil {
		return nil, fmt.Errorf("fetch column for nullability: %w", err)
	}

	// The nth entry determines if the output column described by descs[n] is
	// nullable. plan.Outputs might contain more entries than cols because the
	// plan output also contains information like sort columns.
	nullables := make([]bool, len(descs))
	for i := range nullables {
		nullables[i] = true // assume nullable until proven otherwise
	}
	for i, col := range cols {
		if i == len(outputPlan.Outputs) {
			// plan.Outputs might not have the same output because the top level node
			// joins child outputs like with append.
			break
		}
		nullables[i] = isColNullable(query, outputPlan, outputPlan.Outputs[i], col)
	}
	return nullables, nil
}

func createParamArgs(argCount int) []interface{} {
	args := make([]interface{}, argCount)
	for i := range argCount {
		args[i] = nil
	}
	return args
}

func extractDoc(query *ast.SourceQuery) []string {
	if query.Doc == nil || len(query.Doc.List) <= 1 {
		return nil
	}
	// Drop last line, like: "-- name: Foo :exec"
	lines := make([]string, len(query.Doc.List)-1)
	for i := range lines {
		comment := query.Doc.List[i].Text
		// TrimLeft to remove runs of dashes. TrimPrefix only removes fixed number.
		noDashes := strings.TrimLeft(comment, "-")
		lines[i] = strings.TrimSpace(noDashes)
	}
	return lines
}

func countVoids(outputs []OutputColumn) int {
	n := 0
	for _, out := range outputs {
		if _, ok := out.PgType.(pg.VoidType); ok {
			n++
		}
	}
	return n
}
