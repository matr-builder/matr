package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// Command defines a matr cmd
type Command struct {
	Name       string
	Summary    string
	Doc        string
	Params     []Param
	Returns    []Param
	IsExported bool
	WrapFunc   bool
}

// Param defines a matr HandlerFunc parameter
type Param struct {
	Name string
	Type string
}

// Parse parses a matr file to identify the available handlers
func Parse(file string) ([]Command, error) {
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, file, nil, 4)
	if err != nil {
		panic(err)
	}

	funcs := []Command{}

	if len(f.Comments) == 0 || f.Comments[0].Pos() != 1 || f.Comments[0].Text() != "+build matr\n" {
		return funcs, errors.New("invalid Matrfile: o matr build tag found")
	}

	for _, d := range f.Decls {
		t, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}
		funcs = append(funcs, parseCmd(t))
	}

	return funcs, nil
}

func parseCmd(t *ast.FuncDecl) Command {
	cmd := Command{
		Name:       t.Name.String(),
		Params:     []Param{},
		IsExported: ast.IsExported(t.Name.String()),
	}

	if t.Doc != nil {
		d := []string{}
		for _, ds := range t.Doc.List {

			d = append(d, strings.Replace(ds.Text, "//", "", 1))
		}
		cmd.Summary = d[0]
		cmd.Doc = strings.Join(d, "\n")
	}

	if t.Type == nil || t.Type.Params == nil || len(t.Type.Params.List) == 0 {
		return cmd
	}

	for _, fld := range t.Type.Params.List {
		for _, name := range fld.Names {
			cmd.Params = append(cmd.Params, Param{Name: name.String(), Type: formatType(fld.Type)})
		}
	}

	if t.Type.Results == nil || len(t.Type.Results.List) == 0 {
		return cmd
	}

	for _, fld := range t.Type.Results.List {
		for _, name := range fld.Names {
			cmd.Returns = append(cmd.Returns, Param{Name: name.String(), Type: formatType(fld.Type)})
		}
	}

	return cmd
}

func formatType(t ast.Expr) string {
	switch tt := t.(type) {
	case *ast.StarExpr:
		return "*" + formatType(tt.X)
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", tt.X, tt.Sel)
	case *ast.Ident:
		return tt.String()
	case *ast.Ellipsis:
		return "[]" + formatType(tt.Elt)
	default:
		return ""
	}
}
