# lesiw.io/linelen [![Go Reference](https://pkg.go.dev/badge/lesiw.io/linelen.svg)](https://pkg.go.dev/lesiw.io/linelen)

An `analysis.Analyzer` for line length.
Currently hardcoded to a max line length of 79 and a tab-to-space ratio of 4.

## Usage

```go
package main

import (
	"testing"

	"lesiw.io/checker"
    "lesiw.io/linelen"
)

func TestCheck(t *testing.T) { checker.Run(t, linelen.Analyzer) }
```

## To do
* Configurable line length and tab-to-space ratio.
* [Tests!](https://pkg.go.dev/golang.org/x/tools/go/analysis/analysistest)
