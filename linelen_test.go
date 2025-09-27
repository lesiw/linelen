package linelen

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
	"lesiw.io/checker"
)

func TestCheck(t *testing.T) { checker.Run(t, Analyzer) }

func TestAnalysisTest(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analyzer,
		"basic", "longimportblock", "longimportstmt")
}
