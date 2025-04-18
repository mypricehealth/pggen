package golang

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/mypricehealth/pggen/internal/ast"
	"github.com/mypricehealth/pggen/internal/casing"
	"github.com/mypricehealth/pggen/internal/codegen"
	"github.com/mypricehealth/pggen/internal/codegen/golang/gotype"
	"github.com/mypricehealth/pggen/internal/gomod"
)

// Templater creates query file templates.
type Templater struct {
	caser            casing.Caser
	resolver         TypeResolver
	pkg              string // Go package name
	inlineParamCount int
}

// TemplaterOpts is options to control the template logic.
type TemplaterOpts struct {
	Caser    casing.Caser
	Resolver TypeResolver
	Pkg      string // Go package name
	// How many params to inline when calling querier methods.
	InlineParamCount int
}

func NewTemplater(opts TemplaterOpts) Templater {
	return Templater{
		pkg:              opts.Pkg,
		caser:            opts.Caser,
		resolver:         opts.Resolver,
		inlineParamCount: opts.InlineParamCount,
	}
}

// TemplateAll creates query template files for each codegen.QueryFile.
func (tm Templater) TemplateAll(files []codegen.QueryFile) ([]TemplatedFile, error) {
	goQueryFiles := make([]TemplatedFile, 0, len(files))
	allDeclarers := NewDeclarerSet()

	// Pick leader file to define common structs and interfaces via Declarer.
	firstIndex := -1
	firstName := string(unicode.MaxRune)
	for i, f := range files {
		if f.SourcePath < firstName {
			firstIndex = i
			firstName = f.SourcePath
		}
	}

	for i, queryFile := range files {
		isLeader := i == firstIndex
		goFile, decls, err := tm.templateFile(queryFile, isLeader)
		if err != nil {
			return nil, fmt.Errorf("template query file %s for go: %w", queryFile.SourcePath, err)
		}
		goQueryFiles = append(goQueryFiles, goFile)
		allDeclarers.AddAll(decls.ListAll()...)
	}

	// Add declarers to leader file.
	goQueryFiles[firstIndex].Declarers = allDeclarers.ListAll()

	// Remove unneeded pgconn import if possible.
	for _, file := range goQueryFiles {
		if !file.needsPgconnImport() {
			file.Imports.RemovePackage("github.com/jackc/pgx/v5/pgconn")
		}

		if file.needsPGXImport() {
			file.Imports.AddPackage("github.com/jackc/pgx/v5")
		}

		if file.needsPgtypeImport() {
			file.Imports.AddPackage("github.com/jackc/pgx/v5/pgtype")
		}
	}

	// Remove self imports.
	for _, file := range goQueryFiles {
		selfPkg, err := gomod.GuessPackage(file.SourcePath)
		if err != nil || selfPkg == "" {
			continue // ignore error, assume it's not a self import
		}

		file.Imports.RemovePackage(selfPkg)
	}
	return goQueryFiles, nil
}

