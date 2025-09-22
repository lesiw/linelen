package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"
	"lesiw.io/linelen"
)

func main() { singlechecker.Main(linelen.Analyzer) }
