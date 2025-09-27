# lesiw.io/linelen [![Go Reference](https://pkg.go.dev/badge/lesiw.io/linelen.svg)](https://pkg.go.dev/lesiw.io/linelen)

An `analysis.Analyzer` for line length.
Defaults to a max line length of 79 and a tab-to-space ratio of 4.

## Usage (test)

```go
package main

import (
	"testing"

	"lesiw.io/checker"
    "lesiw.io/linelen"
)

func TestCheck(t *testing.T) { checker.Run(t, linelen.Analyzer) }
```

## Usage (CLI)

```sh
go get -tool lesiw.io/linelen/cmd/linelen
go tool linelen ./... # equivalent to: go tool linelen -len 79 -tab 4 ./...
```
