package linelen

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"lesiw.io/checker"
)

func TestCheck(t *testing.T) { checker.Run(t, Analyzer) }

// The table is written as plain literals so a case reads as the line it
// describes, which puts it over the limit the table itself tests.
//
//ignore:linelen
var checkLineTests = []struct {
	name  string
	line  string
	allow bool
}{
	{"within the limit", `x := "short"`, true},
	{"string is continuous", `x := "the quick brown fox"`, true},
	{"string is not continuous", `x := "short" + abcdefgh`, false},
	{"several strings that each fall short", `x := "abc" + "def" + gh`, false},
	{"url is continuous", `// https://example.com`, true},
	{"url is not continuous", `// see https://e and more`, false},
	{"string and url each fall short", `x := "abc" + https://e`, false},
	{"url inside a string is counted once", `x := "ab://cd" + more`, false},
}

func TestCheckLine(t *testing.T) {
	c := Config{Len: 20}
	for _, tt := range checkLineTests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.CheckLine([]byte(tt.line)); got != tt.allow {
				t.Errorf("CheckLine(%q) = %v, want %v", tt.line, got, tt.allow)
			}
		})
	}
}

func TestAnalysisTest(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analyzer,
		"basic", "longimportblock", "longimportstmt",
		"methodwrap", "continuous",
	)
}