// templateFile creates the data needed to build a Go file for a query file.
// Also returns any declarations needed by this query file. The caller must
// dedupe declarations.
func (tm Templater) templateFile(file codegen.QueryFile, isLeader bool) (TemplatedFile, DeclarerSet, error) {
	imports := NewImportSet()
	imports.AddPackage("context")
	imports.AddPackage("fmt")
	imports.AddPackage("github.com/jackc/pgx/v5/pgconn")
	if isLeader {
		imports.AddPackage("github.com/jackc/pgx/v5/pgtype")
		imports.AddPackage("github.com/jackc/pgx/v5")
		imports.AddPackage("sync")
	}

	pkgPath := ""
	// NOTE: err == nil check
	// Attempt to guess package path. Ignore error if it doesn't work because
	// resolving the package isn't perfect. We'll fall back to an unqualified
	// type which will likely work since the type is probably declared in this
	// package.
	if pkg, err := gomod.GuessPackage(file.SourcePath); err == nil {
		pkgPath = pkg
	}

	queries := make([]TemplatedQuery, 0, len(file.Queries))
	declarers := NewDeclarerSet()
	for _, query := range file.Queries {
		// Build doc string.
		docs := strings.Builder{}
		avgCharsPerLine := 40
		docs.Grow(len(query.Doc) * avgCharsPerLine)
		for i, d := range query.Doc {
			if i > 0 {
				docs.WriteByte('\t') // first line is already indented in the template
			}
			docs.WriteString("// ")
			docs.WriteString(d)
			docs.WriteRune('\n')
		}

		// Build inputs.
		inputs := make([]TemplatedParam, len(query.Inputs))
		for i, input := range query.Inputs {
			goType, err := tm.resolver.Resolve(input.PgType, input.IsOptional, pkgPath, false)
			if err != nil {
				return TemplatedFile{}, nil, err
			}

			imports.AddType(goType)
			inputs[i] = TemplatedParam{
				UpperName: tm.chooseUpperName(input.PgName, "UnnamedParam", i, len(query.Inputs)),
				LowerName: tm.chooseLowerName(input.PgName, "unnamedParam", i, len(query.Inputs)),
				QualType:  gotype.QualifyType(goType, pkgPath),
				Type:      goType,
				RawName:   input,
			}
			ds := FindInputDeclarers(goType).ListAll()
			declarers.AddAll(ds...)
		}

		// Build outputs.
		outputs := make([]TemplatedColumn, len(query.Outputs))
		for i, out := range query.Outputs {
			goType, err := tm.resolver.Resolve(out.PgType, out.Nullable, pkgPath, true)
			if err != nil {
				return TemplatedFile{}, nil, err
			}
			imports.AddType(goType)
			outputs[i] = TemplatedColumn{
				PgName:    out.PgName,
				UpperName: tm.chooseUpperName(out.PgName, "UnnamedColumn", i, len(query.Outputs)),
				LowerName: tm.chooseLowerName(out.PgName, "UnnamedColumn", i, len(query.Outputs)),
				Type:      goType,
				QualType:  gotype.QualifyType(goType, pkgPath),
			}
			ds := FindOutputDeclarers(goType).ListAll()
			declarers.AddAll(ds...)
		}

		nonVoidCols := removeVoidColumns(outputs)
		resultKind := query.ResultKind
		if len(nonVoidCols) == 0 && resultKind != ast.ResultKindString {
			resultKind = ast.ResultKindExec
		}

		var sqlVarName string
		if resultKind == ast.ResultKindString {
			// Export the raw SQL if the result kind is `:string`
			sqlVarName = tm.caser.ToUpperGoIdent(query.Name) + "SQL"
		} else {
			sqlVarName = tm.caser.ToLowerGoIdent(query.Name) + "SQL"
		}

		queries = append(queries, TemplatedQuery{
			Name:             tm.caser.ToUpperGoIdent(query.Name),
			SQLVarName:       sqlVarName,
			ResultKind:       resultKind,
			Doc:              docs.String(),
			PreparedSQL:      query.PreparedSQL,
			Inputs:           inputs,
			Outputs:          nonVoidCols,
			ScanCols:         outputs,
			InlineParamCount: tm.inlineParamCount,
		})
	}

	return TemplatedFile{
		PkgPath:    pkgPath,
		GoPkg:      tm.pkg,
		SourcePath: file.SourcePath,
		Queries:    queries,
		Imports:    imports,
		IsLeader:   isLeader,
	}, declarers, nil
}

// chooseUpperName converts pgName into a capitalized Go identifier name.
// If it's not possible to convert pgName into an identifier, uses fallback with
// a suffix using idx.
func (tm Templater) chooseUpperName(pgName string, fallback string, idx int, numOptions int) string {
	if name := tm.caser.ToUpperGoIdent(pgName); name != "" {
		return name
	}
	suffix := strconv.Itoa(idx)
	if numOptions > 9 {
		suffix = fmt.Sprintf("%2d", idx)
	}
	return fallback + suffix
}

// chooseLowerName converts pgName into an uncapitalized Go identifier name.
// If it's not possible to convert pgName into an identifier, uses fallback with
// a suffix using idx.
func (tm Templater) chooseLowerName(pgName string, fallback string, idx int, numOptions int) string {
	if name := tm.caser.ToLowerGoIdent(pgName); name != "" {
		return name
	}
	suffix := strconv.Itoa(idx)
	if numOptions > 9 {
		suffix = fmt.Sprintf("%2d", idx)
	}
	return fallback + suffix
}

// removeVoidColumns makes a copy of cols with all VoidType columns removed.
// Useful because return types shouldn't contain the void type, but we need
// to use a nil placeholder for void types when scanning a pgx.Row.
func removeVoidColumns(cols []TemplatedColumn) []TemplatedColumn {
	outs := make([]TemplatedColumn, 0, len(cols))
	for _, col := range cols {
		if _, ok := col.Type.(*gotype.VoidType); ok {
			continue
		}
		outs = append(outs, col)
	}
	return outs
}
