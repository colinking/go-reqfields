package reqfields

import (
	"fmt"
	"go/ast"
	"reflect"
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
	// Since we are doing two AST passes in the same run, we will only have
	// access to information about the current package. To support reporting
	// across other packages, we'll need to move this to a separate analyzer
	// run and persist this data as object facts that we can query below.
	structs := collectStructs(pass)
	fmt.Printf("Found the following structs: %+v\n", structs)

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			cl, ok := n.(*ast.CompositeLit)
			if !ok {
				return true
			}

			// Traverse the composite literal's type until we get to an identifier.
			var ident *ast.Ident
			expr := cl.Type
			for ident == nil {
				switch next := expr.(type) {
				case *ast.Ident:
					ident = next
				case *ast.SelectorExpr:
					ident = next.Sel
				case *ast.StarExpr:
					expr = next.X
				case *ast.ArrayType:
					expr = next.Elt
				case *ast.MapType, *ast.StructType:
					// Skip untyped maps/structs.
					return true
				default:
					pass.ReportRangef(cl.Type, "ERROR: unexpected expression type, got %s", reflect.TypeOf(expr))
					return true
				}
			}

			s, ok := structs[ident.Name]
			if !ok {
				pass.ReportRangef(cl.Type, "ERROR: unknown struct: %s", ident.Name)
				return true
			}

			foundFields := map[string]struct{}{}
			for _, e := range cl.Elts {
				kv, ok := e.(*ast.KeyValueExpr)
				if !ok {
					pass.ReportRangef(e, "ERROR: expected a key-value expr")
					continue
				}
				id, ok := kv.Key.(*ast.Ident)
				if !ok {
					pass.ReportRangef(kv.Key, "ERROR: expected kv.Key to be an identifier, got %s", reflect.TypeOf(kv.Key))
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
				pass.ReportRangef(cl.Type, "%s is missing: %s", s.Name, strings.Join(missingFields, ", "))
			}

			return true
		})
	}

	return nil, nil
}

type Struct struct {
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

			fullName := ts.Name.String()
			s := Struct{
				Name: ts.Name.String(),
			}

			for _, field := range st.Fields.List {
				if field.Tag == nil {
					continue
				}

				if tag, err := structtag.Parse(trimQuotes(field.Tag.Value)); err != nil {
					pass.ReportRangef(field.Tag, "unable to parse struct tag (%s): %+v", field.Tag.Value, err)
					continue
				} else if v, _ := tag.Get("require"); v != nil && v.Name == "true" {
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
