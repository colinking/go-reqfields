package main

import (
	"github.com/colinking/go-reqfields"
	"golang.org/x/tools/go/analysis"
)

//go:generate go build -buildmode=plugin main.go
type analyzerPlugin struct{}

// This must be implemented
func (*analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		reqfields.Analyzer,
	}
}

// Export for golangci-lint.
//
// See: https://golangci-lint.run/contributing/new-linters/#how-to-add-a-private-linter-to-golangci-lint
var AnalyzerPlugin analyzerPlugin //nolint:deadcode
