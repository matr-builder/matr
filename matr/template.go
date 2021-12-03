package matr

import (
	"html/template"
	"io"
	"strings"

	"github.com/matr-builder/matr/parser"
)

const defaultTemplate = `//go:build matr

package main

import (
	"context"
	"os"

	"github.com/matr-builder/matr/matr"
)

func main() {
	// Create new Matr instance
	m := matr.New()

	{{- range .}}
	{{if .IsExported }}
	// {{if .Summary}}{{.Summary}}{{else}}{{.Name}}{{end}}
	m.Handle(&matr.Task{
		Name: "{{cmdname .Name}}",
		Summary: "{{trim .Summary}}",
		Doc: ` + "`{{trim .Doc}}`," + `
		Handler: {{.Name}},
	})
	{{- end -}}
	{{- end}}

	// Run Matr
	if err := m.Run(context.Background(), os.Args[1:]...); err != nil {
		os.Stderr.WriteString("ERROR: "+err.Error()+"\n")
	}
}
`

// generate ...
func generate(cmds []parser.Command, w io.Writer) error {
	// Create a new template and parse the letter into it.
	t := template.Must(template.New("matr").Funcs(template.FuncMap{
		"title": strings.Title,
		"trim":  strings.TrimSpace,
		"cmdname": func(name string) string {
			return parser.LowerFirst(parser.CamelToHyphen(name))
		},
	}).Parse(defaultTemplate))
	return t.Execute(w, cmds)
}
