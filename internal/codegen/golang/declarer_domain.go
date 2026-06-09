package golang

import (
	"strconv"
	"strings"

	"github.com/mypricehealth/pggen/internal/codegen/golang/gotype"
)

// DomainTypeDeclarer declares a new string type and the const values to map to a
// Postgres enum.
type DomainTypeDeclarer struct {
	domain *gotype.DomainType
}

func NewDomainTypeDeclarer(enum *gotype.DomainType) DomainTypeDeclarer {
	return DomainTypeDeclarer{domain: enum}
}

func (e DomainTypeDeclarer) DedupeKey() string {
	return "domain_type::" + e.domain.PgDomain.Name
}

func (e DomainTypeDeclarer) Declare(pkgPath string) (string, error) {
	sb := &strings.Builder{}

	sb.WriteString("// ")
	sb.WriteString(e.domain.Name)
	sb.WriteString(" represents the Postgres domain ")
	sb.WriteString(strconv.Quote(e.domain.PgDomain.Name))
	sb.WriteString(".\n")

	sb.WriteString("type ")
	sb.WriteString(e.domain.Name)
	sb.WriteString(" ")
	qualType := gotype.QualifyType(e.domain.Elem, pkgPath)
	sb.WriteString(qualType)
	sb.WriteString("\n\n")

	return sb.String(), nil
}
