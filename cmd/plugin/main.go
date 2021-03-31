package main

import (
	"fmt"

	reqfields "github.com/colinking/go-required"
	"golang.org/x/tools/go/analysis"
)

//go:generate go build -buildmode=plugin main.go
type analyzerPlugin struct{}

// This must be implemented
func (*analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	fmt.Printf("Returning analyzer for go-required!\n")

	return []*analysis.Analyzer{
		reqfields.Analyzer,
	}
}

// Export for golangci-lint.
//
// See: https://golangci-lint.run/contributing/new-linters/#how-to-add-a-private-linter-to-golangci-lint
var AnalyzerPlugin analyzerPlugin
