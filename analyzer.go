package reqfields

import (
	"go/ast"
	"reflect"
	"strings"

	"github.com/fatih/structtag"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "reqfields",
	Doc:  Doc,
	// Requires:   []*analysis.Analyzer{inspect.Analyzer},
	Run:       runf,
	FactTypes: []analysis.Fact{&Struct{}},
}

const Doc = `check for missing required struct fields

By default, all struct fields are assumed to be optional. A field
can be marked as required using a struct tag:

	type Hat struct {
		Style string ` + "`" + `require:"true"` + "`" + `
	}
`

// TODO: better name
type Struct struct {
	RequiredFields []string
}

func (this *Struct) AFact() {}

func runf(pass *analysis.Pass) (interface{}, error) {
	findStructs(pass)
	checkStructs(pass)

	return nil, nil
}

func findStructs(pass *analysis.Pass) {
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

			s := Struct{}
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

			obj, ok := pass.TypesInfo.Defs[ts.Name]
			if !ok {
				pass.ReportRangef(ts, "unknown definition %s", ts.Name.String())
				return true
			}
			if !pass.ImportObjectFact(obj, &s) {
				pass.ExportObjectFact(obj, &s)
			}

			return true
		})
	}
}

func checkStructs(pass *analysis.Pass) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			cl, ok := n.(*ast.CompositeLit)
			if !ok {
				return true
			}

			checkStruct(pass, cl, nil)

			return true
		})
	}
}

func identFromType(pass *analysis.Pass, cl *ast.CompositeLit, parentIdent *ast.Ident) *ast.Ident {
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
		case *ast.MapType, *ast.StructType, *ast.InterfaceType:
			// Skip untyped maps/structs.
			return nil
		case nil:
			if parentIdent != nil {
				return parentIdent
			}
			pass.ReportRangef(cl, "ERROR: got nil type")
		default:
			pass.ReportRangef(cl, "ERROR: unexpected expression type, got %s", reflect.TypeOf(expr))
			return nil
		}
	}

	return ident
}

func checkStruct(pass *analysis.Pass, cl *ast.CompositeLit, parentIdent *ast.Ident) {
	ident := identFromType(pass, cl, parentIdent)
	if ident == nil {
		return
	}

	foundFields := map[string]struct{}{}
	for _, e := range cl.Elts {
		switch v := e.(type) {
		case *ast.KeyValueExpr:
			id, ok := v.Key.(*ast.Ident)
			if !ok {
				pass.ReportRangef(v.Key, "ERROR: expected kv.Key to be an identifier, got %s", reflect.TypeOf(v.Key))
				continue
			}

			foundFields[id.Name] = struct{}{}
		case *ast.CompositeLit:
			// checkStruct(pass, v, ident)
		// TODO: support unnamed struct parameters
		case *ast.SelectorExpr, *ast.UnaryExpr, *ast.CallExpr, *ast.Ident:
			continue
		default:
			pass.ReportRangef(e, "ERROR: expected a key-value expr, got %s", reflect.TypeOf(v))
		}
	}

	obj, ok := pass.TypesInfo.Uses[ident]
	if !ok {
		pass.ReportRangef(cl, "ERROR: unknown identifier: %s", ident.Name)
		return
	}
	// TODO: figure out how to filter out non-struct identifiers:
	switch ident.Name {
	case "uintptr", "string":
		return
	}
	var s Struct
	if !pass.ImportObjectFact(obj, &s) {
		pass.ReportRangef(cl, "ERROR: no fact for: %s (type=%s)", ident.Name, obj.Type())
		return
	}

	missingFields := []string{}
	for _, requiredField := range s.RequiredFields {
		if _, ok := foundFields[requiredField]; !ok {
			missingFields = append(missingFields, requiredField)
		}
	}
	if len(missingFields) > 0 {
		pass.ReportRangef(cl, "%s is missing: %s", ident.Name, strings.Join(missingFields, ", "))
	}
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
