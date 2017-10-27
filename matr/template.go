package matr

const defaultTemplate = `// +build matr

package main

import (
	"context"
	"log"
	"os"

	"github.com/matr-builder/matr/matr"
)

func main() {
	// Create new Matr instance
	m := matr.New()

	{{- range .}}
	{{if .IsExported }}
	// {{if .Summary}}{{.Summary}}{{else}}{{.Name}}{{end}}
	m.Handle("{{cmdname .Name}}", {{.Name}})
	{{- end -}}
	{{- end}}

	// Run Matr
	if err := m.Run(context.Background(), os.Args[1:]...); err != nil {
		log.Fatal(err)
	}
}
`
