package linelen

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "lineLength",
	Doc:  "Reports lines longer than 79 characters",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		if ast.IsGenerated(f) {
			continue
		}
		var pos int
		file := pass.Fset.File(f.FileStart)
		src, err := pass.ReadFile(file.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		for raw := range strings.SplitSeq(string(src), "\n") {
			line := strings.ReplaceAll(raw, "\t", "    ")
			if len(line) > 79 {
				pass.Reportf(
					file.Pos(pos),
					"line is %d characters long, exceeds 79 limit",
					len(line),
				)
			}
			pos += len(raw) + 1
		}
	}
	return nil, nil
}
