# lesiw.io/linelen

[![Go Reference](https://pkg.go.dev/badge/lesiw.io/linelen.svg)](https://pkg.go.dev/lesiw.io/linelen)
[![CI](https://github.com/lesiw/linelen/actions/workflows/main.yml/badge.svg?branch=main)](https://github.com/lesiw/linelen/actions/workflows/main.yml)
[![Release](https://img.shields.io/github/v/tag/lesiw/linelen?sort=semver&label=release)](https://github.com/lesiw/linelen/tags)
[![Go Version](https://img.shields.io/github/go-mod/go-version/lesiw/linelen)](../go.mod)
[![Discord](https://img.shields.io/discord/1145827224516300971?logo=discord&logoColor=white&color=5865F2&label=discord)](https://lesiw.dev/discord)
[![License](https://img.shields.io/github/license/lesiw/linelen)](../LICENSE)

An `analysis.Analyzer` for line length.

Width is measured in runes after expanding tabs, so a multibyte
character counts once and a tab counts as its display width.
Carriage returns are stripped before measuring. The defaults — a
limit of 79 and a tab width of 4 — are flags.

## Checks

### Lines over the limit

Lines wider than the limit are reported with their measured width.
A tab counts four columns by default. Positions where wrapping is
the author's layout are exempt: import declarations, and function,
method, and interface method signatures. A signature stays on one
line or wraps as its author chose, while the body is still
checked. Generated files are skipped, since a generated line
cannot be reflowed by hand.

## Usage

```sh
go get -tool lesiw.io/linelen/cmd/linelen
go tool linelen ./... # equivalent to: go tool linelen -len 79 -tab 4 ./...
```
