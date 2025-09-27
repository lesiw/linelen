package linelen

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var (
	flagLen = 79
	flagTab = 4
)

var Analyzer = &analysis.Analyzer{
	Name: "linelen",
	Doc:  "reports lines longer than the specified character limit",
	Run:  run,
}

func init() {
	Analyzer.Flags.IntVar(&flagLen, "len", flagLen, "maximum line length")
	Analyzer.Flags.IntVar(&flagTab, "tab", flagTab, "tab width in spaces")
}

func run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		if ast.IsGenerated(f) {
			continue
		}
		ignore := make(map[int]struct{})
		ast.Inspect(f, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.ImportSpec:
				ignore[pass.Fset.Position(node.Pos()).Line] = struct{}{}
			case *ast.GenDecl:
				if node.Tok.String() == "import" {
					pos := pass.Fset.Position(node.Pos())
					end := pass.Fset.Position(node.End())
					for line := pos.Line; line <= end.Line; line++ {
						ignore[line] = struct{}{}
					}
				}
			}
			return true
		})
		file := pass.Fset.File(f.FileStart)
		src, err := pass.ReadFile(file.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		var n, pos int
		for raw := range strings.SplitSeq(string(src), "\n") {
			n++
			line := strings.ReplaceAll(raw, "\t", strings.Repeat(" ", flagTab))
			if _, skip := ignore[n]; !skip && len(line) > flagLen {
				pass.Reportf(
					file.Pos(pos),
					"line is %d characters long, exceeds %d limit",
					len(line), flagLen,
				)
			}
			pos += len(raw) + 1 // len("\n")
		}
	}
	return nil, nil
}
