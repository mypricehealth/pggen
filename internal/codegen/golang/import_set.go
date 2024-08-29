package golang

import (
	"github.com/mypricehealth/pggen/internal/codegen/golang/gotype"
	"golang.org/x/exp/maps"
)

// ImportSet contains a set of imports required by one Go file.
type ImportSet struct {
	imports map[string]struct{}
}

func NewImportSet() *ImportSet {
	return &ImportSet{imports: make(map[string]struct{}, 4)}
}

// AddPackage adds a fully qualified package path to the set, like
// "github.com/mypricehealth/pggen/foo".
func (s *ImportSet) AddPackage(p string) {
	s.imports[p] = struct{}{}
}

// AddPackage removes a fully qualified package path from the set, like
// "github.com/mypricehealth/pggen/foo".
func (s *ImportSet) RemovePackage(p string) {
	delete(s.imports, p)
}

// AddType adds all fully qualified package paths needed for type and any child
// types.
func (s *ImportSet) AddType(typ gotype.Type) {
	importPath := typ.Import()
	if importPath != "" {
		s.AddPackage(typ.Import())
	}

	comp, ok := typ.(*gotype.CompositeType)
	if !ok {
		return
	}
	for _, childType := range comp.FieldTypes {
		s.AddType(childType)
	}
}

// Get returns a new slice containing the packages, suitable
// for an import statement.
func (s *ImportSet) Get() []string {
	return maps.Keys(s.imports)
}
