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

## Author

**go-css-lexer** © [Baoshuo](https://baoshuo.ren), Released under the [MIT](./LICENSE) License.

> [Personal Homepage](https://baoshuo.ren) · [Blog](https://blog.baoshuo.ren) · GitHub [@renbaoshuo](https://github.com/renbaoshuo)
