// Package linelen provides an analysis.Analyzer for line length,
// and the line-length rules themselves for other tools to share.
package linelen

import (
	"bytes"
	"cmp"
	"flag"
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// Config configures the line-length rules. The zero value uses the
// defaults: a 79 column limit and 4 column tabs.
type Config struct {
	Len int // maximum line width in visual columns
	Tab int // visual width of a tab
}

// len and tab resolve unset fields to the defaults, so cmp.Or
// appears here and nowhere else.
func (c Config) len() int { return cmp.Or(c.Len, 79) }
func (c Config) tab() int { return cmp.Or(c.Tab, 4) }

// Flags registers the len and tab flags on fs, bound to c, so every
// command sharing the rules exposes identical knobs.
func (c *Config) Flags(fs *flag.FlagSet) {
	fs.IntVar(&c.Len, "len", c.len(), "maximum line length")
	fs.IntVar(&c.Tab, "tab", c.tab(), "tab width in spaces")
}

// Width returns the visual width of line: a tab counts Tab columns,
// every other rune one. A trailing carriage return does not count.
func (c Config) Width(line []byte) (width int) {
	tab := c.tab()
	line = bytes.TrimSuffix(line, []byte("\r"))
	for _, r := range string(line) {
		if r == '\t' {
			width += tab
		} else {
			width++
		}
	}
	return width
}

// Line describes a line of a source file.
type Line struct {
	N     int       // 1-based line number
	Width int       // visual width in columns
	Pos   token.Pos // position of the line's first byte
}

// Check returns f's overlong lines: every line of src whose width
// exceeds the limit and which Checked reports as measured.
func (c Config) Check(fset *token.FileSet, f *ast.File, src []byte) (long []Line) {
	var (
		n, pos  int
		checked = Checked(fset, f)
		file    = fset.File(f.FileStart)
		limit   = c.len()
	)
	for line := range bytes.SplitSeq(src, []byte("\n")) {
		n++
		width := c.Width(line)
		if checked(n) && width > limit {
			long = append(long, Line{N: n, Width: width, Pos: file.Pos(pos)})
		}
		pos += len(line) + 1 // len("\n")
	}
	return long
}

// Checked returns a predicate reporting whether the limit applies
// to a line of f. It does not apply to import declarations, and to
// function, method, and interface method signatures, whose wrapping
// is the author's chosen layout; a function's body is measured.
func Checked(fset *token.FileSet, f *ast.File) func(line int) bool {
	var (
		skip = make(map[int]struct{})
		mark = func(start, end token.Pos) {
			s, e := fset.Position(start).Line, fset.Position(end).Line
			for line := s; line <= e; line++ {
				skip[line] = struct{}{}
			}
		}
	)
	ast.Inspect(f, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.ImportSpec:
			mark(node.Pos(), node.Pos())
		case *ast.GenDecl:
			if node.Tok.String() == "import" {
				mark(node.Pos(), node.End())
			}
		case *ast.FuncDecl:
			// Function and method signatures may wrap across lines;
			// the wrap is the user's chosen layout. The body still
			// gets checked.
			end := node.Type.End()
			if node.Body != nil {
				end = node.Body.Lbrace
			}
			mark(node.Pos(), end)
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
				mark(field.Pos(), field.End())
			}
		}
		return true
	})
	return func(line int) bool {
		_, skipped := skip[line]
		return !skipped
	}
}

var (
	cfg      Config
	Analyzer = &analysis.Analyzer{
		Name: "linelen",
		Doc:  "reports lines longer than the specified character limit",
		Run:  run,
	}
)

func init() {
	cfg.Flags(&Analyzer.Flags)
}

func run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		if ast.IsGenerated(f) {
			continue
		}
		file := pass.Fset.File(f.FileStart)
		src, err := pass.ReadFile(file.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		for _, l := range cfg.Check(pass.Fset, f, src) {
			pass.Reportf(l.Pos,
				"line is %d characters long, exceeds %d limit",
				l.Width, cfg.len(),
			)
		}
	}
	return nil, nil
}
