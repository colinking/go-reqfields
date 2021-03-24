package reqfields

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/fatih/structtag"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "reqfields",
	Doc:  "reports missing required struct fields",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	structs := collectStructs(pass)
	fmt.Printf("Found the following structs: %+v\n", structs)

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			cl, ok := n.(*ast.CompositeLit)
			if !ok {
				return true
			}

			// TODO: map this type to a package name
			snid, ok := cl.Type.(*ast.Ident)
			if !ok {
				// TODO: hide this behind a debug flag
				pass.Reportf(cl.Type.Pos(), "expected an identifier")
				return true
			}

			s, ok := structs[snid.Name]
			if !ok {
				// TODO: hide this behind a debug flag
				pass.Reportf(cl.Type.Pos(), "unknown struct: %s", snid.Name)
				return true
			}

			foundFields := map[string]struct{}{}
			for _, e := range cl.Elts {
				kv, ok := e.(*ast.KeyValueExpr)
				if !ok {
					// TODO: hide this behind a debug flag
					pass.Reportf(e.Pos(), "expected a key-value expr")
					continue
				}
				id, ok := kv.Key.(*ast.Ident)
				if !ok {
					// TODO: hide this behind a debug flag
					pass.Reportf(kv.Key.Pos(), "expected an identifier")
					continue
				}

				foundFields[id.Name] = struct{}{}
			}

			missingFields := []string{}
			for _, requiredField := range s.RequiredFields {
				if _, ok := foundFields[requiredField]; !ok {
					missingFields = append(missingFields, requiredField)
				}
			}
			if len(missingFields) > 0 {
				pass.Reportf(cl.Type.Pos(), "%s is missing: %s", s.Name, strings.Join(missingFields, ", "))
			}

			return true
		})
	}

	return nil, nil
}

type Struct struct {
	Package        string
	Name           string
	RequiredFields []string
}

func collectStructs(pass *analysis.Pass) map[string]Struct {
	structs := map[string]Struct{}

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			ts, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}

			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				return true
			}

			// TODO: associate this struct with a package
			fullName := ts.Name.String()
			s := Struct{
				Package: "todo",
				Name:    ts.Name.String(),
			}

			for _, field := range st.Fields.List {
				if field.Tag == nil {
					continue
				}

				if tag, err := structtag.Parse(trimQuotes(field.Tag.Value)); err != nil {
					// TODO: hide this behind a debug flag
					pass.Reportf(field.Tag.Pos(), "unable to parse struct tag (%s): %+v", field.Tag.Value, err)
					continue
				} else if v, _ := tag.Get("required"); v != nil && v.Name == "true" {
					s.RequiredFields = append(s.RequiredFields, field.Names[0].Name)
				}
			}

			structs[fullName] = s

			return true
		})
	}

	return structs
}

func trimQuotes(s string) string {
	if s[0] == '"' || s[0] == '`' {
		s = s[1:]
	}
	if s[len(s)-1] == '"' || s[len(s)-1] == '`' {
		s = s[:len(s)-1]
	}

	return s
}
