package main

import (
	"github.com/colinking/go-required"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(reqfields.Analyzer)
}
