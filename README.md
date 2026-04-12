# Go CSS Lexer

A lexer for CSS (Cascading Style Sheets) files written in Go.

The library implements a tokenizer algorithm inspired by Blink, closely mirroring the parsing logic used by modern browsers.

Target version of CSS syntax: [CSS Syntax Module Level 3 (W3C Candidate Recommendation Draft; Dec 24, 2021)](https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/), which is the latest stable version of the CSS syntax specification as of July 2025.

## Installation

```bash
go get go.baoshuo.dev/csslexer
```

## API

### Input

You can create an `Input` instance using one of the following constructors:

- `NewInput(input string) *Input`
- `NewInputRunes(runes []rune) *Input`
- `NewInputBytes(input []byte) *Input`
- `NewInputReader(r io.Reader) *Input`

The lexer requires an `Input` instance to read the CSS content.

### Lexer

Create a lexer:

```go
lexer := csslexer.NewLexer(input)
```

Read next token:

```go
token := lexer.Next()
```

The types of tokens can be found in the `csslexer.TokenType` type, and the definition of each token type is available in `token.go`.

## Benchmark

```bash
go test -benchmem -bench=^BenchmarkLexer$ "go.baoshuo.dev/csslexer/bench"
```

Result:

```text
goos: darwin
goarch: arm64
pkg: go.baoshuo.dev/csslexer/bench
cpu: Apple M4 Pro
BenchmarkLexer/960.css.gz-14                 15025       79833 ns/op    119.2 MiB/s      33280 B/op      2718 allocs/op
BenchmarkLexer/animate.css.gz-14              1701      710401 ns/op    127.4 MiB/s     310769 B/op     30864 allocs/op
BenchmarkLexer/blueprint.css.gz-14            8978      141014 ns/op    117.7 MiB/s      68976 B/op      5634 allocs/op
BenchmarkLexer/bootstrap-4.css.gz-14           750     1485375 ns/op    128.4 MiB/s     722180 B/op     62280 allocs/op
BenchmarkLexer/bootstrap-5.css.gz-14           789     1526705 ns/op    127.8 MiB/s     734099 B/op     63325 allocs/op
BenchmarkLexer/font-awesome-all.css.gz-14     2506      480782 ns/op    144.4 MiB/s     252929 B/op     23377 allocs/op
BenchmarkLexer/foundation.css.gz-14            856     1324012 ns/op    120.9 MiB/s     627897 B/op     52686 allocs/op
BenchmarkLexer/inuit.css.gz-14                4260      287830 ns/op    175.3 MiB/s     120128 B/op      8499 allocs/op
BenchmarkLexer/normalize.css.gz-14           29696       41324 ns/op    141.0 MiB/s      13200 B/op       892 allocs/op
BenchmarkLexer/pure.css.gz-14                 6192      190956 ns/op    143.8 MiB/s      84456 B/op      6585 allocs/op
BenchmarkLexer/reset.css.gz-14               95690       12467 ns/op    82.99 MiB/s       5504 B/op       422 allocs/op
BenchmarkLexer/tailwind.css.gz-14               75    15526441 ns/op    146.0 MiB/s    7082370 B/op    622999 allocs/op
```

## Author

**go-css-lexer** © [Baoshuo](https://baoshuo.ren), Released under the [MIT](./LICENSE) License.

> [Personal Homepage](https://baoshuo.ren) · [Blog](https://blog.baoshuo.ren) · GitHub [@renbaoshuo](https://github.com/renbaoshuo)
