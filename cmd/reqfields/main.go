package main

import (
	"github.com/colinking/go-reqfields"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(reqfields.Analyzer)
}
