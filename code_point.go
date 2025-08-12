package csslexer

import (
	"go.baoshuo.dev/cssutil"
)

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#would-start-an-identifier
func (l *Lexer) nextCharsAreIdentifier() bool {
	first := l.r.Peek(0)

	if cssutil.IsIdentStartCodePoint(first) {
		return true
	}

	second := l.r.Peek(1)

	if twoCharsAreValidEscape(first, second) {
		return true
	}

	if first == '-' {
		return cssutil.IsIdentStartCodePoint(second) || second == '-' ||
			twoCharsAreValidEscape(second, l.r.Peek(2))
	}

	return false
}

// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#starts-with-a-number
func (l *Lexer) nextCharsAreNumber() bool {
	first := l.r.Peek(0)

	if cssutil.IsDigit(first) {
		return true
	}

	second := l.r.Peek(1)

	if first == '+' || first == '-' {
		if cssutil.IsDigit(second) {
			return true
		}

		if second == '.' {
			third := l.r.Peek(2)

			if cssutil.IsDigit(third) {
				return true
			}
		}
	}

	if first == '.' {
		return cssutil.IsDigit(second)
	}

	return false
}
