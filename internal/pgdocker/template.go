package pgdocker

type pgTemplate struct {
	PGPass      string
	InitScripts []string
}

const dockerfileTemplate = `
{{- /*gotype: github.com/kirk-anchor/pggen/internal/pgdocker.pgTemplate*/ -}}
{{- define "dockerfile" -}}
FROM postgres:13
{{ range .InitScripts }}
COPY {{.}} /docker-entrypoint-initdb.d/
{{ end }}
{{- end }}
`
