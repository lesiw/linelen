package linelen

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"unicode/utf8"

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
		ignoreRange := func(start, end token.Pos) {
			startLine := pass.Fset.Position(start).Line
			endLine := pass.Fset.Position(end).Line
			for line := startLine; line <= endLine; line++ {
				ignore[line] = struct{}{}
			}
		}
		ast.Inspect(f, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.ImportSpec:
				ignore[pass.Fset.Position(node.Pos()).Line] = struct{}{}
			case *ast.GenDecl:
				if node.Tok.String() == "import" {
					ignoreRange(node.Pos(), node.End())
				}
			case *ast.FuncDecl:
				// Function and method signatures may wrap across lines;
				// the wrap is the user's chosen layout. The body still
				// gets checked.
				end := node.Type.End()
				if node.Body != nil {
					end = node.Body.Lbrace
				}
				ignoreRange(node.Pos(), end)
			case *ast.InterfaceType:
				// Interface methods are *ast.Field entries whose Type
				// is *ast.FuncType. The whole field span is the
				// signature; let it wrap.
				if node.Methods == nil {
					return true
				}
				for _, field := range node.Methods.List {
					if _, ok := field.Type.(*ast.FuncType); !ok {
						continue
					}
					ignoreRange(field.Pos(), field.End())
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
			rawLen := len(raw)
			raw = strings.TrimRight(raw, "\r")
			line := strings.ReplaceAll(raw, "\t", strings.Repeat(" ", flagTab))
			width := utf8.RuneCountInString(line)
			if _, skip := ignore[n]; !skip && width > flagLen {
				pass.Reportf(
					file.Pos(pos),
					"line is %d characters long, exceeds %d limit",
					width, flagLen,
				)
			}
			pos += rawLen + 1 // len("\n")
		}
	}
	return nil, nil
}
