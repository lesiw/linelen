package linelen

import (
	"testing"

	"lesiw.io/checker"
)

func TestCheck(t *testing.T) { checker.Run(t, Analyzer) }
