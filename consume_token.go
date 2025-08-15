package csslexer

import (
	"strings"

	"go.baoshuo.dev/cssutil"
)

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#consume-numeric-token
func (l *Lexer) consumeNumericToken() (TokenType, string) {
	number := l.consumeNumber()

	if l.nextCharsAreIdentifier() {
		unit := l.consumeName()

		return DimensionToken, number + unit
	} else if l.r.Peek(0) == '%' {
		l.r.Move(1) // consume '%'

		return PercentageToken, number + "%"
	}

	return NumberToken, number
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#urange
func (l *Lexer) consumeUnicodeRangeToken() (TokenType, string) {
	// range start
	start_length_remaining := 6
	for next := l.r.Peek(0); start_length_remaining > 0 && next != EOF && cssutil.IsHexDigit(next); next = l.r.Peek(0) {
		l.r.Move(1) // consume the hex digit
		start_length_remaining--
	}

	if start_length_remaining > 0 && l.r.Peek(0) == '?' { // wildcard range
		for start_length_remaining > 0 && l.r.Peek(0) == '?' {
			l.r.Move(1) // consume the '?'
			start_length_remaining--
		}
	} else if l.r.Peek(0) == '-' && cssutil.IsHexDigit(l.r.Peek(1)) { // range end
		l.r.Move(1) // consume the '-'

		end_length_remaining := 6
		for next := l.r.Peek(0); end_length_remaining > 0 && next != EOF && cssutil.IsHexDigit(next); next = l.r.Peek(0) {
			l.r.Move(1) // consume the hex digit
			end_length_remaining--
		}
	}

	return UnicodeRangeToken, l.r.CurrentString()
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#consume-ident-like-token
func (l *Lexer) consumeIdentLikeToken() (TokenType, string) {
	name := l.consumeName()

	if l.r.Peek(0) == '(' {
		l.r.Move(1) // consume the opening parenthesis
		if strings.ToLower(name) == "url" {
			// The spec is slightly different so as to avoid dropping whitespace
			// tokens, but they wouldn't be used and this is easier.
			l.consumeWhitespace()

			next := l.r.Peek(0)
			if next != '"' && next != '\'' {
				return l.consumeURLToken()
			}
		}

		return FunctionToken, name
	}

	return IdentToken, name
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#consume-string-token
func (l *Lexer) consumeStringToken() (TokenType, string) {
	var result strings.Builder

	until := l.r.Peek(0) // the opening quote, already checked valid by the caller
	l.r.Move(1)

	for {
		next := l.r.Peek(0)

		if next == until {
			l.r.Move(1)
			return StringToken, result.String()
		}

		if next == EOF {
			return StringToken, result.String()
		}

		if cssutil.IsNewline(next) {
			return BadStringToken, result.String()
		}

		if next == '\\' {
			next_next := l.r.Peek(1)

			if next_next == EOF {
				l.r.Move(1) // consume the backslash
				continue
			}

			if cssutil.IsNewline(next_next) {
				l.r.Move(1)
				l.consumeSingleWhitespace()
			} else if twoCharsAreValidEscape(next, next_next) {
				l.r.Move(1) // consume the backslash
				result.WriteRune(l.consumeEscape())
			} else {
				result.WriteRune(next)
				l.r.Move(1)
			}
		} else {
			result.WriteRune(next)
			l.r.Move(1) // consume the current rune
		}
	}
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#consume-url-token
func (l *Lexer) consumeURLToken() (TokenType, string) {
	var result strings.Builder

	for {
		next := l.r.Peek(0)

		if next == ')' {
			l.r.Move(1)
			return UrlToken, result.String()
		}

		if next == EOF {
			return UrlToken, result.String()
		}

		if cssutil.IsWhitespace(next) {
			l.consumeWhitespace()

			next_next := l.r.Peek(0)
			if next_next == ')' {
				l.r.Move(1) // consume the closing parenthesis
				return UrlToken, result.String()
			}
			if next_next == EOF {
				return UrlToken, result.String()
			}

			// If the next character is not a closing parenthesis, there's an error and we should mark it as a bad URL token.
			break
		}

		if next == '"' || next == '\'' || next == '(' || cssutil.IsNonPrintableCodePoint(next) {
			l.r.Move(1) // consume the invalid character
			break
		}

		if next == '\\' {
			if twoCharsAreValidEscape(next, l.r.Peek(1)) {
				l.r.Move(1) // consume the backslash
				result.WriteRune(l.consumeEscape())
				continue
			} else {
				break
			}
		}

		result.WriteRune(next)
		l.r.Move(1) // consume the current rune
	}

	l.consumeBadUrlRemnants()
	return BadUrlToken, ""
}
